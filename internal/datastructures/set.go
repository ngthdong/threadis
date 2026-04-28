package datastructures

type Set struct {
	key string 
	dict map[string]struct{}
}

func NewSet(key string) *Set {
	return &Set{
		key: key,
		dict: make(map[string]struct{}),
	}
}

func (s *Set) Add(members ...string) int {
	countNewMemberes := 0
	for _, mem := range members {
		if _, exists := s.dict[mem]; !exists {
			s.dict[mem] = struct{}{}
			countNewMemberes += 1
		}
	}		
	return 	countNewMemberes
}

func (s *Set) Rem(members ...string) int {
	countRemovedMembers := 0
	for _, mem := range members {
		if _, exists := s.dict[mem]; exists {
			delete(s.dict, mem)
			countRemovedMembers += 1
		}
	}		
	return 	countRemovedMembers
}

func (s *Set) IsMember(member string) int {
	if _, exist := s.dict[member]; exist {
		return 1
	}
	return 0
}

func (s *Set) Members() []string {
	members := make([]string, 0, len(s.dict))
	for k, _ := range s.dict {
		members = append(members, k)
	}
	return members 
}