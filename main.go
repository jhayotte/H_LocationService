/*package main implements the Location service. It fetchs driver's location in
NSQ and stores it in REDIS. An endpoint allows to retrieve them for a specific
driver in a specific time duration
*/
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
	"github.com/rubyist/circuitbreaker"
	"gopkg.in/redis.v3"
)

const (
	//NSQconnnection is the connection string to NSQ
	NSQconnnection string = "127.0.0.1:4150"

	//NSQstream is the stream name used in NSQ by Location Service
	NSQstream string = "topic_location"

	//REDISconnection is the connection string to REDIS
	REDISconnection string = "127.0.0.1:6379"
)

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

// DriverLocationHandler retrieves the last location of a customer according the
// time frame given in parameter
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

//RedisRPush inserts the speicified val in the speicified key
func RedisRPush(m DriverLocationResult) {

	key := strconv.Itoa(m.DriverID)
	v, _ := json.Marshal(m.LocationResult)
	val := string(v)

	client := RedisInit()
	err := client.RPush(key, val).Err()
	if err != nil {
		log.Printf("Could not push in Redis.")
		panic(err)
	}
}

//Mapping of DriverLocationRequest and DriverLocationResult
func Mapping(d DriverLocationRequest) DriverLocationResult {
	var m DriverLocationResult
	m.DriverID = d.DriverID
	m.LocationResult.Latitude = d.LocationRequest.Latitude
	m.LocationResult.Longitude = d.LocationRequest.Longitude
	m.LocationResult.UpdatedAt = d.LocationRequest.UpdatedAt.Format(time.RFC3339)
	return m
}

//UnqueueDriversLocation from NSQ
func UnqueueDriversLocation() {
	wg := &sync.WaitGroup{}
	wg.Add(4)

	var message DriverLocationRequest

	config := nsq.NewConfig()
	q, _ := nsq.NewConsumer(NSQstream, "worker_location_service", config)
	q.AddHandler(nsq.HandlerFunc(func(m *nsq.Message) error {
		json.Unmarshal(m.Body, &message)

		//Format the request in the format wanted
		messageFormatted := Mapping(message)

		//Insert in Redis
		RedisRPush(messageFormatted)

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

//GetDriversLocationFromGateway is wrapped in a circuit breaker.
// its role is to fetch messages in NSQ and insert them in Redis
//
func GetDriversLocationFromGateway() {
	// Creates a circuit breaker that will trip after 10 failures
	// using a time out of 5 seconds
	cb := circuit.NewThresholdBreaker(10)

	cb.Call(func() error {
		// This is where you'll do some remote call
		UnqueueDriversLocation()
		// If it fails, return an error
		return nil
	}, time.Second*5) // This will time out after 5 seconds, which counts as a failure
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/drivers/{driverID:[0-9]+}/coordinates", DriverLocationHandler).Methods("GET")
	log.Printf("Server started and listening on port %d.", 8001)
	http.Handle("/", r)

	//Collects messages in NSQ to store them in REDIS
	go GetDriversLocationFromGateway()

	log.Fatal(http.ListenAndServe(":8001", r))
}
