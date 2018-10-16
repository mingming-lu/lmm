all: start

install:
	docker-compose run --rm --no-deps manager bash -c "rm -rf manager/node_modules && npm --prefix /manager install"

start:
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d

restart:
	make stop
	make

stop:
	docker-compose down

logs:
	docker-compose logs -f
