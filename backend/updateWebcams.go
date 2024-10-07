package functions

import (
	"fmt"
	"net/http"

	"github.com/Yeetii/live-weather/lib"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

type CommonInput struct {
	WebcamId string
	Location []float64
	ImageUrl string
}

func init() {
	functions.HTTP("updateWebcams", UpdateWebcams)
}

func UpdateWebcams(w http.ResponseWriter, r *http.Request) {
	inputs := []CommonInput{{WebcamId: "borga", Location: []float64{15.03789571840728, 64.84199155484801}, ImageUrl: "https://www.airviro.com/borga/webcam/latestimg.jpg"},
		{WebcamId: "helags", Location: []float64{12.505582249386759, 62.917014196762445}, ImageUrl: "https://www.airviro.com/helags/webcam/latestimg.jpg"},
		{WebcamId: "ramundberget", Location: []float64{12.37264481898198, 62.69248269325625}, ImageUrl: "https://www.airviro.com/ramundberget/webcam/latestimg.jpg"},
		{WebcamId: "bydalen", Location: []float64{13.75263354936005, 63.10759607237622}, ImageUrl: "https://www.airviro.com/bydalen/webcam/latestimg.jpg"},
		{WebcamId: "trillevallen", Location: []float64{13.206262037974694, 63.25443937020817}, ImageUrl: "https://api.trafikinfo.trafikverket.se/v2/Images/RoadConditionCamera_39635528.Jpeg?type=fullsize&maxage=140"},
		{WebcamId: "gevsjön", Location: []float64{12.702011476723557, 63.36705212550013}, ImageUrl: "https://api.trafikinfo.trafikverket.se/v2/Images/RoadConditionCamera_39635384.Jpeg?type=fullsize&maxage=140"},
		{WebcamId: "handöl", Location: []float64{12.382941421804786, 63.26835926669674}, ImageUrl: "https://api.trafikinfo.trafikverket.se/v2/Images/RoadConditionCamera_39635520.Jpeg?type=fullsize&maxage=140"},
		{WebcamId: "medstugan", Location: []float64{12.407996503920517, 63.519546112666376}, ImageUrl: "https://api.trafikinfo.trafikverket.se/v2/Images/RoadConditionCamera_39626819.Jpeg?type=fullsize&maxage=140"},
		{WebcamId: "storlien", Location: []float64{12.088252196720449, 63.31759038262924}, ImageUrl: "https://api.trafikinfo.trafikverket.se/v2/Images/RoadConditionCamera_39636227.Jpeg?type=fullsize&maxage=140"},
		{WebcamId: "nedalshytta", Location: []float64{12.101315126910368, 62.97826646239796}, ImageUrl: "https://metnet.no/custcams/nedalshytta/laget/webcam_hd.jpg"},
		{WebcamId: "meråker", Location: []float64{11.679622045416139, 63.456829044603644}, ImageUrl: "https://metnet.no/custcams/merakeralpin2/laget/webcam_hd.jpg"},
	}

	for _, input := range inputs {
		var fileName = fmt.Sprintf("webcam-%s.jpg", input.WebcamId)
		var url = fmt.Sprintf(input.ImageUrl)
		lib.UploadToFirebaseStorage(url, fileName, input.Location)
	}
}
