package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

type Weather struct {
	Location struct {
		Name      string `json:"name"`
		State     string `json:"region"`
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

	// loading .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading the .env file.")
	}

	// construct the url query
	baseURL := os.Getenv("BASE_URL")
	apikey := os.Getenv("API_KEY")
	q := url.QueryEscape(os.Getenv("DEFAULT_LOCATION"))
	days := "1"
	aqi := "yes"

	// check for command-line argument
	if len(os.Args) >= 2 {

		// if there is one argument - ex. sunny.exe "atlanta"
		if len(os.Args) == 2 {
			q = url.QueryEscape(os.Args[1])

			// if there is multiple argument, construct as a single query - ex. sunny.exe "los angeles"
		} else {
			q = url.QueryEscape(strings.Join(os.Args[1:], " "))
		}
	}

	fullUrl := baseURL + apikey + "&q=" + q + "&days=" + days + "&aqi=" + aqi

	// fetch a request from weather api
	res, err := http.Get(fullUrl)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		panic("Weather API is not available.")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	// desconstruct response body onto weather struct
	var weather Weather
	err = json.Unmarshal(body, &weather)
	if err != nil {
		panic(err)
	}

	location, current, hours := weather.Location, weather.Current, weather.Forecast.Forecastday[0].Hour

	// print ex. - Norcross, Georgia: 81F, Partly cloudy - Feels like 83F
	curr := fmt.Sprintf(
		"\n%s, %s: %.0fF, %s - Feels like %.0fF\n",
		location.Name, location.State, current.Temp, current.Condition.Text, current.FeelsLike)
	color.Yellow(curr)

	// print ex. - Humidity: 72%, Cloud: 75%, UV Index: 7, AQI: 15
	curr2 := fmt.Sprintf(
		"Humidity: %d%%, Cloud: %d%%, UV Index: %.0f, AQI: %.0f\n",
		current.Humidity, current.Cloud, current.UV, current.AirQuality.PM25)
	color.Cyan(curr2)

	// print ex. dash
	dash := utf8.RuneCountInString(curr2)
	fmt.Println(strings.Repeat("-", dash))

	// print ex. - 05:00PM - 86F, 0%, Overcast
	for _, hour := range hours {
		date := time.Unix(hour.TimeEpoch, 0)

		if date.Before(time.Now()) {
			continue
		}

		message := fmt.Sprintf(
			"%s - %.0fF, %.0f%%, %s\n",
			date.Format("03:04PM"), hour.TempF, hour.ChanceOfRain, hour.Condition.Text)

		if hour.ChanceOfRain < 90 {
			fmt.Print(message)
		} else {
			color.Red(message)
		}
	}

}
