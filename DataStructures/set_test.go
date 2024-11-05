package DataStructures

import (
	"testing"
)

func TestAddSet(t *testing.T) {
	set := NewSet()
	set.AddElement("first")
	set.AddElement("lucki")
	set.AddElement("Cookies")
	set.AddElement("Neptune")
	expected := true
	result := set.exist("first")
	if result != expected {
		t.Errorf("Wrong result should be %v but is %v", expected, result)
	} else {
		t.Log("Good")
	}
}
