# SGCommute 🚕🚌🚇

## Introduction

### What is SGCommute

SGCommute is a **REST API wrapper** for **Singapore Land Transport Authority (Singapore)'s DataMall APIs**. Compatible with all programming languages, it allows simpler and more straightforward access to LTA DataMall's data via HTTP requests.

### What data can SGCommute give me

SGCommute is able to retrieve the following data from LTA DataMall's API

- Regarding buses
  - Bus details
    - Operating company
    - AM/PM peak hours
    - Loop/bi-directional service
    - Location at which the bus service loops (if it is a looped service)
  - Bus routes
    - All bus stops covered by the bus service
    - Weekday/Saturday/Sunday first/last bus service timing at specific bus stop
- Regarding bus stops
  - Road name
  - Landmarks near bus stops to aid in identifying bus stop (if any)
  - Coordinates of bus stop (longitude and latitude)
  - Bus arrival timings at that bus stop

## Get started

### Usage

#### HTTP GET requests

This project is hosted on a server and you may simply send GET requests to `https://sgcommute-287703.appspot.com`.

##### __Example routes for data__

To retrieve all buses that are in operation:

`https://sgcommute-287703.appspot.com/buses`

To retrieve all bus stops in Singapore that are in operation:

`https://sgcommute-287703.appspot.com/busstops`

To retrieve bus arrival timings at a specific bus stop (Eg. 01012):

`https://sgcommute-287703.appspot.com/busstops/01012/arrivals`

For a full list of routes, please visit the wiki page of this project.

### Deploy

#### **On localhost**

##### **Step 1**

`git clone https://github.com/Beebeeoii/SGCommute.git`

##### **Step 2**

Navigate to `sgcommute-localhost` > `private` > `bus-stops` > `bus-stops.go` and `sgcommute-localhost` > `private` > `buses` > `buses.go`

##### **Step 3**

Register for an API Key on the [LTA DataMall website](https://www.mytransport.sg/content/mytransport/home/dataMall/request-for-api.html)

##### **Step 4**

With the API key, replace `yourAPIKey` in `const apiKey = "yourAPIKey"` with the API Key you received via email from Step 3

##### **Step 5**

Build and run `main.go` and the server is accessible from `localhost:8080`

Voila!

#### **On Google App Engine**

##### Step 1

`git clone https://github.com/Beebeeoii/SGCommute.git`

##### Step 2

Navigate to `go-app` > `app.yaml`

##### Step 3

Register for an API Key on the [LTA DataMall website](https://www.mytransport.sg/content/mytransport/home/dataMall/request-for-api.html)

##### Step 4

With the API key, replace `yourAPIKey` in `API_KEY: "yourAPIKey"` with the API Key you received via email from Step 3

##### Step 5

Upload both `go-app` and `gopath` onto Google App Engine via Cloudshell Editor

##### Step 6

Enter Cloudshell and perform `cd go-app/`, `gcloud app deploy`

Voila!

## License

    Copyright 2020 Beebeeoii

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

        http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
