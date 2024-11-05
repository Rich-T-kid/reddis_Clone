package main

import (
	"fmt"
	"rd_clone/DataStructures"
)

// Define a struct to represent the JSON data
func main() {
	reddis := DataStructures.ConfigHashMap("_Storage", "Data")
	defer reddis.Finish()
	/*
		reddis.SetKey("type", "shit")
		reddis.SetKey("rich", "richard")
		reddis.SetKeyTTL("Molly", "richard", time.Second)
		reddis.SetKey("lucki", "best")

		reddis.SetKeyTTL("neptun", "dest", 5*time.Second)
		reddis.SetKeyTTL("KingBeast", "Den of lions", 5*time.Second)
		reddis.SetKeyTTL("futur", "guapo", 12*time.Second)
		fmt.Println(reddis.KeysAndTTL())
		reddis.UpdateTTl("neptun", 4*time.Second)
		fmt.Println(reddis.KeysAndTTL())
		time.Sleep(3 * time.Second)
		fmt.Println(reddis.KeysAndTTL())
		//fmt.Println(reddis.Get("Molly"))
	*/
	fmt.Println(reddis.KeysAndTTL())

}
