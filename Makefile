GOPATH:=$(CURDIR):$(GOPATH)

test:
	go test

test-bench:
	go test -test.bench=".*" 

debug:
	gin --appPort 8089 main.go
