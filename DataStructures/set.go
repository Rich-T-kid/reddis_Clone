package DataStructures

type Set struct {
	size  int
	items map[string]bool // only the key is what matters, booleans are only their to save memory

}

func (s *Set) exist(element string) bool {
	for key := range s.items {
		if key == element {
			return true
		}

	}
	return false
}

func (s *Set) AddElement(element string) error {
	if s.exist(element) {
		// if it already exist dont need to do anything dont increase size
		return nil
	}
	s.items[element] = false
	s.size++

	return nil
}
func (s *Set) RemoveElemnt(element string) error {
	for key := range s.items {
		if element == key {
			s.size--
			delete(s.items, key)
		}
	}
	return nil
}

func (s *Set) Clear() error {
	emptyMap := make(map[string]bool)
	s.items = emptyMap
	s.size = 0
	return nil
}
func (s *Set) isEmpty() bool {
	return s.size == 0
}

func (s *Set) Elements() []string {
	var elements []string
	for key := range s.items {
		elements = append(elements, key)
	}
	return elements
}

func NewSet() *Set {
	return &Set{
		items: make(map[string]bool),
	}
}
