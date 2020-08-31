package buses

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type allBusesDetails struct {
	TotalBuses   int                   `json:"totalBuses"`
	BusesDetails []individualBusDetail `json:"value"`
	APIKeyValid  bool                  `json:"isAPIKeyValid"`
}

type individualBusDetail struct {
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

type allBusesRoutes struct {
	TotalStops  int                  `json:"totalStops"`
	BusesRoutes []individualBusRoute `json:"value"`
	APIKeyValid bool                 `json:"isAPIKeyValid"`
}

type individualBusRoute struct {
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

func retrieveBusDetailsFromLTA(apiKey string) (data allBusesDetails) {
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

		var dataset allBusesDetails
		error = json.Unmarshal([]byte(body), &dataset)

		if len(dataset.BusesDetails) == 0 {
			data.TotalBuses = len(data.BusesDetails)
			data.APIKeyValid = true
			return
		}

		data.BusesDetails = append(data.BusesDetails, dataset.BusesDetails...)

		n += 500
	}
}

// GetAllBusDetails ref by main
func GetAllBusDetails(w http.ResponseWriter, r *http.Request) {
	apiKey := os.Getenv("API_KEY")

	if apiKey == "" {
		w.Write([]byte("API Key not found! Please attach your unique API Key to the Header of GET request"))
		return
	}

	data := retrieveBusDetailsFromLTA(apiKey)

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

func GetSingleBusDetail(w http.ResponseWriter, r *http.Request) {
	apiKey := os.Getenv("API_KEY")

	if apiKey == "" {
		w.Write([]byte("API Key not found! Please attach your unique API Key to the Header of GET request"))
		return
	}

	var vars = mux.Vars(r)
	var requestedBusNumber = vars["busNumber"]

	data := retrieveBusDetailsFromLTA(apiKey)

	if !data.APIKeyValid {
		w.Write([]byte("Invalid API Key!"))
	} else {
		var allBuses = data.BusesDetails

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
}

func retrieveBusRoutesFromLTA(requestedBusNumber string, apiKey string) (data allBusesRoutes) {
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

		var dataset allBusesRoutes
		error = json.Unmarshal([]byte(body), &dataset)

		if len(dataset.BusesRoutes) == 0 {
			data.TotalStops = len(data.BusesRoutes)
			data.APIKeyValid = true
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

func GetSingleBusRoute(w http.ResponseWriter, r *http.Request) {
	apiKey := os.Getenv("API_KEY")

	if apiKey == "" {
		w.Write([]byte("API Key not found! Please attach your unique API Key to the Header of GET request"))
		return
	}

	var vars = mux.Vars(r)
	var requestedBusNumber = vars["busNumber"]

	var allBusesRoutesDetails = retrieveBusRoutesFromLTA(requestedBusNumber, apiKey)

	if !allBusesRoutesDetails.APIKeyValid {
		w.Write([]byte("Invalid API Key!"))
	} else {
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
}
