package main

import (
	"bytes"
	"time"

	"github.com/hdac-io/casperlabs-ee-grpc-go-util/util"

	"github.com/hdac-io/casperlabs-ee-grpc-go-util/grpc"
)

func main() {
	// Connect to ee sock.
	client := grpc.Connect(`/Users/yun/.casperlabs/.casper-node.sock`)

	// laod wasm code
	mintCode, posCode, cntDefCode, cntCallCode, mintInstallCode, posInstallCode := loadWasmCode()

	// validate all wasm code
	mintResult := grpc.Validate(client, mintCode, 1)
	println(mintResult)
	posResult := grpc.Validate(client, posCode, 1)
	println(posResult)
	cntDefResult := grpc.Validate(client, cntDefCode, 1)
	println(cntDefResult)
	cntCallResult := grpc.Validate(client, cntCallCode, 1)
	println(cntCallResult)
	mintInstallResult := grpc.Validate(client, mintInstallCode, 1)
	println(mintInstallResult)
	posInstallResult := grpc.Validate(client, posInstallCode, 1)
	println(posInstallResult)

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

	// Run genesis and commit
	/* Legacy RunGenensis
	parentStateHash, effects := grpc.RunGenensis(client,
		genesisAddress,
		"100",
		0,
		mintCode,
		posCode,
		validates,
		1)
	*/

	parentStateHash, effects := grpc.RunGenensisWithChainSpec(client,
		networkName,
		0,
		1,
		mintInstallCode,
		posInstallCode,
		accounts,
		costs)

	postStateHash, bonds := grpc.Commit(client, rootStateHash, effects, 1)
	if bytes.Equal(postStateHash, parentStateHash) {
		rootStateHash = postStateHash
	}
	println(util.EncodeToHexString(rootStateHash))
	println(bonds[0].String())

	// Run "Counter Define contract"
	effects2 := grpc.Execute(client, rootStateHash, time.Now().Unix(), uint64(10), genesisAddress, cntDefCode, cntDefCode, 1)

	postStateHash2, bonds2 := grpc.Commit(client, rootStateHash, effects2, 1)
	rootStateHash = postStateHash2
	println(util.EncodeToHexString(postStateHash2))
	println(bonds2[0].String())

	// Run "Counter Call contract"
	effects3 := grpc.Execute(client, rootStateHash, time.Now().Unix(), uint64(10), genesisAddress, cntCallCode, cntCallCode, 1)

	postStateHash3, bonds3 := grpc.Commit(client, rootStateHash, effects3, 1)
	rootStateHash = postStateHash3
	println(util.EncodeToHexString(postStateHash3))
	println(bonds3[0].String())

	// Query counter contract.
	path := []string{"counter", "count"}
	queryResult1, queryData1 := grpc.Query(client, rootStateHash, genesisAddress, path, 1)
	println(queryResult1, queryData1.GetIntValue())

	effects4 := grpc.Execute(client, rootStateHash, time.Now().Unix(), uint64(10), genesisAddress, cntCallCode, cntCallCode, 1)

	postStateHash4, bonds4 := grpc.Commit(client, rootStateHash, effects4, 1)
	rootStateHash = postStateHash4
	println(util.EncodeToHexString(postStateHash4))
	println(bonds4[0].String())

	queryResult2, queryData2 := grpc.Query(client, rootStateHash, genesisAddress, path, 1)
	println(queryResult2, queryData2.GetIntValue())
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
