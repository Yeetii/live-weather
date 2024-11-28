package functions

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/PuerkitoBio/goquery"
	"github.com/Yeetii/live-weather/lib"
)

type SkistarWeatherInput struct {
	Id             string
	Url            string
	TopLocation    []float64
	BottomLocation []float64
}

func init() {
	functions.HTTP("updateSkistarWeather", updateSkistarWeather)
}

func updateSkistarWeather(w http.ResponseWriter, r *http.Request) {
	inputs := []SkistarWeatherInput{
		{Id: "areby", Url: "https://www.skistar.com/Lpv/Forecast?lang=sv&area=areby", TopLocation: nil, BottomLocation: []float64{63.403513916879106, 13.059243790348805}},
		{Id: "hogzon", Url: "https://www.skistar.com/Lpv/Forecast?lang=sv&area=hogzon", TopLocation: []float64{63.42746531861163, 13.07798918790169}, BottomLocation: nil},
		{Id: "bjornen", Url: "https://www.skistar.com/Lpv/Forecast?lang=sv&area=bjornen", TopLocation: []float64{63.40397330198576, 13.112439722480248}, BottomLocation: []float64{63.39058903591519, 13.124520371380962}},
		{Id: "duved", Url: "https://www.skistar.com/Lpv/Forecast?lang=sv&area=duved", TopLocation: []float64{63.40925052213198, 12.933974437408123}, BottomLocation: []float64{63.39653432268454, 12.924465636198212}},
		{Id: "vemdalsskalet", Url: "https://www.skistar.com/Lpv/Forecast?lang=sv&area=vemdalsskalet", TopLocation: []float64{62.483387, 13.956566}, BottomLocation: []float64{62.484503, 13.967102}},
		{Id: "bjornrike", Url: "https://www.skistar.com/Lpv/Forecast?lang=sv&area=bjornrike", TopLocation: []float64{62.41864, 13.98688}, BottomLocation: []float64{62.42142, 13.95809}},
		{Id: "klovsjostorhogna", Url: "https://www.skistar.com/Lpv/Forecast?lang=sv&area=klovsjostorhogna", TopLocation: []float64{62.49811, 14.09203}, BottomLocation: []float64{62.49464, 14.11936}},
	}

	var observations []lib.Observation

	for _, input := range inputs {
		var weather = scrapeCurrentWeather(input.Url)
		log.Println(input.Id)
		log.Println(weather)
		if input.TopLocation != nil {
			id := "skistar-" + input.Id + "-top"
			observation := lib.Observation{Id: &id, Latitude: &input.TopLocation[0], Longitude: &input.TopLocation[1], TemperatureC: &weather.TemperatureTop, WindSpeedMs: &weather.WindSpeedTop, WindGustSpeedMs: &weather.GustWindpeedTop}
			observations = append(observations, observation)
		}
		if input.BottomLocation != nil {
			id := "skistar-" + input.Id + "-bottom"
			observation := lib.Observation{Id: &id, Latitude: &input.BottomLocation[0], Longitude: &input.BottomLocation[1], TemperatureC: &weather.TemperatureBottom, WindSpeedMs: &weather.GustWindspeedBottom, WindGustSpeedMs: &weather.GustWindspeedBottom}
			observations = append(observations, observation)
		}
	}
	lib.UploadObservationsToFirestore(observations)
}

type weatherMeasurement struct {
	TemperatureTop      float64
	TemperatureBottom   float64
	WindSpeedTop        float64
	WindSpeedBottom     float64
	GustWindpeedTop     float64
	GustWindspeedBottom float64
}

func scrapeCurrentWeather(url string) weatherMeasurement {
	res, err := http.Get(url)
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

	tempFields := doc.Find(".lpv-info-weather__text")

	topTempElement := tempFields.First()
	bottomTempElement := tempFields.Eq(1)

	topTemp, err := extractFloat(topTempElement)
	if err != nil {
		log.Println(err)
	}
	bottomTemp, err := extractFloat(bottomTempElement)
	if err != nil {
		log.Println(err)
	}

	windContainers := doc.Find(".lpv-info-list__value")

	topWindContainer := windContainers.First()

	topWindSpeed, err := extractFloat(topWindContainer.Children().First())
	if err != nil {
		log.Println(err)
	}
	gustWindpeedTop, err := extractFloat(topWindContainer.Children().Eq(1))
	if err != nil {
		log.Println(err)
	}

	bottomWindContainer := windContainers.Eq(1)

	bottomWindSpeed, err := extractFloat(bottomWindContainer.Children().First())
	if err != nil {
		log.Println(err)
	}
	gustWindspeedBottom, err := extractFloat(bottomWindContainer.Children().Eq(1))
	if err != nil {
		log.Println(err)
	}

	return weatherMeasurement{
		TemperatureTop:      topTemp,
		TemperatureBottom:   bottomTemp,
		WindSpeedTop:        topWindSpeed,
		GustWindpeedTop:     gustWindpeedTop,
		WindSpeedBottom:     bottomWindSpeed,
		GustWindspeedBottom: gustWindspeedBottom,
	}
}

func extractFloat(element *goquery.Selection) (float64, error) {
	text := element.Text()
	re := regexp.MustCompile(`[-]?\d+\.?\d*`)
	match := re.FindString(text)
	if match == "" {
		return 0, error(fmt.Errorf("no float found in string: %s", text))
	}
	return strconv.ParseFloat(match, 64)
}
