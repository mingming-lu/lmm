all: start

start:
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d

go/mod/tidy:
	docker-compose run --rm --no-deps api go mod tidy

build: build-cli

build-cli:
	docker-compose -f docker-compose.yml run --rm --no-deps api bash -c "cd cmd/cli && go install -v"

# example:
# > make cli env=dev commands="hello-world"
cli:
	docker-compose -f docker-compose.yml -f docker-compose.${env}.yml run --rm api cli ${commands}

restart: stop
	make start

stop:
	docker-compose down

test:
	docker-compose -f docker-compose.yml -f docker-compose.test.yml up --abort-on-container-exit

logs:
	docker-compose logs -f
