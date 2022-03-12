package main

import (
	"fmt"
	"time"

	"github.com/Revolyssup/screenlimit/action"
	"github.com/Revolyssup/screenlimit/db"
	"github.com/Revolyssup/screenlimit/server"
)

const PASS = "default"

func main() {
	store, err := db.New(PASS)
	if err != nil {
		fmt.Println("could not initialize store ", err.Error())
		return
	}
	rs := action.NewDialog(10, store)
	go server.Run("1401", store)
	for {
	WAIT:
		ch := make(chan bool, 1)
		fmt.Println("Enter password in 10 seconds or pc will reboot")
		select {
		case <-time.After(10 * time.Second):
			action.RunCron(ch, rs)
			if <-ch {
				goto WAIT
			}
			panic("rebooted")
			// err := exec.Command("reboot").Run()
			// if err != nil {
			// 	fmt.Println("Could not reboot ", err.Error())
			// 	return
			// }
		}
	}

}
