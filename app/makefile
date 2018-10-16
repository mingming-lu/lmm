all: start

install:
	docker-compose run --rm --no-deps app bash -c "rm -rf app/node_modules && npm --prefix /app install"

start:
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d

restart:
	make stop
	make

stop:
	docker-compose down

logs:
	docker-compose logs -f
