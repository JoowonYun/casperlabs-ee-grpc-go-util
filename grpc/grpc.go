// Package grpc 는 Casperlabs의 Execution Engine의 GRPC Client 모듈을 정의한 모듈이다.
package grpc

import (
	"context"
	"fmt"

	"github.com/hdac-io/casperlabs-ee-grpc-go-util/protobuf/io/casperlabs/casper/consensus/state"
	"github.com/hdac-io/casperlabs-ee-grpc-go-util/protobuf/io/casperlabs/ipc"
	"github.com/hdac-io/casperlabs-ee-grpc-go-util/protobuf/io/casperlabs/ipc/transforms"
	"github.com/hdac-io/casperlabs-ee-grpc-go-util/storedvalue"

	"github.com/hdac-io/casperlabs-ee-grpc-go-util/util"

	"google.golang.org/grpc"
)

// Connect 은 Casperlabs의 Execution Engine의 unix socket으로 연결하는 함수.
func Connect(path string) ipc.ExecutionEngineServiceClient {
	path = `unix:////` + path

	conn, e := grpc.Dial(path, grpc.WithInsecure())
	if e != nil {
		panic(e)
	}

	client := ipc.NewExecutionEngineServiceClient(conn)

	return client
}

// RunGenesis 는 Execution Engine을 시작할 때 Genensis정보를 chain에 떄라 초기화하는 함수.
//
// ChainSpec_GenesisConfig 정보를 파라미터로 받아
// RunGenesis 후 결과를 return 받는다.
func RunGenesis(
	client ipc.ExecutionEngineServiceClient, genesisConfig *ipc.ChainSpec_GenesisConfig) (*ipc.GenesisResponse, error) {
	return client.RunGenesis(
		context.TODO(),
		genesisConfig)
}

// Commit 은 Execute한 effects를 적용시킬 때 사용하는 함수.
//
// State Hash, Execute한 effects를 파라미터로 받아,
// Commit 후 state hash 와 현재 Bonding 된 validator의 정보를 return 받는다.
func Commit(client ipc.ExecutionEngineServiceClient,
	prestateHash []byte,
	effects []*transforms.TransformEntry,
	protocolVersion *state.ProtocolVersion) (postStateHash []byte, validators []*ipc.Bond, errMessage string) {
	r, err := client.Commit(
		context.TODO(),
		&ipc.CommitRequest{
			PrestateHash:    prestateHash,
			Effects:         effects,
			ProtocolVersion: protocolVersion})
	if err != nil {
		errMessage = err.Error()
	}

	switch r.GetResult().(type) {
	case *ipc.CommitResponse_Success:
		postStateHash = r.GetSuccess().GetPoststateHash()
		validators = r.GetSuccess().GetBondedValidators()
	case *ipc.CommitResponse_MissingPrestate:
		errMessage = fmt.Sprintf("%s\nMissing prestate : %s", errMessage, util.EncodeToHexString(r.GetMissingPrestate().GetHash()))
	case *ipc.CommitResponse_KeyNotFound:
		errMessage = fmt.Sprintf("%s\nKey not Found ", errMessage)
		var hashValue []byte
		switch r.GetKeyNotFound().GetValue().(type) {
		case *state.Key_Address_:
			errMessage = fmt.Sprintf("%s\n(Address)", errMessage)
			hashValue = r.GetKeyNotFound().GetAddress().GetAccount()
		case *state.Key_Hash_:
			errMessage = fmt.Sprintf("%s\n(Hash)", errMessage)
			hashValue = r.GetKeyNotFound().GetHash().GetHash()
		case *state.Key_Uref:
			errMessage = fmt.Sprintf("%s\n(Uref)", errMessage)
			hashValue = r.GetKeyNotFound().GetUref().GetUref()
		case *state.Key_Local_:
			errMessage = fmt.Sprintf("%s\n(Local)", errMessage)
			hashValue = r.GetKeyNotFound().GetLocal().GetHash()
		}
		errMessage = fmt.Sprintf("%s : %s", errMessage, util.EncodeToHexString(hashValue))
	case *ipc.CommitResponse_TypeMismatch:
		errMessage = fmt.Sprintf("%s\nType missmatch : expected (%s), but (%s)", errMessage, r.GetTypeMismatch().GetExpected(), r.GetTypeMismatch().GetFound())
	case *ipc.CommitResponse_FailedTransform:
		errMessage = fmt.Sprintf("%s\nFailed transform : %s", errMessage, r.GetFailedTransform().GetMessage())
	}

	return postStateHash, validators, errMessage
}

