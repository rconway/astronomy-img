import image_fetcher
import nasa


def main():
    fetcher = nasa.NasaImageFetcher()
    try:
        fetcher.get_image()
    except image_fetcher.ImageFetcherError as e:
        print(f"Problem fetching image: {e}")


if __name__ == "__main__":
    main()
