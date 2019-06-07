all: dev

install: install-api install-app install-manager

install-api:
	cd api && make build

install-app:
	cd app && make install

install-manager:
	cd manager && make install

init: init-api

init-api:
	cd api && make init

dev:
	make start

start:
	make start-gateway
	make start-services -j 4

start-gateway:
	cd gateway && make

start-services: start-api start-app start-manager

start-api:
	cd api && make

start-app:
	cd app && make

start-manager:
	cd manager && make

stop:
	make stop-services -j 4
	make stop-gateway

stop-services: stop-api stop-app stop-manager

stop-api:
	cd api && make stop

stop-app:
	cd app && make stop

stop-manager:
	cd manager && make stop

stop-gateway:
	cd gateway && make stop

restart:
	make stop
	make start
