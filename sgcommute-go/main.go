package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type AllBusesDetails struct {
	TotalBuses   int                   `json:"totalBuses"`
	BusesDetails []IndividualBusDetail `json:"value"`
}

type IndividualBusDetail struct {
	ServiceNo       string `json:"ServiceNo"`
	Operator        string `json:"Operator"`
	Direction       int    `json:"Direction"`
	Category        string `json:"Category"`
	OriginCode      string `json:"OriginCode"`
	DestinationCode string `json:"DestinationCode"`
	AM_Peak_Freq    string `json:"AM_Peak_Freq"`
	AM_Offpeak_freq string `json:"AM_Offpeak_freq"`
	PM_Peak_Freq    string `json:"PM_Peak_Freq"`
	PM_Offpeak_Freq string `json:"PM_Offpeak_Freq"`
	LoopDesc        string `json:"LoopDesc"`
}

type AllBusesRoutes struct {
	TotalStops  int                  `json:"totalStops"`
	BusesRoutes []IndividualBusRoute `json:"value"`
}

type IndividualBusRoute struct {
	ServiceNo    string  `json:"ServiceNo"`
	Operator     string  `json:"Operator"`
	Direction    int     `json:"Direction"`
	StopSequence int     `json:"StopSequence"`
	BusStopCode  string  `json:"BusStopCode"`
	Distance     float32 `json:"Distance"`
	WD_FirstBus  string  `json:"WD_FirstBus"`
	WD_LastBus   string  `json:"WD_LastBus"`
	SAT_FirstBus string  `json:"SAT_FirstBus"`
	SAT_LastBus  string  `json:"SAT_LastBus"`
	SUN_FirstBus string  `json:"SUN_FirstBus"`
	SUN_LastBus  string  `json:"SUN_LastBus"`
}

type AllBusesArrivalDetails struct {
	TotalBuses        int                  `json:"totalBuses"`
	BusArrivalDetails []BusServiceArrivals `json:"Services"`
}

type BusServiceArrivals struct {
	ServiceNo string                      `json:"ServiceNo"`
	Operator  string                      `json:"Operator"`
	NextBus   IndividualBusArrivalDetails `json:"NextBus"`
	NextBus2  IndividualBusArrivalDetails `json:"NextBus2"`
	NextBus3  IndividualBusArrivalDetails `json:"NextBus3"`
}

type IndividualBusArrivalDetails struct {
	OriginCode       string `json:"OriginCode"`
	DestinationCode  string `json:"DestinationCode"`
	EstimatedArrival string `json:"EstimatedArrival"`
	Latitude         string `json:"Latitude"`
	Longitude        string `json:"Longitude"`
	VisitNumber      string `json:"VisitNumber"`
	Load             string `json:"Load"`
	Feature          string `json:"Feature"`
	Type             string `json:"Type"`
}

type AllBusStops struct {
	TotalBusStops  int                  `json:"totalBusStops"`
	BusStopDetails []IndividualBusStops `json:"value"`
}

type IndividualBusStops struct {
	BusStopCode string  `json:"BusStopCode"`
	RoadName    string  `json:"RoadName"`
	Description string  `json:"Description"`
	Latitude    float64 `json:"Latitude"`
	Longitude   float64 `json:"Longitude"`
}

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
	myRouter.HandleFunc("/buses", getAllBusDetails)
	myRouter.HandleFunc("/buses/{busNumber}", getSingleBusDetail)
	myRouter.HandleFunc("/buses/{busNumber}/route", getSingleBusRoute)

	//everything related to bus stops
	myRouter.HandleFunc("/busstops", getAllBusStopsDetails)
	myRouter.HandleFunc("/busstops/{busStopNumber}", getSingleBusStopDetail)
	myRouter.HandleFunc("/busstops/{busStopNumber}/arrivals", getBusArrivals)
	myRouter.HandleFunc("/busstops/{busStopNumber}/{busNumber}", getSpecificBusArrival)
	// myRouter.HandleFunc("/busstops/{busStopNumber}/crowd", getSingleBusStopCrowd)
}

