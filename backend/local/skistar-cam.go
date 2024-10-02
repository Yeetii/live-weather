package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
)

func downloadImage(url string, filePath string) error {
	// Send HTTP GET request to the image URL
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch image: %w", err)
	}
	defer resp.Body.Close()

	// Check if the response status is OK (status code 200)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Create the file where the image will be saved
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Copy the image data from the response body to the file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write image to file: %w", err)
	}

	return nil
}

// ScrapeWebcamInfo scrapes the webcam image and title from the page
func ScrapeWebcamInfo() string {
	// Send an HTTP request to fetch the page
	res, err := http.Get("https://www.skistar.com/sv/vara-skidorter/are/vinter-i-are/vader-och-backar/webbkameror-are/WebCam/?webcamId=46")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("Failed to fetch the webpage: %s", res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the input element with data-range-mapper-value="23" and extract the data-image-url
	var imageUrl string
	doc.Find("input.fn-lpv-image-data-holder").Each(func(i int, s *goquery.Selection) {
		// Check if the element has the desired attribute
		if value, exists := s.Attr("data-range-mapper-value"); exists && value == "23" {
			// Extract the data-image-url attribute
			imageUrl, _ = s.Attr("data-image-url")
			fmt.Printf("Webcam Image URL (range 23): %s\n", imageUrl)
		}
	})
	return imageUrl
}

func main() {
	var url = ScrapeWebcamInfo()
	fmt.Println("Downloading image..." + url)
	downloadImage(url, "webcam.jpg")
}
