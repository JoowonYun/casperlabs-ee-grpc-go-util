package main

import (
	"bytes"
	"context"
	"encoding/hex"
	cons "io/casperlabs/casper/consensus"
	state "io/casperlabs/casper/consensus/state"
	ipc "io/casperlabs/ipc"
	"io/ioutil"
	"time"

	"github.com/golang/protobuf/proto"

	"golang.org/x/crypto/blake2b"
	"google.golang.org/grpc"
)

func main() {
	conn, e := grpc.Dial("unix:////Users/yun/.casperlabs/.casper-node.sock", grpc.WithInsecure())
	if e != nil {
		panic(e)
	}
	defer conn.Close()

	client := ipc.NewExecutionEngineServiceClient(conn)
	mintCode, posCode, cntDefCode, cntCallCode := loadWasmCode()

	emptyStateHash, err := hex.DecodeString("3307a54ca6d5bfbafc0ef1b003f3ec4941c011ee7f79889e44416754de2f091d")

	rootStateHash := emptyStateHash

	genesisAddress, err := hex.DecodeString("d70243dd9d0d646fd6df282a8f7a8fa05a6629bec01d8024c3611eb1c1fb9f84")
	if err != nil {
		panic(err)
	}

	parentStateHash, effects := runGenensis(client, genesisAddress, mintCode, posCode)

	postStateHash, bonds := commit(client, rootStateHash, effects)
	if bytes.Equal(postStateHash, parentStateHash) {
		rootStateHash = postStateHash
	}

	println(bonds[0].String())

	validate(client, mintCode)
	validate(client, posCode)
	validate(client, cntDefCode)
	validate(client, cntCallCode)

	effects2 := execute(client, genesisAddress, rootStateHash, cntDefCode)

	postStateHash2, bonds2 := commit(client, rootStateHash, effects2)
	println()
	rootStateHash = postStateHash2
	println(hex.EncodeToString(postStateHash2))
	println(bonds2[0].String())

	effects3 := execute(client, genesisAddress, rootStateHash, cntCallCode)

	postStateHash3, bonds3 := commit(client, rootStateHash, effects3)
	println()
	println(hex.EncodeToString(rootStateHash))
	rootStateHash = postStateHash3
	println(hex.EncodeToString(rootStateHash))
	println(hex.EncodeToString(postStateHash3))
	println(bonds3[0].String())

	str := []string{"counter", "count"}
	query(client, rootStateHash, genesisAddress, str)

	effects4 := execute(client, genesisAddress, rootStateHash, cntCallCode)

	postStateHash4, bonds4 := commit(client, rootStateHash, effects4)
	rootStateHash = postStateHash4
	println(hex.EncodeToString(postStateHash4))
	println(bonds4[0].String())

	query(client, rootStateHash, genesisAddress, str)
}

func loadWasmCode() (mintCode []byte, posCode []byte, cntDefCode []byte, cntCallCode []byte) {
	mintCode, err1 := ioutil.ReadFile("./contracts/mint_token.wasm")
	if err1 != nil {
		panic(err1)
	}

	posCode, err2 := ioutil.ReadFile("./contracts/pos.wasm")
	if err2 != nil {
		panic(err2)
	}

	cntDefCode, err3 := ioutil.ReadFile("./contracts/counterdefine.wasm")
	if err3 != nil {
		panic(err3)
	}

	cntCallCode, err4 := ioutil.ReadFile("./contracts/countercall.wasm")
	if err4 != nil {
		panic(err4)
	}

	return mintCode, posCode, cntDefCode, cntCallCode
}

func validate(client ipc.ExecutionEngineServiceClient, wasmCode []byte) {
	ctx, cancle := context.WithTimeout(context.Background(), time.Second)
	defer cancle()

	r, err := client.Validate(
		ctx,
		&ipc.ValidateRequest{
			SessionCode: wasmCode,
			PaymentCode: wasmCode})
	if err != nil {
		panic(err)
	}

	print(r.GetSuccess())
}

func runGenensis(client ipc.ExecutionEngineServiceClient, genesisAddress []byte, mintCode []byte, posCode []byte) (parentStateHash []byte, effects *ipc.ExecutionEffect) {
	initialMotes := &state.BigInt{Value: "123", BitWidth: uint32(512)}

	timestamp := uint64(0)

	deployMintCode := &ipc.DeployCode{Code: mintCode}

	deployPosCode := &ipc.DeployCode{Code: posCode}

	validateAddress, err := hex.DecodeString("93236a9263d2ac6198c5ed211774c745d5dc62a910cb84276f8a7c4959208915")
	if err != nil {
		panic(err)
	}

	genesisValidators := []*ipc.Bond{&ipc.Bond{ValidatorPublicKey: validateAddress, Stake: &state.BigInt{Value: "100", BitWidth: 512}}}

	protocolVersion := &state.ProtocolVersion{Value: uint64(1)}

	ctx, cancle := context.WithTimeout(context.Background(), time.Second)
	defer cancle()

	r, err := client.RunGenesis(
		ctx,
		&ipc.GenesisRequest{
			Address:           genesisAddress,
			InitialMotes:      initialMotes,
			Timestamp:         timestamp,
			MintCode:          deployMintCode,
			ProofOfStakeCode:  deployPosCode,
			GenesisValidators: genesisValidators,
			ProtocolVersion:   protocolVersion})
	if err != nil {
		panic(err)
	}

	genesisResult := r.GetSuccess()

	println(`PostStaeHash : ` + hex.EncodeToString(genesisResult.PoststateHash))
	for index, value := range genesisResult.GetEffect().GetOpMap() {
		println(`Key [` + string(index) + `]: ` + value.GetKey().String())
		println(`Operlation ` + string(index) + ` : ` + value.GetOperation().String())
	}

	return genesisResult.PoststateHash, genesisResult.GetEffect()
}

