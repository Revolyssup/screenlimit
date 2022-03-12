package db

import (
	"fmt"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Store struct {
	Role string `gorm:"role"`
	Pass string `gorm:"pass"`
	db   *gorm.DB
	mx   sync.Mutex
}

const role = "ADMIN"

func New(pass string) (*Store, error) {
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&Store{})
	st := &Store{}
	err = db.Where("role = ?", role).First(st).Error
	if err != nil || st.Pass == "" {
		fmt.Println("password ", st.Pass)
		st := &Store{
			Role: role,
			Pass: pass,
			db:   db,
		}
		err = db.Create(st).Error
		return st, err
	}
	if err != nil {
		fmt.Println("not initializing store ", err.Error())
	}
	st.db = db
	fmt.Println("ashish,pass is ", st.Pass)
	return st, err
}
func (s *Store) GetPassword() string {
	s.mx.Lock()
	defer s.mx.Unlock()
	st := s.Pass
	return st
}

func (s *Store) SetPassword(pass string) {
	s.mx.Lock()
	defer s.mx.Unlock()
	fmt.Println("setting new pass")
	s.Pass = pass
	s.Role = role
	err := s.db.Model(&Store{}).Where("role = ?", role).Update("pass", pass).Error
	if err != nil {
		fmt.Println("err ", err.Error())
	}
}
