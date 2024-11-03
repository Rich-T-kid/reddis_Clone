package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Define a struct to represent the JSON data
type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}
type HardDisk struct {
	persitance *os.File
	directory  string
	fileName   string
}
type hashTableStats struct {
	size      int
	keys      []string
	ttl_Times []string // times of TTL for all keys
}
type HashTable struct {
	Collective map[string]interface{} //  user will have their own dict
	ErrorLog   *os.File
	RecoverLog *os.File
	sync.Mutex
	Storage *HardDisk
	Stats   *hashTableStats
}

var (
	instance *HashTable
)

/*
If element does exist return
*/
func (h *HashTable) get(key string) (bool, interface{}) {
	if value, ok := h.Collective[key]; ok {
		return true, value
	}
	return true, nil
}
func (h *HashTable) setKey(key string, value interface{}) error {
	h.Collective[key] = value
	return nil

}
func (h *HashTable) deleteKey(inputkey string) error {
	for key := range h.Collective {
		if key == inputkey {
			delete(h.Collective, key)
			return nil
		}
	}
	return errors.New("Key doesnt exist to delete")
}

func (h *HashTable) exist(key string) bool {
	if _, ok := h.Collective[key]; ok {
		return true
	} else {
		return false
	}

}
func (h *HashTable) updateTTl() error {
	return nil
}

func (h *HashTable) keys() []string {
	return h.Stats.keys
}

func (h *hashTableStats) incrementsize() {
	h.size++
}

func (h *hashTableStats) addkey(key string) {
	h.keys = append(h.keys, key)
}
func (h *hashTableStats) addTTL(key, TTL string) {
	formatedString := fmt.Sprintf("%s:%s", key, TTL)
	h.ttl_Times = append(h.ttl_Times, formatedString)
}

func (h *hashTableStats) removekey(key, TTL string) {

}
func (h *hashTableStats) removeTTL(key, TTL string) {
}

func (h *HashTable) loadData(hashTable map[string]interface{}, storage *HardDisk) {
	filePath := filepath.Join(storage.directory, storage.fileName)
	file, err := os.Open(filePath)
	if err != nil {
		panic("error had occured opening the file")
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&hashTable)
	if err != nil {
		fmt.Println("Error reading JSON data:", err)
		return
	}
	fmt.Println("filled in hashtable")

}

func (h *HashTable) PersistData(HashTable map[string]interface{}, storage *HardDisk) {
	filePath := filepath.Join(storage.directory, storage.fileName)
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic("error had occured opening the file")
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(HashTable); err != nil {
		panic("error writing to stoarge")
	}
	fmt.Println("Done updating json file for persistant storage")

}
func CreateStorage(directory, filename string) *HardDisk {

	err := os.MkdirAll(directory, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating directories:", err)
		panic("Error directory coulndt be made")
	}
	fmt.Println("Directories created successfully at:", directory)
	DataFilePath := filepath.Join(directory, filename)
	DataFile, err := os.OpenFile(DataFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		if os.IsExist(err) {
			DataFile, _ = os.OpenFile(DataFilePath, os.O_RDWR, 0644) // Open existing file
		} else {
			panic("error creating Error Logs")
		}
	}
	return &HardDisk{
		persitance: DataFile,
		directory:  directory,
		fileName:   filename,
	}
}
func createHashTableStat() *hashTableStats {
	return &hashTableStats{
		keys:      make([]string, 1),
		ttl_Times: make([]string, 1),
	}
}
func createTable(storage *HardDisk) *HashTable {
	const dirPath = "info"
	now := time.Now()
	formattedTime := now.Format("2006-01-02 15:04:05")
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating directories:", err)
		panic("Error directory coulndt be made")
	}
	fmt.Println("Directories created successfully at:", dirPath)
	errorFilePath := filepath.Join(dirPath, "InfoLog.txt")
	recoverLogPath := filepath.Join(dirPath, "Recover.txt")
	// Open error log file only if it does not exist
	errorFile, err := os.OpenFile(errorFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		if os.IsExist(err) {
			errorFile, _ = os.OpenFile(errorFilePath, os.O_RDWR, 0644) // Open existing file
		} else {
			panic("error creating Error Logs")
		}
	}

	// Open recovery log file only if it does not exist
	recoverFile, err := os.OpenFile(recoverLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		if os.IsExist(err) {
			recoverFile, _ = os.OpenFile(recoverLogPath, os.O_RDWR, 0644) // Open existing file
		} else {
			panic("error creating recovery file")
		}
	}
	var Logmessage = "Hash Table is being constructed" + formattedTime + "\n"
	// this will allow use to differential different times when the program is ran
	recoverFile.Write([]byte(Logmessage))
	errorFile.Write([]byte(Logmessage))
	return &HashTable{
		Collective: make(map[string]interface{}),
		ErrorLog:   errorFile,
		RecoverLog: recoverFile,
		Storage:    storage,
		Stats:      createHashTableStat(),
	}
}

func gethashtable() *HashTable {
	if instance == nil {
		panic("this should never execute ")
	} else {
		return instance
	}
}
func storeHashTable(input *HashTable) {
	instance = input
}
func init() {
	print("\n starting")
	stroage := CreateStorage("_storage", "test.json")
	table := createTable(stroage)
	storeHashTable(table)
	fmt.Println("Pre Loaded Data", table)
	table.loadData(table.Collective, table.Storage)
	fmt.Println("Post Loaded Data ", table)
	print("\n done")
}

func main() {
	// Directory and file paths
	fmt.Println("")
	fmt.Println("")
	reddis := gethashtable()
	fmt.Println(reddis.Collective)
	fmt.Println("->>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	reddis.setKey("richard", "cooldude")
	fmt.Println(reddis.get("richard"))
	fmt.Println(reddis.Collective)
	reddis.setKey("test2", "value2")
	reddis.setKey("test3", "test3")
	reddis.setKey("test4", "test3")
	fmt.Println("->>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	fmt.Println("->>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	fmt.Println(reddis.Collective)
	fmt.Println("->>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	reddis.PersistData(reddis.Collective, reddis.Storage)
	fmt.Println("->>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")

	//redisClone := hashtable()
	//print(redisClone)

	/*
		dirPath := "data/person"
		filePath := filepath.Join(dirPath, "person.json")

		// Step 1: Create directories
		err := os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating direactories:", err)
			return
		}
		fmt.Println("Directories created successfully at:", dirPath)

		// Step 2: Create and write to the JSON file
		person := Person{
			Name:    "Alice",
			Age:     30,
			Country: "USA",
		}

		file, err := os.Create(filePath)
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		err = encoder.Encode(person)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
		fmt.Println("JSON file created and data written successfully at:", filePath)

		// Step 3: Read the JSON file
		file, err = os.Open(filePath)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer file.Close()

		var personData Person
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&personData)
		if err != nil {
			fmt.Println("Error reading JSON data:", err)
			return
		}
		fmt.Printf("Read JSON data: Name: %s, Age: %d, Country: %s\n", personData.Name, personData.Age, personData.Country)

		// Step 4: Update the JSON file
		personData.Age = 31 // Modify the data (e.g., updating the age)

		// Reopen the file for writing, clear it, and write the updated data
		file, err = os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			fmt.Println("Error reopening file for updating:", err)
			return
		}
		defer file.Close()

		encoder = json.NewEncoder(file)
		err = encoder.Encode(personData)
		if err != nil {
			fmt.Println("Error writing updated data to file:", err)
			return
		}
	*/

}
