package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/hdac-io/casperlabs-ee-grpc-go-util/grpc"
	"github.com/hdac-io/casperlabs-ee-grpc-go-util/protobuf/io/casperlabs/casper/consensus"
	"github.com/hdac-io/casperlabs-ee-grpc-go-util/protobuf/io/casperlabs/ipc"
	"github.com/hdac-io/casperlabs-ee-grpc-go-util/protobuf/io/casperlabs/ipc/transforms"
	"github.com/hdac-io/casperlabs-ee-grpc-go-util/util"
)

func stressTest() {
	// Init variable
	println("This is stress test of EE. Be sure that the running EE is RELEASE BUILD!!!")
	println("Setting genesis account...")
	emptyStateHash := util.DecodeHexString(util.StrEmptyStateHash)
	rootStateHash := emptyStateHash
	genesisAddress := util.DecodeHexString("d70243dd9d0d646fd6df282a8f7a8fa05a6629bec01d8024c3611eb1c1fb9f84")
	systemContract := make([]byte, 32)

	chainName := "hdac"
	costs := map[string]uint32{
		"regular":            1,
		"div-multiplier":     16,
		"mul-multiplier":     4,
		"mem-multiplier":     2,
		"mem-initial-pages":  4096,
		"mem-grow-per-page":  8192,
		"mem-copy-per-byte":  1,
		"max-stack-height":   65536,
		"opcodes-multiplier": 3,
		"opcodes-divisor":    8}

	protocolVersion := util.MakeProtocolVersion(1, 0, 0)

	// Connect to ee sock.
	socketPath := os.Getenv("HOME") + `/.casperlabs/.casper-node.sock`
	client := grpc.Connect(socketPath)

	genesisConfig, err := util.GenesisConfigMock(
		chainName, genesisAddress, "500000000000000", "10000000000", protocolVersion, costs,
		"./example/contracts/mint_install.wasm", "./example/contracts/pos_install.wasm")
	if err != nil {
		fmt.Printf("Bad GenesisConfigMock err : %v", err)
		return
	}

	//var parentStateHash , effects
	response, err := grpc.RunGenesis(client, genesisConfig)
	if err != nil {
		panic(err)
	}

	var parentStateHash []byte
	var effects []*transforms.TransformEntry

	switch response.GetResult().(type) {
	case *ipc.GenesisResponse_Success:
		parentStateHash = response.GetSuccess().GetPoststateHash()
		effects = response.GetSuccess().GetEffect().GetTransformMap()
	case *ipc.GenesisResponse_FailedDeploy:
		panic(response.GetFailedDeploy().GetMessage())
	}

	postStateHash, bonds, errMessage := grpc.Commit(client, rootStateHash, effects, protocolVersion)
	if bytes.Equal(postStateHash, parentStateHash) {
		rootStateHash = postStateHash
	}
	println(util.EncodeToHexString(rootStateHash))
	println(bonds[0].String())

	queryResult10, errMessage := grpc.Query(client, rootStateHash, "address", systemContract, []string{}, protocolVersion)
	proxyHash := queryResult10.GetAccount().GetNamedKeys()[0].GetKey().GetHash().GetHash()
	println(util.EncodeToHexString(genesisAddress), ": ", proxyHash)
	println(errMessage)
	println("Genesis setting complete")

	println("Creating 100,000 accounts")
	numAcc := 100000
	recipientAddressArr := make([][32]byte, numAcc)
	//Create 100_000 accounts
	for idx := 0; idx < numAcc; idx++ {
		byteaddr := sha256.Sum256([]byte(strconv.Itoa(idx)))
		recipientAddressArr = append(recipientAddressArr, byteaddr)
	}

	println("First scenario: 1 deploy per 1 commit. Too slow that just do 1000 txs")
	var startTime, endTime time.Time

	for idx, unitaddr := range recipientAddressArr {
		if idx == 10 {
			startTime = time.Now()
		}
		if idx%100 == 0 {
			println(idx, " deploys creating...")
		}
		if idx == 1000 {
			break
		}
		timestamp := time.Now().Unix()
		deploys := util.MakeInitDeploys()

		deploy, _ := makeSendDeploy(genesisAddress, unitaddr[:], proxyHash, chainName, timestamp)
		deploys = util.AddDeploy(deploys, deploy)
		res, _ := grpc.Execute(client, rootStateHash, timestamp, deploys, protocolVersion)
		effects, err := executeErrorHandler(res)
		if err != nil {
			panic(err)
		}

		postStateHash, _, errMessage := grpc.Commit(client, rootStateHash, effects, protocolVersion)
		rootStateHash = postStateHash
		if errMessage != "" {
			panic(errMessage)
		}
	}
	endTime = time.Now()
	elapsedTime := endTime.Sub(startTime).Seconds()
	println("Total time: ", elapsedTime)
	println("TPS: ", 990/elapsedTime)

	println("Second scenario: Gather deploys and put them into 1 commit.")
	deploys := util.MakeInitDeploys()
	startTime = time.Now()

	for idx, unitaddr := range recipientAddressArr {
		if idx == 1000 {
			break
		}

		if idx%100 == 0 {
			println(idx, " deploys creating...")
		}
		timestamp := time.Now().Unix()
		deploy, _ := makeSendDeploy(genesisAddress, unitaddr[:], proxyHash, chainName, timestamp)
		deploys = util.AddDeploy(deploys, deploy)
	}

	timestamp := time.Now().Unix()
	println("Execute")
	res, _ := grpc.Execute(client, rootStateHash, timestamp, deploys, protocolVersion)
	effects, err = executeErrorHandler(res)
	if err != nil {
		panic(err)
	}

	println("Commit")
	postStateHash, _, errMessage = grpc.Commit(client, rootStateHash, effects, protocolVersion)
	rootStateHash = postStateHash
	if errMessage != "" {
		panic(errMessage)
	}
	endTime = time.Now()
	elapsedTime = endTime.Sub(startTime).Seconds()
	println("Total time: ", elapsedTime)
	println("TPS: ", 1000/elapsedTime)
}

