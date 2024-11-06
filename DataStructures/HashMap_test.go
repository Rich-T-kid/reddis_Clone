package DataStructures

import (
	"testing"
)

const testDir = "_testDir"
const testFile = "testData"

func TestMapinit(t *testing.T) {
	test_map := ConfigHashMap(testDir, testFile)
	defer test_map.persistData(test_map.Collective, test_map.Storage, 0)
	test_map.SetKey("richard", "baah")
	expected := "baah"
	_, result := test_map.Get("richard")
	if expected != result {
		t.Errorf("Wrong output should be %s got %s", expected, result)
	} else {
		t.Log("Good")
	}
}

func setupHashTable() *HashTable {
	return ConfigHashMap("test_data", "test_file")
}

func TestSetAndGetKey(t *testing.T) {
	table := setupHashTable()
	key, value := "testKey", "testValue"

	// Test setting a key
	table.SetKey(key, value)

	// Test getting the key
	found, gotValue := table.Get(key)
	if !found {
		t.Errorf("Expected to find key '%s', but it was not found", key)
	}
	if gotValue != value {
		t.Errorf("Expected value '%v', got '%v'", value, gotValue)
	}
}

func TestDeleteKey(t *testing.T) {
	table := setupHashTable()
	key := "deleteKey"
	table.SetKey(key, "someValue")

	// Test deleting the key
	err := table.DeleteKey(key)
	if err != nil {
		t.Errorf("DeleteKey returned an error: %v", err)
	}

	// Check if the key still exists
	found, _ := table.Get(key)
	if found {
		t.Error("Expected key to be deleted, but it was found")
	}

	// Test deleting a non-existing key
	err = table.DeleteKey("nonExistentKey")
	if err == nil {
		t.Error("Expected an error when deleting a non-existent key, but got nil")
	}
}

func TestIncrement(t *testing.T) {
	table := setupHashTable()
	key := "incrementKey"
	table.SetKey(key, 5)

	// Test incrementing the key's value
	err := table.Increment(key)
	if err != nil {
		t.Errorf("Increment returned an error: %v", err)
	}

	// Verify the incremented value
	_, value := table.Get(key)
	expected := 6
	if value != expected {
		t.Errorf("Expected value %d, got %v", expected, value)
	}
}

func TestDecrement(t *testing.T) {
	table := setupHashTable()
	key := "decrementKey"
	table.SetKey(key, 5)

	// Test decrementing the key's value
	err := table.Decrement(key)
	if err != nil {
		t.Errorf("Decrement returned an error: %v", err)
	}

	// Verify the decremented value
	_, value := table.Get(key)
	expected := 4
	if value != expected {
		t.Errorf("Expected value %d, got %v", expected, value)
	}
}

func TestExist(t *testing.T) {
	table := setupHashTable()
	key := "existKey"
	table.SetKey(key, "someValue")

	// Test key existence
	if !table.Exist(key) {
		t.Error("Expected key to exist but it was not found")
	}

	// Test non-existing key
	if table.Exist("nonExistentKey") {
		t.Error("Expected key to not exist, but it was found")
	}
}

func TestKeys(t *testing.T) {
	table := setupHashTable()
	keys := []string{"key1", "key2", "key3"}
	for _, key := range keys {
		table.SetKey(key, "value")
	}

	// Test retrieving all keys
	allKeys := table.Keys()
	for _, key := range keys {
		found := false
		for _, k := range allKeys {
			if k == key {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected key '%s' to be in the list of keys", key)
		}
	}
}
func TestSadd(t *testing.T) {
	h := setupHashTable()
	err := h.Sadd("testSet", "value1", "value2", "value3")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if h.SCard("testSet") != 3 {
		t.Errorf("Expected 3 elements, got %d", h.SCard("testSet"))
	}
}

func TestIsMember(t *testing.T) {
	h := setupHashTable()
	h.Sadd("testSet", "value1", "value2")
	isMember, err := h.IsMember("testSet", "value1")
	if err != nil || !isMember {
		t.Errorf("Expected value1 to be a member, got error: %v", err)
	}
	isMember, err = h.IsMember("testSet", "nonExistentValue")
	if err != nil || isMember {
		t.Errorf("Expected nonExistentValue not to be a member, got: %v", isMember)
	}
}

func TestSmembers(t *testing.T) {
	h := setupHashTable()
	h.Sadd("testSet", "value1", "value2", "value3")
	members := h.Smembers("testSet")
	expectedMembers := []string{"value1", "value2", "value3"}
	if len(members) != len(expectedMembers) {
		t.Errorf("Expected %d members, got %d", len(expectedMembers), len(members))
	}
	for _, val := range expectedMembers {
		found := false
		for _, member := range members {
			if member == val {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected member %s not found in %v", val, members)
		}
	}
}

func TestSRem(t *testing.T) {
	h := setupHashTable()
	h.Sadd("testSet", "value1", "value2")
	err := h.SRem("testSet", "value1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	isMember, err := h.IsMember("testSet", "value1")
	if err != nil || isMember {
		t.Errorf("Expected value1 to be removed from set, but it still exists")
	}

	err = h.SRem("testSet", "nonExistentValue")
	if err == nil {
		t.Errorf("Expected error for non-existent member, but got none")
	}
}

func TestSCard(t *testing.T) {
	h := setupHashTable()
	if h.SCard("testSet") != 0 {
		t.Errorf("Expected 0 elements, got %d", h.SCard("testSet"))
	}
	h.Sadd("testSet", "value1", "value2")
	if h.SCard("testSet") != 2 {
		t.Errorf("Expected 2 elements, got %d", h.SCard("testSet"))
	}
	h.SRem("testSet", "value1")
	if h.SCard("testSet") != 1 {
		t.Errorf("Expected 1 element after removal, got %d", h.SCard("testSet"))
	}
}

func TestSetExist(t *testing.T) {
	h := setupHashTable()
	if h.SetExist("nonExistentSet") {
		t.Errorf("Expected set not to exist")
	}
	h.Sadd("newSet", "value1")
	if !h.SetExist("newSet") {
		t.Errorf("Expected set to exist after adding elements")
	}
}

/*
func TestUpdateTTL(t *testing.T) {
	table := setupHashTable()
	key := "ttlKey"
	table.SetKey(key, "value")

	// Test updating TTL (dummy behavior in this case)
	err := table.UpdateTTl(key, time.Now().Add(time.Minute))
	if err != nil {
		t.Errorf("UpdateTTL returned an error: %v", err)
	}
}*/
