TAG ?= "latest"

build: clean
	go build -o updater main.go

update:
	go mod tidy
	go mod vendor

container:
	docker build -t updater:$(TAG) .

clean: 
	rm -f updater
