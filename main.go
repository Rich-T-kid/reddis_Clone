package main

import (
	"fmt"
	"rd_clone/DataStructures"
)

// Define a struct to represent the JSON data
func main() {
	reddis := DataStructures.ConfigHashMap("_Storage", "Data")
	defer reddis.PersistData(reddis.Collective, reddis.Storage)
	fmt.Println("reddis clone-> ", reddis)
	reddis.SetKey("type", "shit")
	reddis.SetKey("rich", "richard")
	reddis.Get("rich")
	fmt.Println(reddis.Get("r"))
}
