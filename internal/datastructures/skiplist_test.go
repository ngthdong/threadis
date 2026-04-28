package datastructures_test

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ngthdong/threadis/internal/datastructures"
)

func setupSkiplist() *datastructures.Skiplist {
	rand.Seed(1) // deterministic test
	return datastructures.CreateSkiplist()
}

func TestInsertAndRank(t *testing.T) {
	sl := setupSkiplist()

	sl.Insert(1, "a")
	sl.Insert(2, "b")
	sl.Insert(3, "c")

	assert.Equal(t, uint32(3), sl.Length())

	assert.Equal(t, uint32(1), sl.GetRank(1, "a"))
	assert.Equal(t, uint32(2), sl.GetRank(2, "b"))
	assert.Equal(t, uint32(3), sl.GetRank(3, "c"))
}

func TestInsertDuplicateScoreLexOrder(t *testing.T) {
	sl := setupSkiplist()

	sl.Insert(1, "b")
	sl.Insert(1, "a")
	sl.Insert(1, "c")

	assert.Equal(t, uint32(1), sl.GetRank(1, "a"))
	assert.Equal(t, uint32(2), sl.GetRank(1, "b"))
	assert.Equal(t, uint32(3), sl.GetRank(1, "c"))
}

func TestUpdateScoreNoReposition(t *testing.T) {
	sl := setupSkiplist()

	sl.Insert(1, "a")
	sl.Insert(2, "b")
	sl.Insert(3, "c")

	node := sl.UpdateScore(2, "b", 2.5)

	assert.Equal(t, 2.5, node.Score())

	// vẫn ở vị trí cũ
	assert.Equal(t, uint32(2), sl.GetRank(2.5, "b"))
}

func TestUpdateScoreWithReposition(t *testing.T) {
	sl := setupSkiplist()

	sl.Insert(1, "a")
	sl.Insert(2, "b")
	sl.Insert(3, "c")

	sl.UpdateScore(2, "b", 4)

	assert.Equal(t, uint32(3), sl.GetRank(4, "b"))
}

func TestDeleteSkiplist(t *testing.T) {
	sl := setupSkiplist()

	sl.Insert(1, "a")
	sl.Insert(2, "b")
	sl.Insert(3, "c")

	n := sl.Delete(2, "b")
	assert.Equal(t, 1, n)

	assert.Equal(t, uint32(2), sl.Length())

	assert.Equal(t, uint32(0), sl.GetRank(2, "b"))
	assert.Equal(t, uint32(1), sl.GetRank(1, "a"))
	assert.Equal(t, uint32(2), sl.GetRank(3, "c"))
}

func TestDeleteNonExisting(t *testing.T) {
	sl := setupSkiplist()

	sl.Insert(1, "a")

	n := sl.Delete(2, "b")
	assert.Equal(t, 0, n)
}

func TestMixedOperations(t *testing.T) {
	sl := setupSkiplist()

	sl.Insert(1, "a")
	sl.Insert(5, "e")
	sl.Insert(3, "c")
	sl.Insert(2, "b")
	sl.Insert(4, "d")

	assert.Equal(t, uint32(5), sl.Length())

	// check order
	assert.Equal(t, uint32(1), sl.GetRank(1, "a"))
	assert.Equal(t, uint32(2), sl.GetRank(2, "b"))
	assert.Equal(t, uint32(3), sl.GetRank(3, "c"))
	assert.Equal(t, uint32(4), sl.GetRank(4, "d"))
	assert.Equal(t, uint32(5), sl.GetRank(5, "e"))

	// delete + update
	sl.Delete(3, "c")
	sl.UpdateScore(4, "d", 0)

	assert.Equal(t, uint32(1), sl.GetRank(0, "d"))
}