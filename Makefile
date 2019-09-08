
test:
	# 加上-v就会输出详细信息（包括自定义日志）
	#go test -v ./service ./utils
	#go test -v ./models ./main_test.go -configPath $(CURDIR)/config/dev/
	export CBLOG_CONFIG_PATH=$(CURDIR)/config/dev/ && go test -v ./...

test-bench:
	go test -test.bench=".*" ./... -configPath $(CURDIR)/config/dev/

debug:
	gin --appPort 8089 main.go

run:
	go build -o cblog main.go && ./cblog

build-linux:
	echo "build linux"
	GOOS=linux GOARCH=amd64 go build -o cblog-service main.go

build:
	echo "build local"
	go build -o cblog-service main.go

docs:
	type apidoc || echo "请 npm i -g apidoc"
	apidoc -e ./_ref/ -o ./doc/
