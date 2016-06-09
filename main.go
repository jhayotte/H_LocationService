// main.go
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/bitly/go-nsq"
)

//DriverLocationHandler retrieves the last location of a customer according the time frame given in parameter
func DriverLocationHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Here are your location! :)!\n"))
}

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
