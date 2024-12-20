package functions

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	firebase "firebase.google.com/go"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	geojson "github.com/paulmach/go.geojson"
	"google.golang.org/api/option"
)

func init() {
	functions.HTTP("updateTrafikverket", UpdateTrafikverket)
}

func parseCoordinate(coordinate string) []float64 {
	var lat, lon float64
	fmt.Sscanf(coordinate, "POINT (%f %f)", &lon, &lat)
	return []float64{lon, lat}
}

// Firebase Function to fetch from Trafikverket API and store in Firestore
func UpdateTrafikverket(w http.ResponseWriter, r *http.Request) {
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

	// Initialize Firestore client
	firestoreClient, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("error initializing Firestore client: %v", err)
	}
	defer firestoreClient.Close()

	authKey := os.Getenv("TRAFIKVERKET_AUTH_KEY")
	if authKey == "" {
		http.Error(w, "TRAFIKVERKET_AUTH_KEY not set in environment", http.StatusInternalServerError)
		return
	}

	// Fetch data from Trafikverket API
	apiURL := "https://api.trafikinfo.trafikverket.se/v2/data.json"

	// Define the XML payload
	xmlData := fmt.Sprintf(`
	<REQUEST>
		<LOGIN authenticationkey="%s" />
		<QUERY objecttype="WeatherMeasurepoint" schemaversion="2.1">
			<FILTER>
				<WITHIN name="Geometry.SWEREF99TM" shape="box" value="311863 6858375, 552124 7169867"/>
			</FILTER>
		</QUERY>
	</REQUEST>`, authKey)

	// Make the POST request
	resp, err := http.Post(apiURL, "application/xml", bytes.NewBuffer([]byte(xmlData)))
	if err != nil {
		log.Fatalf("Failed to make the request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read API response", http.StatusInternalServerError)
		return
	}

	var trafikData TrafikverketAPIResponse
	err = json.Unmarshal(body, &trafikData)
	if err != nil {
		log.Printf("Failed to parse API response: %v", err)
		http.Error(w, "Failed to parse API response", http.StatusInternalServerError)
		return
	}

	// Map trafikData into geojson features
	var features []geojson.Feature
	for _, result := range trafikData.RESPONSE.RESULT {
		for _, measurepoint := range result.WeatherMeasurepoint {
			coordinate := parseCoordinate(measurepoint.Geometry.WGS84)
			feature := geojson.NewPointFeature(coordinate)
			feature.ID = "trafikverket-" + measurepoint.ID
			feature.Properties = map[string]interface{}{
				"name":              measurepoint.Name,
				"elevation":         nil,
				"temperature_c":     measurepoint.Observation.Air.Temperature.Value,
				"windSpeed_ms":      measurepoint.Observation.Wind[0].Speed.Value,
				"windDirection_deg": measurepoint.Observation.Wind[0].Direction.Value,
				"windGustSpeed_ms":  nil,
				"humidity_percent":  measurepoint.Observation.Air.RelativeHumidity.Value,
				"newSnow24h_cm":     nil,
				"snowDepth_cm":      nil,
				"visibility_m":      measurepoint.Observation.Air.VisibleDistance.Value,
			}
			features = append(features, *feature)
		}
	}

	// Store each WeatherMeasurepoint in Firestore
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

	// Send a success response
	fmt.Fprintln(w, "Data successfully fetched and stored in Firestore.")

}

