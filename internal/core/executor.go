package core

import (
	"errors"
	"syscall"
)

func cmdPING(args []string) []byte {
	var res []byte
	if len(args) > 1 {
		return Encode(errors.New("ERR wrong number of arguments for 'ping' command"), false)
	}

	if len(args) == 0 {
		res = Encode("PONG", true)
	} else {
		res = Encode(args[0], false)
	}
	return res
}

func ExecuteAndResponse(cmd *Command, connFd int) error {
	var res []byte

	switch cmd.Cmd {
	case "PING":
		res = cmdPING(cmd.Args)
	case "SET":
		res = cmdSET(cmd.Args)
	case "GET":
		res = cmdGET(cmd.Args)
	case "TTL":
		res = cmdTTL(cmd.Args)
	case "EXPIRE":
		res = cmdEXPIRE(cmd.Args)
	case "DEL":
		res = cmdDEL(cmd.Args)
	case "EXISTS":
		res = cmdEXISTS(cmd.Args)
	case "ZADD":
		res = cmdZADD(cmd.Args)
	case "ZSCORE":
		res = cmdZSCORE(cmd.Args)
	case "ZRANK":
		res = cmdZRANK(cmd.Args)
	case "ZRANGE":
		res = cmdZRANGE(cmd.Args)
	case "ZREM":
		res = cmdZREM(cmd.Args)
	case "SADD":
		res = cmdSADD(cmd.Args)
	case "SREM":
		res = cmdSREM(cmd.Args)
	case "SMEMBERS":
		res = cmdSMEMBERS(cmd.Args)
	case "SISMEMBER":
		res = cmdSISMEMBER(cmd.Args)
	default:
		res = []byte("-CMD NOT FOUND\r\n")
	}
	_, err := syscall.Write(connFd, res)
	return err
}
