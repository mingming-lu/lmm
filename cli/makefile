all:

.PHONY: python
python:
	docker-compose build python
	docker-compose run --rm python

resize-asset-photos:
	docker-compose run --rm python sh -c 'python3 resize_image.py --src_dir /project/asset/static/photos/ --dst_dir /project/asset/static/photos/${width}w --width=${width}'
