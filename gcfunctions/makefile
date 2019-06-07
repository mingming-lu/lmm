deploy/on_asset_uploaded:
	cd on_asset_uploaded \
	&& gcloud functions deploy main \
		--memory=128MB \
		--region=asia-northeast1 \
		--runtime=python37 \
		--trigger-event=google.storage.object.finalize \
		--trigger-resource=${SUBSCRIBE_BUCKET}
