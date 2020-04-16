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
	filename := fetcher.Fetch()

	log.Printf("Image retrieved to filename: %v\n", filename)
}
