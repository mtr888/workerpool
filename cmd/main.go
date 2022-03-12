package main

import (
	"fmt"
	"math/rand"
	"time"
	"workerpool/entity"
)

const workerCount, usersCount = 5, 10

func main() {
	rand.Seed(time.Now().Unix())

	startTime := time.Now()

	done := make(chan struct{})
	defer close(done)

	usersChannel := entity.GenerateUsers(done, usersCount)
	errorsChannel := entity.WorkerPool(done, usersChannel, workerCount)

	for err := range errorsChannel {
		fmt.Println("err = ", err)
	}

	fmt.Printf("DONE! Time Elapsed: %.2f seconds\n", time.Since(startTime).Seconds())
	time.Sleep(1 * time.Second)
}
