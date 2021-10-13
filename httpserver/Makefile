export tag=v1.0
root:
    export ROOT=github.com/davshen/k8s_native_camp_0

build:
    echo "building httpserver binary"
    mkdir -p bin/amd64
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/amd64 .

release: build
    echo "building httpserver container"
    docker build -t shendawei/httpserver:${tag} .

push: release
    echo "pushing shendawei/httpserver"
    docker push shendawei/httpserver:v1.0