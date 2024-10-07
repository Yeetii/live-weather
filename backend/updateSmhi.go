package functions

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	firebase "firebase.google.com/go"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	geojson "github.com/paulmach/go.geojson"
	"google.golang.org/api/option"
)

func init() {
	functions.HTTP("updateSmhi", UpdateSmhi)
}

const apiURL = "https://opendata-download-metobs.smhi.se/api/version/1.0/parameter/"

const (
	MinLongitude = 11.91821627146622
	MinLatitude  = 61.72869520035822
	MaxLongitude = 18.493133525180227
	MaxLatitude  = 64.42201973845242
)

func UpdateSmhi(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var opts []option.ClientOption
	if _, err := os.Stat("service-account.json"); err == nil {
		opts = append(opts, option.WithCredentialsFile("service-account.json"))
	}
	conf := &firebase.Config{
		DatabaseURL: "https://live-weather-eefc5.firebaseio.com",
	}
	app, err := firebase.NewApp(ctx, conf, opts...)
	if err != nil {
		log.Fatalf("error initializing app: %v", err)
	}

	firestoreClient, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("error initializing Firestore client: %v", err)
	}
	defer firestoreClient.Close()

	stations, err := getStations(1)
	if err != nil {
		log.Fatalf("error getting stations: %v", err)
		http.Error(w, "error getting stations", http.StatusInternalServerError)
	}
	stations = filterStations(stations)

	var observations []Observation
	for _, station := range stations {
		observation, err := getObservation(1, station.ID)
		if err != nil {
			log.Printf("error getting observation for station %s: %v", station.Name, err)
			continue
		}
		observations = append(observations, observation)
	}

	// TODO: When multiple measurements, combine them into one document before storing.

	var features []geojson.Feature
	for _, observation := range observations {
		feature := geojson.NewPointFeature([]float64{*observation.Longitude, *observation.Latitude})
		feature.ID = *observation.Id
		feature.Properties = map[string]interface{}{
			"name":              observation.Name,
			"elevation":         observation.Elevation,
			"temperature_c":     observation.TemperatureC,
			"windSpeed_ms":      observation.WindSpeedMs,
			"windDirection_deg": observation.WindDirectionDeg,
			"windGustSpeed_ms":  observation.WindGustSpeedMs,
			"humidity_percent":  observation.HumidityPercent,
			"newSnow24h_cm":     nil,
			"snowDepth_cm":      observation.SnowDepthCm,
			"visibility_m":      observation.VisibilityM,
		}
		features = append(features, *feature)

	}

	for _, feature := range features {
		id, ok := feature.ID.(string)
		if !ok {
			log.Fatal("feature.ID is not a string")
		}

		geoJsonBytes, marshalErr := feature.MarshalJSON()
		if marshalErr != nil {
			log.Printf("Failed to marshal feature: %v", marshalErr)
			http.Error(w, "Failed to marshal feature", http.StatusInternalServerError)
			return
		}

		var geoJsonMap map[string]interface{}
		unmarshalErr := json.Unmarshal(geoJsonBytes, &geoJsonMap)
		if unmarshalErr != nil {
			log.Printf("Failed to unmarshal JSON into map: %v", unmarshalErr)
			http.Error(w, "Failed to process JSON", http.StatusInternalServerError)
			return
		}

		_, err := firestoreClient.Collection("weatherObservations").Doc(id).Set(ctx, geoJsonMap)
		if err != nil {
			log.Printf("Failed to store document: %v", err)
			http.Error(w, "Failed to store data in Firestore", http.StatusInternalServerError)
			return
		}
	}

	log.Println("Data successfully fetched and stored in Firestore.")
	fmt.Fprintln(w, "Data successfully fetched and stored in Firestore.")
}

func getZero[T any]() T {
	var result T
	return result
}

func fetchFromApi[T any](url string) (T, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to make the request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read the response body: %v", err)
		return getZero[T](), err
	}

	var object T
	err = json.Unmarshal(body, &object)
	if err != nil {
		log.Printf("Failed to parse API response: %v", err)
		return getZero[T](), err
	}
	return object, nil
}

func setValue(observation *Observation, value *float64, measureMentindex int) {
	switch measureMentindex {
	case 1:
		observation.TemperatureC = value
	case 4:
		observation.WindSpeedMs = value
	case 3:
		observation.WindDirectionDeg = value
	case 21:
		observation.WindGustSpeedMs = value
	case 6:
		observation.HumidityPercent = value
	case 8:
		observation.SnowDepthCm = value
	case 12:
		observation.VisibilityM = value
	default:
		log.Printf("Unknown measurement index: %v", measureMentindex)
		return
	}
}

