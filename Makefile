all: build web

build:
	GO15VENDOREXPERIMENT=1 go build

save-deps:
	GO15VENDOREXPERIMENT=1 godep save ./...

test:
	GO15VENDOREXPERIMENT=1 go test

run:
	./certrotate
