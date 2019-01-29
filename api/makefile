all: start

install:
	docker-compose -f docker-compose.yml run --rm --no-deps api go mod download

start:
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d

build:
	docker-compose -f docker-compose.yml run --rm --no-deps api bash -c "cd cmd/api && go install -v"
	docker-compose -f docker-compose.yml run --rm --no-deps api bash -c "cd cmd/cli && go install -v"

# example:
# > make cli env=dev commands="hello-world"
cli:
	docker-compose -f docker-compose.yml -f docker-compose.${env}.yml run --rm api cli ${commands}

prod:
	docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d

restart:
	make stop
	make

stop:
	docker-compose down

test:
	docker-compose -f docker-compose.yml -f docker-compose.test.yml up --abort-on-container-exit

logs:
	docker-compose logs -f
