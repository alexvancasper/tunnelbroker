package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model        // adds ID, created_at etc.
	Login      string `json:"login" gorm:"unique"`
	Password   string `json:"password"`
	API        string `json:"apikey" gorm:"unique"`
	Tunnels    []Tunnel
}
