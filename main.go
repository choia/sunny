package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

type Weather struct {
	Location struct {
		Name      string `json:"name"`
		State     string `json:"state"`
		Country   string `json:"country"`
		LocalTime string `json:"localtime"`
	} `json:"location"`
	Current struct {
		Temp       float64 `json:"temp_f"`
		Humidity   int16   `json:"humidity"`
		Cloud      int16   `json:"cloud"`
		FeelsLike  float64 `json:"feelslike_f"`
		UV         float64 `json:"uv"`
		AirQuality struct {
			PM25 float64 `json:"pm2_5"`
		} `json:"air_quality"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`
	Forecast struct {
		Forecastday []struct {
			Hour []struct {
				TimeEpoch    int64   `json:"time_epoch"`
				TempF        float64 `json:"temp_f"`
				ChanceOfRain float64 `json:"chance_of_rain"`
				Condition    struct {
					Text string `json:"text"`
				} `json:"condition"`
			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading the .env file.")
	}

	baseURL := os.Getenv("BASE_URL")
	apikey := os.Getenv("API_KEY")
	q := url.QueryEscape(os.Getenv("DEFAULT_LOCATION"))
	days := "1"
	aqi := "yes"

	fullUrl := baseURL + apikey + "&q=" + q + "&days=" + days + "&aqi=" + aqi

	res, err := http.Get(fullUrl)
	if err != nil {
		panic("Weather API is not available.")
	}
	// fmt.Println(fullUrl)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", body)
}
