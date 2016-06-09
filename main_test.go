// main.go (web-server)

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/bitly/go-nsq"
)

func EnqueueTest() {
	fmt.Println("Insert some driver location in NSQ")

	var driverLocations = []DriverLocation{
		DriverLocation{
			DriverID:    1,
			Latitude:    48.8566,
			Longitude:   2.3522,
			CreatedDate: time.Now(),
		},
		DriverLocation{
			DriverID:    1,
			Latitude:    48.8544,
			Longitude:   2.3521,
			CreatedDate: time.Now().Add(time.Second),
		},
		DriverLocation{
			DriverID:    1,
			Latitude:    48.8544,
			Longitude:   2.3520,
			CreatedDate: time.Now().Add(2 * time.Second),
		},
	}

	config := nsq.NewConfig()
	w, _ := nsq.NewProducer(NSQconnnection, config)

	for _, u := range driverLocations {
		location, _ := json.Marshal(u)
		err := w.Publish("topic_location", []byte(location))
		if err != nil {
			log.Panic("Could not connect")
		}
	}

	w.Stop()
}

func UnqueueTest() {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	var m DriverLocation

	config := nsq.NewConfig()
	q, _ := nsq.NewConsumer(NSQstream, "Unqueue_test", config)
	q.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		log.Printf("Got a message: %v", json.Unmarshal(message.Body, &m))
		wg.Done()
		return nil
	}))

	err := q.ConnectToNSQLookupd(NSQconnnection)
	if err != nil {
		log.Panic("Could not connect")
	}
	wg.Wait()
}

func TestInsertDriverLocation(t *testing.T) {
	EnqueueTest()
}
