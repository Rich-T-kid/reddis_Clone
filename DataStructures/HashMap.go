package DataStructures

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	instance *HashTable
	MaxTTL   = 0
)

type HashTable struct {
	Collective map[string]interface{} //  user will have their own dict
	key_TTL    map[string]*Tuple      `json:"-"`
	Sets       map[string]*Set
	ErrorLog   *os.File
	RecoverLog *os.File
	Storage    *HardDisk
	Stats      *hashTableStats
	sync.Mutex
}

func (h *HashTable) loadHashTable(hashTable map[string]interface{}, storage *HardDisk, idx int) {
	filePath := filepath.Join(storage.directory, storage.fileNames[0])
	file, err := os.Open(filePath)
	if err != nil {
		panic("error had occured opening the file")
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&hashTable)
	if err != nil {
		if err == io.EOF {
			fmt.Println("The current hardDrive for the cache is empty")
			return
		}
		fmt.Println("Error reading JSON data:", err)
		return
	}
	fmt.Println("filled in hashtable")

}
func (h *HashTable) loadTuples(tupleMap map[string]*Tuple, storage *HardDisk, idx int) {
	filePath := filepath.Join(storage.directory, storage.fileNames[idx])
	file, err := os.Open(filePath)
	if err != nil {
		panic("error had occured opening the file")
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&tupleMap)
	if err != nil {
		if err == io.EOF {
			fmt.Println("The current hardDrive for the cache is empty")
			return
		}
		fmt.Println("Error reading JSON data:", err)
		return
	}
	fmt.Println("filled in Tuple table")
}
func (h *HashTable) loadSets(Sets map[string]*Set, storage *HardDisk, idx int) {
	filePath := filepath.Join(storage.directory, storage.fileNames[idx])
	file, err := os.Open(filePath)
	if err != nil {
		panic("error had occured opening the file")
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&Sets)
	if err != nil {
		if err == io.EOF {
			fmt.Println("The current hardDrive for the cache is empty")
			return
		}
		fmt.Println("Error reading JSON data:", err)
		return
	}
	fmt.Println("filled in Set Structs table")
}

