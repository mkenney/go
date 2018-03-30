package main

import (
    "bytes"
    "fmt"
    "golang.org/x/tour/tree"
)

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
    rWalk(t, ch)
    close(ch)
}

func rWalk(t *tree.Tree, ch chan int) {
    if t.Left != nil {
        rWalk(t.Left, ch)
    }
    if t.Right != nil {
        rWalk(t.Right, ch)
    }
    ch <- t.Value
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
    var ret_val bool = false

    ch1 := make(chan int)
    ch2 := make(chan int)

    go Walk(t1, ch1)
    go Walk(t2, ch2)

    var a_tour []int
    var b_tour []int
    select {
        case a := <- ch1:
            a_tour[]
        case a := <- ch2:
            b_tour.WriteString(a)
    }

    ret_val = (a_tour == b_tour)

    return ret_val
}

func main() {

    fmt.PrintF("%v", Same(tree.New(1), tree.New(1)))

    walk := make(chan int)
    go Walk(tree.New(1), walk)
    for a := range walk {
        fmt.Println(a)
    }

}
