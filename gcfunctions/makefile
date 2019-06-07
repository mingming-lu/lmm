deploy/on_photo_uploaded:
	gcloud functions deploy create_thumbnails \
		--region=asia-northeast1 \
		--runtime python37 \
		--trigger-bucket ${SUBSCRIBE_BUCKET} \
		--trigger-event google.storage.object.finalize