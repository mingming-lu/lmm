.PHONY: all
all:
	make run -j

.PHONY: install
install:
	go get -u github.com/akinaru-lu/elesion
	go get -u github.com/go-sql-driver/mysql
	go get -u github.com/google/uuid
	go get -u github.com/stretchr/testify/assert
	npm --prefix app install
	npm --prefix manager install

.PHONY: run
run: app api image manager

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

