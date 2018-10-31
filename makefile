all: dev

dev:
	make start

start:
	cd gateway && make
	cd messaging && make
	cd api && make
	cd app && make
	cd asset && make
	cd manager && make
	cd docs && make

stop:
	cd api && make stop
	cd app && make stop
	cd asset && make stop
	cd manager && make stop
	cd docs && make stop
	cd messaging && make stop
	cd gateway && make stop

restart:
	make stop
	make start

install:
	cd api && make install
	cd app && make install
	cd manager && make install
