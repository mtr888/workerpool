package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
	"workerpool/entity"
)

const workerCount, usersCount = 5, 1000

func main() {
	rand.Seed(time.Now().Unix())

	startTime := time.Now()
	userCh := make(chan entity.User)
	wg := &sync.WaitGroup{}

	go entity.GenerateUsers(usersCount, userCh)

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go entity.Worker(i+1, userCh, wg)
	}
	wg.Wait()
	fmt.Printf("DONE! Time Elapsed: %.2f seconds\n", time.Since(startTime).Seconds())
}
