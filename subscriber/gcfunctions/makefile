deploy/on_asset_uploaded:
	gcloud functions deploy on_asset_uploaded \
		--memory=256MB \
		--region=asia-northeast1 \
		--retry \
		--runtime=python37 \
		--trigger-event=google.storage.object.finalize \
		--trigger-resource=${SUBSCRIBE_BUCKET}

deploy/on_user_registered:
	gcloud functions deploy on_user_registered \
		--memory=128MB \
		--region=asia-northeast1 \
		--retry \
		--runtime=python37 \
		--trigger-event=providers/cloud.pubsub/eventTypes/topic.publish \
		--trigger-resource=UserRegistered
