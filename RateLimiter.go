package ratelimiter

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type job func() // тип функция

func worker(id int, jobs <-chan job, wg *sync.WaitGroup, cPS *uint64, rPS uint64) {
	defer wg.Done()
	for fn := range jobs {
		for atomic.LoadUint64(cPS) > rPS {
			fmt.Println("Sleeping")
			time.Sleep(time.Millisecond * 20)
		}
		atomic.AddUint64(cPS, 1)
		//fmt.Println("worker", id, "started de job")
		fn()
		//fmt.Println("worker", id, "finished da job")
	}
	//fmt.Println("Done - worker", id )
}

func RateLimiter(jobs chan job, totalWorkers int, rPM uint64, wg *sync.WaitGroup) {
	// rate received as rPM (rate per minute), but calculations actually in
	// rPS = rPM/60 for smoother operation
	defer wg.Done()
	var cpS uint64 // counterPerS

	go func() {
		for {
			time.Sleep(time.Second)
			fmt.Println(cpS)
			atomic.CompareAndSwapUint64(&cpS, cpS, 0)
		}
	}() // resetting the counter

	for w := 1; w <= totalWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, wg, &cpS, rPM/60)
	}
}
