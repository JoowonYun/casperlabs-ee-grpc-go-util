package main

import (
	"bytes"
	"math/big"
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
	mintCode, posCode, cntDefCode, cntCallCode, mintInstallCode, posInstallCode, transferToAccountCode, standardPaymentCode := loadWasmCode()

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
	effects2, errMessage := grpc.Execute(client, rootStateHash, time.Now().Unix(), uint64(10), genesisAddress, cntDefCode, []byte{}, cntDefCode, []byte{}, protocolVersion)

	postStateHash2, bonds2, errMessage := grpc.Commit(client, rootStateHash, effects2, protocolVersion)
	rootStateHash = postStateHash2
	println(util.EncodeToHexString(postStateHash2))
	println(bonds2[0].String())

	// Run "Counter Call contract"
	effects3, errMessage := grpc.Execute(client, rootStateHash, time.Now().Unix(), uint64(10), genesisAddress, cntCallCode, []byte{}, cntCallCode, []byte{}, protocolVersion)

	postStateHash3, bonds3, errMessage := grpc.Commit(client, rootStateHash, effects3, protocolVersion)
	rootStateHash = postStateHash3
	println(util.EncodeToHexString(postStateHash3))
	println(bonds3[0].String())

	// Query counter contract.
	path := []string{"counter", "count"}
	queryResult1, errMessage := grpc.Query(client, rootStateHash, genesisAddress, path, protocolVersion)
	println(queryResult1.GetIntValue())
	println(errMessage)

	effects4, errMessage := grpc.Execute(client, rootStateHash, time.Now().Unix(), uint64(10), genesisAddress, cntCallCode, []byte{}, cntCallCode, []byte{}, protocolVersion)

	postStateHash4, bonds4, errMessage := grpc.Commit(client, rootStateHash, effects4, protocolVersion)
	rootStateHash = postStateHash4
	println(util.EncodeToHexString(postStateHash4))
	println(bonds4[0].String())

	queryResult2, errMessage := grpc.Query(client, rootStateHash, genesisAddress, path, protocolVersion)
	println(queryResult2.GetIntValue())
	println(errMessage)

	// Run "Send transaction"
	sessionAbi := util.MakeArgsTransferToAccount("93236a9263d2ac6198c5ed211774c745d5dc62a910cb84276f8a7c4959208915", uint64(10))
	paymentAbi := util.MakeArgsStandardPayment(new(big.Int).SetUint64(1))
	effects5, errMessage := grpc.Execute(client, rootStateHash, time.Now().Unix(), uint64(10), genesisAddress, standardPaymentCode, paymentAbi, transferToAccountCode, sessionAbi, protocolVersion)

	postStateHash5, bonds5, errMessage := grpc.Commit(client, rootStateHash, effects5, protocolVersion)
	rootStateHash = postStateHash5
	println(util.EncodeToHexString(postStateHash5))
	println(bonds5[0].String())

	// Upgrade costs data..
	costs["regular"] = 2
	nextProtocolVersion := util.MakeProtocolVersion(2, 0, 0)
	postStateHash6, effects6, errMessage := grpc.Upgrade(client, parentStateHash, cntDefCode, costs, protocolVersion, nextProtocolVersion)
	postStateHash7, bonds6, errMessage := grpc.Commit(client, postStateHash6, effects6, nextProtocolVersion)
	if bytes.Equal(postStateHash6, postStateHash7) {
		rootStateHash = postStateHash6
		protocolVersion = nextProtocolVersion
	}
	println(util.EncodeToHexString(rootStateHash))
	println(bonds6[0].String())
}

func loadWasmCode() (mintCode []byte, posCode []byte, cntDefCode []byte, cntCallCode []byte, mintInstallCode []byte, posInstallCode []byte, transferToAccountCode []byte, standardPaymentCode []byte) {
	mintCode = util.LoadWasmFile("./example/contracts/mint_token.wasm")

	posCode = util.LoadWasmFile("./example/contracts/pos.wasm")

	cntDefCode = util.LoadWasmFile("./example/contracts/counterdefine.wasm")

	cntCallCode = util.LoadWasmFile("./example/contracts/countercall.wasm")

	mintInstallCode = util.LoadWasmFile("./example/contracts/mint_install.wasm")

	posInstallCode = util.LoadWasmFile("./example/contracts/pos_install.wasm")

	transferToAccountCode = util.LoadWasmFile("./example/contracts/transfer_to_account.wasm")

	standardPaymentCode = util.LoadWasmFile("./example/contracts/standard_payment.wasm")

	return mintCode, posCode, cntDefCode, cntCallCode, mintInstallCode, posInstallCode, transferToAccountCode, standardPaymentCode
}
