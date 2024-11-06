/*
find . -type f ! -name "*.go" -exec rm -f {} +    # Remove non-.go files
find . -type d ! -name "." ! -name "*.go" -exec rm -rf {} +   # Remove empty directories
Commands to clean up after the test execute
*/
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

func TestNewSet(t *testing.T) {
	set := NewSet()
	if set == nil {
		t.Error("NewSet() returned nil")
	}
	if !set.isEmpty() {
		t.Error("New set should be empty")
	}
}

func TestAddElement(t *testing.T) {
	set := NewSet()
	err := set.AddElement("first")
	if err != nil {
		t.Errorf("AddElement() returned an error: %v", err)
	}
	if !set.exist("first") {
		t.Error("Element 'first' should exist after adding")
	}
	if set.Size != 1 {
		t.Errorf("Expected size 1, got %d", set.Size)
	}

	// Test adding an existing element
	err = set.AddElement("first")
	if err != nil {
		t.Errorf("AddElement() returned an error when adding an existing element: %v", err)
	}
	if set.Size != 1 {
		t.Errorf("Size should not change when adding an existing element; expected 1, got %d", set.Size)
	}
}

func TestRemoveElement(t *testing.T) {
	set := NewSet()
	set.AddElement("first")
	err := set.RemoveElemnt("first")
	if err != nil {
		t.Errorf("RemoveElemnt() returned an error: %v", err)
	}
	if set.exist("first") {
		t.Error("Element 'first' should not exist after removing")
	}
	if set.Size != 0 {
		t.Errorf("Expected size 0, got %d", set.Size)
	}

	// Test removing a non-existing element
	err = set.RemoveElemnt("nonexistent")
	if err != nil {
		t.Errorf("RemoveElemnt() returned an error when removing a non-existent element: %v", err)
	}
}

func TestClear(t *testing.T) {
	set := NewSet()
	set.AddElement("first")
	set.AddElement("second")
	err := set.Clear()
	if err != nil {
		t.Errorf("Clear() returned an error: %v", err)
	}
	if !set.isEmpty() {
		t.Error("Set should be empty after Clear()")
	}
}

func TestIsEmpty(t *testing.T) {
	set := NewSet()
	if !set.isEmpty() {
		t.Error("New set should be empty")
	}
	set.AddElement("first")
	if set.isEmpty() {
		t.Error("Set should not be empty after adding an element")
	}
	set.Clear()
	if !set.isEmpty() {
		t.Error("Set should be empty after Clear()")
	}
}

func TestElements(t *testing.T) {
	set := NewSet()
	set.AddElement("first")
	set.AddElement("second")
	elements := set.Elements()
	if len(elements) != 2 {
		t.Errorf("Expected 2 elements, got %d", len(elements))
	}
	if !(elements[0] == "first" || elements[1] == "first") {
		t.Error("Expected 'first' to be in elements")
	}
	if !(elements[0] == "second" || elements[1] == "second") {
		t.Error("Expected 'second' to be in elements")
	}
}
