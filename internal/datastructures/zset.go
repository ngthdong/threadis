package datastructures

type ZSet struct {
	zskiplist *Skiplist
	dict      map[string]float64
}

func CreateZSet() *ZSet {
	return &ZSet{
		zskiplist: CreateSkiplist(),
		dict:      map[string]float64{},
	}
}

func (zs *ZSet) Add(score float64, element string) int {
	if len(element) == 0 {
		return 0
	}

	if curScore, exists := zs.dict[element]; exists {
		if curScore != score {
			znode := zs.zskiplist.UpdateScore(curScore, element, score)
			zs.dict[element] = znode.score
		}
		// Return 0 for existing elements (count only new elements)
		return 0
	}

	znode := zs.zskiplist.Insert(score, element)
	zs.dict[element] = znode.score
	return 1
}

func (zs *ZSet) GetRank(ele string, reverse bool) (rank int64, score float64) {
	setSize := zs.zskiplist.length
	score, exist := zs.dict[ele]
	if !exist {
		return -1, 0
	}

	rank = int64(zs.zskiplist.GetRank(score, ele))
	if reverse {
		rank = int64(setSize) - rank
	} else {
		rank--
	}

	return rank, score
}

func (zs *ZSet) GetScore(ele string) (int, float64) {
	score, exist := zs.dict[ele]
	if !exist {
		return -1, 0
	}

	return 1, score
}

func (zs *ZSet) Length() int {
	return len(zs.dict)
}

// Range returns elements in the range [start, stop] with optional scores
// Negative indices are supported (-1 = last element)
// Returns slice of members and slice of scores (if withScores is true)
func (zs *ZSet) Range(start, stop int, withScores bool) ([]string, []float64) {
	size := int(zs.zskiplist.length)

	if size == 0 {
		return []string{}, []float64{}
	}

	// Handle negative indices
	if start < 0 {
		start = size + start
	}
	if stop < 0 {
		stop = size + stop
	}

	// Clamp to valid range
	if start < 0 {
		start = 0
	}
	if start >= size {
		return []string{}, []float64{}
	}
	if stop < 0 {
		return []string{}, []float64{}
	}
	if stop >= size {
		stop = size - 1
	}

	// Ensure start <= stop
	if start > stop {
		return []string{}, []float64{}
	}

	members := []string{}
	scores := []float64{}

	// Traverse skiplist from head
	node := zs.zskiplist.head.levels[0].forward
	for i := 0; i < start && node != nil; i++ {
		node = node.levels[0].forward
	}

	// Collect elements in range
	for i := start; i <= stop && node != nil; i++ {
		members = append(members, node.element)
		if withScores {
			scores = append(scores, node.score)
		}
		node = node.levels[0].forward
	}

	return members, scores
}

// Remove removes one or more members from the sorted set
// Returns the number of members removed
func (zs *ZSet) Remove(members ...string) int {
	removed := 0
	for _, member := range members {
		if score, exists := zs.dict[member]; exists {
			zs.zskiplist.Delete(score, member)
			delete(zs.dict, member)
			removed++
		}
	}
	return removed
}
