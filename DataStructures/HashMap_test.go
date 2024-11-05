package DataStructures

import (
	"testing"
)

const testDir = "_testDir"
const testFile = "testData"

func TestMapinit(t *testing.T) {
	test_map := ConfigHashMap(testDir, testFile)
	defer test_map.PersistData(test_map.Collective, test_map.Storage)
	test_map.SetKey("richard", "baah")
	expected := "baah"
	_, result := test_map.Get("richard")
	if expected != result {
		t.Errorf("Wrong output should be %s got %s", expected, result)
	} else {
		t.Log("Good")
	}
}
