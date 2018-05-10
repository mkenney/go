package main

import (
	"fmt"

	"github.com/mkenney/go/cli"
)

func main() {
	fmt.Printf("%v\n", "string")
	fmt.Printf("%v\n", true)
	fmt.Printf("%v\n", 1)
	fmt.Printf("%v\n", 1.1)

	//var a interface{}
	//b := ""
	//a = "asdf"

	//switch a.(type) {
	//case b:
	//	fmt.Println("got here")
	//default:
	//	fmt.Println("didn't work")
	//}

	cli.Parse()
}
