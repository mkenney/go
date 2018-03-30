package main

import (
	"fmt"
	"time"
)

func testMap() {
	testMap := make(map[int]int)
	var start time.Time
	var end time.Time

	start = time.Now()
	for a := 0; a < length; a++ {
		testMap[a] = a * 100
	}
	end = time.Now()
	fmt.Println("map set time: ", end.Sub(start))

	start = time.Now()
	for a := 0; a < length; a++ {
		val, found := testMap[a]
		if a*100 != val {
			panic(fmt.Sprintf("val %d = %v (%v)", a, val, found))
		}
	}
	end = time.Now()
	fmt.Println("map get time: ", end.Sub(start))

	start = time.Now()
	for a := 0; a < length; a++ {
		delete(testMap, a)
	}
	end = time.Now()
	fmt.Println("map del time: ", end.Sub(start))
}
