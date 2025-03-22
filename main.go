package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
)

const (
	defaultCity = "Dhaka"
	apiKey = "5638256207439e5235995ff5438158dc"
	apiUrl = "https://api.openweathermap.org/data/2.5/forecast"
	geoCoderUrl = "http://api.openweathermap.org/geo/1.0/direct?q={city name},{state code},{country code}&appid={API key}"
)

type Main struct {
    Temp      float64 `json:"temp"`
    Humidity  int     `json:"humidity"`
    // Add other fields as needed
}
type Weather struct {
    Description string `json:"description"`
}

type ForecastEntry struct {
    Dt      int64    `json:"dt"`
    Main    Main     `json:"main"`
    Weather []Weather `json:"weather"`
    DtTxt   string   `json:"dt_txt"`
}

type Coord struct {
    Lat float64 `json:"lat"`
    Lon float64 `json:"lon"`
}

type City struct {
    ID      int    `json:"id"`
    Name    string `json:"name"`
    Coord   Coord  `json:"coord"`
    Country string `json:"country"`
}

type ForecastResponse struct {
    Cod     string          `json:"cod"`
    Message int             `json:"message"`
    Cnt     int             `json:"cnt"`
    List    []ForecastEntry `json:"list"`
    City    City            `json:"city"`
}

type Location struct {
	Name string `json:"name"`
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
	Country string `json:"country"`
	State string `json:"state"`
	LocalNames map[string]string `json:"local_names,omitempty"`
}

func getGeoCodeInfo(cityName string) string  {
	url := fmt.Sprintf("http://api.openweathermap.org/geo/1.0/direct?q=%s&appid=%s", cityName, apiKey)

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching geo code information:", err)
		return ""
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		fmt.Println("Error reading response body:", err)
		return ""
	}

	var locations []Location
	if err := json.Unmarshal(body, &locations); err != nil {
		fmt.Println("Error unmarshalling response body:", err)
		return ""
	}

	if len(locations) == 0 {
		fmt.Println("No location found for the city:", cityName)
		return ""
	}

	weatherAPI := fmt.Sprintf("https://api.openweathermap.org/data/2.5/forecast?lat=%f&lon=%f&appid=%s", locations[0].Lat, locations[0].Lon, apiKey)

	return weatherAPI
	
}

func main() {
	
	cityName := flag.String("city", defaultCity, "City name for weather")
	flag.Parse()


	weatherAPI := getGeoCodeInfo(*cityName)

	response, err := http.Get(weatherAPI)
	if err != nil {
		fmt.Println("Error fetching weather forecast information:", err)
		return
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		fmt.Println("Error reading response body:", err)		
		return
	}

	var forecast ForecastResponse
	if err := json.Unmarshal(body, &forecast); err != nil {
		fmt.Println("Error unmarshalling response body:", err)
		return
	}

	for _, entry := range forecast.List {
		fmt.Printf("Time: %s, Temp: %.2fK, Condition: %s\n",
			entry.DtTxt, entry.Main.Temp, entry.Weather[0].Description)
	}

	

}
