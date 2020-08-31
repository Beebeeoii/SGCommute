package busstops

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const apiKey = "yourAPIKey"

type allBusStops struct {
	TotalBusStops  int                  `json:"totalBusStops"`
	BusStopDetails []individualBusStops `json:"value"`
	APIKeyValid    bool                 `json:"isAPIKeyValid"`
}

type individualBusStops struct {
	BusStopCode string  `json:"BusStopCode"`
	RoadName    string  `json:"RoadName"`
	Description string  `json:"Description"`
	Latitude    float64 `json:"Latitude"`
	Longitude   float64 `json:"Longitude"`
}

func retrieveBusStopsFromLTA(apiKey string) (data allBusStops) {
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

		request.Header.Set("AccountKey", apiKey)
		request.Header.Set("accept", "application/json")

		response, error := httpClient.Do(request)

		if response.Status == "401 UNAUTHORIZED" {
			data.APIKeyValid = false
			return
		}

		if error != nil {
			log.Fatalln(error)
		}
		defer response.Body.Close()

		body, error := ioutil.ReadAll(response.Body)

		var dataset allBusStops
		error = json.Unmarshal([]byte(body), &dataset)

		if len(dataset.BusStopDetails) == 0 {
			data.TotalBusStops = len(data.BusStopDetails)
			data.APIKeyValid = true
			return
		}

		data.BusStopDetails = append(data.BusStopDetails, dataset.BusStopDetails...)

		n += 500
	}
}

func GetAllBusStopsDetails(w http.ResponseWriter, r *http.Request) {
	if apiKey == "" {
		w.Write([]byte("API Key not found! Please fill in API Key"))
		return
	}

	data := retrieveBusStopsFromLTA(apiKey)

	if !data.APIKeyValid {
		w.Write([]byte("Invalid API Key!"))
	} else {
		jsonResponse, error := json.Marshal(data)

		if error != nil {
			http.Error(w, error.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}
}

func GetSingleBusStopDetail(w http.ResponseWriter, r *http.Request) {
	if apiKey == "" {
		w.Write([]byte("API Key not found! Please fill in API Key"))
		return
	}

	var vars = mux.Vars(r)
	var requestedBusStopCode = vars["busStopNumber"]

	data := retrieveBusStopsFromLTA(apiKey)

	if !data.APIKeyValid {
		w.Write([]byte("Invalid API Key!"))
	} else {
		var busStops = data.BusStopDetails

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

		w.Write([]byte("No such bus stop number"))
	}
}

type allBusesArrivalDetails struct {
	TotalBuses        int                  `json:"totalBuses"`
	BusArrivalDetails []busServiceArrivals `json:"Services"`
	APIKeyValid       bool                 `json:"isAPIKeyValid"`
}

type busServiceArrivals struct {
	ServiceNo string                      `json:"ServiceNo"`
	Operator  string                      `json:"Operator"`
	NextBus   individualBusArrivalDetails `json:"NextBus"`
	NextBus2  individualBusArrivalDetails `json:"NextBus2"`
	NextBus3  individualBusArrivalDetails `json:"NextBus3"`
}

type individualBusArrivalDetails struct {
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

func retrieveBusArrivalsFromLTA(requestedBusStopCode string, apiKey string) (data allBusesArrivalDetails) {
	const BusDetailsAddress = "http://datamall2.mytransport.sg/ltaodataservice/BusArrivalv2"

	var httpClient = http.Client{}

	var request, error = http.NewRequest("GET", BusDetailsAddress, nil)
	if error != nil {
		log.Fatalln(error)
	}
	var query = request.URL.Query()
	query.Add("BusStopCode", requestedBusStopCode)
	request.URL.RawQuery = query.Encode()

	request.Header.Set("AccountKey", apiKey)
	request.Header.Set("accept", "application/json")

	response, error := httpClient.Do(request)

	if response.Status == "401 UNAUTHORIZED" {
		data.APIKeyValid = false
		return
	}

	if error != nil {
		log.Fatalln(error)
	}
	defer response.Body.Close()

	body, error := ioutil.ReadAll(response.Body)

	error = json.Unmarshal([]byte(body), &data)

	data.TotalBuses = len(data.BusArrivalDetails)
	data.APIKeyValid = true
	return
}

func GetBusArrivals(w http.ResponseWriter, r *http.Request) {
	if apiKey == "" {
		w.Write([]byte("API Key not found! Please fill in API Key"))
		return
	}

	var vars = mux.Vars(r)
	var requestedBusStopCode = vars["busStopNumber"]

	data := retrieveBusArrivalsFromLTA(requestedBusStopCode, apiKey)

	if !data.APIKeyValid {
		w.Write([]byte("Invalid API Key!"))
	} else {
		jsonResponse, error := json.Marshal(data)

		if error != nil {
			http.Error(w, error.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}
}

func retrieveSpecificBusArrivalFromLTA(requestedBusStopCode string, requestedBusNumber string, apiKey string) (data allBusesArrivalDetails) {
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

	request.Header.Set("AccountKey", apiKey)
	request.Header.Set("accept", "application/json")

	response, error := httpClient.Do(request)

	if response.Status == "401 UNAUTHORIZED" {
		data.APIKeyValid = false
		return
	}

	if error != nil {
		log.Fatalln(error)
	}
	defer response.Body.Close()

	body, error := ioutil.ReadAll(response.Body)

	error = json.Unmarshal([]byte(body), &data)

	data.TotalBuses = len(data.BusArrivalDetails)
	data.APIKeyValid = true
	return
}

func GetSpecificBusArrival(w http.ResponseWriter, r *http.Request) {
	if apiKey == "" {
		w.Write([]byte("API Key not found! Please fill in API Key"))
		return
	}

	var vars = mux.Vars(r)
	var requestedBusStopCode = vars["busStopNumber"]
	var requestedBusNumber = vars["busNumber"]

	data := retrieveSpecificBusArrivalFromLTA(requestedBusStopCode, requestedBusNumber, apiKey)

	if !data.APIKeyValid {
		w.Write([]byte("Invalid API Key!"))
	} else {
		jsonResponse, error := json.Marshal(data)

		if error != nil {
			http.Error(w, error.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}
}
