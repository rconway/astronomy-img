package main

import (
	"log"

	"github.com/rconway/astronomy-img/imagefetcher"
)

func init() {
	log.Println("...astronomy-img...")
}

func main() {
	var fetcher imagefetcher.ImageFetcher
	fetcher = &imagefetcher.NasaImageFetcher{}
	filename, err := fetcher.Fetch()

	if err == nil {
		log.Printf("Image retrieved to filename: %v\n", filename)
	} else {
		log.Printf("Error fetching image: %v\n", err)
	}
}
