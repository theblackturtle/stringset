package stringset

import (
	"fmt"
	"log"
	"strings"

	"github.com/projectdiscovery/hmap/store/hybrid"
)

type Set struct {
	hm *hybrid.HybridMap
}

// Deduplicate utilizes the Set type to generate a unique list of strings from the input slice.
func Deduplicate(input []string) []string {
	return New(input...).Slice()
}

// New returns a Set containing the values provided in the arguments.
func New(initial ...string) Set {
	s := Set{}
	hm, err := hybrid.New(hybrid.DefaultHybridOptions)
	if err != nil {
		log.Fatal(err)
	}
	s.hm = hm
	for _, v := range initial {
		s.Insert(v)
	}

	return s
}

// Has returns true if the receiver Set already contains the element string argument.
func (s Set) Has(element string) bool {
	// _, exists := s[strings.ToLower(element)]
	_, exists := s.hm.Get(strings.ToLower(element))
	return exists
}

// Insert adds the element string argument to the receiver Set.
func (s Set) Insert(element string) {
	s.hm.Set(strings.ToLower(element), []byte{})
}

// InsertMany adds all the elements strings into the receiver Set.
func (s Set) InsertMany(elements ...string) {
	for _, i := range elements {
		s.Insert(i)
	}
}

// Remove will delete the element string from the receiver Set.
func (s Set) Remove(element string) {
	e := strings.ToLower(element)
	// delete(s, e)
	s.hm.Del(e)
}

// Slice returns a string slice that contains all the elements in the Set.
func (s Set) Slice() []string {
	var i uint64

	k := make([]string, s.Len())
	s.hm.Scan(func(b1, b2 []byte) error {
		k[i] = string(b1)
		i++
		return nil
	})

	return k
}

// Union adds all the elements from the other Set argument into the receiver Set.
func (s Set) Union(other Set) {
	other.hm.Scan(func(b1, b2 []byte) error {
		s.Insert(string(b1))
		return nil
	})
	// for k := range other {
	// 	s.Insert(k)
	// }
}

// Len returns the number of elements in the receiver Set.
func (s Set) Len() int {
	size := 0
	s.hm.Scan(func(b1, b2 []byte) error {
		size++
		return nil
	})
	return size
}

// Subtract removes all elements in the other Set argument from the receiver Set.
func (s Set) Subtract(other Set) {
	other.hm.Scan(func(b1, b2 []byte) error {
		s.hm.Del(string(b1))
		return nil
	})
	// for item := range other {
	// 	s.Remove(item)
	// }
}

// Intersect causes the receiver Set to only contain elements also found in the
// other Set argument.
func (s Set) Intersect(other Set) {
	// for item := range s {
	// 	e := strings.ToLower(item)
	// 	if _, exists := other[e]; !exists {
	// 		delete(s, e)
	// 	}
	// }
	s.hm.Scan(func(b1, b2 []byte) error {
		e := strings.ToLower(string(b1))
		if _, exists := other.hm.Get(e); !exists {
			s.hm.Del(e)
		}
		return nil
	})
}

// Set implements the flag.Value interface.
func (s *Set) String() string {
	return strings.Join(s.Slice(), ",")
}

// Set implements the flag.Value interface.
func (s *Set) Set(input string) error {
	if input == "" {
		return fmt.Errorf("String parsing failed")
	}

	items := strings.Split(input, ",")
	for _, item := range items {
		s.Insert(strings.TrimSpace(item))
	}
	return nil
}
