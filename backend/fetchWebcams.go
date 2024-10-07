package functions

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	firebase "firebase.google.com/go"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	geojson "github.com/paulmach/go.geojson"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

var firebaseApp *firebase.App

func init() {
	var err error
	ctx := context.Background()
	conf := &firebase.Config{StorageBucket: "live-weather-eefc5.appspot.com"}

	var opts []option.ClientOption
	if _, err := os.Stat("service-account.json"); err == nil {
		opts = append(opts, option.WithCredentialsFile("service-account.json"))
	}
	firebaseApp, err = firebase.NewApp(ctx, conf, opts...)
	if err != nil {
		log.Fatalf("error initializing Firebase app: %v", err)
	}

	functions.HTTP("fetchWebcams", FetchWebcams)
}

type FileInfo struct {
	URL      string          `json:"url"`
	Location geojson.Feature `json:"location"`
}

func FetchWebcams(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Get storage client
	client, err := firebaseApp.Storage(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("error getting Firebase storage client: %v", err), http.StatusInternalServerError)
		return
	}

	// Get the default bucket
	bucket, err := client.DefaultBucket()
	if err != nil {
		http.Error(w, fmt.Sprintf("error getting default Firebase storage bucket: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Println()

	// List files in the bucket
	it := bucket.Objects(ctx, nil)
	var files []FileInfo

	for {
		objectAttrs, err := it.Next()
		if err == iterator.Done {
			break // No more items in the bucket
		}
		if err != nil {
			http.Error(w, fmt.Sprintf("error listing files: %v", err), http.StatusInternalServerError)
			return
		}

		fileName := objectAttrs.Name
		url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucket.BucketName(), fileName)
		location := objectAttrs.Metadata["location"]

		var coords []float64

		erra := json.Unmarshal([]byte(location), &coords)
		if erra != nil {
			coords = []float64{0, 0}
		}

		geojson := geojson.NewFeature(geojson.NewPointGeometry(coords))
		geojson.SetProperty("url", url)
		geojson.ID = fileName

		// Add file information to the list
		files = append(files, FileInfo{URL: url, Location: *geojson})
	}

	// Set response header as JSON
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Convert the files list to JSON and write it to the response
	json.NewEncoder(w).Encode(files)
}
