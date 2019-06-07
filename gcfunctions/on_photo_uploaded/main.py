import io

from google.cloud import storage
from os import path
from PIL import Image


def create_thumbnails(event, context):
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
    client = storage.Client()

    filename, buckername = event['name'], event['bucket']
    bucket = client.get_bucket(buckername)
    src = bucket.get_blob(filename)
    buffer = io.BytesIO()
    src.download_to_file(buffer)

    with Image.open(buffer) as image:
        for width in (320, 640, 960, 1280):
            dst = _create_photo_thumbnail(image, width)
            dst.seek(0)

            name, ext = path.splitext(filename)
            thumbnail = bucket.blob(f"{name}_{width}{ext}")
            thumbnail.cache_control = src.cache_control
            thumbnail.acl = src.acl
            thumbnail.upload_from_file(
              dst,
              content_type=src.content_type,
            )


def _create_photo_thumbnail(image: Image.Image, width: int) -> Image.Image:
    img = photo.copy()

    if img.size[0] < width:
        print(f'skip resizing: the width {img.size[0]} is less than {width}')
        return img
    ratio = width / img.size[0]
    height = int(img.size[1] * ratio)
    img.thumbnail((width, height), Image.ANTIALIAS)
    return img
