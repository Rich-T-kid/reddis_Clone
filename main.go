package main

import (
	"fmt"
	"rd_clone/DataStructures"
)

// Define a struct to represent the JSON data
func main() {
	reddis := DataStructures.ConfigHashMap("_Storage", "Data")
	defer reddis.PersistData(reddis.Collective, reddis.Storage)
	reddis.SetKey("type", "shit")
	reddis.SetKey("rich", "richard")
	reddis.Get("rich")
	reddis.DeleteKey("ri")
	reddis.DeleteKey("rich")
	reddis.SetKey("rich", 5)
	fmt.Println(reddis.Get("rich"))
	for i := 0; i < 3; i++ {
		reddis.Increment("rich")
	}
	reddis.Decrement("rich")
	reddis.SetKey("bob", "square")
	fmt.Println(reddis.Decrement("bob"))
	reddis.Increment("tob")
	fmt.Println(reddis.Get("rich"))

	fmt.Println(reddis.Keys())

}
