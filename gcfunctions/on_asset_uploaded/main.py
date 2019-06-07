import io
import re

from google.cloud import storage
from os import path
from PIL import Image


def main(event, context):
    """
    event: dict of json: {
      "kind": "storage#object",
      "id": string,
      "selfLink": string,
      "name": string,
      "bucket": string,
      "generation": long,
      "metageneration": long,
      "contentType": string,
      "timeCreated": datetime,
      "updated": datetime,
      "timeDeleted": datetime,
      "temporaryHold": boolean,
      "eventBasedHold": boolean,
      "retentionExpirationTime": datetime,
      "storageClass": string,
      "timeStorageClassUpdated": datetime,
      "size": unsigned long,
      "md5Hash": string,
      "mediaLink": string,
      "contentEncoding": string,
      "contentDisposition": string,
      "contentLanguage": string,
      "cacheControl": string,
      "metadata": {
        (key): string
      }
    }
    """

    if event['contentType'].startswith('image/'):

        if event['name'].startswith('thumbnail/'):
            print(f"skip thumbnail: {event['name']}")
            return

        return create_thumbnails(event['bucket'], event['name'])

    print(f'unhandled event: {event}')


def create_thumbnails(bucketname, filename):
    client = storage.Client()

    bucket = client.get_bucket(bucketname)

    src = bucket.get_blob(filename)
    buffer = io.BytesIO()
    src.download_to_file(buffer)

    format = src.content_type.replace('image/', '')

    with Image.open(buffer) as image:
        for width in (320, 640, 960, 1280):
            thumbnail = _create_thumbnail(image, width)
            data = io.BytesIO()
            thumbnail.save(data, format=format)

            dst = bucket.blob(f"thumbnail/w{width}/{filename}")
            dst.cache_control = src.cache_control
            dst.upload_from_file(
                data,
                content_type=src.content_type,
                predefined_acl='publicRead',
                rewind=True,
            )


def _create_thumbnail(image: Image.Image, width: int) -> Image.Image:
    img = image.copy()

    if img.size[0] < width:
        print(f'skip resizing: the width {img.size[0]} is less than {width}')
        return img
    ratio = width / img.size[0]
    height = int(img.size[1] * ratio)
    img.thumbnail((width, height), Image.ANTIALIAS)
    return img
