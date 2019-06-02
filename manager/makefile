all: start

install:
	docker-compose run --rm --no-deps manager sh -c "rm -rf manager/node_modules/* && npm --prefix /manager install"

npm/update:
	docker-compose run --rm --no-deps manager npm --prefix /manager update

start:
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d

restart:
	make stop
	make

stop:
	docker-compose down

logs:
	docker-compose logs -f

build:
	docker-compose -f docker-compose.yml -f docker-compose.prod.yml run --rm --no-deps manager sh -c "npm --prefix /manager run build"
