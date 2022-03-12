package action

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Revolyssup/screenlimit/db"
)

type Cron interface {
	Run() bool
}

func RunCron(c Cron) bool {
	return c.Run()
}

type Restarter struct {
	timer int
	store *db.Store
}

func NewRestarter(timer int, store *db.Store) *Restarter {
	return &Restarter{
		timer: timer,
		store: store,
	}
}
func (r *Restarter) Run() bool {
	temp := make(chan bool, 1)
	go func() {
		for i := 0; i < 3; i++ {
			fmt.Println(3-i, " attempts left")
			reader := bufio.NewReader(os.Stdin)
			fmt.Println("PASSWORD>")
			text, _ := reader.ReadString('\n')
			// convert CRLF to LF
			text = strings.Replace(text, "\n", "", -1)
			fmt.Println("text is ", text)
			if text == r.store.GetPassword() {
				temp <- true
			}
		}
		temp <- false
	}()
	fmt.Println("Enter password in ", r.timer, " seconds")
	select {
	case <-time.After(time.Second * 10):
		return false
	case this := <-temp:
		return this
	}
}
