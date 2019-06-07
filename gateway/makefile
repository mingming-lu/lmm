all: start

start:
	docker-compose up -d

restart:
	make stop
	make

stop:
	docker-compose down

logs:
	docker-compose logs # no follow and display nginx access log
	docker-compose run --rm nginx sh -c "tail /var/log/nginx/access.log"

reload-nginx:
	docker-compose kill -s SIGHUP nginx
