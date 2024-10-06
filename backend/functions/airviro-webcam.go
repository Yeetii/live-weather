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
	inputs := []Input{{WebcamId: "borga", Location: []float64{15.03789571840728, 64.84199155484801}},
		{WebcamId: "helags", Location: []float64{12.505582249386759, 62.917014196762445}},
		{WebcamId: "ramundberget", Location: []float64{12.37264481898198, 62.69248269325625}},
		{WebcamId: "bydalen", Location: []float64{13.75263354936005, 63.10759607237622}},
	}

	for _, input := range inputs {
		var fileName = fmt.Sprintf("airviro-webcam-%s.jpg", input.WebcamId)
		var url = fmt.Sprintf("https://www.airviro.com/%s/webcam/latestimg.jpg", input.WebcamId)
		lib.UploadToFirebaseStorage(url, fileName, input.Location)
	}
}
