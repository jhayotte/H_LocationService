package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

//DriverLocationHandler retrieves the last location of a customer according the time frame given in parameter
func DriverLocationHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Here are your location! :)!\n"))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", DriverLocationHandler)

	//r.HandleFunc("/drivers/{driverID}/{id:[0-9]+}", DriverLocationHandler)

	///drivers/:id/coordinates?minutes=5
	http.Handle("/", r)

	// vars := mux.Vars(request)
	// driverID := vars["driverID"]

	log.Fatal(http.ListenAndServe(":8000", r))
}
