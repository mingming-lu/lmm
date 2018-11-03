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

prod:
	docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
