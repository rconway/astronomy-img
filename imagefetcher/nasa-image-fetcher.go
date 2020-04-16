package imagefetcher

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var nasaURL = "https://api.nasa.gov/planetary/apod?api_key=DEMO_KEY"

// NasaImageFetcher Fetch images from NASA website
type NasaImageFetcher struct{}

func init() {
	log.Println("...nasa-image-fetcher...")
}

func (f *NasaImageFetcher) check(err error) {
	if err != nil {
		log.Panicln(err)
	}
}

// get image URL by requestng json metadata from site
func (f *NasaImageFetcher) getImageURLFromMetadata(siteURL string) (string, error) {
	myClient := &http.Client{Timeout: 10 * time.Second}

	// get metadata from server
	resp, err := myClient.Get(siteURL)
	f.check(err)
	defer resp.Body.Close()
	fmt.Println("Response status:", resp.Status)

	// parse the json response
	type responseData struct {
		URL *string `json:"url"`
	}
	data := &responseData{}
	err = json.NewDecoder(resp.Body).Decode(data)
	f.check(err)

	// error if url is missing from json response
	if data.URL == nil {
		return "", fmt.Errorf("'url' missing from json response")
	}

	// all good - return the image url
	return *data.URL, nil
}

func (f *NasaImageFetcher) getImageFromURL(url string, fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	myClient := &http.Client{Timeout: 10 * time.Second}
	resp, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)

	return err
}

// get the filename from the end of the url
func (f *NasaImageFetcher) getFilenameFromURL(fullURLFile string) (string, error) {
	var filename = ""

	fileURL, err := url.Parse(fullURLFile)
	if err != nil {
		return filename, err
	}

	path := fileURL.Path
	segments := strings.Split(path, "/")
	filename = segments[len(segments)-1]

	return filename, nil
}

// Fetch Fetch the image from the provider
func (f *NasaImageFetcher) Fetch() string {
	// image URL
	imageURLStr, err := f.getImageURLFromMetadata(nasaURL)
	f.check(err)
	log.Printf("Image URL is %v\n", imageURLStr)

	// image itself
	imageFilename, err := f.getFilenameFromURL(imageURLStr)
	f.check(err)
	f.getImageFromURL(imageURLStr, imageFilename)

	return imageFilename
}
