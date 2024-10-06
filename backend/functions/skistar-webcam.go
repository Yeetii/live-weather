package functions

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Yeetii/live-weather/lib"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/PuerkitoBio/goquery"
)

type Input struct {
	WebcamId string
	Location []float64
}

func init() {
	functions.HTTP("fetchSkiStarWebcams", main)
}

// Tege berg = 61
// Tväråvalvet = 62
// Fjällgård = 63
// Stjärntorget = 77
// Kabin = 44
// Sadel = 60
// VM-platå = 45
// Förberget = 49

func main(w http.ResponseWriter, r *http.Request) {
	inputs := []Input{{WebcamId: "46", Location: []float64{12.832031, 64.628695}},
		{WebcamId: "61", Location: []float64{0, 0}},
		{WebcamId: "62", Location: []float64{0, 0}},
		{WebcamId: "63", Location: []float64{0, 0}},
		{WebcamId: "77", Location: []float64{0, 0}},
		{WebcamId: "44", Location: []float64{0, 0}},
		{WebcamId: "60", Location: []float64{0, 0}},
		{WebcamId: "45", Location: []float64{0, 0}},
		{WebcamId: "49", Location: []float64{0, 0}},
	}

	for _, input := range inputs {
		var url = scrapeWebcamUrl(input.WebcamId)
		var fileName = fmt.Sprintf("skistar-webcam-%s.jpg", input.WebcamId)
		lib.UploadToFirebaseStorage(url, fileName, input.Location)
	}
}

func scrapeWebcamUrl(webcamId string) string {
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
