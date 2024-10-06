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
	inputs := []Input{{WebcamId: "46", Location: []float64{13.061854, 63.386158}},
		{WebcamId: "61", Location: []float64{12.97008963763537, 63.410562604101635}},
		{WebcamId: "62", Location: []float64{13.063774, 63.436330}},
		{WebcamId: "63", Location: []float64{13.112216429602679, 63.40863152963157}},
		{WebcamId: "77", Location: []float64{13.07677168424639, 63.402590779054385}},
		{WebcamId: "44", Location: []float64{13.079073, 63.427271}},
		{WebcamId: "60", Location: []float64{13.113218812480724, 63.40418386515296}},
		{WebcamId: "45", Location: []float64{13.063259360877352, 63.41586041513909}},
		{WebcamId: "49", Location: []float64{13.182920769760132, 63.38747245909455}},
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
