
GOPATH:=$(shell go env GOPATH)
MODIFY=Mgithub.com/micro/go-micro/api/proto/api.proto=github.com/micro/go-micro/v2/api/proto

.PHONY: proto
proto:
    
	protoc --proto_path=. --micro_out=${MODIFY}:. --go_out=${MODIFY}:. proto/rooms/rooms.proto
    

.PHONY: build
build: proto

	go build -o rooms-service *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t rooms-service:latest
