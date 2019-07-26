all: start

install:
	docker-compose run --rm --no-deps manager sh -c "rm -rf manager/node_modules/* && npm --prefix /manager install"

npm/update:
	docker-compose run --rm --no-deps manager npm --prefix /manager update

npm/audit:
	docker-compose run --rm --no-deps manager npm --prefix /manager audit

npm/audit-fix:
	docker-compose run --rm --no-deps manager npm --prefix /manager audit fix

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
	docker-compose -f docker-compose.yml run --rm --no-deps manager sh -c "npm --prefix /manager run build"
