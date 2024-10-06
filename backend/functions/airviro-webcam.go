package functions

import (
	"fmt"
	"net/http"

	"github.com/Yeetii/live-weather/lib"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("fetchAirviroWebcams", collectAirviroWebcams)
}

func collectAirviroWebcams(w http.ResponseWriter, r *http.Request) {
	inputs := []Input{{WebcamId: "borga", Location: []float64{0, 0}},
		{WebcamId: "helags", Location: []float64{0, 0}},
		{WebcamId: "ramundberget", Location: []float64{0, 0}},
		{WebcamId: "bydalen", Location: []float64{0, 0}},
	}

	for _, input := range inputs {
		var fileName = fmt.Sprintf("airviro-webcam-%s.jpg", input.WebcamId)
		var url = fmt.Sprintf("https://www.airviro.com/%s/webcam/latestimg.jpg", input.WebcamId)
		lib.UploadToFirebaseStorage(url, fileName, input.Location)
	}
}
