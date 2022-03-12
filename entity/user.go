package entity

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

var Actions = []string{
	"logged in",
	"logged out",
	"create record",
	"delete record",
	"update recodr",
}

type User struct {
	Id    int
	Email string
	Logs  []logItem
}

type logItem struct {
	Action    string
	Timestamp time.Time
}

func (u *User) GetActivityInfo() string {
	output := fmt.Sprintf("UID: %d; Email: %s;\nActivity Log:\n", u.Id, u.Email)
	for index, item := range u.Logs {
		output += fmt.Sprintf("%d. [%s] at %s\n", index, item.Action, item.Timestamp.Format(time.RFC3339))
	}

	return output
}

func (u *User) SaveUserInfo() error {
	fmt.Printf("WRITING FILE FOR UID %d\n", u.Id)

	filename := fmt.Sprintf("../logs/uid%d.txt", u.Id)
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	file.WriteString(u.GetActivityInfo())
	return nil
}

func GenerateUsers(done <-chan struct{}, count int) <-chan User {
	usersChannel := make(chan User)
	go func(count int) {
		defer func() {
			close(usersChannel)
			fmt.Println("generateUsers is done")
		}()
		for i := 0; i < count; i++ {

			select {
			case <-done:
				return
			case usersChannel <- User{
				Id:    i + 1,
				Email: fmt.Sprintf("user%d@company.com", i+1),
				Logs:  generateLogs(rand.Intn(2)),
			}:
			}
		}
	}(count)
	return usersChannel
}

func generateLogs(count int) []logItem {
	logs := make([]logItem, count)
	for i := 0; i < count; i++ {
		logs[i] = logItem{
			Action:    Actions[rand.Intn(len(Actions)-1)],
			Timestamp: time.Now(),
		}
	}
	return logs
}

func WorkerPool(done chan struct{}, usersChannel <-chan User, workerCount int) <-chan error {
	errorChannel := make(chan error)
	wg := sync.WaitGroup{}
	wg.Add(workerCount)

	defer fmt.Println("workerpool is done")

	for i := 0; i < workerCount; i++ {
		i := i
		go func() {
			defer func() {
				wg.Done()
				fmt.Printf("worker #%d is done\n", i+1)
			}()

			for user := range usersChannel {
				err := user.SaveUserInfo()
				if err != nil {
					errorChannel <- err
				}

			}

			// for {
			// 	select {
			// 	case <-done:
			// 		return
			// 	case user := <-usersChannel:
			// 		select {
			// 		case <-done:
			// 			return
			// 		case errorChannel <- func() error {
			// 			return user.SaveUserInfo()
			// 		}():
			// 		}
			// 	}
			// }
		}()
	}

	go func() {
		wg.Wait()
		close(errorChannel)
	}()

	return errorChannel
}
