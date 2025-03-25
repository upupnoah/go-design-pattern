package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int, 2)
	ch <- 1
	ch <- 2
	go func() {
		val := <-ch
		fmt.Println(val)
	}()
	ch <- 3
	time.Sleep(1 * time.Second)
	fmt.Println(<-ch)
}
