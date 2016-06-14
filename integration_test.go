// +build integration

package main

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/bitly/go-nsq"
)

//NSQstream is the stream name used in NSQ by Location Service
const NSQstream string = "topic_location"

//NSQconnnection is the connection string to NSQ
const NSQconnnection string = "127.0.0.1:4150"

//REDISconnection is the connection string to REDIS
const REDISconnection string = "127.0.0.1:6379"

//TestInsertDriverLocationInNSQ inserts driver in NSQ
func TestInsertDriverLocationInNSQ(t *testing.T) {
	var driverLocations = []DriverLocation{
		DriverLocation{
			DriverID:  1,
			Latitude:  48.8566,
			Longitude: 2.3522,
			UpdatedAt: time.Now(),
		},
		DriverLocation{
			DriverID:  1,
			Latitude:  48.8544,
			Longitude: 2.3521,
			UpdatedAt: time.Now().Add(5 * time.Second),
		},
		DriverLocation{
			DriverID:  1,
			Latitude:  48.8544,
			Longitude: 2.3520,
			UpdatedAt: time.Now().Add(10 * time.Second),
		},
	}

	config := nsq.NewConfig()
	w, _ := nsq.NewProducer(NSQconnnection, config)

	for _, u := range driverLocations {
		location, _ := json.Marshal(u)
		err := w.Publish("topic_location", []byte(location))

		if err != nil {
			fmt.Println("Could not publish in NSQ")
			panic(err)
		}
	}

	w.Stop()
}
