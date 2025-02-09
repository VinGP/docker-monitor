package model

import "time"

type ContainerStatus struct {
	PingTime    *float64   `json:"ping_time"                        db:"ping_time"`
	LastSuccess *time.Time `json:"last_success"                     db:"last_success"`
	IPAddress   string     `json:"ip_address"   validate:"required" db:"ip"`
}