/*
writes the current hash table data to the Secondary storages to perist inbetween program exeuction
*/
func (h *HashTable) persistData(sourceData interface{}, storage *HardDisk, idx int) {
	filePath := filepath.Join(storage.directory, storage.fileNames[idx])
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic("error had occured opening the file")
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(sourceData); err != nil {
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
	filePointers := []*os.File{}
	filePointers = append(filePointers, DataFile)

	// Initialize a slice of strings for file names
	files := []string{}
	files = append(files, filename)
	return &HardDisk{
		persitance: filePointers,
		directory:  directory,
		fileNames:  files,
	}
}

func addFile(filename string, Storage *HardDisk) {
	// Generate the full file path with .json extension
	DataFilePath := filepath.Join(Storage.directory, filename)

	// Attempt to open or create the file
	DataFile, err := os.OpenFile(DataFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic("error creating or opening file")
	}

	// Append the file name and file pointer to Storage
	Storage.fileNames = append(Storage.fileNames, filename)
	Storage.persitance = append(Storage.persitance, DataFile)
	fmt.Printf("created %s at location %s\n", filename, DataFilePath)
}

func createHashTableStat() *hashTableStats {
	return &hashTableStats{
		keys:      make([]string, 1),
		ttl_Times: make([]string, 1),
	}
}
func (h *HashTable) Finish() {
	h.persistData(h.Collective, h.Storage, 0)
	h.persistData(h.key_TTL, h.Storage, 1)
	h.persistData(h.Sets, h.Storage, 2)
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
	fmt.Println("Directories created successfully at: ", dirPath)
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
		key_TTL:    make(map[string]*Tuple),
		Sets:       make(map[string]*Set),
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

func ConfigHashMap(directory, filename string) *HashTable {
	stroage := CreateStorage(directory+"_storage", filename+".json")
	table := createTable(stroage)
	storeHashTable(table)
	addFile("TTL_times.json", stroage)
	addFile("Set_structs.json", stroage)
	table.loadHashTable(table.Collective, table.Storage, 0)
	table.loadTuples(table.key_TTL, table.Storage, 1)
	table.loadSets(table.Sets, table.Storage, 2)
	fmt.Println(table.Sets["myset"])
	return table
}

func (h *HashTable) Get(key string) (bool, interface{}) {
	var RecoveryMessage = fmt.Sprintf("Get %s\n", key)
	h.RecoverLog.Write([]byte(RecoveryMessage))
	// check for existance and check for ttl
	if !h.Exist(key) {
		return false, nil
	}
	tuple := h.key_TTL[key]
	if time.Now().After(tuple.ExpiresAt) {
		h.DeleteKey(key)
		return false, nil
	} else {
		return true, h.Collective[key]
	}
}

func (h *HashTable) SetKeyTTL(key string, value interface{}, ttl time.Duration) {
	var RecoveryMessage = fmt.Sprintf("Set %s:%v TTL:%v \n", key, value, ttl.Seconds())
	h.RecoverLog.Write([]byte(RecoveryMessage))
	tuple := NewTuple(value, ttl)
	h.key_TTL[key] = tuple
	h.Collective[key] = value
}

func (h *HashTable) SetKey(key string, value interface{}) {
	var RecoveryMessage = fmt.Sprintf("Set %s:%v\n", key, value)
	h.RecoverLog.Write([]byte(RecoveryMessage))
	NeverExires := time.Hour * 500000
	tuple := NewTuple(value, NeverExires)
	h.key_TTL[key] = tuple
	h.Collective[key] = value
}

func (h *HashTable) DeleteKey(inputkey string) error {
	var RecoveryMessage = fmt.Sprintf("Remove %s\n", inputkey)
	h.RecoverLog.Write([]byte(RecoveryMessage))
	for key := range h.Collective {
		if key == inputkey {
			delete(h.Collective, key)
			delete(h.key_TTL, key)
			return nil
		}
	}
	return errors.New("key doesnt exist to delete")
}

func (h *HashTable) Exist(key string) bool {
	if _, ok := h.Collective[key]; ok {
		return true
	} else {
		return false
	}

}
func (h *HashTable) UpdateTTl(key string, newTime time.Duration) error {
	var RecoveryMessage = fmt.Sprintf("Updating TTL %s:%v\n", key, newTime.Seconds())
	h.RecoverLog.Write([]byte(RecoveryMessage))
	if !h.Exist(key) {
		return errors.New("key doesnt exist to update TTL")
	}
	tuple := NewTuple(h.key_TTL[key], newTime)
	h.key_TTL[key] = tuple
	return nil
}

func (h *HashTable) Keys() []string {
	var RecoveryMessage = fmt.Sprintf("All Keys %d\n", time.Now().Second())
	h.RecoverLog.Write([]byte(RecoveryMessage))
	h.clearDeadKeys()
	var res []string
	for k := range h.Collective {
		res = append(res, k)
	}
	return res
}
func (h *HashTable) KeysAndTTL() []string {
	var RecoveryMessage = fmt.Sprintf("All Keys %d\n", time.Now().Second())
	h.RecoverLog.Write([]byte(RecoveryMessage))
	h.clearDeadKeys()
	var res []string
	for key, tuple := range h.key_TTL {
		value := tuple.Value
		expiresAt := tuple.ExpiresAt.Format("2006-01-02 15:04:05") // Custom time format for readability
		formattedString := fmt.Sprintf("%s -- Value: %v, Expires at: %s\n", key, value, expiresAt)
		res = append(res, formattedString)
	}
	return res
}
func (h *HashTable) clearDeadKeys() {

	var RecoveryMessage = fmt.Sprintf("Clearing Dead Keys %s", time.Now().Format("2006-01-02 15:04:05"))
	h.RecoverLog.Write([]byte(RecoveryMessage))
	// check for existance and check for ttl
	for key := range h.key_TTL {
		tuple := h.key_TTL[key]
		if time.Now().After(tuple.ExpiresAt) {
			h.DeleteKey(key)
		}
	}
}
func (h *HashTable) Increment(key string) error {
	if h.Exist(key) && isNumber(h.Collective[key]) {
		switch v := h.Collective[key].(type) {
		case int:
			h.Collective[key] = v + 1
		case int8:
			h.Collective[key] = v + 1
		case int16:
			h.Collective[key] = v + 1
		case int32:
			h.Collective[key] = v + 1
		case int64:
			h.Collective[key] = v + 1
		case uint:
			h.Collective[key] = v + 1
		case uint8:
			h.Collective[key] = v + 1
		case uint16:
			h.Collective[key] = v + 1
		case uint32:
			h.Collective[key] = v + 1
		case uint64:
			h.Collective[key] = v + 1
		case float32:
			h.Collective[key] = v + 1.0
		case float64:
			h.Collective[key] = v + 1.0
		default:
			return errors.New("unsupported numeric type")
		}
	} else {
		return errors.New("invalid key to increment. either it doesn't exist or isn't of a numeric type")
	}
	return nil
}

func (h *HashTable) Decrement(key string) error {
	if h.Exist(key) && isNumber(h.Collective[key]) {
		switch v := h.Collective[key].(type) {
		case int:
			h.Collective[key] = v - 1
		case int8:
			h.Collective[key] = v - 1
		case int16:
			h.Collective[key] = v - 1
		case int32:
			h.Collective[key] = v - 1
		case int64:
			h.Collective[key] = v - 1
		case uint:
			h.Collective[key] = v - 1
		case uint8:
			h.Collective[key] = v - 1
		case uint16:
			h.Collective[key] = v - 1
		case uint32:
			h.Collective[key] = v - 1
		case uint64:
			h.Collective[key] = v - 1
		case float32:
			h.Collective[key] = v - 1.0
		case float64:
			h.Collective[key] = v - 1.0
		default:
			return errors.New("unsupported numeric type")
		}
	} else {
		return errors.New("invalid key to increment. Either it doesn't exist or isn't of a numeric type")
	}
	return nil
}

const SetDoesntExist = "set doesnt exist"

// SADD setname value1, value2, value3, ... ect
func (h *HashTable) Sadd(setkey string, elements ...string) error {
	if !h.SetExist(setkey) {
		h.Sets[setkey] = NewSet()
	}
	set := h.Sets[setkey]
	for _, element := range elements {
		set.AddElement(element)
	}
	return nil
}

// is element a member
func (h *HashTable) IsMember(setkey, element string) (bool, error) {
	if !h.SetExist(setkey) {
		return false, errors.New(SetDoesntExist)
	}
	set := h.Sets[setkey]
	return set.exist(element), nil

}

// all members
func (h *HashTable) Smembers(setkey string) []string {
	var result []string
	if !h.SetExist(setkey) {
		return result
	}
	set := h.Sets[setkey]
	for k := range set.Items {
		result = append(result, k)
	}
	return result

}

// remove member / members
func (h *HashTable) SRem(setkey, member string) error {
	if !h.SetExist(setkey) {
		return errors.New(SetDoesntExist)
	}
	set := h.Sets[setkey]
	if set.exist(member) {
		delete(h.Sets, member)
		return nil
	}
	return errors.New("key doesnt exist in set")
}

// # of elements in set
func (h *HashTable) SCard(setkey string) int {
	if !h.SetExist(setkey) {
		return 0
	}
	set := h.Sets[setkey]
	return set.Size
}
func (h *HashTable) SetExist(setkey string) bool {
	_, ok := h.Sets[setkey]
	return ok
}

func isNumber(value interface{}) bool {
	switch value.(type) {
	case int, int8, int16, int32, int64, // Signed integers
		uint, uint8, uint16, uint32, uint64, // Unsigned integers
		float32, float64: // Floating-point numbers
		return true
	default:
		return false
	}
}
