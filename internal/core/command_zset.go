package core

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/ngthdong/threadis/internal/constant"
	"github.com/ngthdong/threadis/internal/datastructures"
)

func cmdZADD(args []string) []byte {
	if len(args) < 3 {
		return Encode(errors.New("ERR wrong number of arguments for 'zadd' command"), false)
	}
	key := args[0]
	scoreIndex := 1

	numScoreEleArgs := len(args) - scoreIndex
	if numScoreEleArgs%2 == 1 || numScoreEleArgs == 0 {
		return Encode(errors.New("ERR wrong number of arguments for 'zadd' command"), false)
	}

	zset, exist := zsetStore[key]
	if !exist {
		zset = datastructures.CreateZSet()
		zsetStore[key] = zset
	}

	count := 0
	for i := scoreIndex; i < len(args); i += 2 {
		score, err := strconv.ParseFloat(args[i], 64)
		if err != nil {
			return Encode(errors.New("ERR value is not a valid float"), false)
		}
		member := args[i+1]

		// Add returns 1 if new element, 0 if existing (even if updated)
		if zset.Add(score, member) == 1 {
			count++
		}
	}

	return Encode(int64(count), false)
}

func cmdZSCORE(args []string) []byte {
	if len(args) != 2 {
		return Encode(errors.New("ERR wrong number of arguments for 'zscore' command"), false)
	}

	key, member := args[0], args[1]
	zset, exist := zsetStore[key]
	if !exist {
		return constant.RespNil
	}

	ret, score := zset.GetScore(member)
	if ret == -1 {
		return constant.RespNil
	}

	return Encode(fmt.Sprintf("%f", score), false)
}

func cmdZRANK(args []string) []byte {
	if len(args) != 2 {
		return Encode(errors.New("ERR wrong number of arguments for 'zrank' command"), false)
	}

	key, member := args[0], args[1]
	zset, exist := zsetStore[key]
	if !exist {
		return constant.RespNil
	}

	rank, _ := zset.GetRank(member, false)
	return Encode(rank, false)
}

func cmdZRANGE(args []string) []byte {
	if len(args) < 3 {
		return Encode(errors.New("ERR wrong number of arguments for 'zrange' command"), false)
	}

	key := args[0]
	start, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return Encode(errors.New("ERR value is not an integer or out of range"), false)
	}
	stop, err := strconv.ParseInt(args[2], 10, 64)
	if err != nil {
		return Encode(errors.New("ERR value is not an integer or out of range"), false)
	}

	withScores := false
	if len(args) > 3 {
		if len(args) != 4 {
			return Encode(errors.New("ERR syntax error"), false)
		}
		if strings.ToUpper(args[3]) == "WITHSCORES" {
			withScores = true
		} else {
			return Encode(errors.New("ERR syntax error"), false)
		}
	}

	zset, exist := zsetStore[key]
	if !exist {
		// Return empty array for non-existing key
		return Encode([]interface{}{}, false)
	}

	members, scores := zset.Range(int(start), int(stop), withScores)

	if !withScores {
		// Return array of members
		result := make([]interface{}, len(members))
		for i, member := range members {
			result[i] = member
		}
		return Encode(result, false)
	}

	// Return array of [member1, score1, member2, score2, ...]
	result := make([]interface{}, len(members)*2)
	for i, member := range members {
		result[i*2] = member
		result[i*2+1] = fmt.Sprintf("%f", scores[i])
	}
	return Encode(result, false)
}

func cmdZREM(args []string) []byte {
	if len(args) < 2 {
		return Encode(errors.New("ERR wrong number of arguments for 'zrem' command"), false)
	}

	key := args[0]
	members := args[1:]

	zset, exist := zsetStore[key]
	if !exist {
		return Encode(int64(0), false)
	}

	count := zset.Remove(members...)
	return Encode(int64(count), false)
}
