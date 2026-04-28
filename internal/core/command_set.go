package core

import (
	"errors"

	"github.com/ngthdong/threadis/internal/datastructures"
)

func cmdSADD(args []string) []byte {
	if len(args) < 2 {
		return Encode(errors.New("ERR wrong number of arguments for 'sadd' command"), false)
	}

	key := args[0] 
	set, exist := setStore[key]
	if !exist {
		set = datastructures.NewSet(key)
		setStore[key] = set
	}

	count := set.Add(args[1:]...)
	return Encode(count, false)
}

func cmdSREM(args []string) []byte {
	if len(args) < 2 {
		return Encode(errors.New("ERR wrong number of arguments for 'srem' command"), false)
	}

	key := args[0]
	set, exist := setStore[key]
	if !exist {
		set = datastructures.NewSet(key)
		setStore[key] = set
	}

	count := set.Rem(args[1:]...)
	return Encode(count, false)
}

func cmdSMEMBERS(args []string) []byte {
	if len(args) != 1 {
		return Encode(errors.New("ERR wrong number of arguments for 'smembers' command"), false)
	}

	key := args[0]
	set, exist := setStore[key]
	if !exist {
		return Encode(make([]string, 0), false)
	}

	return Encode(set.Members(), false)
}

func cmdSISMEMBER(args []string) []byte {
	if len(args) != 2 {
		return Encode(errors.New("ERR wrong number of arguments for 'sismember' command"), false)
	}

	key := args[0]
	set, exist := setStore[key]
	if !exist {
		return Encode(0, false)
	}

	return Encode(set.IsMember(args[1]), false)
}
