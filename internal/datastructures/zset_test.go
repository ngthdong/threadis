package datastructures_test

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ngthdong/threadis/internal/datastructures"
)

func setupZSet() *datastructures.ZSet {
	rand.Seed(1) // deterministic
	return datastructures.CreateZSet()
}

func TestZSetAddAndLength(t *testing.T) {
	zs := setupZSet()

	n := zs.Add(1, "a")
	assert.Equal(t, 1, n)

	n = zs.Add(2, "b")
	assert.Equal(t, 1, n)

	n = zs.Add(3, "c")
	assert.Equal(t, 1, n)

	assert.Equal(t, 3, zs.Length())
}

func TestZSetAddDuplicate(t *testing.T) {
	zs := setupZSet()

	zs.Add(1, "a")

	// same score → no change
	n := zs.Add(1, "a")
	assert.Equal(t, 1, n)

	// different score → update
	n = zs.Add(2, "a")
	assert.Equal(t, 1, n)

	ok, score := zs.GetScore("a")
	assert.Equal(t, 1, ok)
	assert.Equal(t, 2.0, score)
}

func TestZSetRank(t *testing.T) {
	zs := setupZSet()

	zs.Add(1, "a")
	zs.Add(2, "b")
	zs.Add(3, "c")

	rank, _ := zs.GetRank("a", false)
	assert.Equal(t, int64(0), rank)

	rank, _ = zs.GetRank("b", false)
	assert.Equal(t, int64(1), rank)

	rank, _ = zs.GetRank("c", false)
	assert.Equal(t, int64(2), rank)
}

func TestZSetRankReverse(t *testing.T) {
	zs := setupZSet()

	zs.Add(1, "a")
	zs.Add(2, "b")
	zs.Add(3, "c")

	rank, _ := zs.GetRank("a", true)
	assert.Equal(t, int64(2), rank)

	rank, _ = zs.GetRank("c", true)
	assert.Equal(t, int64(0), rank)
}

func TestZSetGetScore(t *testing.T) {
	zs := setupZSet()

	zs.Add(10, "x")

	ok, score := zs.GetScore("x")
	assert.Equal(t, 1, ok)
	assert.Equal(t, 10.0, score)

	ok, score = zs.GetScore("y")
	assert.Equal(t, -1, ok)
	assert.Equal(t, 0.0, score)
}

func TestZSetNonExistingRank(t *testing.T) {
	zs := setupZSet()

	zs.Add(1, "a")

	rank, score := zs.GetRank("not_exist", false)
	assert.Equal(t, int64(-1), rank)
	assert.Equal(t, 0.0, score)
}

func TestZSetUpdateReorder(t *testing.T) {
	zs := setupZSet()

	zs.Add(1, "a")
	zs.Add(2, "b")
	zs.Add(3, "c")

	// update b 
	zs.Add(5, "b")

	rank, _ := zs.GetRank("b", false)
	assert.Equal(t, int64(2), rank)

	rank, _ = zs.GetRank("b", true)
	assert.Equal(t, int64(0), rank)
}

func TestZSetEmptyElement(t *testing.T) {
	zs := setupZSet()

	n := zs.Add(1, "")
	assert.Equal(t, 0, n)

	assert.Equal(t, 0, zs.Length())
}