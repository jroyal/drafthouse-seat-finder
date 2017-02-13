package utils

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStringSet(t *testing.T) {
	s := NewStringSet()
	assert.NotNil(t, s)
	assert.Equal(t, 0, len(s.elements))
}

func TestAdd(t *testing.T) {
	s := NewStringSet()

	s.Add("elem")
	assert.Equal(t, 1, len(s.elements))

	s.Add("test")
	assert.Equal(t, 2, len(s.elements))

	s.Add("elem")
	assert.Equal(t, 2, len(s.elements))
}

func TestAddSlice(t *testing.T) {
	s := NewStringSet()

	s.Add("elem")
	assert.Equal(t, 1, len(s.elements))

	slice := []string{"James", "Dude", "elem"}
	s.AddSlice(slice)
	assert.Equal(t, 3, len(s.elements))
}

func TestContains(t *testing.T) {
	s := NewStringSet()

	s.Add("elem")
	assert.Equal(t, 1, len(s.elements))

	assert.True(t, s.Contains("elem"))
	assert.False(t, s.Contains("test"))
}

func TestToSlice(t *testing.T) {
	s := NewStringSet()

	s.Add("elem")

	slice := []string{"elem", "James", "Dude"}
	s.AddSlice(slice)
	assert.Equal(t, 3, len(s.elements))

	result := s.ToSlice()

	sort.Strings(slice)
	sort.Strings(result)
	assert.Equal(t, slice, result)
}
