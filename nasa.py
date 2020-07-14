import image_fetcher
import requests

nasa_URL = "https://api.nasa.gov/planetary/apod?api_key=DEMO_KEY"


class NasaMetadataError(image_fetcher.ImageFetcherError):
    def __init__(self, reason):
        super().__init__(reason)


class NasaImageFetcher(image_fetcher.ImageFetcher):
    def __init__(self):
        super().__init__()

    def get_image_URL_from_metadata(self, site_url):
        try:
            r = requests.get(nasa_URL, timeout=5)
            r.raise_for_status()
            json_metadata = r.json()
        except requests.exceptions.ConnectTimeout:
            raise NasaMetadataError(
                f"Timeout retrieving metadata from URL {nasa_URL}")
        except requests.HTTPError:
            raise NasaMetadataError(
                f"Could not retrieve metadata from URL {nasa_URL}")
        except ValueError:
            raise NasaMetadataError("Invalid metadata json")
        except Exception as e:
            raise NasaMetadataError(
                f"Could not retrieve metadata due to exception: {type(e)} - {e}")

        return json_metadata['url']

    def get_image(self):
        # Today's image URL
        image_url = self.get_image_URL_from_metadata(nasa_URL)
        print(f"Image URL: {image_url}")

        # Extract filename for saving
        filename = image_url.split('/')[-1]
        print(f"Image filename: {filename}")

        # Request the image
        response = requests.get(image_url)

        # Save image from response
        with open(filename, 'wb') as f:
            f.write(response.content)

        # Dump HTTP meta-data
        print(response.status_code)
        print(response.headers['content-type'])
        print(response.encoding)
