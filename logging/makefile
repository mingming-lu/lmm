all: dev

dev:
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d

prod:
	docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d

stop:
	docker-compose -f docker-compose.yml down

restart:
	make stop
	make dev

logs:
	docker-compose -f docker-compose.yml logs -f
