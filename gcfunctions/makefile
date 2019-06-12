deploy/on_asset_uploaded:
	cd on_asset_uploaded \
	&& gcloud functions deploy on_asset_uploaded \
		--memory=256MB \
		--region=asia-northeast1 \
		--retry \
		--runtime=python37 \
		--trigger-event=google.storage.object.finalize \
		--trigger-resource=${SUBSCRIBE_BUCKET}
