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
