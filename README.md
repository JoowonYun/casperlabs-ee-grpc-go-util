#   Casperlabs-ee-grpc-go-util  ![Build Status](https://travis-ci.com/hdac-io/casperlabs-ee-grpc-go-util.svg?token=GLBy52vf7tteTZjpJNu9&branch=master)


This module is for using GRPC and UTIL in GO to utilize Casperlabs Execution-Engine.


## Install
- Download to module in $GOPATH
```bash
$ go get -u github.com/hdac-io/casperlabs-ee-grpc-go-util
```

## Example
- Running to Casperlabs Execution-Engine 
```bash
$ git clone https://github.com/hdac-io/CasperLabs.git

$ cd execution-engine

$ cargo build

$ ./target/debug/casperlabs-engine-grpc-server ~/.casperlabs/.casper-node.sock
```
- Install "casperlabs-ee-grpc-go-util" module

- Import "casperlabs-ee-grpc-go-util"
```go
import (
    "github.com/hdac-io/casperlabs-ee-grpc-go-util/util"
	"github.com/hdac-io/casperlabs-ee-grpc-go-util/grpc"
)

func main() {
    client := grpc.Connect(`/.casperlabs/.casper-node.sock`)

	mintTokenCode := util.LoadWasmFile("./example/contracts/mint_token.wasm")
}
```

## Integration test
- Running casperlabs-engine-grpc-server
```bash
$ $EXECUTION_ENGINE_HOME/target/debug/casperlabs-engine-grpc-server ~/.casperlabs/.casper-node.sock
```
- Run
```bash
$ make integration-test
```

## Unit test
```bash
$ make test
```

## Develop
- Install protobuf
- Generate proto file
```bash
$ make proto
```
