package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

// expensive function that we want to avoid duplicate calls to
func getDataFromDB(key string) (string, error) {
	log.Printf("DB query executing for key: %s\n", key)
	// Simulate an expensive operation
	time.Sleep(100 * time.Millisecond)
	return fmt.Sprintf("Data for key: %s", key), nil
}

func main() {
	// Create a new singleflight group
	var g singleflight.Group

	// Create a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Simulate multiple concurrent requests for the same data
	for i := range 10 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			key := "user-123" // Same key used by all goroutines

			// Using singleflight to deduplicate calls
			result, err, shared := g.Do(key, func() (interface{}, error) {
				return getDataFromDB(key)
			})

			if err != nil {
				log.Printf("Error: %v", err)
				return
			}

			data := result.(string)
			log.Printf("Goroutine %d got result: %s (shared: %v)\n", id, data, shared)
		}(i)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	// Demonstrate DoChan which returns a channel
	fmt.Println("\nDemonstrating DoChan:")
	wg.Add(3)

	for i := range 3 {
		go func(id int) {
			defer wg.Done()
			key := "product-456"

			ch := g.DoChan(key, func() (interface{}, error) {
				log.Printf("DoChan function executing for key: %s\n", key)
				time.Sleep(200 * time.Millisecond)
				return fmt.Sprintf("Product data for: %s", key), nil
			})

			// Wait for the result
			result := <-ch
			log.Printf("DoChan goroutine %d got result: %s (shared: %v)\n",
				id, result.Val.(string), result.Shared)
		}(i)
	}

	wg.Wait()

	// Demonstrate Forget which drops cached values
	fmt.Println("\nDemonstrating Forget:")

	// First call
	result1, _, _ := g.Do("cache-key", func() (any, error) {
		log.Println("Computing value for cache-key")
		return "cached value", nil
	})
	log.Printf("First call result: %s\n", result1)

	// This call should return the cached result without executing the function
	result2, _, _ := g.Do("cache-key", func() (any, error) {
		log.Println("This shouldn't execute because the value is in flight")
		return "new value", nil
	})
	log.Printf("Second call result (should be cached): %s\n", result2)

	// Forget the key
	g.Forget("cache-key")
	log.Println("Forgot cache-key")

	// This should execute the function again
	result3, _, _ := g.Do("cache-key", func() (any, error) {
		log.Println("Computing value for cache-key again after Forget")
		return "new cached value", nil
	})
	log.Printf("Third call result (after Forget): %s\n", result3)
}
