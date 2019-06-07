import io

from google.cloud import storage
from os import path
from PIL import Image


def create_photo_thumbnails(event, context):
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

    original = event
    bucket = client.get_bucket(original['bucket'])
    src = bucket.get_blob(original['name'])
    buffer = io.BytesIO()
    blob.download_to_file(buffer)

    with Image.open(buffer) as image:
        for width in (320, 640, 960, 1280):
            dst = _create_photo_thumbnail(image, width)
            name, ext = path.splitext(original['name'])
            bucket.blob(f"{name}_{width}{ext}").upload_from_file(dst)


def _create_photo_thumbnail(image: Image.Image, width: int) -> Image.Image:
    img = photo.copy()

    if img.size[0] < width:
        print(f'skip resizing: the width {img.size[0]} is less than {width}')
        return img
    ratio = width / img.size[0]
    height = int(img.size[1] * ratio)
    img.thumbnail((width, height), Image.ANTIALIAS)
    return img
