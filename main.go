package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bitly/go-nsq"
	"github.com/gorilla/mux"
	"gopkg.in/redis.v3"
)

//DriverLocationHandler retrieves the last location of a customer according the time frame given in parameter
func DriverLocationHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("\t%s",
		r.RequestURI)

	// Read route parameter
	vars := mux.Vars(r)
	driverID := vars["driverID"]
	queryMinutes := r.URL.Query().Get("minutes")
	if queryMinutes == "" {
		log.Printf("queryMinutes missing")
		w.WriteHeader(http.StatusInternalServerError)
	}

	minutes, err := strconv.ParseInt(queryMinutes, 10, 16)

	if err != nil {
		log.Printf("Received bad query Minutes parameter: %s.", queryMinutes)
		w.WriteHeader(http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")
	response := []byte(GetDriverLocation(driverID, minutes))

	w.Write(response)
}

//GetDriverLocation returns the location of customer in the last N minutes
func GetDriverLocation(key string, minutes int64) string {
	client := RedisInit()

	if minutes < 1 {
		minutes = 1
	}
	//According to specification Drivers are pushing their locations every 5 seconds so in order to retrieves them according minutes, we do minutes * 12.
	take := minutes * 12

	r, err := client.LRange(key, 0, take).Result()
	if err != nil {
		log.Printf("Could not get from Redis.")
		panic(err)
	}
	result := "[" + strings.Join(r, ",") + "]"
	return result
}

//RedisInit connects to Redis
func RedisInit() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     REDISconnection,
		Password: "",
		DB:       0,
	})
	return client
}

//StoreLocation insert driver location in Redis
func StoreLocation(d DriverLocationRequest) {
	client := RedisInit()

	var m DriverLocationResult
	m.DriverID = d.DriverID
	m.LocationResult.Latitude = d.LocationRequest.Latitude
	m.LocationResult.Longitude = d.LocationRequest.Longitude
	m.LocationResult.UpdatedAt = d.LocationRequest.UpdatedAt.Format(time.RFC3339)

	key := strconv.Itoa(m.DriverID)
	v, _ := json.Marshal(m.LocationResult)
	val := string(v)

	err := client.RPush(key, val).Err()
	if err != nil {
		log.Printf("Could not push in Redis.")
		panic(err)
	}
}

//GetLocationFromGateway fetchs messages in NSQ and insert them in Redis
func GetLocationFromGateway() {
	wg := &sync.WaitGroup{}
	wg.Add(4)

	var message DriverLocationRequest

	config := nsq.NewConfig()
	q, _ := nsq.NewConsumer(NSQstream, "worker_location_service", config)
	q.AddHandler(nsq.HandlerFunc(func(m *nsq.Message) error {
		json.Unmarshal(m.Body, &message)

		//Store messages in Redis
		StoreLocation(message)

		log.Printf(string(m.Body))
		return nil
	}))

	err := q.ConnectToNSQD(NSQconnnection)
	if err != nil {
		log.Printf("Could not connect to NSQ.")
		panic(err)
	}
	wg.Wait()
	wg.Done()
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/drivers/{driverID:[0-9]+}/coordinates", DriverLocationHandler).Methods("GET")
	log.Printf("Server started and listening on port %d.", 8001)
	http.Handle("/", r)

	//Collects messages in NSQ to store them in REDIS
	go GetLocationFromGateway()

	log.Fatal(http.ListenAndServe(":8001", r))
}
