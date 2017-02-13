package utils

// StringSet is a barebones set implementation for strings
type StringSet struct {
	elements map[string]struct{}
}

// NewStringSet creates a new StringSet with an empty map
func NewStringSet() *StringSet {
	elem := make(map[string]struct{})
	s := StringSet{elem}
	return &s
}

// Add will add an element to the set if the set doesn't currently contain the element
func (s *StringSet) Add(elem string) {
	s.elements[elem] = struct{}{}
}

// AddSlice is a convenience function to loop through a string slice and add the elements of
// the slice to the set
func (s *StringSet) AddSlice(elems []string) {
	for _, elem := range elems {
		s.elements[elem] = struct{}{}
	}
}

// Contains will determine if a given element belongs to the set
func (s *StringSet) Contains(elem string) bool {
	_, exists := s.elements[elem]
	return exists
}

// ToSlice will convert the StringSet to a string slice
func (s *StringSet) ToSlice() []string {
	i := 0
	elements := make([]string, len(s.elements))
	for k := range s.elements {
		elements[i] = k
		i++
	}
	return elements
}
