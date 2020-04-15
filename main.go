package main

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

func init() {
	log.Println("...astronomy-img...")
}

func check(err error) {
	if err != nil {
		log.Panicln(err)
	}
}

// get image URL by requestng json metadata from site
func getImageURL(siteURL string) (string, error) {
	myClient := &http.Client{Timeout: 10 * time.Second}

	// get metadata from server
	resp, err := myClient.Get(siteURL)
	check(err)
	defer resp.Body.Close()
	fmt.Println("Response status:", resp.Status)

	// parse the json response
	type responseData struct {
		URL *string `json:"url"`
	}
	data := &responseData{}
	err = json.NewDecoder(resp.Body).Decode(data)
	check(err)

	// error is url is missing from json response
	if data.URL == nil {
		return "", fmt.Errorf("'url' missing from json response")
	}

	// all good - return the image url
	return *data.URL, nil
}

func getImage(url string, fileName string) error {
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
func buildFileName(fullURLFile string) (string, error) {
	var fileName = ""

	fileURL, err := url.Parse(fullURLFile)
	if err != nil {
		return fileName, err
	}

	path := fileURL.Path
	segments := strings.Split(path, "/")
	fileName = segments[len(segments)-1]

	return fileName, nil
}

func main() {
	// image URL
	imageURLStr, err := getImageURL(nasaURL)
	check(err)
	log.Printf("Image URL is %v\n", imageURLStr)

	// image itself
	imageFilename, err := buildFileName(imageURLStr)
	check(err)
	getImage(imageURLStr, imageFilename)
}
