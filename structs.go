//structs.go

package main

import "time"

//Location contains a location at a specific time
type Location struct {
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	UpdatedAt time.Time `json:"updated_at"`
}

//DriverLocation contains the position of one driver
type DriverLocation struct {
	DriverID int `json:"driverID"`
	Location Location
}