func retrieveBusDetailsFromLTA() (data AllBusesDetails) {
	var n = 0
	const BusDetailsAddress = "http://datamall2.mytransport.sg/ltaodataservice/BusServices"

	for {
		var httpClient = http.Client{}

		var request, error = http.NewRequest("GET", BusDetailsAddress, nil)
		if error != nil {
			log.Fatalln(error)
		}
		var query = request.URL.Query()
		query.Add("$skip", fmt.Sprint(n))
		request.URL.RawQuery = query.Encode()

		request.Header.Set("AccountKey", "7WS52+iASeigG9HsbWDM6Q==")
		request.Header.Set("accept", "application/json")

		response, error := httpClient.Do(request)
		if error != nil {
			log.Fatalln(error)
		}
		defer response.Body.Close()

		body, error := ioutil.ReadAll(response.Body)

		var dataset AllBusesDetails
		error = json.Unmarshal([]byte(body), &dataset)

		if len(dataset.BusesDetails) == 0 {
			data.TotalBuses = len(data.BusesDetails)
			return
		}

		data.BusesDetails = append(data.BusesDetails, dataset.BusesDetails...)

		n += 500
	}
}

func retrieveBusRoutesFromLTA(requestedBusNumber string) (data AllBusesRoutes) {
	var n = 0
	const BusDetailsAddress = "http://datamall2.mytransport.sg/ltaodataservice/BusRoutes"

	for {
		var httpClient = http.Client{}

		var request, error = http.NewRequest("GET", BusDetailsAddress, nil)
		if error != nil {
			log.Fatalln(error)
		}
		var query = request.URL.Query()
		query.Add("$skip", fmt.Sprint(n))
		request.URL.RawQuery = query.Encode()

		request.Header.Set("AccountKey", "7WS52+iASeigG9HsbWDM6Q==")
		request.Header.Set("accept", "application/json")

		response, error := httpClient.Do(request)
		if error != nil {
			log.Fatalln(error)
		}
		defer response.Body.Close()

		body, error := ioutil.ReadAll(response.Body)

		var dataset AllBusesRoutes
		error = json.Unmarshal([]byte(body), &dataset)

		if len(dataset.BusesRoutes) == 0 {
			data.TotalStops = len(data.BusesRoutes)
			return
		}

		for _, n := range dataset.BusesRoutes {
			if n.ServiceNo == requestedBusNumber {
				data.BusesRoutes = append(data.BusesRoutes, n)
			}
		}
		n += 500
	}
}

func retrieveBusArrivalsFromLTA(requestedBusStopCode string) (data AllBusesArrivalDetails) {
	const BusDetailsAddress = "http://datamall2.mytransport.sg/ltaodataservice/BusArrivalv2"

	var httpClient = http.Client{}

	var request, error = http.NewRequest("GET", BusDetailsAddress, nil)
	if error != nil {
		log.Fatalln(error)
	}
	var query = request.URL.Query()
	query.Add("BusStopCode", requestedBusStopCode)
	request.URL.RawQuery = query.Encode()

	request.Header.Set("AccountKey", "7WS52+iASeigG9HsbWDM6Q==")
	request.Header.Set("accept", "application/json")

	response, error := httpClient.Do(request)
	if error != nil {
		log.Fatalln(error)
	}
	defer response.Body.Close()

	body, error := ioutil.ReadAll(response.Body)

	error = json.Unmarshal([]byte(body), &data)

	data.TotalBuses = len(data.BusArrivalDetails)
	return
}

func retrieveSpecificBusArrivalFromLTA(requestedBusStopCode string, requestedBusNumber string) (data AllBusesArrivalDetails) {
	const BusDetailsAddress = "http://datamall2.mytransport.sg/ltaodataservice/BusArrivalv2"

	var httpClient = http.Client{}

	var request, error = http.NewRequest("GET", BusDetailsAddress, nil)
	if error != nil {
		log.Fatalln(error)
	}
	var query = request.URL.Query()
	query.Add("BusStopCode", requestedBusStopCode)
	query.Add("ServiceNo", requestedBusNumber)
	request.URL.RawQuery = query.Encode()

	request.Header.Set("AccountKey", "7WS52+iASeigG9HsbWDM6Q==")
	request.Header.Set("accept", "application/json")

	response, error := httpClient.Do(request)
	if error != nil {
		log.Fatalln(error)
	}
	defer response.Body.Close()

	body, error := ioutil.ReadAll(response.Body)

	error = json.Unmarshal([]byte(body), &data)

	data.TotalBuses = len(data.BusArrivalDetails)
	return
}

