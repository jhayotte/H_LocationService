//structs.go

package main

import "time"

//LocationResult contains a location at a specific time
type LocationResult struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	UpdatedAt string  `json:"updated_at"`
}

//DriverLocationResult contains the position of one driver
type DriverLocationResult struct {
	DriverID       int `json:"driverID"`
	LocationResult LocationResult
}

//LocationRequest contains a location at a specific time
type LocationRequest struct {
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	UpdatedAt time.Time `json:"updated_at"`
}

//DriverLocationRequest contains the position of one driver
type DriverLocationRequest struct {
	DriverID        int `json:"driverID"`
	LocationRequest LocationRequest
}
