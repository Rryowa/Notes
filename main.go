package main

import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

type CLI struct {
	nWorkers       uint64
	maxConcurrency uint64 // Allowed to run at the same time
	sem            chan struct{}
}

//***How to DYNAMICALLY change num of goroutines running in paraller***

// NOTE: IF USE (c *CLI) receiver, 'c.sem <- struct{}' will cause a deadlock, after changing max!
// NOTE: channels passed by reference.
// NOTE: if you use make, new or &, you can pass it to another function without copying the underlying data.
func worker(id uint64, sem chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	// Acquire semaphore
	fmt.Printf("Worker %d: Waiting to acquire semaphore\n", id)
	sem <- struct{}{}

	// Do work
	//log.Printf("Worker %d: Semaphore acquired, running\n", id)
	time.Sleep(5 * time.Second)

	// Release semaphore
	<-sem
	log.Printf("Worker %d: Semaphore released\n", id)
}

func main() {
	c := &CLI{nWorkers: 10, maxConcurrency: 2}

	// Create a buffered channel with a capacity of maxConcurrency
	c.sem = make(chan struct{}, c.maxConcurrency)

	var wg sync.WaitGroup

	// We start 10 goroutines but only 2 of them will run in parallel
	for i := uint64(1); i <= atomic.LoadUint64(&c.nWorkers); i++ {
		if i == 6 {
			time.Sleep(5 * time.Second)
			atomic.StoreUint64(&c.maxConcurrency, 1)
			c.sem = make(chan struct{}, 1)
			log.Println("Max set: ", cap(c.sem), "Worker:", i)
		}
		wg.Add(1)

		//NOTE ABOUT COPYING c.sem CHANNEL IN worker
		go worker(i, c.sem, &wg)
	}
	wg.Wait()
	close(c.sem) // Remember to close the channel once done
	fmt.Println("All workers have completed")
}
