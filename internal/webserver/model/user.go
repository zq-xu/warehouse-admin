package model

import (
	"gorm.io/gorm"

	"zq-xu/warehouse-admin/pkg/store"
)

const (
	UserTableName = "user"
)

// description:"default the column ID is the primaryKey
type User struct {
	gorm.Model

	Name     string
	Alias    string
	Password string
	Comment  string
	Role     int
	Status   int
}

func (u *User) TableName() string {
	return UserTableName
}

func init() {
	store.RegisterTable(&User{})
}
