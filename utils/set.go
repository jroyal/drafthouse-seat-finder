package utils

type StringSet struct {
	elements map[string]struct{}
}

func NewStringSet() *StringSet {
	elem := make(map[string]struct{})
	s := StringSet{elem}
	return &s
}

func (s *StringSet) Add(elem string) {
	s.elements[elem] = struct{}{}
}

func (s *StringSet) AddSlice(elems []string) {
	for _, elem := range elems {
		s.elements[elem] = struct{}{}
	}
}

func (s *StringSet) Contains(elem string) bool {
	_, exists := s.elements[elem]
	return exists
}

func (s *StringSet) ToSlice() []string {
	i := 0
	elements := make([]string, len(s.elements))
	for k := range s.elements {
		elements[i] = k
		i++
	}
	return elements
}