func makeSendDeploy(senderAddr, recipientAddr, proxyHash []byte, chainName string, timestamp int64) (*ipc.DeployItem, error) {
	sessionArgs := []*consensus.Deploy_Arg{
		&consensus.Deploy_Arg{
			Name: "method",
			Value: &consensus.Deploy_Arg_Value{
				Value: &consensus.Deploy_Arg_Value_StringValue{
					StringValue: "transfer_to_account"}}},
		&consensus.Deploy_Arg{
			Name: "address",
			Value: &consensus.Deploy_Arg_Value{
				Value: &consensus.Deploy_Arg_Value_BytesValue{
					BytesValue: recipientAddr}}},
		&consensus.Deploy_Arg{
			Name: "amount",
			Value: &consensus.Deploy_Arg_Value{
				Value: &consensus.Deploy_Arg_Value_LongValue{
					LongValue: int64(10)}}},
	}

	paymentArgs := []*consensus.Deploy_Arg{
		&consensus.Deploy_Arg{
			Name: "method",
			Value: &consensus.Deploy_Arg_Value{
				Value: &consensus.Deploy_Arg_Value_StringValue{
					StringValue: "standard_payment",
				},
			},
		},
		&consensus.Deploy_Arg{
			Name: "fee",
			Value: &consensus.Deploy_Arg_Value{
				Value: &consensus.Deploy_Arg_Value_IntValue{
					IntValue: 100000000,
				},
			},
		},
	}

	deploy, err := util.MakeDeploy(senderAddr, util.HASH, proxyHash, sessionArgs, util.HASH, proxyHash, paymentArgs, uint64(10), timestamp, chainName)
	return deploy, err
}
