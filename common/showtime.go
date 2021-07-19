package common

import (
	"golang.org/x/sync/semaphore"
	"log"
	"math/big"
	"sync"
	"time"
)

func ShowTime(c int, t int64, f func(s *semaphore.Weighted, wg *sync.WaitGroup, idx int)) {
	start := time.Now()
	defer func() {
		sub := time.Now().Sub(start)
		log.Printf("耗时：%v", sub)
		var r big.Float
		r.Quo(big.NewFloat(float64(c)), big.NewFloat(sub.Seconds()))
		log.Printf(" hits/sec  %v", r.String())
	}()

	var sem *semaphore.Weighted
	if t > 0 {
		sem = semaphore.NewWeighted(t)
	}
	count := c
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(count)
	for i := 0; i < count; i++ {
		go f(sem, &waitGroup, i)
	}

	waitGroup.Wait()
}
