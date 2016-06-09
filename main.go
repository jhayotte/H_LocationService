// main.go
package main

import (
	"encoding/json"
	"log"
	"strconv"
	//"log"
	//"net/http"
	"sync"

	"github.com/bitly/go-nsq"
)

//DriverLocationHandler retrieves the last location of a customer according the time frame given in parameter
// func DriverLocationHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("Here are your location! :)!\n"))
// }

func main() {
	// r := mux.NewRouter()
	// r.HandleFunc("/", DriverLocationHandler)
	//
	// //r.HandleFunc("/drivers/{driverID}/{id:[0-9]+}", DriverLocationHandler)
	//
	// ///drivers/:id/coordinates?minutes=5
	// http.Handle("/", r)
	//
	// // vars := mux.Vars(request)
	// // driverID := vars["driverID"]
	//
	// log.Fatal(http.ListenAndServe(":8000", r))

	Unqueue()

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
		log.Printf("Driver %v was in latitude: %v and longitude: %v at %v", strconv.Itoa(m.DriverID), strconv.FormatFloat(m.Latitude, 'f', 3, 32), strconv.FormatFloat(m.Longitude, 'f', 3, 32), m.CreatedDate)
		return nil
	}))

	err := q.ConnectToNSQD(NSQconnnection)
	if err != nil {
		log.Panic("Could not connect")
	}
	wg.Wait()
	wg.Done()

}
