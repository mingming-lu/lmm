import io
import re

from datetime import datetime, timezone
from dateutil import parser
from google.cloud import storage
from os import path
from PIL import Image

import on_assert_updated
import on_user_registered


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
    on_assert_updated.run(event, context)


def on_user_registered(event, context):
    on_user_registered.run(event, context)
