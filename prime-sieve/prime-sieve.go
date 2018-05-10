// A concurrent prime sieve
package main

import (
	"fmt"
	"sync"
)

var idmux sync.Mutex
var idint = 0

func getID() int {
	idmux.Lock()
	idint++
	idmux.Unlock()
	return idint
}

// Generate sends the sequence 2, 3, 4, ... to channel 'ch'.
func Generate(ch chan<- int) {
	//id := fmt.Sprintf("Generate(%d)", getID())
	for i := 2; ; i++ {
		//fmt.Printf("\n---------------------------------\n%s: 1. sending %d\n", id, i)
		ch <- i // Send 'i' to channel 'ch'.
	}
}

// Filter copies the values from channel 'in' to channel 'out',
// removing those divisible by 'prime'.
func Filter(num <-chan int, result chan<- int, prime int) {
	//id := fmt.Sprintf("Filter(%d)", getID())
	//fmt.Printf("\n	%s: 2. executing filter() on %d\n", id, prime)
	for {
		i := <-num // Receive value from 'in'.
		//fmt.Printf("\n		%s: 3. received %d\n", id, i)
		//fmt.Printf("\n			%s: 4. %d %% %d = %d\n", id, i, prime, i%prime)
		if i%prime != 0 {
			//fmt.Printf("\n			%s: - using %d\n", id, i)
			result <- i // Send 'i' to 'out'.
		} else {
			//fmt.Printf("\n			%s: - discarding %d", id, i)
		}
		//fmt.Printf("\n")
	}
}

// The prime sieve: Daisy-chain Filter processes.
func main() {
	num := make(chan int) // Create a new channel.
	go Generate(num)      // Launch Generate goroutine.
	for i := 0; i < 500; i++ {
		prime := <-num
		//fmt.Printf("\n				5. prime: %d\n", prime)
		fmt.Printf("\n%d: %d\n", i+1, prime)
		if prime > 2100 {
			break
		}
		result := make(chan int)
		go Filter(num, result, prime)
		num = result
	}
}
