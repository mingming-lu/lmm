.PHONY: all
all:
	make run -j
.PHONY: install
install:
	go get -v github.com/akinaru-lu/elesion
	go get -v github.com/go-sql-driver/mysql
	go get -v github.com/google/uuid
	go get -v github.com/stretchr/testify/assert
	go get -v github.com/stretchr/testify/mock
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

.PHONY: test
test: test-api

.PHONY: test-api
test-api:
	go test -v lmm/api/context/account/domain/model
	go test -v lmm/api/context/account/domain/repository
	go test -v lmm/api/context/account/ui
	go test -v lmm/api/context/account/appservice
	go test -v lmm/api/testing
