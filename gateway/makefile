all: start

start:
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d

restart:
	make stop
	make

stop:
	docker-compose down

logs:
	docker-compose logs # no follow and display nginx access log
	docker-compose run --rm nginx bash -c "tail /var/log/nginx/access.log"

prod:
	docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
