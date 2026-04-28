package datastructures_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ngthdong/threadis/internal/datastructures"
)

func TestAdd(t *testing.T) {
	s := datastructures.NewSet("myset")

	n := s.Add("a", "b", "c")
	assert.Equal(t, 3, n)

	// add duplicate
	n = s.Add("a", "b")
	assert.Equal(t, 0, n)

	assert.Equal(t, 1, s.IsMember("a"))
	assert.Equal(t, 1, s.IsMember("b"))
	assert.Equal(t, 1, s.IsMember("c"))
}

func TestRem(t *testing.T) {
	s := datastructures.NewSet("myset")

	s.Add("a", "b", "c")

	n := s.Rem("a", "b")
	assert.Equal(t, 2, n)

	assert.Equal(t, 0, s.IsMember("a"))
	assert.Equal(t, 0, s.IsMember("b"))
	assert.Equal(t, 1, s.IsMember("c"))

	// remove non-existing
	n = s.Rem("x")
	assert.Equal(t, 0, n)
}

func TestIsMember(t *testing.T) {
	s := datastructures.NewSet("myset")

	s.Add("a")

	assert.Equal(t, 1, s.IsMember("a"))
	assert.Equal(t, 0, s.IsMember("b"))
}

func TestMembers(t *testing.T) {
	s := datastructures.NewSet("myset")

	s.Add("a", "b", "c")

	members := s.Members()

	assert.Len(t, members, 3)
	assert.ElementsMatch(t, []string{"a", "b", "c"}, members)
}

func TestAddAfterRemove(t *testing.T) {
	s := datastructures.NewSet("myset")

	s.Add("a")
	s.Rem("a")

	assert.Equal(t, 0, s.IsMember("a"))

	n := s.Add("a")
	assert.Equal(t, 1, n)
	assert.Equal(t, 1, s.IsMember("a"))
}

func TestEmptySet(t *testing.T) {
	s := datastructures.NewSet("myset")

	assert.Equal(t, 0, s.IsMember("a"))

	members := s.Members()
	assert.Len(t, members, 0)

	n := s.Rem("a")
	assert.Equal(t, 0, n)
}