export tag=v1.1
root:
	export ROOT=github.com/cncamp/golang

build:
	echo "building httpserver binary"
	mkdir -p bin/amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/amd64 .

release: build
	echo "building httpserver container"
	docker build -t shendawei/httpserver:${tag} .

push: release
	echo "pushing shendawei/httpserver"
	docker push shendawei/httpserver:${tag}