func getObservation(measurementIndex int, stationID int) (Observation, error) {
	fetchUrl := fmt.Sprintf("%v%v/station/%v/period/latest-hour/data.json", apiURL, measurementIndex, stationID)
	measurement, err := fetchFromApi[Measurement](fetchUrl)
	if err != nil {
		log.Printf("Failed to fetch measurement on %v: %v", fetchUrl, err)
		return Observation{}, err
	}

	var elevation float64
	var lat float64
	var lon float64

	if len(measurement.Position) > 0 {
		elevation = measurement.Position[0].Height
		lat = measurement.Position[0].Latitude
		lon = measurement.Position[0].Longitude
	}

	id := "smhi-" + measurement.Station.Key
	name := measurement.Station.Name

	observation := Observation{
		Id:        &id,
		Elevation: &elevation,
		Latitude:  &lat,
		Longitude: &lon,
		Name:      &name,
	}

	if len(measurement.Value) > 0 {
		value := measurement.Value[0].Value
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			log.Printf("Failed to parse value: %v", err)
			return Observation{}, err
		}
		setValue(&observation, &floatValue, measurementIndex)
	}

	return observation, nil
}

func getStations(measurementIndex int) ([]Station, error) {
	resp, err := http.Get(fmt.Sprintf("%v%v.json", apiURL, measurementIndex))
	if err != nil {
		log.Fatalf("Failed to make the request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read the response body: %v", err)
		return []Station{}, err
	}

	var tempStations Stations
	err = json.Unmarshal(body, &tempStations)
	if err != nil {
		log.Printf("Failed to parse API response: %v", err)
		return []Station{}, err
	}
	return tempStations.Station, nil
}

func withinBoundingBox(station Station) bool {
	return station.Longitude >= MinLongitude &&
		station.Longitude <= MaxLongitude &&
		station.Latitude >= MinLatitude &&
		station.Latitude <= MaxLatitude
}

func filterStations(stations []Station) []Station {
	var filteredStations []Station
	for _, station := range stations {
		if station.Active && withinBoundingBox(station) {
			filteredStations = append(filteredStations, station)
		}
	}
	return filteredStations
}

type Observation struct {
	Id               *string  `json:"id"`
	Name             *string  `json:"name"`
	Latitude         *float64 `json:"latitude"`
	Longitude        *float64 `json:"longitude"`
	Elevation        *float64 `json:"elevation"`
	TemperatureC     *float64 `json:"temperature_c"`
	WindSpeedMs      *float64 `json:"windSpeed_ms"`
	WindDirectionDeg *float64 `json:"windDirection_deg"`
	WindGustSpeedMs  *float64 `json:"windGustSpeed_ms"`
	HumidityPercent  *float64 `json:"humidity_percent"`
	NewSnow24hCm     *float64 `json:"newSnow24h_cm"`
	SnowDepthCm      *float64 `json:"snowDepth_cm"`
	VisibilityM      *float64 `json:"visibility_m"`
}

type Station struct {
	Key     string `json:"key"`
	Updated int64  `json:"updated"`
	Title   string `json:"title"`
	Summary string `json:"summary"`
	Link    []struct {
		Href string `json:"href"`
		Rel  string `json:"rel"`
		Type string `json:"type"`
	} `json:"link"`
	Name              string  `json:"name"`
	Owner             string  `json:"owner"`
	OwnerCategory     string  `json:"ownerCategory"`
	MeasuringStations string  `json:"measuringStations"`
	ID                int     `json:"id"`
	Height            float64 `json:"height"`
	Latitude          float64 `json:"latitude"`
	Longitude         float64 `json:"longitude"`
	Active            bool    `json:"active"`
	From              int64   `json:"from"`
	To                int64   `json:"to"`
}

type Stations struct {
	Key       string `json:"key"`
	Updated   int64  `json:"updated"`
	Title     string `json:"title"`
	Summary   string `json:"summary"`
	Unit      string `json:"unit"`
	ValueType string `json:"valueType"`
	Link      []struct {
		Href string `json:"href"`
		Rel  string `json:"rel"`
		Type string `json:"type"`
	} `json:"link"`
	StationSet []struct {
		Key     string `json:"key"`
		Updated int64  `json:"updated"`
		Title   string `json:"title"`
		Summary string `json:"summary"`
		Link    []struct {
			Href string `json:"href"`
			Rel  string `json:"rel"`
			Type string `json:"type"`
		} `json:"link"`
	} `json:"stationSet"`
	Station []Station `json:"station"`
}

type Measurement struct {
	Updated   int64 `json:"updated"`
	Parameter struct {
		Key     string `json:"key"`
		Name    string `json:"name"`
		Summary string `json:"summary"`
		Unit    string `json:"unit"`
	} `json:"parameter"`
	Station struct {
		Key               string  `json:"key"`
		Name              string  `json:"name"`
		Owner             string  `json:"owner"`
		OwnerCategory     string  `json:"ownerCategory"`
		MeasuringStations string  `json:"measuringStations"`
		Height            float64 `json:"height"`
	} `json:"station"`
	Period struct {
		Key      string `json:"key"`
		From     int64  `json:"from"`
		To       int64  `json:"to"`
		Summary  string `json:"summary"`
		Sampling string `json:"sampling"`
	} `json:"period"`
	Position []struct {
		From      int64   `json:"from"`
		To        int64   `json:"to"`
		Height    float64 `json:"height"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"position"`
	Link []struct {
		Href string `json:"href"`
		Rel  string `json:"rel"`
		Type string `json:"type"`
	} `json:"link"`
	Value []struct {
		Date    int64  `json:"date"`
		Value   string `json:"value"`
		Quality string `json:"quality"`
	} `json:"value"`
}
