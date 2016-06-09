//structs.go

package main

import "time"

//DriverLocation contain the position of one driver
type DriverLocation struct {
	DriverID    int       `json:"driverID"`
	Latitude    float32   `json:"latitude"`
	Longitude   float32   `json:"longitude"`
	CreatedDate time.Time `json:"createdDate"`
}
