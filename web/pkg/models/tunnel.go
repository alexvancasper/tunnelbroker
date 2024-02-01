package models

import "gorm.io/gorm"

type Tunnel struct {
	gorm.Model
	Configured bool   `json:"configured"`
	UserID     uint   `json:"userid"`
	P2P        string `json:"p2p"`
	PD         string `json:"pd"`
	IPv4Remote string `json:"ipv4remote"`
	IPv4Local  string `json:"ipv4local"`
	TunnelName string `json:"tunnelname"`
}
