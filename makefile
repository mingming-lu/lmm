all: dev

install:
	cd api && make install
	cd app && make install
	cd manager && make install

dev:
	make start

start:
	make start-gateway
	make start-services -j 4

start-gateway:
	cd gateway && make

start-services: start-api start-app start-asset start-manager

start-api:
	cd api && make

start-app:
	cd app && make

start-asset:
	cd asset && make

start-manager:
	cd manager && make

stop:
	make stop-services -j 4
	make stop-gateway

stop-services: stop-api stop-app stop-asset stop-manager

stop-api:
	cd api && make stop

stop-app:
	cd app && make stop

stop-asset:
	cd asset && make stop

stop-manager:
	cd manager && make stop

stop-gateway:
	cd gateway && make stop

restart:
	make stop
	make start

go-build:
	cd api && make build
	cd asset && make build

scale-api:
	docker-compose -f api/docker-compose.yml -f api/docker-compose.${env}.yml up -d --scale api=${n}