// Query 는 특정 state 에서 해당 Key의 path에 대한 정보를 조회해주는 함수.
//
// State hash, Key type, Key Data, path를 파라미터로 받아
// Query 후 결과를 return 해준다.
func Query(client ipc.ExecutionEngineServiceClient,
	stateHash []byte,
	keyType string,
	keyData []byte,
	path []string,
	protocolVersion *state.ProtocolVersion) (result []byte, errMessage string) {

	var key *state.Key
	switch keyType {
	case "address":
		key = &state.Key{Value: &state.Key_Address_{Address: &state.Key_Address{Account: keyData}}}
	case "local":
		key = &state.Key{Value: &state.Key_Local_{Local: &state.Key_Local{Hash: keyData}}}
	case "uref":
		key = &state.Key{Value: &state.Key_Uref{Uref: &state.Key_URef{Uref: keyData}}}
	case "hash":
		key = &state.Key{Value: &state.Key_Hash_{Hash: &state.Key_Hash{Hash: keyData}}}
	}

	r, err := client.Query(
		context.TODO(),
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

// Execute 는 deploys를 실행할떄의 effects를 받아오는 함수.
//
// state hash, timestamp, deploys를 파라미터로 받아
// Execute 후 전체 response return 해준다.
func Execute(client ipc.ExecutionEngineServiceClient,
	parentStateHash []byte,
	int64timestamp int64,
	deploys []*ipc.DeployItem,
	protocolVersion *state.ProtocolVersion) (response *ipc.ExecuteResponse, err error) {

	timestamp := uint64(int64timestamp)

	return client.Execute(
		context.TODO(),
		&ipc.ExecuteRequest{
			ParentStateHash: parentStateHash,
			BlockTime:       timestamp,
			Deploys:         deploys,
			ProtocolVersion: protocolVersion})
}

// Upgrade 는 Wasm 코드나 Cost를 변경하여 Protocol Version을 Upgrade할 때 활용
//
// State hash, 변경할 Insatll Wasm코드, Cost, 현재 protocol version, 다음 protocol version을 파라미터로 받으며,
// Install wasm 코드를 변경할지, Cost를 변경할지는 옵션으로 가능하며 Upgrade 를 통해 변경한 후
// 변경될 state hash, effects를 return 해준다.
func Upgrade(client ipc.ExecutionEngineServiceClient,
	parentStateHash []byte,
	wasmCode []byte,
	mapCosts map[string]uint32,
	currentProtocolVersion *state.ProtocolVersion,
	nextProtocolVersion *state.ProtocolVersion) (postStateHash []byte, effects []*transforms.TransformEntry, errMessage string) {

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

	r, err := client.Upgrade(
		context.TODO(),
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
		effects = r.GetSuccess().GetEffect().GetTransformMap()
	case *ipc.UpgradeResponse_FailedDeploy:
		errMessage = r.GetFailedDeploy().GetMessage()
	}

	return postStateHash, effects, errMessage
}

// QueryBalance 는 address의 balance를 조회할 때 사용하는 함수.
//
// 조회할 state hash와 address를 파라미터로 받아, key를 address로 Query한다.
// name key에서 name이 mint인 uref를 추출하여 hex string로 변환하고 purse Id를 abi로 변환한 후 hex string으로 변환하여 붙인다.
// 해당 값을 blake2b256을 하면 local bytes 값이 추출된다. 이 값을 key를 local로 하여 Query한다.
// 받아온 uref값을 Key로 하여 Query하면 BigInt 형태의 blanace를 return 해준다.
func QueryBalance(client ipc.ExecutionEngineServiceClient,
	stateHash []byte,
	address []byte,
	protocolVersion *state.ProtocolVersion) (balance string, errMessage string) {

	res, errMessage := Query(client, stateHash, "address", address, []string{}, protocolVersion)
	if errMessage != "" {
		return balance, errMessage
	}

	var storedValue storedvalue.StoredValue
	storedValue, err, _ := storedValue.FromBytes(res)
	if err != nil {
		return balance, err.Error()
	}
	account := storedValue.Account
	purseID := account.PurseId.Address
	var mintUref []byte
	for _, namedKey := range account.NamedKeys {
		if namedKey.Name == "mint" {
			mintUref = namedKey.Key.Uref.Address
			break
		}
	}

	resBytes := append(mintUref, purseID...)
	localBytes := util.Blake2b256(resBytes)

	res, errMessage = Query(client, stateHash, "local", localBytes, []string{}, protocolVersion)
	if errMessage != "" {
		return balance, errMessage
	}

	storedValue, err, _ = storedValue.FromBytes(res)
	if err != nil {
		return balance, err.Error()
	}
	uref := storedValue.ClValue.ToStateValues().GetKey().GetUref().GetUref()

	res, errMessage = Query(client, stateHash, "uref", uref, []string{}, protocolVersion)
	if errMessage != "" {
		return balance, errMessage
	}

	storedValue, err, _ = storedValue.FromBytes(res)
	balance = storedValue.ClValue.ToStateValues().GetBigInt().GetValue()

	return balance, errMessage
}
