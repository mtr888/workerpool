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
	// time.Sleep(time.Millisecond * 1)
	return nil
}

func GenerateUsers(count int, userCh chan User) {
	defer close(userCh)

	for i := 0; i < count; i++ {
		userCh <- User{
			Id:    i + 1,
			Email: fmt.Sprintf("user%d@company.com", i+1),
			Logs:  generateLogs(rand.Intn(1000)),
		}
		fmt.Printf("generated user %d\n", i+1)
		// time.Sleep(time.Millisecond * 1)
	}
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

func Worker(id int, userCh <-chan User, wg *sync.WaitGroup) {
	defer wg.Done()
	for user := range userCh {
		err := user.SaveUserInfo()
		if err != nil {
			fmt.Println(err)
		}
	}
}
