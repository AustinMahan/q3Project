package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

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
	router.HandleFunc("/start/address/{address1}/end/address/{address2}", getAddresses).Methods("GET")
	router.HandleFunc("/start/lat/{startLat}/long/{startLong}/end/address/{address}", curLocationStart).Methods("GET")

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))
}

func helloName(res http.ResponseWriter, req *http.Request) {
	name := mux.Vars(req)["name"]
	fmt.Fprintln(res, "Hello "+name)
}

func getLatLong(res http.ResponseWriter, req *http.Request) {
	startLat := mux.Vars(req)["startLat"]
	startLong := mux.Vars(req)["startLong"]

	url := fmt.Sprintf("http://api.openchargemap.io/v2/poi/?output=json&countrycode=US&distance=100&maxresults=50&latitude=%s&longitude=%s", startLat, startLong)
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
	allBetween := getStationsBetween(startLat, startLong, endLat, endLong, datas, num)
	output := StationToJsonWStartEnd(allBetween, startLat, startLong, endLat, endLong)

	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprint(res, output)
}

func getAddresses(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	address1, _ := mux.Vars(req)["address1"]
	address1 = strings.Replace(address1, " ", "+", -1)
	address1 = fmt.Sprintf("https://maps.googleapis.com/maps/api/place/textsearch/json?query=%s&key=%s", address1, "AIzaSyDMbb_u0IslPPN7gmPfFbPCoqAz1Mmussc")

	address1Lat, address1Lng := latLngGetter(address1)

	address2, _ := mux.Vars(req)["address2"]
	address2 = strings.Replace(address2, " ", "+", -1)
	address2 = fmt.Sprintf("https://maps.googleapis.com/maps/api/place/textsearch/json?query=%s&key=%s", address2, "AIzaSyDMbb_u0IslPPN7gmPfFbPCoqAz1Mmussc")
	address2Lat, address2Lng := latLngGetter(address2)
	if address1Lat == 0 || address1Lng == 0 || address2Lat == 0 || address2Lng == 0 {
		fmt.Fprint(res, "[]")
	} else {
		num := getDisanceBetween(address1Lat, address1Lng, address2Lat, address2Lng)
		maxStations := getMaxStations(num)

		url := fmt.Sprintf("http://api.openchargemap.io/v2/poi/?output=json&countrycode=US&latitude=%s&longitude=%s&distance=%s&&maxresults=%s", toString(address1Lat), toString(address1Lng), toString(num), maxStations)
		datas := urlGetter(url)
		allBetween := getStationsBetween(address1Lat, address1Lng, address2Lat, address2Lng, datas, num)
		output := StationToJsonWStartEnd(allBetween, address1Lat, address1Lng, address2Lat, address2Lng)

		fmt.Fprint(res, output)
	}
}
func curLocationStart(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	startLat, _ := strconv.ParseFloat(mux.Vars(req)["startLat"], 64)
	startLong, _ := strconv.ParseFloat(mux.Vars(req)["startLong"], 64)

	address, _ := mux.Vars(req)["address"]
	address = strings.Replace(address, " ", "+", -1)
	address = fmt.Sprintf("https://maps.googleapis.com/maps/api/place/textsearch/json?query=%s&key=%s", address, "AIzaSyDMbb_u0IslPPN7gmPfFbPCoqAz1Mmussc")
	addressLat, addressLng := latLngGetter(address)
	if startLat == 0 || startLong == 0 || addressLat == 0 || addressLng == 0 {
		fmt.Fprint(res, "[]")
	} else {
		num := getDisanceBetween(startLat, startLong, addressLat, addressLng)
		maxStations := getMaxStations(num)

		url := fmt.Sprintf("http://api.openchargemap.io/v2/poi/?output=json&countrycode=US&latitude=%s&longitude=%s&distance=%s&&maxresults=%s", toString(startLat), toString(startLong), toString(num), maxStations)
		datas := urlGetter(url)
		allBetween := getStationsBetween(startLat, startLong, addressLat, addressLng, datas, num)
		output := StationToJsonWStartEnd(allBetween, startLat, startLong, addressLat, addressLng)

		fmt.Fprint(res, output)
	}
}
