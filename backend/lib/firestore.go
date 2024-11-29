package lib

import (
	"context"
	"encoding/json"
	"log"
	"os"

	firebase "firebase.google.com/go"
	geojson "github.com/paulmach/go.geojson"
	"google.golang.org/api/option"
)

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
	NewSnow72hCm     *float64 `json:"newSnow72h_cm"`
	SnowDepthCm      *float64 `json:"snowDepth_cm"`
	VisibilityM      *float64 `json:"visibility_m"`
}

func UploadObservationsToFirestore(observations []Observation) error {
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
			"newSnow24h_cm":     observation.NewSnow24hCm,
			"newSnow72h_cm":     observation.NewSnow72hCm,
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
			return marshalErr
		}

		var geoJsonMap map[string]interface{}
		unmarshalErr := json.Unmarshal(geoJsonBytes, &geoJsonMap)
		if unmarshalErr != nil {
			log.Printf("Failed to unmarshal JSON into map: %v", unmarshalErr)
			return unmarshalErr
		}

		_, err := firestoreClient.Collection("weatherObservations").Doc(id).Set(ctx, geoJsonMap)
		if err != nil {
			log.Printf("Failed to store document: %v", err)
			return err
		}
	}
	return nil
}
