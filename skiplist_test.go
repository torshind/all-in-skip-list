package skiplist

import (
	"fmt"
	"math"
	"math/rand"
	"slices"
	"sync"
	"testing"
)

func TestGetRandomLevel(t *testing.T) {
	list := NewSkipList[int, float64](16)
	data := make([]int, 10e6)

	for i := range data {
		data[i] = list.getRandomLevel()
	}

	n := slices.Max(data)

	// Compute the observed frequency distribution
	observed := make([]int, n+1)
	for _, v := range data {
		observed[v]++
	}
	slices.Reverse(observed)
	fmt.Println(observed)

	var totalError, totalRelativeError float64

	for i := 1; i < len(observed); i++ {
		// Calculate the expected value (double the previous element)
		expected := 2 * observed[i-1]

		if expected != 0 && observed[i] != 0 {
			// Calculate the absolute error
			error := math.Abs(float64(observed[i] - expected))

			// Calculate the relative error
			relativeError := error / float64(observed[i-1])

			// Accumulate the errors
			totalError += error
			totalRelativeError += relativeError

			// Output information for each pair
			fmt.Printf("Element %d: %d, Expected: %d, Absolute Error: %.6e, Relative Error: %.6e\n", i, observed[i], expected, error, relativeError)
		}
	}

	// Calculate the average relative error
	averageRelativeError := totalRelativeError / float64(len(observed)-1)

	// Output the total error and average relative error
	fmt.Printf("\nTotal Error: %.6e\n", totalError)
	fmt.Printf("Average Relative Error: %.6e\n", averageRelativeError)
}

func TestInsert(t *testing.T) {
	list := NewSkipList[int, float64](32)

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
	list := NewSkipList[int, float64](32)
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
