package main

import (
	"fmt"
	"sync"
)

var syncOnce sync.Once

func work() {
	fmt.Println("work")
}

func main() {
	// syncOnce.Do(work)
	// syncOnce.Do(work)
	wg := sync.WaitGroup{}
	for range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			syncOnce.Do(work)
		}()
	}
	wg.Wait()
}
