// main.go (web-server)

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/bitly/go-nsq"
)

const NSQconnnection string = "127.0.0.1:4150"

//DriverLocation contain the position of one driver
type DriverLocation struct {
	DriverID    string    `json:"driverID"`
	Latitude    string    `json:"latitude"`
	Longitude   string    `json:"longitude"`
	CreatedDate time.Time `json:"createdDate"`
}

func LoadTest() {
	fmt.Println("Insert some driver location in NSQ")

	var driverLocations = []DriverLocation{
		DriverLocation{
			DriverID:    "1",
			Latitude:    "48.8566",
			Longitude:   "2.3522",
			CreatedDate: time.Now(),
		},
		DriverLocation{
			DriverID:    "1",
			Latitude:    "48.8544",
			Longitude:   "2.3521",
			CreatedDate: time.Now().Add(time.Second),
		},
		DriverLocation{
			DriverID:    "1",
			Latitude:    "48.8544",
			Longitude:   "2.3520",
			CreatedDate: time.Now().Add(2 * time.Second),
		},
	}

	config := nsq.NewConfig()
	w, _ := nsq.NewProducer(NSQconnnection, config)

	for _, u := range driverLocations {
		location, _ := json.Marshal(u)
		err := w.Publish("write_test", []byte(location))
		if err != nil {
			log.Panic("Could not connect")
		}
	}

	w.Stop()
}

func TestInsertDriverLocation(t *testing.T) {
	LoadTest()
}
