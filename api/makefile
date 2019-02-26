all: start

install:
	docker-compose -f docker-compose.yml run --rm --no-deps api go mod download

start:
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d

build: build-api build-cli

build-api:
	docker-compose -f docker-compose.yml run --rm --no-deps api bash -c "cd cmd/api && go install -v"

build-cli:
	docker-compose -f docker-compose.yml run --rm --no-deps api bash -c "cd cmd/cli && go install -v"

# example:
# > make cli env=dev commands="hello-world"
cli:
	docker-compose -f docker-compose.yml -f docker-compose.${env}.yml run --rm api cli ${commands}

prod:
	docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d

restart: stop
	make start

stop:
	docker-compose down

test:
	docker-compose -f docker-compose.yml -f docker-compose.test.yml up --abort-on-container-exit

logs:
	docker-compose logs -f
