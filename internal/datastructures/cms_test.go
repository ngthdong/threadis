package datastructures_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ngthdong/threadis/internal/datastructures"
)

func setupCMS() *datastructures.CMS {
	return datastructures.NewCMS(1000, 5)
}

func TestCMS_IncrAndCount(t *testing.T) {
	cms := setupCMS()

	cms.IncrBy("apple", 1)
	cms.IncrBy("apple", 2)
	cms.IncrBy("banana", 1)

	apple := cms.Count("apple")
	banana := cms.Count("banana")

	assert.GreaterOrEqual(t, apple, uint32(3))
	assert.GreaterOrEqual(t, banana, uint32(1))
}

func TestCMS_MultipleKeys(t *testing.T) {
	cms := setupCMS()

	cms.IncrBy("a", 1)
	cms.IncrBy("b", 2)
	cms.IncrBy("c", 3)

	assert.GreaterOrEqual(t, cms.Count("a"), uint32(1))
	assert.GreaterOrEqual(t, cms.Count("b"), uint32(2))
	assert.GreaterOrEqual(t, cms.Count("c"), uint32(3))
}

func TestCMS_ZeroValue(t *testing.T) {
	cms := setupCMS()

	count := cms.Count("not_exist")

	assert.Equal(t, uint32(0), count)
}

func TestCMS_IncrByReturnValue(t *testing.T) {
	cms := setupCMS()

	v := cms.IncrBy("x", 5)
	assert.GreaterOrEqual(t, v, uint32(5))

	v = cms.IncrBy("x", 5)
	assert.GreaterOrEqual(t, v, uint32(10))
}

func TestCMS_NoUnderestimate(t *testing.T) {
	cms := setupCMS()

	trueCount := uint32(0)

	for i := 0; i < 100; i++ {
		cms.IncrBy("key", 1)
		trueCount++
	}

	estimate := cms.Count("key")

	assert.GreaterOrEqual(t, estimate, trueCount)
}

func TestCMS_OverflowProtection(t *testing.T) {
	cms := setupCMS()

	cms.IncrBy("big", math.MaxUint32-10)

	cms.IncrBy("big", 100)

	count := cms.Count("big")

	assert.Equal(t, uint32(math.MaxUint32), count)
}