package grpc

import (
	"context"
	"time"

	cons "github.com/joowonyun/casperlabs-ee-grpc-go-util/protobuf/io/casperlabs/casper/consensus"
	state "github.com/joowonyun/casperlabs-ee-grpc-go-util/protobuf/io/casperlabs/casper/consensus/state"
	ipc "github.com/joowonyun/casperlabs-ee-grpc-go-util/protobuf/io/casperlabs/ipc"

	"github.com/joowonyun/casperlabs-ee-grpc-go-util/util"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
)

func Connect(path string) ipc.ExecutionEngineServiceClient {
	path = `unix:////` + path

	conn, e := grpc.Dial(path, grpc.WithInsecure())
	if e != nil {
		panic(e)
	}

	client := ipc.NewExecutionEngineServiceClient(conn)

	return client
}

func RunGenensis(
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
			&ipc.Bond{ValidatorPublicKey: util.DecodeHexString(address), Stake: &state.BigInt{Value: stake, BitWidth: 512}})
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := client.RunGenesis(
		ctx,
		&ipc.GenesisRequest{
			Address:           util.DecodeHexString(genesisAddress),
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

func Commit(client ipc.ExecutionEngineServiceClient,
	prestateHash []byte,
	effects *ipc.ExecutionEffect) (postStateHash []byte, validators []*ipc.Bond) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

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

func Validate(client ipc.ExecutionEngineServiceClient, wasmCode []byte) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

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

func Query(client ipc.ExecutionEngineServiceClient,
	stateHash []byte,
	genensisAddress string,
	path []string) (bool, *state.Value) {
	key := &state.Key{Value: &state.Key_Address_{Address: &state.Key_Address{Account: util.DecodeHexString(genensisAddress)}}}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

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

func Execute(client ipc.ExecutionEngineServiceClient,
	parentStateHash []byte,
	timestamp int64,
	gasPrice uint64,
	strGenensisAddress string,
	paymentWasmCode []byte,
	sessionWasmCode []byte,
	motesTransferredInPayment uint64,
	protocolVersion int) (effects *ipc.ExecutionEffect) {

	u64Timestamp := uint64(timestamp)
	genensisAddress := util.DecodeHexString(strGenensisAddress)

	deployBody := &cons.Deploy_Body{
		Session: &cons.Deploy_Code{Contract: &cons.Deploy_Code_Wasm{Wasm: sessionWasmCode}},
		Payment: &cons.Deploy_Code{Contract: &cons.Deploy_Code_Wasm{Wasm: paymentWasmCode}}}

	marshalDeployBody, err := proto.Marshal(deployBody)
	bodyHash := util.Blake2b256(marshalDeployBody)

	deployHeader := &cons.Deploy_Header{
		AccountPublicKey: genensisAddress,
		Timestamp:        u64Timestamp,
		GasPrice:         gasPrice,
		BodyHash:         bodyHash}

	marshalDeployHeader, err := proto.Marshal(deployHeader)
	headerHash := util.Blake2b256(marshalDeployHeader)

	deploys := []*ipc.DeployItem{
		&ipc.DeployItem{
			Address:                   genensisAddress,
			Session:                   &ipc.DeployPayload{Payload: &ipc.DeployPayload_DeployCode{DeployCode: &ipc.DeployCode{Code: sessionWasmCode}}},
			Payment:                   &ipc.DeployPayload{Payload: &ipc.DeployPayload_DeployCode{DeployCode: &ipc.DeployCode{Code: paymentWasmCode}}},
			MotesTransferredInPayment: motesTransferredInPayment,
			GasPrice:                  gasPrice,
			AuthorizationKeys:         [][]byte{genensisAddress},
			DeployHash:                headerHash}}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

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
