.PHONY: all
all:
	make run -j

.PHONY: install
install:
	go get -u github.com/akinaru-lu/elesion
	go get -u github.com/go-sql-driver/mysql
	npm --prefix app install

.PHONY: run
run: app api image

.PHONY: app
app: app/package.json
	npm --prefix app run dev

.PHONY: api
api: api/main.go
	go run api/main.go

.PHONY: image
image: image/main.go
	go run image/main.go

