package db

import (
	"fmt"
	"sync"
)

type Role struct {
	Role string `gorm:"role"`
	Pass string `gorm:"pass"`
}
type RoleStore struct {
	db *DB
	mx sync.Mutex
}

const ADMIN = "admin"

func NewRoleStore(db *DB) *RoleStore {
	return &RoleStore{db: db}
}
func (s *RoleStore) GetAdminPassword() string {
	s.mx.Lock()
	defer s.mx.Unlock()
	role, err := s.GetRole(ADMIN)
	if err != nil {
		return ""
	}
	return role.Pass
}
func (s *RoleStore) SetAdminPassword(pass string) {
	s.mx.Lock()
	defer s.mx.Unlock()
	_, err := s.SetRole(ADMIN, pass)
	if err != nil {
		return
	}
}
func (s *RoleStore) GetRole(role string) (*Role, error) {
	s.db.mx.Lock()
	defer s.db.mx.Unlock()
	r := Role{}
	err := s.db.Where("role = ?", role).Find(&r).Error
	return &r, err
}
func (s *RoleStore) GetRoles() ([]Role, error) {
	s.db.mx.Lock()
	defer s.db.mx.Unlock()
	r := []Role{}
	err := s.db.Find(&r).Error
	return r, err
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
