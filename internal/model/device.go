package model

import "time"

type DeviceState string

const (
	StateAvailable DeviceState = "available"
	StateInUse     DeviceState = "in-use"
	StateInactive  DeviceState = "inactive"
)

type Device struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Brand     string      `json:"brand"`
	State     DeviceState `json:"state"`
	CreatedAt time.Time   `json:"created_at"`
}