type TrafikverketAPIResponse struct {
	RESPONSE struct {
		RESULT []struct {
			WeatherMeasurepoint []struct {
				ID       string `json:"Id"`
				Name     string `json:"Name"`
				Geometry struct {
					SWEREF99TM string `json:"SWEREF99TM"`
					WGS84      string `json:"WGS84"`
				} `json:"Geometry"`
				Observation struct {
					Sample  time.Time `json:"Sample"`
					Weather struct {
						Precipitation string `json:"Precipitation"`
					} `json:"Weather"`
					Air struct {
						Temperature struct {
							Origin      string  `json:"Origin"`
							SensorNames string  `json:"SensorNames"`
							Value       float64 `json:"Value"`
						} `json:"Temperature"`
						Dewpoint struct {
							SensorNames string  `json:"SensorNames"`
							Value       float64 `json:"Value"`
						} `json:"Dewpoint"`
						RelativeHumidity struct {
							Origin      string  `json:"Origin"`
							SensorNames string  `json:"SensorNames"`
							Value       float64 `json:"Value"`
						} `json:"RelativeHumidity"`
						VisibleDistance struct {
							Origin      string  `json:"Origin"`
							SensorNames string  `json:"SensorNames"`
							Value       float64 `json:"Value"`
						} `json:"VisibleDistance"`
					} `json:"Air"`
					Wind []struct {
						Height float64 `json:"Height"`
						Speed  struct {
							Origin      string  `json:"Origin"`
							SensorNames string  `json:"SensorNames"`
							Value       float64 `json:"Value"`
						} `json:"Speed"`
						Direction struct {
							Origin      string  `json:"Origin"`
							SensorNames string  `json:"SensorNames"`
							Value       float64 `json:"Value"`
						} `json:"Direction"`
					} `json:"Wind"`
					Aggregated5Minutes struct {
						Precipitation struct {
							Rain    bool `json:"Rain"`
							Snow    bool `json:"Snow"`
							RainSum struct {
								Origin      string  `json:"Origin"`
								SensorNames string  `json:"SensorNames"`
								Value       float64 `json:"Value"`
							} `json:"RainSum"`
							SnowSum struct {
								Solid struct {
									Value float64 `json:"Value"`
								} `json:"Solid"`
								WaterEquivalent struct {
									Origin      string  `json:"Origin"`
									SensorNames string  `json:"SensorNames"`
									Value       float64 `json:"Value"`
								} `json:"WaterEquivalent"`
							} `json:"SnowSum"`
							TotalWaterEquivalent struct {
								Value float64 `json:"Value"`
							} `json:"TotalWaterEquivalent"`
						} `json:"Precipitation"`
					} `json:"Aggregated5minutes"`
					Aggregated10Minutes struct {
						Wind struct {
							SpeedMax struct {
								Origin      string  `json:"Origin"`
								SensorNames string  `json:"SensorNames"`
								Value       float64 `json:"Value"`
							} `json:"SpeedMax"`
							SpeedAverage struct {
								Origin      string  `json:"Origin"`
								SensorNames string  `json:"SensorNames"`
								Value       float64 `json:"Value"`
							} `json:"SpeedAverage"`
						} `json:"Wind"`
						Precipitation struct {
							Rain    bool `json:"Rain"`
							Snow    bool `json:"Snow"`
							RainSum struct {
								Origin      string  `json:"Origin"`
								SensorNames string  `json:"SensorNames"`
								Value       float64 `json:"Value"`
							} `json:"RainSum"`
							SnowSum struct {
								Solid struct {
									Value float64 `json:"Value"`
								} `json:"Solid"`
								WaterEquivalent struct {
									Origin      string  `json:"Origin"`
									SensorNames string  `json:"SensorNames"`
									Value       float64 `json:"Value"`
								} `json:"WaterEquivalent"`
							} `json:"SnowSum"`
							TotalWaterEquivalent struct {
								Value float64 `json:"Value"`
							} `json:"TotalWaterEquivalent"`
						} `json:"Precipitation"`
					} `json:"Aggregated10minutes"`
					Aggregated30Minutes struct {
						Wind struct {
							SpeedMax struct {
								Origin      string  `json:"Origin"`
								SensorNames string  `json:"SensorNames"`
								Value       float64 `json:"Value"`
							} `json:"SpeedMax"`
							SpeedAverage struct {
								Origin      string  `json:"Origin"`
								SensorNames string  `json:"SensorNames"`
								Value       float64 `json:"Value"`
							} `json:"SpeedAverage"`
						} `json:"Wind"`
						Precipitation struct {
							Rain    bool `json:"Rain"`
							Snow    bool `json:"Snow"`
							RainSum struct {
								Origin      string  `json:"Origin"`
								SensorNames string  `json:"SensorNames"`
								Value       float64 `json:"Value"`
							} `json:"RainSum"`
							SnowSum struct {
								Solid struct {
									Value float64 `json:"Value"`
								} `json:"Solid"`
								WaterEquivalent struct {
									Origin      string  `json:"Origin"`
									SensorNames string  `json:"SensorNames"`
									Value       float64 `json:"Value"`
								} `json:"WaterEquivalent"`
							} `json:"SnowSum"`
							TotalWaterEquivalent struct {
								Value float64 `json:"Value"`
							} `json:"TotalWaterEquivalent"`
						} `json:"Precipitation"`
					} `json:"Aggregated30minutes"`
					ID string `json:"Id"`
				} `json:"Observation"`
				Deleted      bool      `json:"Deleted"`
				ModifiedTime time.Time `json:"ModifiedTime"`
			} `json:"WeatherMeasurepoint"`
		} `json:"RESULT"`
	} `json:"RESPONSE"`
}
