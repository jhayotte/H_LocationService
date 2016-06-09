// main.go (web-server)

package main

import (
	"fmt"
	"log"
	"testing"

	"github.com/bitly/go-nsq"
)

const NSQconnnection string = "127.0.0.1:4150"

//DriverLocation contain the position of one driver
type DriverLocation struct {
	DriverID  string `json:"driverID"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

func LoadTest() {
	fmt.Println("Insert some driver location in NSQ")

	var driverLocations = []DriverLocation{
		DriverLocation{
			DriverID:  "1",
			Latitude:  "48.8566",
			Longitude: "2.3522",
		},
		DriverLocation{
			DriverID:  "1",
			Latitude:  "48.8544",
			Longitude: "2.3521",
		},
		DriverLocation{
			DriverID:  "1",
			Latitude:  "48.8544",
			Longitude: "2.3520",
		},
	}

	config := nsq.NewConfig()
	w, _ := nsq.NewProducer(NSQconnnection, config)

	for _, u := range driverLocations {
		err := w.Publish("write_test", []byte("{\"driver\":"+u.DriverID+",\"latitude\":"+u.Latitude+",\"longitude\": "+u.Longitude+"}"))
		if err != nil {
			log.Panic("Could not connect")
		}
	}

	w.Stop()
}

func TestInsertDriverLocation(t *testing.T) {
	LoadTest()
}
