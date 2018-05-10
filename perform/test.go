package main

import (
	"fmt"
	"time"
)

func main() {
	a := 1
	start := time.Now()
	for i := 1; i < 100000000; i++ {
		a += i
		_ = a % i
	}
	fmt.Printf("Go took %f seconds\n", time.Now().Sub(start).Seconds())
}
