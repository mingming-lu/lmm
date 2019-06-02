all: start

build:
	docker-compose -f docker-compose.yml run --rm --no-deps asset bash -c "cd cmd/asset && go install -v"

start:
	docker-compose -f docker-compose.yml up -d

restart:
	make stop
	make

stop:
	docker-compose down

logs:
	docker-compose logs -f

prod:
	docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
