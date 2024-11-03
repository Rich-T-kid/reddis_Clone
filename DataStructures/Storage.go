package DataStructures

import (
	"fmt"
	"os"
)

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

func (h *hashTableStats) addkey(key string) {
	h.keys = append(h.keys, key)
}
func (h *hashTableStats) addTTL(key, TTL string) {
	formatedString := fmt.Sprintf("%s:%s", key, TTL)
	h.ttl_Times = append(h.ttl_Times, formatedString)
}

func (h *hashTableStats) removekey(key string) {
	h.size--
	for i := range h.keys {
		if h.keys[i] == key {
			h.keys = append(h.keys[:i], h.keys[i+1:]...)
			break // Exit the loop once the key is found and removed
		}
	}
	fmt.Println("element doesnt exist to be removed from keys list for table stats")
}
func (h *hashTableStats) removeTTL(key, TTL string) {
	formatedString := fmt.Sprintf("%s:%s", key, TTL)
	for i, element := range h.ttl_Times {
		if element == formatedString {
			h.ttl_Times = append(h.ttl_Times[:i], h.ttl_Times[i+1:]...)
			break
		}
	}
	fmt.Println(formatedString + " doest exist to be deleted from ttl list for table stats")
}
