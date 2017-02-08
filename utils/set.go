package utils

type StringSet struct {
	Elements map[string]struct{}
}

func NewStringSet() *StringSet {
	elem := make(map[string]struct{})
	s := StringSet{elem}
	return &s
}

func (s *StringSet) Add(elem string) {
	s.Elements[elem] = struct{}{}
}

func (s *StringSet) AddSlice(elems []string) {
	for _, elem := range elems {
		s.Elements[elem] = struct{}{}
	}
}

func (s *StringSet) ToSlice() []string {
	i := 0
	elements := make([]string, len(s.Elements))
	for k := range s.Elements {
		elements[i] = k
		i++
	}
	return elements
}
