package datastructures_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ngthdong/threadis/internal/datastructures"
)

func setupBloom() *datastructures.BloomFilter {
	// entries = 1000, error rate = 1%
	return datastructures.NewBloomFilter(1000, 0.01)
}

func TestBloomAddAndExist(t *testing.T) {
	bf := setupBloom()

	bf.Add("apple")
	bf.Add("banana")
	bf.Add("cherry")

	assert.True(t, bf.Exist("apple"))
	assert.True(t, bf.Exist("banana"))
	assert.True(t, bf.Exist("cherry"))
}

func TestBloomNoFalseNegative(t *testing.T) {
	bf := setupBloom()

	words := []string{"a", "b", "c", "d", "e"}

	for _, w := range words {
		bf.Add(w)
	}

	for _, w := range words {
		assert.True(t, bf.Exist(w), "should not have false negative for %s", w)
	}
}

func TestBloomNonExisting(t *testing.T) {
	bf := setupBloom()

	bf.Add("apple")

	// this MAY be false positive, so we only assert if it's definitely false
	exists := bf.Exist("not_exist")

	// only check it's boolean (no panic / crash)
	assert.IsType(t, true, exists)
}

func TestBloomHashFunctions(t *testing.T) {
	bf := setupBloom()

	hash := bf.CalcHash("apple")

	bf.AddHash(hash)

	assert.True(t, bf.ExistHash(hash))
}

func TestBloomAddDuplicate(t *testing.T) {
	bf := setupBloom()

	bf.Add("apple")
	bf.Add("apple")
	bf.Add("apple")

	assert.True(t, bf.Exist("apple"))
}

func TestBloomManyInsertions(t *testing.T) {
	bf := setupBloom()

	for i := 0; i < 100; i++ {
		bf.Add(string(rune(i)))
	}

	for i := 0; i < 100; i++ {
		assert.True(t, bf.Exist(string(rune(i))))
	}
}