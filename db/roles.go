package db

import (
	"fmt"
)

type Roles struct {
	Role string `gorm:"role"`
	Pass string `gorm:"pass"`
	db   *DB
}

func NewRole(role string, pass string, db *DB) (*Roles, error) {
	st := &Roles{}
	err := db.Where("role = ?", role).First(st).Error
	if err != nil || st.Pass == "" {
		fmt.Println("password ", st.Pass)
		st := &Roles{
			Role: role,
			Pass: pass,
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
func (s *Roles) GetPassword() string {
	s.db.mx.Lock()
	defer s.db.mx.Unlock()
	st := s.Pass
	return st
}

func (s *Roles) SetPassword(role string, pass string) {
	s.db.mx.Lock()
	defer s.db.mx.Unlock()
	fmt.Println("setting new pass")
	s.Pass = pass
	s.Role = role
	err := s.db.Model(&Roles{}).Where("role = ?", role).Update("pass", pass).Error
	if err != nil {
		fmt.Println("err ", err.Error())
	}
}
