package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

/*
	goal:
	- process data from an input channel using multiple goroutines
	- each goroutine should have a timeout of 5 seconds to process each item
	- if processing takes longer than 5 seconds, the goroutine should stop processing and move to the next item
	- collect results in an output channel
*/

func processData(data int) int {
	// Simulate processing time
	seconds := rand.Intn(10)
	fmt.Println("Processing", data, "for", seconds, "seconds")
	time.Sleep(time.Duration(seconds) * time.Second)
	return data * 2
}

func main() {
	in := make(chan int)
	out := make(chan int)

	go func() {
		for i := 0; i < 100; i++ {
			in <- i
		}
		close(in)
	}()

	now := time.Now()
	parallelProcess(in, out, 5)
	for data := range out {
		fmt.Println(data)
	}
	fmt.Println("Processing took:", time.Since(now))
}

func worker(in <-chan int, out chan<- int, wg *sync.WaitGroup) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	defer wg.Done()

	for {
		select {
		case data, ok := <-in:
			if !ok {
				return
			}

			resultChan := make(chan int)
			go func() {
				resultChan <- processData(data)
				close(resultChan)
			}()

			select {
			case <-ctx.Done():
				fmt.Println("Worker timed out while sending result")
				return
			case result := <-resultChan:
				out <- result
				return
			}
		case <-ctx.Done():
			fmt.Println("Worker timed out")
			return
		}
	}
}

// stop after 5 seconds
func parallelProcess(in <-chan int, out chan<- int, numWorkers int) {
	var wg sync.WaitGroup
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go worker(in, out, &wg)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
}
