package grpc

import (
	"context"
	"time"

	cons "github.com/hdac-io/casperlabs-ee-grpc-go-util/protobuf/io/casperlabs/casper/consensus"
	state "github.com/hdac-io/casperlabs-ee-grpc-go-util/protobuf/io/casperlabs/casper/consensus/state"
	ipc "github.com/hdac-io/casperlabs-ee-grpc-go-util/protobuf/io/casperlabs/ipc"

	"github.com/hdac-io/casperlabs-ee-grpc-go-util/util"

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
	protocolVersion *state.ProtocolVersion) (parentStateHash []byte, effects *ipc.ExecutionEffect) {
	initialMotes := &state.BigInt{Value: strInitialMotes, BitWidth: uint32(512)}

	deployMintCode := &ipc.DeployCode{Code: mintCode}

	deployPosCode := &ipc.DeployCode{Code: posCode}

	genesisValidators := []*ipc.Bond{}

	for address, stake := range validators {
		genesisValidators = append(genesisValidators,
			&ipc.Bond{ValidatorPublicKey: util.DecodeHexString(address), Stake: &state.BigInt{Value: stake, BitWidth: 512}})
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
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
			ProtocolVersion:   protocolVersion})
	if err != nil {
		panic(err)
	}

	genesisResult := r.GetSuccess()

	return genesisResult.PoststateHash, genesisResult.GetEffect()
}

func RunGenensisWithChainSpec(client ipc.ExecutionEngineServiceClient,
	name string,
	timestamp int64,
	protocolVersion *state.ProtocolVersion,
	mintInstallCode []byte,
	posInstallCode []byte,
	mapAccounts map[string][]string,
	mapCosts map[string]uint32) (parentStateHash []byte, effects *ipc.ExecutionEffect, errMessage string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	accounts := []*ipc.ChainSpec_GenesisAccount{}

	for address, strAccount := range mapAccounts {
		accounts = append(accounts, &ipc.ChainSpec_GenesisAccount{PublicKey: util.DecodeHexString(address), Balance: &state.BigInt{Value: strAccount[0], BitWidth: 512}, BondedAmount: &state.BigInt{Value: strAccount[1], BitWidth: 512}})
	}

	costs := &ipc.ChainSpec_CostTable{
		Wasm: &ipc.ChainSpec_CostTable_WasmCosts{
			Regular:        mapCosts["regular"],
			Div:            mapCosts["div-multiplier"],
			Mul:            mapCosts["mul-multiplier"],
			Mem:            mapCosts["mem-multiplier"],
			InitialMem:     mapCosts["mem-initial-pages"],
			GrowMem:        mapCosts["mem-grow-per-page"],
			Memcpy:         mapCosts["mem-copy-per-byte"],
			MaxStackHeight: mapCosts["max-stack-height"],
			OpcodesMul:     mapCosts["opcodes-multiplier"],
			OpcodesDiv:     mapCosts["opcodes-divisor"]}}

	r, err := client.RunGenesisWithChainspec(
		ctx,
		&ipc.ChainSpec_GenesisConfig{
			Name:            name,
			Timestamp:       uint64(timestamp),
			ProtocolVersion: protocolVersion,
			MintInstaller:   mintInstallCode,
			PosInstaller:    posInstallCode,
			Accounts:        accounts,
			Costs:           costs})

	if err != nil {
		errMessage = err.Error()
	}

	switch r.GetResult().(type) {
	case *ipc.GenesisResponse_Success:
		parentStateHash = r.GetSuccess().GetPoststateHash()
		effects = r.GetSuccess().GetEffect()
	case *ipc.GenesisResponse_FailedDeploy:
		errMessage += r.GetFailedDeploy().GetMessage()
	}

	return parentStateHash, effects, errMessage
}

