package main

import (
	"fmt"
	"os"
	"runtime/trace"
	"sync"
)

var mapMux = sync.Mutex{}
var length = 10

func main() {
	fmt.Println("slicemap iterations: ", length)
	trace.Start(os.Stderr) // gotrace
	defer trace.Stop()     // gotrace
	testMap()

	testSliceMap()

	testModMap()

	testGoMap()
}
