all: start

start:
	docker-compose up -d

restart:
	make stop
	make

stop:
	docker-compose down

logs:
	docker-compose logs -f