func Commit(client ipc.ExecutionEngineServiceClient,
	prestateHash []byte,
	effects *ipc.ExecutionEffect,
	protocolVersion *state.ProtocolVersion) (postStateHash []byte, validators []*ipc.Bond, errMessage string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	r, err := client.Commit(
		ctx,
		&ipc.CommitRequest{
			PrestateHash:    prestateHash,
			Effects:         effects.GetTransformMap(),
			ProtocolVersion: protocolVersion})
	if err != nil {
		errMessage = err.Error()
	}

	switch r.GetResult().(type) {
	case *ipc.CommitResponse_Success:
		postStateHash = r.GetSuccess().GetPoststateHash()
		validators = r.GetSuccess().GetBondedValidators()
	case *ipc.CommitResponse_MissingPrestate:
		errMessage = "Missing prestate : " + util.EncodeToHexString(r.GetMissingPrestate().GetHash())
	case *ipc.CommitResponse_KeyNotFound:
		errMessage += "Key not Found "
		var hashValue []byte
		switch r.GetKeyNotFound().GetValue().(type) {
		case *state.Key_Address_:
			errMessage += "(Address)"
			hashValue = r.GetKeyNotFound().GetAddress().GetAccount()
		case *state.Key_Hash_:
			errMessage += "(Hash)"
			hashValue = r.GetKeyNotFound().GetHash().GetHash()
		case *state.Key_Uref:
			errMessage += "(Uref)"
			hashValue = r.GetKeyNotFound().GetUref().GetUref()
		case *state.Key_Local_:
			errMessage += "(Local)"
			hashValue = r.GetKeyNotFound().GetLocal().GetHash()
		}
		errMessage += " : " + util.EncodeToHexString(hashValue)
	case *ipc.CommitResponse_TypeMismatch:
		errMessage += "Type missmatch : expected (" + r.GetTypeMismatch().GetExpected() + "), but (" + r.GetTypeMismatch().GetFound() + ")"
	case *ipc.CommitResponse_FailedTransform:
		errMessage += "Failed transform : " + r.GetFailedTransform().GetMessage()
	}

	return postStateHash, validators, errMessage
}

func Validate(client ipc.ExecutionEngineServiceClient, wasmCode []byte, protocolVersion *state.ProtocolVersion) (result bool, errMessage string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := client.Validate(
		ctx,
		&ipc.ValidateRequest{
			WasmCode:        wasmCode,
			ProtocolVersion: protocolVersion})
	if err != nil {
		errMessage = err.Error()
	}

	switch r.GetResult().(type) {
	case *ipc.ValidateResponse_Success:
		result = true
	case *ipc.ValidateResponse_Failure:
		result = false
		errMessage = r.GetFailure()
	}

	return result, errMessage
}

func Query(client ipc.ExecutionEngineServiceClient,
	stateHash []byte,
	keyType string,
	keyData string,
	path []string,
	protocolVersion *state.ProtocolVersion) (result *state.Value, errMessage string) {

	var key *state.Key
	switch keyType {
	case "address":
		key = &state.Key{Value: &state.Key_Address_{Address: &state.Key_Address{Account: util.DecodeHexString(keyData)}}}
	case "local":
		key = &state.Key{Value: &state.Key_Local_{Local: &state.Key_Local{Hash: util.DecodeHexString(keyData)}}}
	case "uref":
		key = &state.Key{Value: &state.Key_Uref{Uref: &state.Key_URef{Uref: util.DecodeHexString(keyData)}}}
	case "hash":
		key = &state.Key{Value: &state.Key_Hash_{Hash: &state.Key_Hash{Hash: util.DecodeHexString(keyData)}}}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	r, err := client.Query(
		ctx,
		&ipc.QueryRequest{
			StateHash:       stateHash,
			BaseKey:         key,
			Path:            path,
			ProtocolVersion: protocolVersion})
	if err != nil {
		errMessage = err.Error()
	}

	switch r.GetResult().(type) {
	case *ipc.QueryResponse_Success:
		result = r.GetSuccess()
	case *ipc.QueryResponse_Failure:
		errMessage = r.GetFailure()
	}

	return result, errMessage
}

func Execute(client ipc.ExecutionEngineServiceClient,
	parentStateHash []byte,
	timestamp int64,
	gasPrice uint64,
	strGenensisAddress string,
	paymentWasmCode []byte,
	paymentArgs []byte,
	sessionWasmCode []byte,
	sessionArgs []byte,
	protocolVersion *state.ProtocolVersion) (effects *ipc.ExecutionEffect, errMessage string) {

	u64Timestamp := uint64(timestamp)
	genensisAddress := util.DecodeHexString(strGenensisAddress)

	deployBody := &cons.Deploy_Body{
		Session: &cons.Deploy_Code{Contract: &cons.Deploy_Code_Wasm{Wasm: sessionWasmCode}, AbiArgs: sessionArgs},
		Payment: &cons.Deploy_Code{Contract: &cons.Deploy_Code_Wasm{Wasm: paymentWasmCode}, AbiArgs: paymentArgs}}

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
			Address:           genensisAddress,
			Session:           &ipc.DeployPayload{Payload: &ipc.DeployPayload_DeployCode{DeployCode: &ipc.DeployCode{Code: sessionWasmCode, Args: sessionArgs}}},
			Payment:           &ipc.DeployPayload{Payload: &ipc.DeployPayload_DeployCode{DeployCode: &ipc.DeployCode{Code: paymentWasmCode, Args: paymentArgs}}},
			GasPrice:          gasPrice,
			AuthorizationKeys: [][]byte{genensisAddress},
			DeployHash:        headerHash}}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	r, err := client.Execute(
		ctx,
		&ipc.ExecuteRequest{
			ParentStateHash: parentStateHash,
			BlockTime:       u64Timestamp,
			Deploys:         deploys,
			ProtocolVersion: protocolVersion})
	if err != nil {
		errMessage = err.Error()
	}

	switch r.GetResult().(type) {
	case *ipc.ExecuteResponse_Success:
		effects = r.GetSuccess().GetDeployResults()[0].GetExecutionResult().GetEffects()
	case *ipc.ExecuteResponse_MissingParent:
		errMessage = "Missing parentstate : " + util.EncodeToHexString(r.GetMissingParent().GetHash())
	}

	return effects, errMessage
}

