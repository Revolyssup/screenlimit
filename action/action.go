package action

import (
	"fmt"
	"time"

	"github.com/Revolyssup/screenlimit/db"
)

type Cron interface {
	Run(chan bool) //Will send a true/false signal
}

func RunCron(ch chan bool, c Cron) {
	c.Run(ch)
}

type Dialog struct {
	timer int
	store *db.RoleStore
}

func NewDialog(timer int, store *db.RoleStore) *Dialog {
	return &Dialog{
		timer: timer,
		store: store,
	}
}

func (r *Dialog) Run(ch chan bool) {
	fmt.Println("going to open dialog")
	pass := make(chan info, 10)
	go func() {
		fmt.Println("Enter password in ", r.timer, " seconds")
		select {
		case <-time.After(time.Second * time.Duration(r.timer)):
			ch <- false
			return
		case pswrd := <-pass:
			role, err := r.store.GetRole("child")
			if err != nil || role == nil {
				fmt.Println(err.Error())
				ch <- false
				return
			}
			if pswrd.pass == role.Pass {
				fmt.Println("password is correct which is ", pswrd.pass)
				ch <- true
			} else {
				fmt.Printf("password is incorrect.expected %s got %s ", role.Pass, pswrd.pass)
				ch <- false
			}
			return
		}
	}()

	RunApp(pass)
}
