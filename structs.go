//structs.go

package main

import "time"

//DriverLocation contain the position of one driver
type DriverLocation struct {
	DriverID    string    `json:"driverID"`
	Latitude    string    `json:"latitude"`
	Longitude   string    `json:"longitude"`
	CreatedDate time.Time `json:"createdDate"`
}
