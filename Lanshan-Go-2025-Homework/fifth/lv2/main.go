package main

import (
	"fmt"
	"sync"
	"taskpool-demo/mypool"
)

func main() {
	var (
		counter int
		mu      sync.Mutex
	)

	pool := mypool.New(5)

	for i := 0; i < 1000; i++ {
		pool.Submit(func() {
			mu.Lock()
			counter++
			mu.Unlock()
		})
	}

	pool.Wait()
	pool.Close()

	fmt.Println("最终结果:", counter)
}
