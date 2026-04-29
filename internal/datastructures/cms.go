package datastructures

import (
	"math"
	"github.com/spaolacci/murmur3"
)

const Log10PointFive = -0.30102999566

type CMS struct {
	width   uint32
	depth   uint32
	counter [][]uint32
}

func NewCMS(width uint32, depth uint32) *CMS {
	cms := &CMS{
		width: width,
		depth: depth,
	}

	cms.counter = make([][]uint32, depth)
	for i := uint32(0); i < depth; i++ {
		cms.counter[i] = make([]uint32, width)
	}

	return cms
}

func CalcCMSDim(errRate float64, errProb float64) (uint32, uint32) {
	width := uint32(math.Ceil(2.0 / errRate))
	depth := uint32(math.Ceil(math.Log10(errProb) / Log10PointFive))
	return width, depth
}

func (cms *CMS) calcHash(item string, seed uint32) uint32 {
	hasher := murmur3.New32WithSeed(seed)
	hasher.Write([]byte(item))
	return hasher.Sum32()
}


func (cms *CMS) IncrBy(item string, value uint32) uint32 {
	var minCount uint32 = math.MaxUint32

	for i := uint32(0); i < cms.depth; i++ {
		hash := cms.calcHash(item, i)
		j := hash % cms.width

		// Safely add the value to prevent overflow.
		if math.MaxUint32 - cms.counter[i][j] < value {
			cms.counter[i][j] = math.MaxUint32
		} else {
			cms.counter[i][j] += value
		}

		if cms.counter[i][j] < minCount {
			minCount = cms.counter[i][j]
		}
	}

	return minCount
}


func (cms *CMS) Count(item string) uint32 {
	var minCount uint32 = math.MaxUint32

	for i := uint32(0); i < cms.depth; i++ {
		hash := cms.calcHash(item, i)
		j := hash % cms.width

		if cms.counter[i][j] < minCount {
			minCount = cms.counter[i][j]
		}
	}

	return minCount
}