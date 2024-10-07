package lib

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func UploadToFirebaseStorage(url string, fileName string, location []float64) error {
	fmt.Println("Uploading image to Firebase Storage...")
	// Initialize Firebase app
	ctx := context.Background()
	conf := &firebase.Config{StorageBucket: "live-weather-eefc5.appspot.com"}

	var opts []option.ClientOption
	if _, err := os.Stat("service-account.json"); err == nil {
		opts = append(opts, option.WithCredentialsFile("service-account.json"))
	}
	app, err := firebase.NewApp(ctx, conf, opts...)
	if err != nil {
		return fmt.Errorf("error initializing Firebase app: %w", err)
	}

	// Get storage client
	client, err := app.Storage(ctx)
	if err != nil {
		return fmt.Errorf("error getting Firebase storage client: %w", err)
	}

	// Get the default bucket
	bucket, err := client.DefaultBucket()
	if err != nil {
		return fmt.Errorf("error getting default Firebase storage bucket: %w", err)
	}

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	object := bucket.Object(fileName)
	writer := object.NewWriter(ctx)
	defer writer.Close()

	writer.ObjectAttrs.Metadata = map[string]string{
		"location": fmt.Sprintf(`[%f, %f]`, location[0], location[1]),
	}

	if _, err := io.Copy(writer, resp.Body); err != nil {
		return fmt.Errorf("error writing to Firebase storage: %w", err)
	}

	log.Printf("Successfully uploaded %s to Firebase Storage\n", fileName)
	return nil
}
