import image_fetcher
import requests

nasa_URL = "https://api.nasa.gov/planetary/apod?api_key=DEMO_KEY"


class NasaMetadataError(image_fetcher.ImageFetcherError):
    def __init__(self, reason):
        self.reason = reason


def get_image_URL_from_metadata(site_url):
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
            f"Could not retrieve metadata due to exception: {type(e)}")
    else:
        print(f"Image URL: {json_metadata['url']}")


def get_image():
    get_image_URL_from_metadata(nasa_URL)
