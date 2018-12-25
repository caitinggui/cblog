GOPATH:=$(CURDIR):$(GOPATH)

test:
	go test

test-bench:
	go test -test.bench=".*" 

debug:
	gin --appPort 8089 main.go

run:
	go run main.go

build-linux:
	echo "build linux"
	GOOS=linux GOARCH=amd64 go build -o cblog-service main.go

build:
	echo "build local"
	go build -o cblog-service main.go
