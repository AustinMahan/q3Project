package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"
)

func StationToJson(stations []Station) string {
	byte, err := json.Marshal(stations)
	if err != nil {
		fmt.Println(err)
	}

	return string(byte)
}

func urlGetter(url string) []Station {
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println("Bad")
		panic(err)
	}

	jsonDataFromHttp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Bad")
		panic(err)
	}
	var jsonData []Station
	err = json.Unmarshal([]byte(jsonDataFromHttp), &jsonData)
	if err != nil {
		fmt.Println("Bad")
		panic(err)
	}
	return jsonData
}

func latLngGetter(url string) (lat float64, lng float64) {
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println("Bad")
		panic(err)
	}

	jsonDataFromHttp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Bad")
		panic(err)
	}
	var jsonData addressToLatLng
	err = json.Unmarshal([]byte(jsonDataFromHttp), &jsonData)
	if err != nil {
		fmt.Println("Bad")
		panic(err)
	}
	return jsonData.Results[0].Geometry.Location.Lat, jsonData.Results[0].Geometry.Location.Lng
}

func keepLines(s string, n int) string {
	result := strings.Join(strings.Split(s, "\n")[:n], "\n")
	return strings.Replace(result, "\r", "", -1)
}

type Station struct {
	AddressInfo struct {
		AddressLine1    string  `json:"AddressLine1"`
		Latitude        float64 `json:"Latitude"`
		Longitude       float64 `json:"Longitude"`
		Postcode        string  `json:"Postcode"`
		StateOrProvince string  `json:"StateOrProvince"`
		Town            string  `json:"Town"`
	} `json:"AddressInfo"`
}

type addressToLatLng struct {
	Results []struct {
		Geometry struct {
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
		} `json:"geometry"`
	} `json:"results"`
}

func getDisanceBetween(lat1 float64, long1 float64, lat2 float64, long2 float64) float64 {
	var dLat float64 = deg2rad(lat2 - lat1)
	var dLon float64 = deg2rad(long2 - long1)
	var a float64 = math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(deg2rad(lat1))*math.Cos(deg2rad(lat2))*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	var c = 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	var d = 3961 * c
	return d
}

func deg2rad(deg float64) float64 {
	return deg * (math.Pi / 180)
}

func getStationsBetween(lat1 float64, long1 float64, lat2 float64, long2 float64, stations []Station, endDistance float64) []Station {
	var output []Station

	num := len(stations)
	for i := 0; i < num; i++ {
		stationDist := getDisanceBetween(lat2, long2, stations[i].AddressInfo.Latitude, stations[i].AddressInfo.Longitude)
		startToStation := getDisanceBetween(lat1, long1, stations[i].AddressInfo.Latitude, stations[i].AddressInfo.Longitude)
		if endDistance > stationDist && (stationDist+startToStation) < endDistance*1.07 {
			output = append(output, stations[i])
		}
	}
	return output
}

func getMaxStations(num float64) string {
	return "9999"
}

func toString(num float64) string {
	return strconv.FormatFloat(num, 'f', -1, 64)
}
