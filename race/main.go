package main

import (
	"fmt"
	"sync"
)

func main() {
	results := make(map[int]int)
	wg := sync.WaitGroup{}

	m := sync.Mutex{}

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			m.Lock()
			results[n] = n * n
			m.Unlock()
		}(i)
	}

	wg.Wait()

	fmt.Println(results)
}
