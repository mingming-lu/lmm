all: dev

docker-compose: docker
	cd docker && docker-compose -f docker-compose.yml $(args) $(cmd)

build:
	make docker-compose cmd=build

install:
	make install-app
	make install-api
	make install-manager

install-app:
	cd docker && docker-compose run --rm app bash -c "rm -rf app/node_modules && npm --prefix /app install"

install-api:
	cd docker && docker-compose run --rm api bash -c "go get -u -v github.com/golang/dep/cmd/dep && cd /go/src/lmm/api && rm -rf vendor && dep ensure -v"

install-manager:
	cd docker && docker-compose run --rm manager bash -c "rm -rf manager/node_modules && npm --prefix /manager install"

prod:
	make docker-compose cmd=up

dev:
	make prod args="-f docker-compose.dev.yml"

test:
	make docker-compose args="-f docker-compose.test.yml" cmd="run $(target)"

test-api:
	make test target=api

cli: script
	make docker-compose cmd="run cli python $(target)"

restart:
	make stop
	make

stop:
	make docker-compose cmd=down
