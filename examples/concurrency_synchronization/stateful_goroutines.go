package main

import (
	"fmt"
	"math/rand"
	"sync/atomic"
	"time"
)

/* STATEFUL GOROUTINES
Stateful goroutine encapsulated state (data) and manages it internally (within a goroutine), allowing safe access without requiring locks or mutexes.
- Instead of directly sharing state, pass messages or tasks to the goroutine via channels, which handles the state internally
func statefulWorker(jobs <-chan int, results chan<- int) {
    var state int
    for job := range jobs {
        state += job
        results <- state
    }
}

jobs := make(chan int, 10)
results := make(chan int, 10)
go statefulWorker(jobs, results)

- By isolating the state within a goroutine it avoids the need for synchronization mechanisms like mutexes, providing a simpler concurrency model.
- Manages counters, session data, or any mutable state safely across multiple requests or tasks.
*/

type readOp struct {
	key  int
	resp chan int
}
type writeOp struct {
	key  int
	val  int
	resp chan bool
}

func stateful_goroutines() {

	var readOps uint64
	var writeOps uint64

	reads := make(chan readOp)
	writes := make(chan writeOp)

	go func() {
		var state = make(map[int]int)
		for {
			select {
			case read := <-reads:
				read.resp <- state[read.key]
			case write := <-writes:
				state[write.key] = write.val
				write.resp <- true
			}
		}
	}()

	for r := 0; r < 100; r++ {
		go func() {
			for {
				read := readOp{
					key:  rand.Intn(5),
					resp: make(chan int)}
				reads <- read
				<-read.resp
				atomic.AddUint64(&readOps, 1)
				time.Sleep(time.Millisecond)
			}
		}()
	}

	for w := 0; w < 10; w++ {
		go func() {
			for {
				write := writeOp{
					key:  rand.Intn(5),
					val:  rand.Intn(100),
					resp: make(chan bool)}
				writes <- write
				<-write.resp
				atomic.AddUint64(&writeOps, 1)
				time.Sleep(time.Millisecond)
			}
		}()
	}

	time.Sleep(time.Second)

	readOpsFinal := atomic.LoadUint64(&readOps)
	fmt.Println("readOps:", readOpsFinal)
	writeOpsFinal := atomic.LoadUint64(&writeOps)
	fmt.Println("writeOps:", writeOpsFinal)
}
