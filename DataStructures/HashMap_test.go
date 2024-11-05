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
