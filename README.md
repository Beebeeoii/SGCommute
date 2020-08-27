# SGCommute 🚕🚌🚇

## Introduction

### What is SGCommute

SGCommute is a **REST API wrapper** for **Singapore Land Transport Authority (Singapore)'s DataMall APIs**. Compatible with all programming languages, it allows for simple and straightforward access to LTA DataMall's data via HTTP requests.

On top of wrapping pure LTA DataMall APIs, SGCommute has additional high-level functions that will aid you in retrieving data.

📣 SGCommute is still in development stage.

### What data can SGCommute give me

Currently, SGCommute is able to retrieve the following data from LTA DataMall's API

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

### Getting an API Key

LTA DataMall API requires an API Key for data to be retrieved. As such, SGCommute requires you to provide an API key for it to work.

Registering for an API Key is free on the [LTA DataMall website](https://www.mytransport.sg/content/mytransport/home/dataMall/request-for-api.html).

### Usage

Attach your API Key to the **header** as a *key-value pair* of **HTTP GET requests** as such:

`API_KEY:{YOUR-API-KEY}`

where `API_KEY` is the key and `{YOUR-API-KEY}` is the value.

#### Getting all routes

To view all available routes of SGCommute, please read the documentation or access the [home page of SGCommute](https://sgcommute-287703.appspot.com/buses)
