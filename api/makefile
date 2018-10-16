all: start

start:
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d

restart:
	make stop
	make

stop:
	docker-compose down

test:
	docker-compose -f docker-compose.yml -f docker-compose.test.yml up --abort-on-container-exit

logs:
	docker-compose logs -f
