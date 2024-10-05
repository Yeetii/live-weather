package functions

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	firebase "firebase.google.com/go"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/PuerkitoBio/goquery"
	"google.golang.org/api/option"
)

func init() {
	functions.HTTP("fetchSkiStarWebcams", main)
}

func main(w http.ResponseWriter, r *http.Request) {
	var input = struct {
		WebcamId string
		Location []float64
	}{WebcamId: "46", Location: []float64{12.832031, 64.628695}}
	var url = ScrapeWebcamUrl(input.WebcamId)
	var fileName = fmt.Sprintf("skistar-webcam-%s.jpg", input.WebcamId)
	uploadToFirebaseStorage(url, fileName, input.Location)
}

func uploadToFirebaseStorage(url string, fileName string, location []float64) error {
	fmt.Println("Uploading image to Firebase Storage...")
	// Initialize Firebase app
	ctx := context.Background()
	conf := &firebase.Config{StorageBucket: "live-weather-eefc5.appspot.com"}
	app, err := firebase.NewApp(ctx, conf, option.WithCredentialsFile("service-account.json"))
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
		"location": fmt.Sprintf(`{"type": "Point", "coordinates": [%f, %f]}`, location[0], location[1]),
	}

	if _, err := io.Copy(writer, resp.Body); err != nil {
		return fmt.Errorf("error writing to Firebase storage: %w", err)
	}

	log.Printf("Successfully uploaded %s to Firebase Storage\n", fileName)
	return nil
}

func ScrapeWebcamUrl(webcamId string) string {
	res, err := http.Get("https://www.skistar.com/sv/vara-skidorter/are/vinter-i-are/vader-och-backar/webbkameror-are/WebCam/?webcamId=" + webcamId)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("Failed to fetch the webpage: %s", res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the input element with data-range-mapper-value="23" and extract the data-image-url.
	// The site shows the last 24h, with the 23rd image being the latest.
	var imageUrl string
	doc.Find("input.fn-lpv-image-data-holder").Each(func(i int, s *goquery.Selection) {
		// Check if the element has the desired attribute
		if value, exists := s.Attr("data-range-mapper-value"); exists && value == "23" {
			imageUrl, _ = s.Attr("data-image-url")
			fmt.Printf("Webcam Image URL (range 23): %s\n", imageUrl)
		}
	})
	var largeImageUrl = imageUrl[:len(imageUrl)-5]
	return largeImageUrl
}
