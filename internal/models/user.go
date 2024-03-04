package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model         // adds ID, created_at etc.
	Login       string `json:"login" gorm:"unique"`
	Password    string `json:"password"`
	API         string `json:"apikey" gorm:"unique"`
	TunnelCount int    `json:"tunnelcount" gorm:"default:0"`
	TunnelLimit int    `json:"tunnellimit" gorm:"default:3"`
	Tunnels     []Tunnel
}
