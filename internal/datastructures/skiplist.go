package datastructures

import (
	"math"
	"math/rand"
	"strings"
)

const MAX_LEVEL = 32

type Node struct {
	element  string
	score    float64
	backward *Node
	levels   []SkiplistLevel
}

type SkiplistLevel struct {
	forward *Node
	span    uint32
}

type Skiplist struct {
	head   *Node
	tail   *Node
	length uint32
	level  int
}

func (sl *Skiplist) randomLevel() int {
	level := 1
	for rand.Intn(2) == 1 && level < MAX_LEVEL {
		level += 1
	}
	return level
}

func (sl *Skiplist) CreateNode(level int, score float64, element string) *Node {
	return &Node{
		element:  element,
		score:    score,
		backward: nil,
		levels:   make([]SkiplistLevel, level),
	}
}

func CreateSkiplist() *Skiplist {
	sl := Skiplist{
		length: 0,
		level:  1,
	}
	sl.head = sl.CreateNode(MAX_LEVEL, math.Inf(-1), "")
	sl.head.backward = nil
	sl.tail = nil
	return &sl
}

func (sl *Skiplist) Length() uint32 {
	return sl.length
}

func (n *Node) Score() float64 {
	return n.score
}

func (sl *Skiplist) Insert(score float64, element string) *Node {
	update := [MAX_LEVEL]*Node{}
	rank := [MAX_LEVEL]uint32{}
	head := sl.head

	for i := sl.level - 1; i >= 0; i-- {
		if i == sl.level-1 {
			rank[i] = 0
		} else {
			rank[i] = rank[i+1]
		}

		for head.levels[i].forward != nil &&
			(head.levels[i].forward.score < score || (head.levels[i].forward.score == score &&
				strings.Compare(head.levels[i].forward.element, element) == -1)) {
			rank[i] += head.levels[i].span
			head = head.levels[i].forward
		}

		update[i] = head
	}

	level := sl.randomLevel()

	if level > sl.level {
		for i := sl.level; i < level; i++ {
			rank[i] = 0
			update[i] = sl.head
			update[i].levels[i].span = sl.length
		}
		sl.level = level
	}

	newNode := sl.CreateNode(level, score, element)

	for i := 0; i < level; i++ {
		newNode.levels[i].forward = update[i].levels[i].forward
		update[i].levels[i].forward = newNode
		newNode.levels[i].span = update[i].levels[i].span - (rank[0] - rank[i])
		update[i].levels[i].span = rank[0] - rank[i] + 1
	}

	for i := level; i < sl.level; i++ {
		update[i].levels[i].span++
	}

	if update[0] == sl.head {
		newNode.backward = nil
	} else {
		newNode.backward = update[0]
	}

	if newNode.levels[0].forward != nil {
		newNode.levels[0].forward.backward = newNode
	} else {
		sl.tail = newNode
	}

	sl.length++

	return newNode
}

func (sl *Skiplist) GetRank(score float64, element string) uint32 {
	x := sl.head
	var rank uint32 = 0
	for i := sl.level - 1; i >= 0; i-- {
		for x.levels[i].forward != nil && (x.levels[i].forward.score < score ||
			(x.levels[i].forward.score == score &&
				strings.Compare(x.levels[i].forward.element, element) <= 0)) {
			rank += x.levels[i].span
			x = x.levels[i].forward
		}
		if x.score == score && strings.Compare(x.element, element) == 0 {
			return rank
		}
	}
	return 0
}

func (sl *Skiplist) UpdateScore(curScore float64, element string, newScore float64) *Node {
	update := [MAX_LEVEL]*Node{}
	head := sl.head

	for i := sl.level - 1; i >= 0; i-- {
		for head.levels[i].forward != nil &&
			(head.levels[i].forward.score < curScore || (head.levels[i].forward.score == curScore &&
				strings.Compare(head.levels[i].forward.element, element) == -1)) {
			head = head.levels[i].forward
		}
		update[i] = head
	}

	head = head.levels[0].forward
	if (head.backward == nil || head.backward.score < newScore) &&
		(head.levels[0].forward == nil || head.levels[0].forward.score > newScore) {
		head.score = newScore
		return head
	}

	sl.DeleteNode(head, update)
	newNode := sl.Insert(newScore, element)
	return newNode
}

func (sl *Skiplist) DeleteNode(delNode *Node, update [MAX_LEVEL]*Node) {
	for i := 0; i < sl.level; i++ {
		if update[i].levels[i].forward == delNode {
			update[i].levels[i].span += delNode.levels[i].span - 1
			update[i].levels[i].forward = delNode.levels[i].forward
		} else {
			update[i].levels[i].span--
		}
	}

	if delNode.levels[0].forward != nil {
		delNode.levels[0].forward.backward = delNode.backward
	} else {
		sl.tail = delNode.backward
	}

	for sl.level > 1 && sl.head.levels[sl.level-1].forward == nil {
		sl.level--
	}

	sl.length--
}

func (sl *Skiplist) Delete(score float64, element string) int {
	update := [MAX_LEVEL]*Node{}
	head := sl.head

	for i := sl.level - 1; i >= 0; i-- {
		for head.levels[i].forward != nil && (head.levels[i].forward.score < score ||
			(head.levels[i].forward.score == score &&
				strings.Compare(head.levels[i].forward.element, element) == -1)) {
			head = head.levels[i].forward
		}
		update[i] = head
	}

	head = head.levels[0].forward
	if head != nil && head.score == score && strings.Compare(head.element, element) == 0 {
		sl.DeleteNode(head, update)
		return 1
	}

	return 0
}
