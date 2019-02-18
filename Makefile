GOPATH:=$(CURDIR):$(GOPATH)

test:
	# 加上-v就会输出详细信息（包括自定义日志）
	go test -v ./... -configPath $(CURDIR)/config/dev/

test-bench:
	go test -test.bench=".*" ./... -configPath $(CURDIR)/config/dev/

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

docs:
	type apidoc || echo "请 npm i -g apidoc"
	apidoc -e ./_ref/ -o ./doc/
