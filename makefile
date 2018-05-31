.PHONY: all
all:
	make run -j
.PHONY: install
install:
	go get -u github.com/golang/dep/cmd/dep
	cd api && dep ensure
	cd image && dep ensure
	rm -rf manager/node_modules
	rm -rf app/node_modules
	npm --prefix app install
	npm --prefix manager install

.PHONY: run
run: app api image manager docs

.PHONY: app
app: app/package.json
	npm --prefix app run dev

.PHONY: manager
manager: manager/package.json
	npm --prefix manager run dev

.PHONY: api
api: api/main.go
	go run api/main.go

.PHONY: image
image: image/main.go
	go run image/main.go

.PHONY: docker
docker: docker
	cd docker && docker-compose up

.PHONY: test
test: test-api

.PHONY: cli
cli: script
	cd docker && docker-compose run cli bash

.PHONY: test-api
test-api:
	go test -v lmm/api/context/account/appservice
	go test -v lmm/api/context/account/domain/model
	go test -v lmm/api/context/account/domain/repository
	go test -v lmm/api/context/account/domain/service
	go test -v lmm/api/context/account/ui
	go test -v lmm/api/testing
	go test -v lmm/api/usecase/auth
