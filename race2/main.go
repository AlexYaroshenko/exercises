package main

import (
	"fmt"
	"sync"
)

func main() {
	type res struct {
		key int
		val int
	}

	results := make(map[int]int)
	wg := sync.WaitGroup{}
	resCh := make(chan res)

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			resCh <- res{key: n, val: n * n}
		}(i)
	}

	go func() {
		wg.Wait()
		close(resCh)
	}()

	for r := range resCh {
		results[r.key] = r.val
	}

	fmt.Println(results)
}
