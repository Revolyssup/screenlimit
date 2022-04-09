package db

import (
	"fmt"

	"github.com/Revolyssup/screenlimit/db/actions"
)

type Events struct {
	Timestamp  string       `gorm:"timestamp"`
	Action     string       `gorm:"action"`
	ActionType actions.Type `gorm:"type"`
	db         *DB
}

func NewEvents(time string, db *DB) *Events {
	ev := Events{
		db: db,
	}
	ev.Add(time, "initialized the events database", actions.Initialize)
	return &ev
}

func (e *Events) Add(time string, action string, at actions.Type) error {
	e.db.mx.Lock()
	defer e.db.mx.Unlock()
	fmt.Println("starting to add")
	err := e.db.Create(&Events{
		Timestamp:  time,
		Action:     action,
		ActionType: at,
	}).Error
	if err != nil {
		fmt.Println("err ", err.Error())
		return err
	}
	return nil
}

func (e *Events) Get(pagesize int, offset int, t actions.Type) (ev []Events, err error) {
	if t == "" {
		e.db.mx.Lock()
		err = e.db.Limit(pagesize).Offset(offset).Find(&ev).Error
		e.db.mx.Unlock()
	} else {
		event := []Events{}
		e.db.mx.Lock()
		err = e.db.Where("type = ?", t).Limit(pagesize).Offset(offset).Find(&event).Error
		e.db.mx.Unlock()
	}
	return
}
