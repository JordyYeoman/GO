package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Weather struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`
	Current struct {
		TempC     float64 `json:"temp_c"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`
	Forecast struct {
		Forecastday []struct {
			Hour []struct {
				TimeEpoch int64   `json:"time_epoch"`
				TempC     float64 `json:"temp_c"`
				Condition struct {
					Text string `json:"text"`
				} `json:"condition"`
				ChanceOfRain float64 `json:"chance_of_rain"`
			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

func getCurrentWeather() string {
	WEATHER_API_KEY := os.Getenv("WEATHER_API_KEY")
	LAT_LONG := os.Getenv("LAT_LONG")

	res, err := http.Get("http://api.weatherapi.com/v1/forecast.json?key=" + WEATHER_API_KEY + "&q=" + LAT_LONG + "&days=1&aqi=no&alerts=no")
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		panic("Weather API not available")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var weather Weather
	err = json.Unmarshal(body, &weather)
	if err != nil {
		panic(err)
	}

	fmt.Println(weather)
	tempStr := strconv.FormatFloat(weather.Current.TempC, 'f', -1, 64)

	return "The temperature is " + tempStr + " and " + weather.Current.Condition.Text
}

func getWakeUpMsg() string {
	t := time.Now()
	weatherStr := getCurrentWeather()

	return "Good morning, the time is " + t.Format(time.Kitchen) + ". " + weatherStr
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	fmt.Println("Yo dawg whats good")

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/good-morning", func(w http.ResponseWriter, r *http.Request) {
		msg := getWakeUpMsg()
		_, err := w.Write([]byte(msg))
		if err != nil {
			log.Fatal("Unable to get a formatted startup message")
		}
	})
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Server Good!"))
		if err != nil {
			log.Fatal("Unable to get a formatted startup message")
		}
	})

	err = http.ListenAndServe("localhost:8080", r)
	if err != nil {
		log.Fatal("Unable to start server")
	}

	// Format docs
}