func Upgrade(client ipc.ExecutionEngineServiceClient,
	parentStateHash []byte,
	wasmCode []byte,
	mapCosts map[string]uint32,
	currentProtocolVersion *state.ProtocolVersion,
	nextProtocolVersion *state.ProtocolVersion) (postStateHash []byte, effects *ipc.ExecutionEffect, errMessage string) {

	costs := &ipc.ChainSpec_CostTable{
		Wasm: &ipc.ChainSpec_CostTable_WasmCosts{
			Regular:        mapCosts["regular"],
			Div:            mapCosts["div-multiplier"],
			Mul:            mapCosts["mul-multiplier"],
			Mem:            mapCosts["mem-multiplier"],
			InitialMem:     mapCosts["mem-initial-pages"],
			GrowMem:        mapCosts["mem-grow-per-page"],
			Memcpy:         mapCosts["mem-copy-per-byte"],
			MaxStackHeight: mapCosts["max-stack-height"],
			OpcodesMul:     mapCosts["opcodes-multiplier"],
			OpcodesDiv:     mapCosts["opcodes-divisor"]}}

	upgradePoint := &ipc.ChainSpec_UpgradePoint{
		ActivationPoint:  &ipc.ChainSpec_ActivationPoint{Rank: uint64(1)},
		ProtocolVersion:  nextProtocolVersion,
		UpgradeInstaller: &ipc.DeployCode{Code: wasmCode},
		NewCosts:         costs}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	r, err := client.Upgrade(
		ctx,
		&ipc.UpgradeRequest{
			ParentStateHash: parentStateHash,
			UpgradePoint:    upgradePoint,
			ProtocolVersion: currentProtocolVersion})
	if err != nil {
		errMessage = err.Error()
	}

	switch r.GetResult().(type) {
	case *ipc.UpgradeResponse_Success:
		postStateHash = r.GetSuccess().GetPostStateHash()
		effects = r.GetSuccess().GetEffect()
	case *ipc.UpgradeResponse_FailedDeploy:
		errMessage = r.GetFailedDeploy().GetMessage()
	}

	return postStateHash, effects, errMessage
}

func QueryBlanace(client ipc.ExecutionEngineServiceClient,
	stateHash []byte,
	address string,
	protocolVersion *state.ProtocolVersion) (balance string, errMessage string) {

	res, errMessage := Query(client, stateHash, "address", address, []string{}, protocolVersion)
	if errMessage != "" {
		return balance, errMessage
	}

	purseID := res.GetAccount().GetPurseId().GetUref()
	namedKeys := res.GetAccount().GetNamedKeys()
	var mintUref []byte
	for _, value := range namedKeys {
		if value.GetName() == "mint" {
			mintUref = value.GetKey().GetUref().GetUref()
			break
		}
	}

	localSrc := util.EncodeToHexString(mintUref) + util.EncodeToHexString(util.AbiBytesToBytes(purseID))
	localBytes := util.Blake2b256(util.DecodeHexString(localSrc))

	res, errMessage = Query(client, stateHash, "local", util.EncodeToHexString(localBytes), []string{}, protocolVersion)
	if errMessage != "" {
		return balance, errMessage
	}

	uref := res.GetKey().GetUref().GetUref()
	res, errMessage = Query(client, stateHash, "uref", util.EncodeToHexString(uref), []string{}, protocolVersion)
	if errMessage != "" {
		return balance, errMessage
	}

	return res.GetBigInt().GetValue(), errMessage
}
