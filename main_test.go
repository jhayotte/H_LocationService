// main.go (web-server)

package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/bitly/go-nsq"
	"gopkg.in/redis.v3"
)

func EnqueueTest() {
	fmt.Println("Insert some driver's locations in NSQ")

	var driverLocations = []DriverLocationRequest{
		DriverLocationRequest{
			DriverID: 1,
			LocationRequest: LocationRequest{
				Latitude:  48.8566,
				Longitude: 2.3522,
				UpdatedAt: time.Now(),
			},
		},
		DriverLocationRequest{
			DriverID: 1,
			LocationRequest: LocationRequest{
				Latitude:  48.8544,
				Longitude: 2.3521,
				UpdatedAt: time.Now().Add(5 * time.Second),
			},
		},
		DriverLocationRequest{
			DriverID: 1,
			LocationRequest: LocationRequest{
				Latitude:  48.8544,
				Longitude: 2.3520,
				UpdatedAt: time.Now().Add(10 * time.Second),
			},
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

func TestInsertDriverLocation(t *testing.T) {
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
		fmt.Println("Could not connect to REDIS")
		t.Fail()
	}
	fmt.Println(pong)
}

func TestAddLocationInRedis(t *testing.T) {
	timeRFC := time.Now().Format(time.RFC3339)

	d := DriverLocationResult{
		DriverID: 1,
		LocationResult: LocationResult{Latitude: 48.8566,
			Longitude: 2.3522,
			UpdatedAt: timeRFC,
		},
	}

	client := redis.NewClient(&redis.Options{
		Addr:     REDISconnection,
		Password: "",
		DB:       0,
	})

	key := strconv.Itoa(d.DriverID)
	v, _ := json.Marshal(d.LocationResult)
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
