package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/Revolyssup/screenlimit/action"
	"github.com/Revolyssup/screenlimit/db"
	"github.com/Revolyssup/screenlimit/policy"
	"github.com/Revolyssup/screenlimit/server"
	"github.com/Revolyssup/screenlimit/sysstats"
)

const PASS = "default"
const PORT = "1401"
const role = "child"

var appsToMonitor = []string{"brave", "slack"}

func main() {
	fmt.Println("Starting server...")
	//initialize the database
	database, err := db.NewDB()
	if err != nil {
		fmt.Println("could not initialize store ", err.Error())
		return
	}

	//initialize the role store
	roleStore := db.NewRoleStore(database)
	_, err = roleStore.SetRole(role, PASS)
	if err != nil {
		fmt.Println("could not set password for children due to: ", err.Error())
		return
	}

	//initialize the event store which will be used by any thing that wants to log events into the database
	eventStore := db.NewEventsStore(fmt.Sprintf("%s", time.Now().Format("01-02-2006 15:04:05")), database)
	// start the aggregator
	stats := sysstats.New(&appsToMonitor, eventStore)
	go stats.Run()
	go server.Run(PORT, roleStore, eventStore, PASS, &appsToMonitor, stats) //start the server

	// start the cron job that prompt after a certain time
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		database.Close()
	}()
	action.RunCron(10, roleStore, eventStore, &policy.PolicyRequest{
		Action: policy.Default,
		Type:   policy.RESTART,
	})
}
