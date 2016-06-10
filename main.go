package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/bitly/go-nsq"
	"github.com/gorilla/mux"
	"gopkg.in/redis.v3"
)

//DriverLocationHandler retrieves the last location of a customer according the time frame given in parameter
func DriverLocationHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	driverID := vars["driverID"]
	minutes := r.URL.Query().Get("minutes")

	result := GetDriverLocation(driverID)

	w.Write([]byte("Here are your location! in the last " + minutes + " \n" + result + "\n"))
}

//AddLocationInRedis insert driver location in Redis
func AddLocationInRedis(d DriverLocation) {
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
}

//GetDriverLocation returns the location of customer
func GetDriverLocation(key string) string {
	client := redis.NewClient(&redis.Options{
		Addr:     REDISconnection,
		Password: "",
		DB:       0,
	})
	redisLLenResult, err := client.LLen(key).Result()
	if err != nil {
		panic(err)
	}
	r := client.LRange(key, 0, redisLLenResult).String()

	return r
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/drivers/{driverID:[0-9]+}/coordinates", DriverLocationHandler).Methods("GET")

	http.Handle("/", r)

	go Unqueue()

	log.Fatal(http.ListenAndServe(":8001", r))
}

//Unqueue message in NSQ
func Unqueue() {
	wg := &sync.WaitGroup{}
	wg.Add(4)

	var m DriverLocation

	config := nsq.NewConfig()
	q, _ := nsq.NewConsumer(NSQstream, "Worker_test", config)
	q.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		json.Unmarshal(message.Body, &m)
		AddLocationInRedis(m)
		log.Printf(string(message.Body))
		return nil
	}))

	err := q.ConnectToNSQD(NSQconnnection)
	if err != nil {
		log.Panic("Could not connect")
	}
	wg.Wait()
	wg.Done()
}
