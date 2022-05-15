package db

import (
	"fmt"
	"sync"

	"github.com/Revolyssup/screenlimit/db/actions"
)

type Event struct {
	Timestamp  string       `gorm:"timestamp"`
	Action     string       `gorm:"action"`
	ActionType actions.Type `gorm:"type"`
}

type EventStore struct {
	db *DB
}

func NewEventsStore(time string, db *DB) *EventStore {
	ev := EventStore{
		db: db,
	}
	ev.Add(time, "initialized the events database", actions.Initialize, "")
	return &ev
}

type seenpids struct {
	s  map[string]bool
	mx sync.Mutex
}

var seenpidsingleton = seenpids{
	s: make(map[string]bool),
}

func (s *seenpids) add(pid string) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.s[pid] = true
}
func (s *seenpids) isThere(pid string) bool {
	s.mx.Lock()
	defer s.mx.Unlock()
	return s.s[pid]
}

//Pass pid as empty, it is not stored in the database and only used to uniquely identify the events
func (e *EventStore) Add(time string, action string, at actions.Type, pid string) {
	if pid != "" && seenpidsingleton.isThere(pid) {
		return
	}
	seenpidsingleton.add(pid)
	e.db.mx.Lock()
	defer e.db.mx.Unlock()
	fmt.Println("starting to add")
	err := e.db.Create(&Event{
		Timestamp:  time,
		Action:     action,
		ActionType: at,
	}).Error
	if err != nil {
		fmt.Println("err ", err.Error())
		return
	}
	fmt.Println("added event ", action)
	return
}

func (e *EventStore) Get(pagesize int, offset int, t actions.Type) (ev []Event, err error) {
	if t == "" {
		e.db.mx.Lock()
		err = e.db.Limit(pagesize).Offset(offset).Find(&ev).Error
		e.db.mx.Unlock()
	} else {
		event := []Event{}
		e.db.mx.Lock()
		err = e.db.Where("type = ?", t).Limit(pagesize).Offset(offset).Find(&event).Error
		e.db.mx.Unlock()
	}
	return
}
