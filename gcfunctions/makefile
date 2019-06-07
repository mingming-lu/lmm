deploy/on_photo_uploaded:
	gcloud functions deploy create_thumbnails \
		--region=asia-northeast1 \
		--runtime=python37 \
		--trigger-event=google.storage.object.finalize \
		--trigger-resource=${SUBSCRIBE_BUCKET} \