func commit(client ipc.ExecutionEngineServiceClient, prestateHash []byte, effects *ipc.ExecutionEffect) (postStateHash []byte, validators []*ipc.Bond) {
	ctx, cancle := context.WithTimeout(context.Background(), time.Second)
	defer cancle()

	r, err := client.Commit(
		ctx,
		&ipc.CommitRequest{
			PrestateHash: prestateHash,
			Effects:      effects.GetTransformMap()})
	if err != nil {
		panic(err)
	}

	commitResult := r.GetSuccess()

	// println(`Poststatehash : ` + hex.EncodeToString(commitResult.PoststateHash))
	// for index, value := range commitResult.BondedValidators {
	// 	println(`Bond [` + string(index) + `]`)
	// 	println(`Validator public key : ` + hex.EncodeToString(value.GetValidatorPublicKey()))
	// 	println(`Stake : ` + value.Stake.GetValue() + ` / ` + string(value.GetStake().BitWidth))
	// }

	return commitResult.GetPoststateHash(), commitResult.GetBondedValidators()
}

func execute(client ipc.ExecutionEngineServiceClient, genensisAddress []byte, parentStateHash []byte, wasmCode []byte) (effects *ipc.ExecutionEffect) {
	timestamp := uint64(time.Now().Nanosecond())
	gasPrice := uint64(100)

	deployBody := &cons.Deploy_Body{
		Session: &cons.Deploy_Code{Contract: &cons.Deploy_Code_Wasm{Wasm: wasmCode}},
		Payment: &cons.Deploy_Code{Contract: &cons.Deploy_Code_Wasm{Wasm: wasmCode}}}

	hash1, err := blake2b.New256(nil)
	if err != nil {
		panic(err)
	}
	marshalDeployBody, err := proto.Marshal(deployBody)
	hash1.Write(marshalDeployBody)
	bodyHash := hash1.Sum(nil)

	deployHeader := &cons.Deploy_Header{
		AccountPublicKey: genensisAddress,
		Timestamp:        timestamp,
		GasPrice:         gasPrice,
		BodyHash:         bodyHash}

	marshalDeployHeader, err := proto.Marshal(deployHeader)
	hash2, err := blake2b.New256(nil)
	hash2.Write(marshalDeployHeader)
	headerHash := hash2.Sum(nil)

	deploys := []*ipc.DeployItem{
		&ipc.DeployItem{
			Address:                   genensisAddress,
			Session:                   &ipc.DeployPayload{Payload: &ipc.DeployPayload_DeployCode{DeployCode: &ipc.DeployCode{Code: wasmCode}}},
			Payment:                   &ipc.DeployPayload{Payload: &ipc.DeployPayload_DeployCode{DeployCode: &ipc.DeployCode{Code: wasmCode}}},
			MotesTransferredInPayment: uint64(0),
			GasPrice:                  gasPrice,
			AuthorizationKeys:         [][]byte{genensisAddress},
			DeployHash:                headerHash}}

	protocolVersion := &state.ProtocolVersion{Value: uint64(1)}

	ctx, cancle := context.WithTimeout(context.Background(), time.Second)
	defer cancle()

	r, err := client.Execute(
		ctx,
		&ipc.ExecuteRequest{
			ParentStateHash: parentStateHash,
			BlockTime:       timestamp,
			Deploys:         deploys,
			ProtocolVersion: protocolVersion})
	if err != nil {
		panic(err)
	}

	executeResult := r.GetSuccess()

	return executeResult.GetDeployResults()[0].GetExecutionResult().Effects
}

func query(client ipc.ExecutionEngineServiceClient, stateHash []byte, genensisAddress []byte, path []string) {
	key := &state.Key{Value: &state.Key_Address_{Address: &state.Key_Address{Account: genensisAddress}}}

	ctx, cancle := context.WithTimeout(context.Background(), time.Second)
	defer cancle()

	r, err := client.Query(
		ctx,
		&ipc.QueryRequest{
			StateHash: stateHash,
			BaseKey:   key,
			Path:      path})
	if err != nil {
		panic(err)
	}

	queryResult := r.GetSuccess()
	println(queryResult.String())

}
