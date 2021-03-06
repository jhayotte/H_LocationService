/*package main implements the Location service. It fetchs driver's location in
NSQ and stores it in REDIS. An endpoint allows to retrieve them for a specific
driver in a specific time duration
*/
package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bitly/go-nsq"
	"github.com/gorilla/mux"
	consul "github.com/hashicorp/consul/api"
	"gopkg.in/redis.v3"
)

//DriverLocationResponse contains the position of one driver
type DriverLocationResponse struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	UpdatedAt string  `json:"updated_at"`
}

//DriverLocation contains the position of one driver
type DriverLocation struct {
	DriverID  int       `json:"driverID"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	UpdatedAt time.Time `json:"updated_at"`
}

var (
	handlers     = flag.Int("handlers", 4, "Number of concurrenct handlers")
	redisClient  *redis.Client
	consulClient *consul.Client
)

func main() {
	var err error
	if consulClient, err = initConsul(); err != nil {
		log.Fatal(err)
	}
	err = register(consulClient, "location", "172.17.0.1", 8001)
	if err != nil {
		log.Fatal(err)
	}
	//NSQstream is the stream name used in NSQ by Location Service
	NSQstream := "topic_location"
	//NSQconnection is the connection string to NSQ
	NSQconnection := "172.17.0.1:4150"
	//REDISconnection is the connection string to REDIS
	REDISconnection := "172.17.0.1:6379"

	redisClient, err = RedisInit(REDISconnection)
	if err != nil {
		log.Fatal(err)
	}
	r := mux.NewRouter()
	r.HandleFunc("/drivers/{driverID:[0-9]+}/coordinates", DriverLocationHandler).Methods("GET")
	log.Printf("Server started and listening on port %d.", 8001)
	http.Handle("/", r)

	//Collects messages in NSQ to store them in REDIS
	GetDriversLocationFromGateway(redisClient, NSQconnection, NSQstream)

	log.Println(http.ListenAndServe(":8001", r))
	unregister(consulClient, "location")
}

//RedisInit connects to Redis
func RedisInit(connection string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     connection,
		Password: "",
		DB:       0,
	})
	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}

// GetDriversLocationFromGateway fetchs driver's location and stores them
func GetDriversLocationFromGateway(redisClient *redis.Client,
	NSQConnection, NSQStream string) {

	config := nsq.NewConfig()
	consumer, errConsumer := nsq.NewConsumer(NSQStream, "worker_location_service", config)
	if errConsumer != nil {
		log.Fatal("Could not create a consumer for nsq. Quit.")
	}

	w := Worker{}

	consumer.ChangeMaxInFlight(2)
	consumer.AddConcurrentHandlers(w, *handlers)

	log.Println("Unqueuing NSQ")

	if err := consumer.ConnectToNSQD(NSQConnection); err != nil {
		panic(err)
	}
}

//RedisRPush inserts the speicified val in the speicified key
func RedisRPush(redisClient *redis.Client, driverID int, m DriverLocationResponse) error {
	key := strconv.Itoa(driverID)
	v, err := json.Marshal(m)
	if err != nil {
		return err
	}
	val := string(v)

	err = redisClient.RPush(key, val).Err()
	if err != nil {
		return err
	}
	return nil
}

// DriverLocationHandler returns a json response with all driver's coordinates
// during the last N minutes
func DriverLocationHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("\t%s", r.RequestURI)

	// Read route parameter
	vars := mux.Vars(r)
	driverID := vars["driverID"]
	queryMinutes := r.URL.Query().Get("minutes")

	//ensures that the queryMinutes is well fulfilled
	if queryMinutes == "" {
		log.Printf("queryMinutes is missing")
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

//GetDriverLocation returns driver's location stored in a json format.
func GetDriverLocation(key string, minutes int64) string {
	if minutes < 1 {
		return "[]"
	}

	// According to specification Drivers are pushing their coordinates every 5
	// seconds so in order to retrieves them according minutes,
	// we do minutes * 12.
	take := minutes * 12

	r, err := redisClient.LRange(key, 0, take).Result()
	if err != nil {
		log.Printf("Could not get from Redis.")
		panic(err)
	}
	result := "[" + strings.Join(r, ",") + "]"
	return result
}

//Mapping of DriverLocationRequest and DriverLocationResult
func Mapping(d DriverLocation) DriverLocationResponse {
	m := DriverLocationResponse{}
	m.Latitude = d.Latitude
	m.Longitude = d.Longitude
	m.UpdatedAt = d.UpdatedAt.Format(time.RFC3339)
	return m
}
