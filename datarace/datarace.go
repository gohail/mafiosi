package main

import "sync"

// Запустите с помощью "go run -race datarace.go".

var globalX int

func main() {
	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			for i := 0; i < 1000000; i++ {
				globalX++
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
