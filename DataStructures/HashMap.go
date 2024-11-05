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
)

type HashTable struct {
	Collective map[string]interface{} //  user will have their own dict
	ErrorLog   *os.File
	RecoverLog *os.File
	sync.Mutex
	Storage *HardDisk
	Stats   *hashTableStats
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
		if err == io.EOF {
			fmt.Println("The current hardDrive for the cache is empty")
			return
		}
		fmt.Println("Error reading JSON data:", err)
		return
	}
	fmt.Println("filled in hashtable")

}

/*
writes the current hash table data to the Secondary storages to perist inbetween program exeuction
*/
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
		panic("error writing to 	stoarge")
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
	table.loadData(table.Collective, table.Storage)
	return table
}

func (h *HashTable) Get(key string) (bool, interface{}) {
	var RecoveryMessage = fmt.Sprintf("Get %s\n", key)
	h.RecoverLog.Write([]byte(RecoveryMessage))
	if value, ok := h.Collective[key]; ok {
		return true, value
	}
	return false, nil
}
func (h *HashTable) SetKey(key string, value interface{}) {
	var RecoveryMessage = fmt.Sprintf("Set %s:%v\n", key, value)
	h.RecoverLog.Write([]byte(RecoveryMessage))
	if h.Exist(key) {
		h.Collective[key] = value
		return
	} else {
		h.Stats.addkey(key)
		h.Collective[key] = value
		return
	}
}

func (h *HashTable) DeleteKey(inputkey string) error {
	var RecoveryMessage = fmt.Sprintf("Remove %s\n", inputkey)
	h.RecoverLog.Write([]byte(RecoveryMessage))
	for key := range h.Collective {
		if key == inputkey {
			delete(h.Collective, key)
			h.Stats.removekey(key)
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
func (h *HashTable) UpdateTTl(key string, newTime time.Time) error {
	var RecoveryMessage = fmt.Sprintf("Updating TTL %s:%d\n", key, newTime.Second())
	h.RecoverLog.Write([]byte(RecoveryMessage))
	return nil
}

func (h *HashTable) Keys() []string {
	var RecoveryMessage = fmt.Sprintf("All Keys %d\n", time.Now().Second())
	h.RecoverLog.Write([]byte(RecoveryMessage))
	var res []string
	for k := range h.Collective {
		res = append(res, k)
	}
	return res
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
