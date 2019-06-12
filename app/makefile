all: start

install:
	docker-compose run --rm --no-deps app sh -c "rm -rf app/node_modules/* && npm --prefix /app install"

npm/update:
	docker-compose run --rm --no-deps app npm --prefix /app update

start:
	docker-compose up -d

restart:
	make stop
	make

stop:
	docker-compose down

logs:
	docker-compose logs -f

build:
	docker-compose -f docker-compose.yml run --rm --no-deps app sh -c "npm --prefix /app run build"
