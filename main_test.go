// main.go (web-server)

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/bitly/go-nsq"
	"gopkg.in/redis.v3"
)

func EnqueueTest() {
	fmt.Println("Insert some driver location in NSQ")

	var driverLocations = []DriverLocation{
		DriverLocation{
			DriverID: 1,
			Location: Location{
				Latitude:  48.8566,
				Longitude: 2.3522,
				UpdatedAt: time.Now(),
			},
		},
		DriverLocation{
			DriverID: 1,
			Location: Location{
				Latitude:  48.8544,
				Longitude: 2.3521,
				UpdatedAt: time.Now().Add(time.Second),
			},
		},
		DriverLocation{
			DriverID: 1,
			Location: Location{
				Latitude:  48.8544,
				Longitude: 2.3520,
				UpdatedAt: time.Now().Add(2 * time.Second),
			},
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

func TestInsertDriverLocation(t *testing.T) {
	EnqueueTest()
	EnqueueTest()
	EnqueueTest()
}

func TestRedisPingPong(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     REDISconnection,
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping().Result()
	if err != nil {
		t.Fail()
	}
	fmt.Println(pong)
}

func TestAddLocationInRedis(t *testing.T) {
	d := DriverLocation{
		DriverID: 1,
		Location: Location{Latitude: 48.8566,
			Longitude: 2.3522,
			UpdatedAt: time.Now(),
		},
	}

	client := redis.NewClient(&redis.Options{
		Addr:     REDISconnection,
		Password: "",
		DB:       0,
	})

	key := strconv.Itoa(d.DriverID)
	v, _ := json.Marshal(d.Location)
	val := string(v)

	err := client.RPush(key, val).Err()
	if err != nil {
		panic(err)
	}
	t.Log("LPush success")

	redisLLenResult, err := client.LLen(key).Result()
	if err != nil {
		panic(err)
	}
	client.LRange(key, 0, redisLLenResult)

	t.Log("Get success: " + "9999")
}
