package main

import (
	"context"
	cons "io/casperlabs/casper/consensus"
	state "io/casperlabs/casper/consensus/state"
	ipc "io/casperlabs/ipc"
	"time"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
)

func connect(path string) ipc.ExecutionEngineServiceClient {
	path = `unix:////` + path

	conn, e := grpc.Dial(path, grpc.WithInsecure())
	if e != nil {
		panic(e)
	}

	client := ipc.NewExecutionEngineServiceClient(conn)

	return client
}

func runGenensis(
	client ipc.ExecutionEngineServiceClient,
	genesisAddress string,
	strInitialMotes string,
	timestamp int64,
	mintCode []byte,
	posCode []byte,
	validators map[string]string,
	protocolVersion int) (parentStateHash []byte, effects *ipc.ExecutionEffect) {
	initialMotes := &state.BigInt{Value: strInitialMotes, BitWidth: uint32(512)}

	deployMintCode := &ipc.DeployCode{Code: mintCode}

	deployPosCode := &ipc.DeployCode{Code: posCode}

	genesisValidators := []*ipc.Bond{}

	for address, stake := range validators {
		genesisValidators = append(genesisValidators,
			&ipc.Bond{ValidatorPublicKey: decodeHexString(address), Stake: &state.BigInt{Value: stake, BitWidth: 512}})
	}

	ctx, cancle := context.WithTimeout(context.Background(), time.Second)
	defer cancle()

	r, err := client.RunGenesis(
		ctx,
		&ipc.GenesisRequest{
			Address:           decodeHexString(genesisAddress),
			InitialMotes:      initialMotes,
			Timestamp:         uint64(timestamp),
			MintCode:          deployMintCode,
			ProofOfStakeCode:  deployPosCode,
			GenesisValidators: genesisValidators,
			ProtocolVersion:   &state.ProtocolVersion{Value: uint64(protocolVersion)}})
	if err != nil {
		panic(err)
	}

	genesisResult := r.GetSuccess()

	return genesisResult.PoststateHash, genesisResult.GetEffect()
}

func commit(client ipc.ExecutionEngineServiceClient,
	prestateHash []byte,
	effects *ipc.ExecutionEffect) (postStateHash []byte, validators []*ipc.Bond) {
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

	return commitResult.GetPoststateHash(), commitResult.GetBondedValidators()
}

func validate(client ipc.ExecutionEngineServiceClient, wasmCode []byte) bool {
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

	return r.GetFailure() == ""
}

func query(client ipc.ExecutionEngineServiceClient,
	stateHash []byte,
	genensisAddress string,
	path []string) (bool, *state.Value) {
	key := &state.Key{Value: &state.Key_Address_{Address: &state.Key_Address{Account: decodeHexString(genensisAddress)}}}

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

	return r.GetFailure() == "", queryResult
}

func execute(client ipc.ExecutionEngineServiceClient,
	parentStateHash []byte,
	timestamp int64,
	gasPrice int,
	strGenensisAddress string,
	paymentWasmCode []byte,
	sessionWasmCode []byte,
	protocolVersion int) (effects *ipc.ExecutionEffect) {

	u64Timestamp := uint64(timestamp)
	u64GasPrice := uint64(gasPrice)
	genensisAddress := decodeHexString(strGenensisAddress)

	deployBody := &cons.Deploy_Body{
		Session: &cons.Deploy_Code{Contract: &cons.Deploy_Code_Wasm{Wasm: sessionWasmCode}},
		Payment: &cons.Deploy_Code{Contract: &cons.Deploy_Code_Wasm{Wasm: paymentWasmCode}}}

	marshalDeployBody, err := proto.Marshal(deployBody)
	bodyHash := blake2b256(marshalDeployBody)

	deployHeader := &cons.Deploy_Header{
		AccountPublicKey: genensisAddress,
		Timestamp:        u64Timestamp,
		GasPrice:         u64GasPrice,
		BodyHash:         bodyHash}

	marshalDeployHeader, err := proto.Marshal(deployHeader)
	headerHash := blake2b256(marshalDeployHeader)

	deploys := []*ipc.DeployItem{
		&ipc.DeployItem{
			Address:                   genensisAddress,
			Session:                   &ipc.DeployPayload{Payload: &ipc.DeployPayload_DeployCode{DeployCode: &ipc.DeployCode{Code: sessionWasmCode}}},
			Payment:                   &ipc.DeployPayload{Payload: &ipc.DeployPayload_DeployCode{DeployCode: &ipc.DeployCode{Code: paymentWasmCode}}},
			MotesTransferredInPayment: uint64(0),
			GasPrice:                  u64GasPrice,
			AuthorizationKeys:         [][]byte{genensisAddress},
			DeployHash:                headerHash}}

	ctx, cancle := context.WithTimeout(context.Background(), time.Second)
	defer cancle()

	r, err := client.Execute(
		ctx,
		&ipc.ExecuteRequest{
			ParentStateHash: parentStateHash,
			BlockTime:       u64Timestamp,
			Deploys:         deploys,
			ProtocolVersion: &state.ProtocolVersion{Value: uint64(protocolVersion)}})
	if err != nil {
		panic(err)
	}

	executeResult := r.GetSuccess()

	return executeResult.GetDeployResults()[0].GetExecutionResult().Effects
}
