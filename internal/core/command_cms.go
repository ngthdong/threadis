package core

import (
	"errors"
	"math"
	"strconv"

	"github.com/ngthdong/threadis/internal/constant"
	"github.com/ngthdong/threadis/internal/datastructures"
)

func cmdCMSINITBYDIM(args []string) []byte {
	if len(args) != 3 {
		return Encode(errors.New("ERR wrong number of arguments for 'CMS.INITBYDIM' command"), false)
	}

	key := args[0]

	width, err := strconv.ParseUint(args[1], 10, 32)
	if err != nil {
		return Encode(errors.New("ERR width must be an integer"), false)
	}

	depth, err := strconv.ParseUint(args[2], 10, 32)
	if err != nil {
		return Encode(errors.New("ERR depth must be an integer"), false)
	}

	if _, exist := cmsStore[key]; exist {
		return Encode(errors.New("ERR CMS key already exists"), false)
	}

	cmsStore[key] = datastructures.NewCMS(uint32(width), uint32(depth))
	return constant.RespOk
}

func cmdCMSINITBYPROB(args []string) []byte {
	if len(args) != 3 {
		return Encode(errors.New("ERR wrong number of arguments for 'CMS.INITBYPROB' command"), false)
	}

	key := args[0]

	errRate, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		return Encode(errors.New("ERR invalid error rate"), false)
	}
	if errRate <= 0 || errRate >= 1 {
		return Encode(errors.New("ERR error rate must be between 0 and 1"), false)
	}

	prob, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		return Encode(errors.New("ERR invalid probability"), false)
	}
	if prob <= 0 || prob >= 1 {
		return Encode(errors.New("ERR probability must be between 0 and 1"), false)
	}

	if _, exist := cmsStore[key]; exist {
		return Encode(errors.New("ERR CMS key already exists"), false)
	}

	w, d := datastructures.CalcCMSDim(errRate, prob)
	cmsStore[key] = datastructures.NewCMS(w, d)

	return constant.RespOk
}

func cmdCMSINCRBY(args []string) []byte {
	if len(args) < 3 || len(args)%2 == 0 {
		return Encode(errors.New("ERR wrong number of arguments for 'CMS.INCRBY' command"), false)
	}

	key := args[0]

	cms, exist := cmsStore[key]
	if !exist {
		return Encode(errors.New("ERR CMS key does not exist"), false)
	}

	var res []string

	for i := 1; i < len(args); i += 2 {
		item := args[i]

		value, err := strconv.ParseUint(args[i+1], 10, 32)
		if err != nil {
			return Encode(errors.New("ERR increment must be a non-negative integer"), false)
		}

		count := cms.IncrBy(item, uint32(value))

		if count == math.MaxUint32 {
			res = append(res, "ERR CMS overflow")
			continue
		}

		res = append(res, strconv.FormatUint(uint64(count), 10))
	}

	return Encode(res, false)
}

func cmdCMSQUERY(args []string) []byte {
	if len(args) < 2 {
		return Encode(errors.New("ERR wrong number of arguments for 'CMS.QUERY' command"), false)
	}

	key := args[0]

	cms, exist := cmsStore[key]
	if !exist {
		return Encode(errors.New("ERR CMS key does not exist"), false)
	}

	var res []string

	for i := 1; i < len(args); i++ {
		res = append(res, strconv.FormatUint(uint64(cms.Count(args[i])), 10))
	}

	return Encode(res, false)
}
