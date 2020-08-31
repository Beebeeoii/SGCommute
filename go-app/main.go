package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"

	busstops "private/bus-stops"
	buses "private/buses"

	"github.com/gorilla/mux"
)

var myRouter = mux.NewRouter().StrictSlash(true)

func main() {
	http.Handle("/", myRouter)
	handleRequests()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, myRouter); err != nil {
		log.Fatal(err)
	}
}

func handleRequests() {
	myRouter.HandleFunc("/", homePage)

	//everything related to buses
	myRouter.HandleFunc("/buses", buses.GetAllBusDetails)
	myRouter.HandleFunc("/buses/{busNumber}", buses.GetSingleBusDetail)
	myRouter.HandleFunc("/buses/{busNumber}/route", buses.GetSingleBusRoute)

	//everything related to bus stops
	myRouter.HandleFunc("/busstops", busstops.GetAllBusStopsDetails)
	myRouter.HandleFunc("/busstops/{busStopNumber}", busstops.GetSingleBusStopDetail)
	myRouter.HandleFunc("/busstops/{busStopNumber}/arrivals", busstops.GetBusArrivals)
	myRouter.HandleFunc("/busstops/{busStopNumber}/{busNumber}", busstops.GetSpecificBusArrival)
	// myRouter.HandleFunc("/busstops/{busStopNumber}/crowd", getSingleBusStopCrowd)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	type guideMessage struct {
		AllBusesDetailsURL                               string `json:"all_buses_details_url"`
		SingleBusDetailsURL                              string `json:"single_bus_details_url"`
		SingleBusRouteDetailsURL                         string `json:"single_bus_route_details_url"`
		AllBusStopsDetailsURL                            string `json:"all_bus_stops_details_url"`
		SingleBusStopDetailsURL                          string `json:"single_bus_stop_details_url"`
		SingleBusStopArrivalsDetailsURL                  string `json:"single_bus_stop_arrivals_details_url"`
		SingleBusStopSpecificBusServiceArrivalDetailsURL string `json:"single_bus_stop_specific_bus_service_arrival_details_url"`
	}

	message := guideMessage{
		"https://sgcommute-287703.appspot.com/buses",
		"https://sgcommute-287703.appspot.com/buses/{busNumber}",
		"https://sgcommute-287703.appspot.com/buses/{busNumber}/route",
		"https://sgcommute-287703.appspot.com/busstops",
		"https://sgcommute-287703.appspot.com/busstops/{busStopNumber}",
		"https://sgcommute-287703.appspot.com/busstops/{busStopNumber}/arrivals",
		"https://sgcommute-287703.appspot.com/busstops/{busStopNumber}/{busNumber}",
	}

	jsonResponse, error := json.Marshal(message)

	if error != nil {
		http.Error(w, error.Error(), http.StatusInternalServerError)
		return
	}
	jsonResponse = bytes.ReplaceAll(jsonResponse, []byte(","), []byte(",\n"))

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
