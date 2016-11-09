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

func rotate(lat1, long1, lat2, long2, radius float64) (float64, float64, float64, float64) {
	half := radius / 2
	long1 += half / (math.Cos(lat1) * 69.172)
	lat1 += half / 69
	long2 -= half / (math.Cos(lat2) * 69.172)
	lat2 -= half / 69
	return lat1, long1, lat2, long2
}

func getStationsBetween(lat1 float64, long1 float64, lat2 float64, long2 float64, stations []Station) []Station {
	var output []Station

	endDistance := getDisanceBetween(lat1, long1, lat2, long2)
	num := len(stations)
	for i := 0; i < num; i++ {
		stationDist := getDisanceBetween(lat2, long2, stations[i].AddressInfo.Latitude, stations[i].AddressInfo.Longitude)
		if endDistance > stationDist {
			output = append(output, stations[i])
		}
	}
	return output
}

func getMaxStations(num float64) string {
	if num > 1000 {
		return "1000"
	} else if num < 1000 && num >= 900 {
		return "900"
	} else if num < 900 && num >= 800 {
		return "800"
	} else if num < 800 && num >= 700 {
		return "700"
	} else if num < 700 && num >= 600 {
		return "600"
	} else if num < 600 && num >= 500 {
		return "500"
	} else if num < 500 && num >= 400 {
		return "400"
	} else if num < 400 && num >= 300 {
		return "300"
	} else if num < 300 && num >= 200 {
		return "250"
	} else {
		return "100"
	}
}

func toString(num float64) string {
	return strconv.FormatFloat(num, 'f', -1, 64)
}
