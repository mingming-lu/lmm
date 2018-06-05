all: dev

docker-compose: docker
	cd docker && docker-compose -f docker-compose.yml $(args) $(cmd)

build:
	make docker-compose cmd=build

# TODO move to Dockerfile
install:
	cd docker && docker-compose run --rm api bash -c "go get -u github.com/golang/dep/cmd/dep && cd /go/src/lmm/api && dep ensure"
	cd docker && docker-compose run --rm image bash -c "go get -u github.com/golang/dep/cmd/dep && cd /go/src/lmm/image && dep ensure"
	rm -rf app/node_modules
	cd docker && docker-compose run --rm app bash -c "npm i npm@latest -g && npm --prefix /app install"
	rm -rf manager/node_modules
	cd docker && docker-compose run --rm manager bash -c "npm i npm@latest -g && npm --prefix /manager install"

prod:
	make docker-compose cmd=up

dev:
	make prod args="-f docker-compose.dev.yml"

test:
	make docker-compose args="-f docker-compose.test.yml" cmd="run $(target)"

cli: script
	make docker-compose cmd="run cli python $(target)"

stop:
	make docker-compose cmd=down
