package main

import (
	"bytes"
	"encoding/hex"
	"time"
)

func main() {
	client := connect(`/Users/yun/.casperlabs/.casper-node.sock`)

	mintCode, posCode, cntDefCode, cntCallCode := loadWasmCode()

	mintResult := validate(client, mintCode)
	println(mintResult)
	posResult := validate(client, posCode)
	println(posResult)
	cntDefResult := validate(client, cntDefCode)
	println(cntDefResult)
	cntCallResult := validate(client, cntCallCode)
	println(cntCallResult)

	emptyStateHash := decodeHexString(strEmptyStateHash)

	rootStateHash := emptyStateHash

	genesisAddress := "d70243dd9d0d646fd6df282a8f7a8fa05a6629bec01d8024c3611eb1c1fb9f84"
	validateAddress1 := "93236a9263d2ac6198c5ed211774c745d5dc62a910cb84276f8a7c4959208915"
	validates := map[string]string{
		validateAddress1: "100",
	}

	parentStateHash, effects := runGenensis(client,
		genesisAddress,
		"100",
		0,
		mintCode,
		posCode,
		validates,
		1)

	postStateHash, bonds := commit(client, rootStateHash, effects)
	if bytes.Equal(postStateHash, parentStateHash) {
		rootStateHash = postStateHash
	}

	println(bonds[0].String())

	effects2 := execute(client, rootStateHash, time.Now().Unix(), 10, genesisAddress, cntDefCode, cntDefCode, 1)

	postStateHash2, bonds2 := commit(client, rootStateHash, effects2)
	rootStateHash = postStateHash2
	println(hex.EncodeToString(postStateHash2))
	println(bonds2[0].String())

	effects3 := execute(client, rootStateHash, time.Now().Unix(), 10, genesisAddress, cntCallCode, cntCallCode, 1)

	postStateHash3, bonds3 := commit(client, rootStateHash, effects3)
	rootStateHash = postStateHash3
	println(hex.EncodeToString(postStateHash3))
	println(bonds3[0].String())

	path := []string{"counter", "count"}
	queryResult1, queryData1 := query(client, rootStateHash, genesisAddress, path)
	println(queryResult1, queryData1.GetIntValue())

	effects4 := execute(client, rootStateHash, time.Now().Unix(), 10, genesisAddress, cntCallCode, cntCallCode, 1)

	postStateHash4, bonds4 := commit(client, rootStateHash, effects4)
	rootStateHash = postStateHash4
	println(hex.EncodeToString(postStateHash4))
	println(bonds4[0].String())

	queryResult2, queryData2 := query(client, rootStateHash, genesisAddress, path)
	println(queryResult2, queryData2.GetIntValue())
}

func loadWasmCode() (mintCode []byte, posCode []byte, cntDefCode []byte, cntCallCode []byte) {
	mintCode = loadWasmFile("./contracts/mint_token.wasm")

	posCode = loadWasmFile("./contracts/pos.wasm")

	cntDefCode = loadWasmFile("./contracts/counterdefine.wasm")

	cntCallCode = loadWasmFile("./contracts/countercall.wasm")

	return mintCode, posCode, cntDefCode, cntCallCode
}
