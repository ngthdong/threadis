package core

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

const CRLF string = "\r\n"

var respNil = []byte("$-1\r\n")

// +OK\r\n => OK, 5
func readSimpleString(data []byte) (string, int, error) {
	pos := 1
	for data[pos] != '\r' {
		pos += 1
	}
	return string(data[1:pos]), pos + 2, nil
}

// :123\r\n => 123
func readInt64(data []byte) (int64, int, error) {
	var res int64 = 0
	pos := 1
	sign := int64(1)

	if data[pos] == '-' {
		sign = -sign
		pos += 1
	} else if data[pos] == '+' {
		pos += 1
	}

	for pos < len(data) && data[pos] != '\r' {
		res = res*10 + int64(data[pos]-'0')
		pos += 1
	}

	if pos == len(data) {
		return -1, pos, fmt.Errorf("invalid number")
	}

	return res * sign, pos + 2, nil
}

func readError(data []byte) (string, int, error) {
	return readSimpleString(data)
}

// $5\r\nhello\r\n => 5, 4
func readPrefixLength(data []byte) (int, int, error) {
	res, pos, err := readInt64(data)
	if err != nil {
		return int(res), pos, fmt.Errorf("failed to read prefix length %-v", err)
	}
	return int(res), pos, nil
}

// $5\r\nhello\r\n => "hello"
func readBulkString(data []byte) (string, int, error) {
	length, pos, err := readPrefixLength(data)
	if err != nil {
		return "", pos, fmt.Errorf("failed to read bulk string %-v", err)
	}
	return string(data[pos : pos+length]), pos + length + 2, nil
}

// *2\r\n$5\r\nhello\r\n$5\r\nworld\r\n => {"hello", "world"}
func readArray(data []byte) (interface{}, int, error) {
	arraySize, pos, err := readPrefixLength(data)
	if err != nil {
		return nil, pos, fmt.Errorf("failed to read array %-v", err)
	}

	res := make([]interface{}, arraySize)

	for i := 0; i < arraySize; i++ {
		ans, localPos, err := decodeOne(data[pos:])
		if err != nil {
			return ans, pos, fmt.Errorf("failed to decode data %-v", err)
		}

		res[i] = ans
		pos += localPos
	}

	return res, pos, nil
}

func decodeOne(data []byte) (interface{}, int, error) {
	if len(data) == 0 {
		return nil, 0, fmt.Errorf("empty data")
	}

	switch data[0] {
	case '+':
		return readSimpleString(data)
	case ':':
		return readInt64(data)
	case '-':
		return readError(data)
	case '$':
		return readBulkString(data)
	case '*':
		return readArray(data)
	}

	return nil, 0, fmt.Errorf("unknown RESP type: %q", data[0])
}

func Decode(data []byte) (interface{}, error) {
	res, _, err := decodeOne(data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode data %-v", err)
	}

	return res, nil
}

func encodeStringArray(sa []string) []byte {
	var buf bytes.Buffer

	buf.WriteByte('*')
	buf.WriteString(strconv.Itoa(len(sa)))
	buf.WriteString("\r\n")

	for _, s := range sa {
		buf.WriteByte('$')
		buf.WriteString(strconv.Itoa(len(s)))
		buf.WriteString("\r\n")
		buf.WriteString(s)
		buf.WriteString("\r\n")
	}

	return buf.Bytes()
}

func Encode(value interface{}, isSimpleString bool) []byte {
	switch v := value.(type) {
	case string:
		if isSimpleString {
			var buf bytes.Buffer
			buf.WriteByte('+')
			buf.WriteString(v)
			buf.WriteString(CRLF)
			return buf.Bytes()
		}
		var buf bytes.Buffer
		buf.WriteByte('$')
		buf.WriteString(strconv.Itoa(len(v)))
		buf.WriteString(CRLF)
		buf.WriteString(v)
		buf.WriteString(CRLF)
		return buf.Bytes()

	case int64, int32, int16, int8, int:
		return []byte(fmt.Sprintf(":%d\r\n", v))

	case error:
		var buf bytes.Buffer
		buf.WriteByte('-')
		buf.WriteString(v.Error())
		buf.WriteString(CRLF)
		return buf.Bytes()

	case []string:
		return encodeStringArray(v)

	case [][]string:
		var buf bytes.Buffer

		buf.WriteByte('*')
		buf.WriteString(strconv.Itoa(len(v)))
		buf.WriteString(CRLF)

		for _, sa := range v {
			buf.Write(encodeStringArray(sa))
		}
		return buf.Bytes()

	case []interface{}:
		var buf bytes.Buffer

		buf.WriteByte('*')
		buf.WriteString(strconv.Itoa(len(v)))
		buf.WriteString(CRLF)

		for _, x := range v {
			buf.Write(Encode(x, false))
		}
		return buf.Bytes()

	default:
		return respNil
	}
}

func ParseCmd(data []byte) (*Command, error) {
	value, err := Decode(data)
	if err != nil {
		return nil, err
	}

	array := value.([]interface{})
	tokens := make([]string, len(array))
	for i := range tokens {
		tokens[i] = array[i].(string)
	}
	res := &Command{Cmd: strings.ToUpper(tokens[0]), Args: tokens[1:]}
	return res, nil
}
