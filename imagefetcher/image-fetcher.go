package imagefetcher

// ImageFetcher Interface for types that can fetch an image
type ImageFetcher interface {
	Fetch() (string, error)
}
