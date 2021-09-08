TAG ?= "latest"

build:
	go build -o updater main.go

update:
	go mod tidy
	go mod vendor

container:
	docker build -t updater:$(TAG) .
