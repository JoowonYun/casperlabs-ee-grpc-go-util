package main

import (
	"bytes"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/hdac-io/casperlabs-ee-grpc-go-util/grpc"
	"github.com/hdac-io/casperlabs-ee-grpc-go-util/protobuf/io/casperlabs/ipc"
	"github.com/hdac-io/casperlabs-ee-grpc-go-util/protobuf/io/casperlabs/ipc/transforms"
	"github.com/hdac-io/casperlabs-ee-grpc-go-util/util"
)

func main() {
	// Init variable
	emptyStateHash := util.DecodeHexString(util.StrEmptyStateHash)
	rootStateHash := emptyStateHash
	genesisAddress := "d70243dd9d0d646fd6df282a8f7a8fa05a6629bec01d8024c3611eb1c1fb9f84"
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

	// laod wasm code
	mintCode, posCode, cntDefCode, cntCallCode, mintInstallCode, posInstallCode, transferToAccountCode, standardPaymentCode, bondingCode, unbondingCode := loadWasmCode()

	// validate all wasm code
	mintResult, errMessage := grpc.Validate(client, mintCode, protocolVersion)
	println(mintResult, len(mintCode))
	println(errMessage)
	posResult, errMessage := grpc.Validate(client, posCode, protocolVersion)
	println(posResult, len(posCode))
	println(errMessage)
	cntDefResult, errMessage := grpc.Validate(client, cntDefCode, protocolVersion)
	println(cntDefResult, len(cntDefCode))
	println(errMessage)
	cntCallResult, errMessage := grpc.Validate(client, cntCallCode, protocolVersion)
	println(cntCallResult, len(cntCallCode))
	println(errMessage)
	mintInstallResult, errMessage := grpc.Validate(client, mintInstallCode, protocolVersion)
	println(mintInstallResult, len(mintInstallCode))
	println(errMessage)
	posInstallResult, errMessage := grpc.Validate(client, posInstallCode, protocolVersion)
	println(posInstallResult, len(posInstallCode))
	println(errMessage)
	transferToAccountResult, errMessage := grpc.Validate(client, transferToAccountCode, protocolVersion)
	println(transferToAccountResult, len(transferToAccountCode))
	println(errMessage)
	standardPaymentResult, errMessage := grpc.Validate(client, standardPaymentCode, protocolVersion)
	println(standardPaymentResult, len(standardPaymentCode))
	println(errMessage)
	bondingResult, errMessage := grpc.Validate(client, bondingCode, protocolVersion)
	println(bondingResult, len(bondingCode))
	println(errMessage)
	unbondingResult, errMessage := grpc.Validate(client, unbondingCode, protocolVersion)
	println(unbondingResult, len(unbondingCode))
	println(errMessage)

	genesisConfig, err := util.GenesisConfigMock(
		chainName, genesisAddress, "500000000", "1000000", protocolVersion, costs,
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

	queryResult, errMessage := grpc.QueryBlanace(client, rootStateHash, genesisAddress, protocolVersion)
	println(genesisAddress, ": ", queryResult)
	println(errMessage)

	// Run "Counter Define contract"
	timestamp := time.Now().Unix()
	paymentAbi := util.MakeArgsStandardPayment(new(big.Int).SetUint64(200000000))
	deploy := util.MakeDeploy(genesisAddress, cntDefCode, []byte{}, standardPaymentCode, paymentAbi, uint64(10), timestamp, chainName)
	deploys := util.MakeInitDeploys()
	deploys = util.AddDeploy(deploys, deploy)

	effects2, errMessage := grpc.Execute(client, rootStateHash, timestamp, deploys, protocolVersion)

	postStateHash2, bonds2, errMessage := grpc.Commit(client, rootStateHash, effects2, protocolVersion)
	rootStateHash = postStateHash2
	println(util.EncodeToHexString(postStateHash2))
	println(bonds2[0].String())

	queryResult, errMessage = grpc.QueryBlanace(client, rootStateHash, genesisAddress, protocolVersion)
	println(genesisAddress, ": ", queryResult)
	println(errMessage)

	// Run "Counter Call contract"
	timestamp = time.Now().Unix()
	deploy = util.MakeDeploy(genesisAddress, cntCallCode, []byte{}, standardPaymentCode, paymentAbi, uint64(10), timestamp, chainName)
	deploys = util.MakeInitDeploys()
	deploys = util.AddDeploy(deploys, deploy)
	effects3, errMessage := grpc.Execute(client, rootStateHash, timestamp, deploys, protocolVersion)

	postStateHash3, bonds3, errMessage := grpc.Commit(client, rootStateHash, effects3, protocolVersion)
	rootStateHash = postStateHash3
	println(util.EncodeToHexString(postStateHash3))
	println(bonds3[0].String())

	// Query counter contract.
	path := []string{"counter", "count"}
	queryResult1, errMessage := grpc.Query(client, rootStateHash, "address", genesisAddress, path, protocolVersion)
	println(queryResult1.GetIntValue())
	println(errMessage)

	queryResult, errMessage = grpc.QueryBlanace(client, rootStateHash, genesisAddress, protocolVersion)
	println(genesisAddress, ": ", queryResult)
	println(errMessage)

	timestamp = time.Now().Unix()
	deploy = util.MakeDeploy(genesisAddress, cntCallCode, []byte{}, standardPaymentCode, paymentAbi, uint64(10), timestamp, chainName)
	deploys = util.MakeInitDeploys()
	deploys = util.AddDeploy(deploys, deploy)
	effects4, errMessage := grpc.Execute(client, rootStateHash, timestamp, deploys, protocolVersion)

	postStateHash4, bonds4, errMessage := grpc.Commit(client, rootStateHash, effects4, protocolVersion)
	rootStateHash = postStateHash4
	println(util.EncodeToHexString(postStateHash4))
	println(bonds4[0].String())

	queryResult2, errMessage := grpc.Query(client, rootStateHash, "address", genesisAddress, path, protocolVersion)
	println(queryResult2.GetIntValue())
	println(errMessage)

	queryResult3, errMessage := grpc.QueryBlanace(client, rootStateHash, genesisAddress, protocolVersion)
	println(genesisAddress, ": ", queryResult3)
	println(errMessage)

	// Run "Send transaction"
	timestamp = time.Now().Unix()
	sessionAbi := util.MakeArgsTransferToAccount("93236a9263d2ac6198c5ed211774c745d5dc62a910cb84276f8a7c4959208915", uint64(10))
	deploy = util.MakeDeploy(genesisAddress, transferToAccountCode, sessionAbi, standardPaymentCode, paymentAbi, uint64(10), timestamp, chainName)
	deploys = util.MakeInitDeploys()
	deploys = util.AddDeploy(deploys, deploy)
	effects5, errMessage := grpc.Execute(client, rootStateHash, timestamp, deploys, protocolVersion)

	postStateHash5, bonds5, errMessage := grpc.Commit(client, rootStateHash, effects5, protocolVersion)
	rootStateHash = postStateHash5
	println(util.EncodeToHexString(postStateHash5))
	println(bonds5[0].String())

	queryResult4, errMessage := grpc.QueryBlanace(client, rootStateHash, genesisAddress, protocolVersion)
	println(genesisAddress, ": ", queryResult4)
	println(errMessage)

	queryResult5, errMessage := grpc.QueryBlanace(client, rootStateHash, "93236a9263d2ac6198c5ed211774c745d5dc62a910cb84276f8a7c4959208915", protocolVersion)
	println("93236a9263d2ac6198c5ed211774c745d5dc62a910cb84276f8a7c4959208915: ", queryResult5)
	println(errMessage)

	// bonding
	timestamp = time.Now().Unix()
	bondingAbi := util.MakeArgsBonding(uint64(10))
	deploy = util.MakeDeploy(genesisAddress, bondingCode, bondingAbi, standardPaymentCode, paymentAbi, uint64(10), timestamp, chainName)
	deploys = util.MakeInitDeploys()
	deploys = util.AddDeploy(deploys, deploy)
	effects6, errMessage := grpc.Execute(client, rootStateHash, timestamp, deploys, protocolVersion)
	postStateHash6, bonds6, errMessage := grpc.Commit(client, rootStateHash, effects6, protocolVersion)
	rootStateHash = postStateHash6
	println(util.EncodeToHexString(rootStateHash))
	println(bonds6[0].String())

	// unbonding
	timestamp = time.Now().Unix()
	ubbondingAbi := util.MakeArgsUnBonding(uint64(100))
	deploy = util.MakeDeploy(genesisAddress, unbondingCode, ubbondingAbi, standardPaymentCode, paymentAbi, uint64(10), timestamp, chainName)
	deploys = util.MakeInitDeploys()
	deploys = util.AddDeploy(deploys, deploy)
	effects7, errMessage := grpc.Execute(client, rootStateHash, timestamp, deploys, protocolVersion)
	postStateHash7, bonds7, errMessage := grpc.Commit(client, rootStateHash, effects7, protocolVersion)
	rootStateHash = postStateHash7
	println(util.EncodeToHexString(rootStateHash))
	println(bonds7[0].String())

	// Upgrade costs data..
	costs["regular"] = 2
	nextProtocolVersion := util.MakeProtocolVersion(2, 0, 0)
	postStateHash8, effects8, errMessage := grpc.Upgrade(client, rootStateHash, cntDefCode, costs, protocolVersion, nextProtocolVersion)
	postStateHash9, bonds8, errMessage := grpc.Commit(client, rootStateHash, effects8, nextProtocolVersion)
	if bytes.Equal(postStateHash8, postStateHash9) {
		rootStateHash = postStateHash8
		protocolVersion = nextProtocolVersion
	}
	println(util.EncodeToHexString(rootStateHash))
	println(bonds8[0].String())
}

func loadWasmCode() (mintCode []byte, posCode []byte, cntDefCode []byte, cntCallCode []byte, mintInstallCode []byte, posInstallCode []byte, transferToAccountCode []byte, standardPaymentCode []byte, bondingCode []byte, unbondingCode []byte) {
	mintCode = util.LoadWasmFile("./example/contracts/mint_token.wasm")

	posCode = util.LoadWasmFile("./example/contracts/pos.wasm")

	cntDefCode = util.LoadWasmFile("./example/contracts/counter_define.wasm")

	cntCallCode = util.LoadWasmFile("./example/contracts/counter_call.wasm")

	mintInstallCode = util.LoadWasmFile("./example/contracts/mint_install.wasm")

	posInstallCode = util.LoadWasmFile("./example/contracts/pos_install.wasm")

	transferToAccountCode = util.LoadWasmFile("./example/contracts/transfer_to_account.wasm")

	standardPaymentCode = util.LoadWasmFile("./example/contracts/standard_payment.wasm")

	bondingCode = util.LoadWasmFile("./example/contracts/bonding.wasm")

	unbondingCode = util.LoadWasmFile("./example/contracts/unbonding.wasm")

	return mintCode, posCode, cntDefCode, cntCallCode, mintInstallCode, posInstallCode, transferToAccountCode, standardPaymentCode, bondingCode, unbondingCode
}
