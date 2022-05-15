package action

import (
	"fmt"
	"sync"
	"time"

	"github.com/Revolyssup/screenlimit/db"
	"github.com/Revolyssup/screenlimit/db/actions"
)

//Run the cronjob after every 10 seconds
func RunCron(t int, rs *db.RoleStore, es *db.EventStore) {
	ch := make(chan bool, 1)
	pass := make(chan info, 10)
	d := newDialog(t, rs, ch)
	time.Sleep(100 * time.Second)
	var wg sync.WaitGroup
	go func() {
		wg.Add(1)
		for {
			fmt.Println("Enter password in 10 seconds or pc will reboot")
			select {
			case <-time.After(60 * time.Second):
				fmt.Println("10 sec over")
				d.run(pass)
				fmt.Println("that was done")
				if <-ch {
					es.Add(time.Now().GoString(), "Child entered the password succesfully", actions.Child, "")
					continue
				}
				fmt.Println("but the wait was never over")
				es.Add(time.Now().GoString(), "Child entered the incorrect password", actions.Child, "")
				time.Sleep(time.Second * 2)
				es.Add(time.Now().GoString(), "System rebooted", actions.System, "")
				panic("rebooted")
				// err := exec.Command("reboot").Run()
				// if err != nil {
				// 	fmt.Println("Could not reboot ", err.Error())
				// 	return
				// }
			}
		}
	}()
	getWindowSingleton(pass).ShowAndRun()
	wg.Done()
}

type dialog struct {
	timer int
	store *db.RoleStore
	ch    chan bool
}

func newDialog(timer int, store *db.RoleStore, ch chan bool) *dialog {
	return &dialog{
		timer: timer,
		store: store,
		ch:    ch,
	}
}

func (r *dialog) run(pass chan info) {
	w := getWindowSingleton(pass)
	go func() {
		fmt.Println("Enter password in ", r.timer, " seconds")
		select {
		case <-time.After(time.Second * time.Duration(r.timer)):
			w.Hide()
			fmt.Println("sending false")
			r.ch <- false
			return
		case pswrd := <-pass:
			w.Hide()
			role, err := r.store.GetRole("child")
			if err != nil || role == nil {
				fmt.Println(err.Error())
				r.ch <- false
				return
			}
			if pswrd.pass == role.Pass {
				fmt.Println("password is correct which is ", pswrd.pass)
				r.ch <- true
			} else {
				fmt.Printf("password is incorrect.expected %s got %s ", role.Pass, pswrd.pass)
				r.ch <- false
			}
			return
		}
	}()
	w.Show()
}
