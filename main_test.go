// main.go (web-server)

package main

import (
	"testing"
	"time"
)

//TestMapping ensures that Mapping method returns the correct result
func TestMapping(t *testing.T) {
	var d = DriverLocationRequest{
		DriverID: 1,
		LocationRequest: LocationRequest{
			Latitude:  48.8566,
			Longitude: 2.3522,
			UpdatedAt: time.Now(),
		}}

	r := Mapping(d)

	//Time parsing uses the same layout values as `Format`.
	_, e := time.Parse(
		time.RFC3339,
		r.LocationResult.UpdatedAt)

	if e != nil {
		t.Fail()
	}
}
