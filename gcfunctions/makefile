deploy/on_photo_uploaded:
	gcloud functions deploy create_thumbnails \
		--ignore-file=makefile
		--memory=128MB \
		--region=asia-northeast1 \
		--runtime=python37 \
		--trigger-event=google.storage.object.finalize \
		--trigger-resource=${SUBSCRIBE_BUCKET} \
