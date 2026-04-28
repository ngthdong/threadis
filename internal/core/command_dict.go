package core

import (
	"errors"
	"strconv"
	"time"

	"github.com/ngthdong/threadis/internal/constant"
)

func cmdSET(args []string) []byte {
	if len(args) < 2 || len(args) == 3 || len(args) > 4 {
		return Encode(errors.New("ERR wrong number of arguments for 'set' command"), false)
	}

	var key, value string
	var ttlMs int64 = -1

	key, value = args[0], args[1]
	if len(args) == 4 {
		ttlSec, err := strconv.ParseInt(args[3], 10, 64)
		if err != nil {
			return Encode(errors.New("ERR value is not an integer or out of range"), false)
		}
		ttlMs = ttlSec * 1000
	}

	dictStore.Set(key, dictStore.NewObj(key, value, ttlMs))
	return constant.RespOk
}

func cmdGET(args []string) []byte {
	if len(args) != 1 {
		return Encode(errors.New("ERR wrong number of arguments for 'get' command"), false)
	}

	key := args[0]
	obj := dictStore.Get(key)
	if obj == nil {
		return constant.RespNil
	}

	if dictStore.HasExpired(key) {
		return constant.RespNil
	}

	return Encode(obj.Value, false)
}

func cmdTTL(args []string) []byte {
	if len(args) != 1 {
		return Encode(errors.New("ERR wrong number of arguments for 'ttl' command"), false)
	}
	key := args[0]
	obj := dictStore.Get(key)
	if obj == nil {
		return constant.TtlKeyNotExist
	}

	exp, isExpirySet := dictStore.GetExpiry(key)
	if !isExpirySet {
		return constant.TtlKeyExistNoExpire
	}

	now := uint64(time.Now().UnixMilli())

	if exp <= now {
		return constant.TtlKeyNotExist
	}

	remainMs := exp - now

	return Encode(int64(remainMs/1000), false)
}

func cmdEXPIRE(args []string) []byte {
	if len(args) != 2 {
		return Encode(errors.New("ERR wrong number of arguments for 'expire' command"), false)
	}

	key := args[0]
	ttlSec, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return Encode(errors.New("ERR value is not an integer or out of range"), false)
	}

	if dictStore.Get(key) == nil {
		return Encode(int64(0), false)
	}

	dictStore.SetExpiry(key, ttlSec*1000)
	return Encode(int64(1), false)
}

func cmdDEL(args []string) []byte {
	if len(args) < 1 {
		return Encode(errors.New("ERR wrong number of arguments for 'del' command"), false)
	}

	var deleted int64 = 0
	for _, k := range args {
		if dictStore.Del(k) {
			deleted++
		}
	}
	return Encode(deleted, false)
}

func cmdEXISTS(args []string) []byte {
	if len(args) < 1 {
		return Encode(errors.New("ERR wrong number of arguments for 'exists' command"), false)
	}

	var count int64 = 0
	for _, k := range args {
		if dictStore.Get(k) != nil {
			count++
		}
	}
	return Encode(count, false)
}
