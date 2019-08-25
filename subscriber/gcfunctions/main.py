import io
import re

from datetime import datetime, timezone
from dateutil import parser
from google.cloud import storage
from os import path
from PIL import Image

from on_asset_uploaded import run as run_on_asset_uploaded
from on_user_registered import run as run_on_user_registered


def on_asset_uploaded(event, context):
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
    run_on_asset_uploaded(event, context)


def on_user_registered(event, context):
    run_on_user_registered(event, context)
