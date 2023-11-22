package skiplist

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
)

func TestInsert(t *testing.T) {
	list := NewSkipList[int, float64]()

	list.Insert(1, 0.2)
	list.Insert(0, 0.1)
	list.Insert(2, 0.3)
	list.Insert(3, 0.4)
	list.Insert(5, 0.6)
	list.Insert(6, 0.7)
	list.Insert(7, 0.8)
	list.Insert(8, 0.9)
	list.Insert(9, 1.0)
	list.Insert(9, 1.1)
	list.Insert(9, 1.2)
	list.Insert(6, 6.6)

	fmt.Println("list")
	fmt.Println(list)
}

func TestInsertParallel(t *testing.T) {
	list := NewSkipList[int, float64]()
	var mutex sync.Mutex

	insertValues := func(numElements int, done chan<- error) {
		// Generate random key

		for j := 0; j < numElements; j++ {
			mutex.Lock()
			list.Insert(rand.Int(), rand.Float64())
			mutex.Unlock()
		}

		// Signal completion without error
		done <- nil
	}

	numGoroutines := 50
	numElements := 50
	done := make(chan error)

	// Insert values using goroutines
	for i := 0; i < numGoroutines; i++ {
		go insertValues(numElements, done)
	}

	// Wait for goroutines to finish
	for i := 0; i < numGoroutines; i++ {
		if err := <-done; err != nil {
			fmt.Printf("Error in goroutine %d: %v\n", i, err)
		}
	}

	fmt.Println("Final list:", list)
	fmt.Println(list)
}
