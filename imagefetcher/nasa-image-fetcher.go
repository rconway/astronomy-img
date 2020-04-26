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

// Module initialiser
func init() {
	log.Println("...nasa-image-fetcher...")
}

// URL for Nasa images
var nasaURL = "https://api.nasa.gov/planetary/apod?api_key=DEMO_KEY"

// NasaImageFetcher Fetch images from NASA website
type NasaImageFetcher struct {
	webClient *http.Client
}

// NasaMetadataResponse json structure for metadata response
type NasaMetadataResponse struct {
	URL *string `json:"url"`
}

// Instantiate (if required) and return the web client
func (f *NasaImageFetcher) getWebClient() *http.Client {
	if f.webClient == nil {
		f.webClient = &http.Client{Timeout: 3 * time.Second}
	}
	return f.webClient
}

// get image URL by requestng json metadata from site
func (f *NasaImageFetcher) getImageURLFromMetadata(siteURL string) (string, error) {
	// get metadata from server
	resp, err := f.getWebClient().Get(siteURL)
	if err != nil {
		return "", fmt.Errorf("Problem retrieving metadata: %w", err)
	}
	defer resp.Body.Close()
	log.Println("Response status:", resp.Status)

	// parse the json response
	data := &NasaMetadataResponse{}
	err = json.NewDecoder(resp.Body).Decode(data)
	if err != nil {
		return "", fmt.Errorf("Problem parsing metadata: %w", err)
	}

	// error if url is missing from json response
	if data.URL == nil {
		return "", fmt.Errorf("'url' missing from json response")
	}

	// all good - return the image url
	return *data.URL, nil
}

func (f *NasaImageFetcher) getImageFromURL(url string, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("Problem creating output file '%v': %w", filename, err)
	}
	defer file.Close()

	resp, err := f.getWebClient().Get(url)
	if err != nil {
		return fmt.Errorf("Problem retrieving image from URL `%v`: %w", url, err)
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
		return filename, fmt.Errorf("Could not deduce filename from URL '%v': %w", fullURLFile, err)
	}

	path := fileURL.Path
	segments := strings.Split(path, "/")
	filename = segments[len(segments)-1]

	return filename, nil
}

// Fetch Fetch the image from the provider
func (f *NasaImageFetcher) Fetch() (string, error) {
	imageFilename := ""
	err := error(nil)

	// image URL
	imageURLStr, err := f.getImageURLFromMetadata(nasaURL)
	if err == nil {
		log.Printf("Image URL is %v\n", imageURLStr)
		// image itself
		imageFilename, err = f.getFilenameFromURL(imageURLStr)
	}

	if err == nil {
		err = f.getImageFromURL(imageURLStr, imageFilename)
	}

	if err != nil {
		imageFilename = ""
	}

	return imageFilename, err
}
