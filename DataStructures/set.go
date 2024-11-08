package DataStructures

type Set struct {
	Items map[string]bool `json:"items"` // only the key is what matters, booleans are only their to save memory

}

func (s *Set) exist(element string) bool {
	for key := range s.Items {
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
	s.Items[element] = false

	return nil
}
func (s *Set) RemoveElemnt(element string) error {
	for key := range s.Items {
		if element == key {
			delete(s.Items, key)
		}
	}
	return nil
}

func (s *Set) Clear() error {
	emptyMap := make(map[string]bool)
	s.Items = emptyMap
	return nil
}
func (s *Set) isEmpty() bool {
	return s.Size() == 0
}

func (s *Set) Elements() []string {
	var elements []string
	for key := range s.Items {
		elements = append(elements, key)
	}
	return elements
}
func (s *Set) Size() int {
	return len(s.Items)
}
func NewSet() *Set {
	return &Set{
		Items: make(map[string]bool),
	}
}
