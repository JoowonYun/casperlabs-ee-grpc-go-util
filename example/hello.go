package main

import (
	"bytes"
	"os"
	"time"

	"github.com/hdac-io/casperlabs-ee-grpc-go-util/util"

	"github.com/hdac-io/casperlabs-ee-grpc-go-util/grpc"
)

func main() {
	// Init variable
	emptyStateHash := util.DecodeHexString(util.StrEmptyStateHash)
	rootStateHash := emptyStateHash
	genesisAddress := "d70243dd9d0d646fd6df282a8f7a8fa05a6629bec01d8024c3611eb1c1fb9f84"
	/* For run_genesis
	validateAddress1 := "93236a9263d2ac6198c5ed211774c745d5dc62a910cb84276f8a7c4959208915"
	validates := map[string]string{
		validateAddress1: "100",}
	*/
	networkName := "hdac"
	accounts := map[string][]string{
		genesisAddress: []string{"50000000", "1000000"}}

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
	mintCode, posCode, cntDefCode, cntCallCode, mintInstallCode, posInstallCode := loadWasmCode()

	// validate all wasm code
	mintResult, errMessage := grpc.Validate(client, mintCode, protocolVersion)
	println(mintResult)
	println(errMessage)
	posResult, errMessage := grpc.Validate(client, posCode, protocolVersion)
	println(posResult)
	println(errMessage)
	cntDefResult, errMessage := grpc.Validate(client, cntDefCode, protocolVersion)
	println(cntDefResult)
	println(errMessage)
	cntCallResult, errMessage := grpc.Validate(client, cntCallCode, protocolVersion)
	println(cntCallResult)
	println(errMessage)
	mintInstallResult, errMessage := grpc.Validate(client, mintInstallCode, protocolVersion)
	println(mintInstallResult)
	println(errMessage)
	posInstallResult, errMessage := grpc.Validate(client, posInstallCode, protocolVersion)
	println(posInstallResult)
	println(errMessage)

	// Run genesis and commit
	/* Legacy RunGenensis
	parentStateHash, effects, errMesage := grpc.RunGenensis(client,
		genesisAddress,
		"100",
		0,
		mintCode,
		posCode,
		validates,
		protocolVersion)
	*/

	parentStateHash, effects, errMessage := grpc.RunGenensisWithChainSpec(client,
		networkName,
		0,
		protocolVersion,
		mintInstallCode,
		posInstallCode,
		accounts,
		costs)

	postStateHash, bonds, errMessage := grpc.Commit(client, rootStateHash, effects, protocolVersion)
	if bytes.Equal(postStateHash, parentStateHash) {
		rootStateHash = postStateHash
	}
	println(util.EncodeToHexString(rootStateHash))
	println(bonds[0].String())

	// Run "Counter Define contract"
	effects2, errMessage := grpc.Execute(client, rootStateHash, time.Now().Unix(), uint64(10), genesisAddress, cntDefCode, cntDefCode, protocolVersion)

	postStateHash2, bonds2, errMessage := grpc.Commit(client, rootStateHash, effects2, protocolVersion)
	rootStateHash = postStateHash2
	println(util.EncodeToHexString(postStateHash2))
	println(bonds2[0].String())

	// Run "Counter Call contract"
	effects3, errMessage := grpc.Execute(client, rootStateHash, time.Now().Unix(), uint64(10), genesisAddress, cntCallCode, cntCallCode, protocolVersion)

	postStateHash3, bonds3, errMessage := grpc.Commit(client, rootStateHash, effects3, protocolVersion)
	rootStateHash = postStateHash3
	println(util.EncodeToHexString(postStateHash3))
	println(bonds3[0].String())

	// Query counter contract.
	path := []string{"counter", "count"}
	queryResult1, errMessage := grpc.Query(client, rootStateHash, genesisAddress, path, protocolVersion)
	println(queryResult1.GetIntValue())
	println(errMessage)

	effects4, errMessage := grpc.Execute(client, rootStateHash, time.Now().Unix(), uint64(10), genesisAddress, cntCallCode, cntCallCode, protocolVersion)

	postStateHash4, bonds4, errMessage := grpc.Commit(client, rootStateHash, effects4, protocolVersion)
	rootStateHash = postStateHash4
	println(util.EncodeToHexString(postStateHash4))
	println(bonds4[0].String())

	queryResult2, errMessage := grpc.Query(client, rootStateHash, genesisAddress, path, protocolVersion)
	println(queryResult2.GetIntValue())
	println(errMessage)

	// Upgrade costs data..
	costs["regular"] = 2
	nextProtocolVersion := util.MakeProtocolVersion(2, 0, 0)
	postStateHash5, effects5, errMessage := grpc.Upgrade(client, parentStateHash, cntDefCode, costs, protocolVersion, nextProtocolVersion)
	postStateHash6, bonds6, errMessage := grpc.Commit(client, postStateHash5, effects5, nextProtocolVersion)
	if bytes.Equal(postStateHash5, postStateHash6) {
		rootStateHash = postStateHash5
		protocolVersion = nextProtocolVersion
	}
	println(util.EncodeToHexString(rootStateHash))
	println(bonds6[0].String())
}

func loadWasmCode() (mintCode []byte, posCode []byte, cntDefCode []byte, cntCallCode []byte, mintInstallCode []byte, posInstallCode []byte) {
	mintCode = util.LoadWasmFile("./example/contracts/mint_token.wasm")

	posCode = util.LoadWasmFile("./example/contracts/pos.wasm")

	cntDefCode = util.LoadWasmFile("./example/contracts/counterdefine.wasm")

	cntCallCode = util.LoadWasmFile("./example/contracts/countercall.wasm")

	mintInstallCode = util.LoadWasmFile("./example/contracts/mint_install.wasm")

	posInstallCode = util.LoadWasmFile("./example/contracts/pos_install.wasm")

	return mintCode, posCode, cntDefCode, cntCallCode, mintInstallCode, posInstallCode
}
