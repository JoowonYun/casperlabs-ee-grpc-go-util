test:
	go test ./util
	#go test ./grpc

clean:
	go clean ./...

install:
	go install ./grpc ./util

proto:
	protoc -I protobuf protobuf/io/casperlabs/casper/consensus/state/state.proto --go_out=plugins=grpc:$$GOPATH/src
	protoc -I protobuf protobuf/io/casperlabs/casper/consensus/consensus.proto --go_out=plugins=grpc:$$GOPATH/src
	protoc -I protobuf protobuf/io/casperlabs/ipc/transforms.proto --go_out=plugins=grpc:$$GOPATH/src
	protoc -I protobuf protobuf/io/casperlabs/ipc/ipc.proto --go_out=plugins=grpc:$$GOPATH/src

example:
	go run ./example/hello.go