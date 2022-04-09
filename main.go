package main

import (
	"fmt"
	"time"

	"github.com/Revolyssup/screenlimit/action"
	"github.com/Revolyssup/screenlimit/db"
	"github.com/Revolyssup/screenlimit/db/actions"
	"github.com/Revolyssup/screenlimit/server"
)

const PASS = "default"
const PORT = "1401"
const role = "child"

func main() {
	database, err := db.NewDB()
	if err != nil {
		fmt.Println("could not initialize store ", err.Error())
		return
	}
	store := db.NewRoleStore(role, PASS, database)
	_, err = store.SetRole(role, PASS)
	if err != nil {
		fmt.Println("could not set password for children due to: ", err.Error())
		return
	}
	events := db.NewEvents(time.Now().GoString(), database)
	rs := action.NewDialog(10, store)
	go server.Run(PORT, store, events)
	for {
		ch := make(chan bool, 1)
		fmt.Println("Enter password in 10 seconds or pc will reboot")
		select {
		case <-time.After(10 * time.Second):
			fmt.Println("10 sec over")
			action.RunCron(ch, rs)
			if <-ch {
				events.Add(time.Now().GoString(), "Child entered the password succesfully", actions.Child)
				continue
			}
			events.Add(time.Now().GoString(), "Child entered the incorrect password", actions.Child)
			time.Sleep(time.Second * 2)
			events.Add(time.Now().GoString(), "System rebooted", actions.System)
			panic("rebooted")
			// err := exec.Command("reboot").Run()
			// if err != nil {
			// 	fmt.Println("Could not reboot ", err.Error())
			// 	return
			// }
		}
	}

}
