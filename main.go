package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

// 39.733501, -104.992597 school
// 39.916591, -104.930168 home

//39.503001, -104.755049 random

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/hello/{name}", helloName).Methods("GET")
	router.HandleFunc("/lat/{startLat}/long/{startLong}", getLatLong).Methods("GET")
	router.HandleFunc("/start/lat/{startLat}/long/{startLong}/end/lat/{endLat}/long/{endLong}", getBetween).Methods("GET")

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))
}

func helloName(res http.ResponseWriter, req *http.Request) {
	name := mux.Vars(req)["name"]
	fmt.Fprintln(res, "Hello "+name)
}

func getLatLong(res http.ResponseWriter, req *http.Request) {
	startLat := mux.Vars(req)["startLat"]
	startLong := mux.Vars(req)["startLong"]

	url := fmt.Sprintf("http://api.openchargemap.io/v2/poi/?output=json&countrycode=US&latitude=%d&longitude=%d", startLat, startLong)
	datas := urlGetter(url)
	output := StationToJson(datas)

	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprintln(res, output)
}

func getBetween(res http.ResponseWriter, req *http.Request) {
	startLat, _ := strconv.ParseFloat(mux.Vars(req)["startLat"], 64)
	startLong, _ := strconv.ParseFloat(mux.Vars(req)["startLong"], 64)
	endLat, _ := strconv.ParseFloat(mux.Vars(req)["endLat"], 64)
	endLong, _ := strconv.ParseFloat(mux.Vars(req)["endLong"], 64)

	num := getDisanceBetween(startLat, startLong, endLat, endLong)
	maxStations := getMaxStations(num)
	fmt.Println(num)
	url := fmt.Sprintf("http://api.openchargemap.io/v2/poi/?output=json&countrycode=US&latitude=%s&longitude=%s&distance=%s&&maxresults=%s", toString(startLat), toString(startLong), toString(num), maxStations)
	datas := urlGetter(url)
	startLat, startLong, endLat, endLong = rotate(startLat, startLong, endLat, endLong, num)
	allBetween := getStationsBetween(startLat, startLong, endLat, endLong, datas)
	output := StationToJson(allBetween)

	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprint(res, output)
}
