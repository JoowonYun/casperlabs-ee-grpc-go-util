.PHONY: test
test:
	go test ./util
	#go test ./grpc

.PHONY: clean
clean:
	go clean ./...

.PHONY: install
install:
	go install ./grpc ./util

.PHONY: proto
proto:
	protoc -I protobuf protobuf/io/casperlabs/casper/consensus/state/state.proto --gogo_out=plugins=grpc:$$GOPATH/src
	protoc -I protobuf protobuf/io/casperlabs/casper/consensus/consensus.proto --gogo_out=plugins=grpc:$$GOPATH/src
	protoc -I protobuf protobuf/io/casperlabs/ipc/transforms/transforms.proto --gogo_out=plugins=grpc:$$GOPATH/src
	protoc -I protobuf protobuf/io/casperlabs/ipc/ipc.proto --gogo_out=plugins=grpc:$$GOPATH/src

.PHONY: example
example:
	go run ./example/hello.go
