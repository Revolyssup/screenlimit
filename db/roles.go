package db

import (
	"fmt"
)

type Role struct {
	Role string `gorm:"role"`
	Pass string `gorm:"pass"`
}
type RoleStore struct {
	db *DB
}

func NewRoleStore(role string, pass string, db *DB) *RoleStore {
	return &RoleStore{db: db}
}
func (s *RoleStore) GetRole(role string) (*Role, error) {
	s.db.mx.Lock()
	defer s.db.mx.Unlock()
	r := Role{}
	err := s.db.Where("role = ?", role).Find(&r).Error
	return &r, err
}

func (s *RoleStore) SetRole(role string, pass string) (*Role, error) {
	s.db.mx.Lock()
	defer s.db.mx.Unlock()
	st := &Role{}
	err := s.db.Where("role = ?", role).First(st).Error
	if err != nil || st.Pass == "" {
		fmt.Println("password ", st.Pass)
		st := &Role{
			Role: role,
			Pass: pass,
		}
		err = s.db.Create(st).Error
		return st, err
	}
	if err != nil {
		fmt.Println("not initializing store ", err.Error())
		return nil, err
	}
	fmt.Println("ashish,pass is ", st.Pass)
	return st, nil
}