func retrieveBusStopsFromLTA() (data AllBusStops) {
	var n = 0
	const BusDetailsAddress = "http://datamall2.mytransport.sg/ltaodataservice/BusStops"

	for {
		var httpClient = http.Client{}

		var request, error = http.NewRequest("GET", BusDetailsAddress, nil)
		if error != nil {
			log.Fatalln(error)
		}
		var query = request.URL.Query()
		query.Add("$skip", fmt.Sprint(n))
		request.URL.RawQuery = query.Encode()

		request.Header.Set("AccountKey", "7WS52+iASeigG9HsbWDM6Q==")
		request.Header.Set("accept", "application/json")

		response, error := httpClient.Do(request)
		if error != nil {
			log.Fatalln(error)
		}
		defer response.Body.Close()

		body, error := ioutil.ReadAll(response.Body)

		var dataset AllBusStops
		error = json.Unmarshal([]byte(body), &dataset)

		if len(dataset.BusStopDetails) == 0 {
			data.TotalBusStops = len(data.BusStopDetails)
			return
		}

		data.BusStopDetails = append(data.BusStopDetails, dataset.BusStopDetails...)

		n += 500
	}
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
		" https://sgcommute-287703.appspot.com/buses",
		" https://sgcommute-287703.appspot.com/buses/{busNumber}",
		" https://sgcommute-287703.appspot.com/buses/{busNumber}/route",
		" https://sgcommute-287703.appspot.com/busstops",
		" https://sgcommute-287703.appspot.com/busstops/{busStopNumber}",
		" https://sgcommute-287703.appspot.com/busstops/{busStopNumber}/arrivals",
		" https://sgcommute-287703.appspot.com/busstops/{busStopNumber}/{busNumber}",
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

func getAllBusDetails(w http.ResponseWriter, r *http.Request) {
	jsonResponse, error := json.Marshal(retrieveBusDetailsFromLTA())

	if error != nil {
		http.Error(w, error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func getSingleBusDetail(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	var requestedBusNumber = vars["busNumber"]

	var allBuses = retrieveBusDetailsFromLTA().BusesDetails

	for _, n := range allBuses {
		if requestedBusNumber == n.ServiceNo {
			jsonResponse, error := json.Marshal(n)
			if error != nil {
				http.Error(w, error.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonResponse)
			return
		}
	}

	w.Write([]byte("No such bus service"))

}

func getSingleBusRoute(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	var requestedBusNumber = vars["busNumber"]

	var allBusesRoutesDetails = retrieveBusRoutesFromLTA(requestedBusNumber)

	if allBusesRoutesDetails.TotalStops == 0 {
		w.Write([]byte("No such bus service"))
		return
	}

	jsonResponse, error := json.Marshal(allBusesRoutesDetails)
	if error != nil {
		http.Error(w, error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func getAllBusStopsDetails(w http.ResponseWriter, r *http.Request) {
	jsonResponse, error := json.Marshal(retrieveBusStopsFromLTA())

	if error != nil {
		http.Error(w, error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func getSingleBusStopDetail(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	var requestedBusStopCode = vars["busStopNumber"]

	var busStops = retrieveBusStopsFromLTA().BusStopDetails

	for _, n := range busStops {
		if requestedBusStopCode == n.BusStopCode {
			jsonResponse, error := json.Marshal(n)
			if error != nil {
				http.Error(w, error.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonResponse)
			return
		}
	}

	w.Write([]byte("No such bus service"))
}

func getBusArrivals(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	var requestedBusStopCode = vars["busStopNumber"]

	jsonResponse, error := json.Marshal(retrieveBusArrivalsFromLTA(requestedBusStopCode))

	if error != nil {
		http.Error(w, error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func getSpecificBusArrival(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	var requestedBusStopCode = vars["busStopNumber"]
	var requestedBusNumber = vars["busNumber"]

	jsonResponse, error := json.Marshal(retrieveSpecificBusArrivalFromLTA(requestedBusStopCode, requestedBusNumber))

	if error != nil {
		http.Error(w, error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
