package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	results := fanIn(
		booring("Amy"),
		booring("Joe"),
	)
	for {
		select {
		case result := <-results:
			fmt.Println(result)
		}
	}
	fmt.Println("end.")
}

func booring(name string) chan string {
	c := make(chan string)
	go func() {
		for a := 0; a < 10; a++ {
			sleepTime := float32(rand.Intn(1000))
			time.Sleep(time.Millisecond * time.Duration(sleepTime))
			c <- fmt.Sprintf("%s: %d took %f seconds\n", name, a, (sleepTime / 1000))
		}
	}()
	return c
}

func fanIn(inputs ...chan string) chan string {
	c := make(chan string)
	for _, input := range inputs {
		chanInput := input
		go func() {
			for {
				faster := time.After(time.Millisecond * 10)
				select {
				case c <- <-chanInput:
				case <-faster:
					fmt.Println("faster...")
				}
			}
		}()
	}
	return c
}
