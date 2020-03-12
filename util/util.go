// Package util 은 Casperlabs의 Execution Engine과 연동시 필요한 모듈을 정의한 모듈이다.
package util

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"github.com/hdac-io/casperlabs-ee-grpc-go-util/protobuf/io/casperlabs/casper/consensus"
	"github.com/hdac-io/casperlabs-ee-grpc-go-util/protobuf/io/casperlabs/casper/consensus/state"
	"github.com/hdac-io/casperlabs-ee-grpc-go-util/protobuf/io/casperlabs/ipc"
	"github.com/hdac-io/casperlabs-ee-grpc-go-util/storedvalue"
	"golang.org/x/crypto/blake2b"
)

// StrEmptyStateHash 는 비어있는 trie의 state Hash 값으로 초기 state hash 값
const StrEmptyStateHash = "3307a54ca6d5bfbafc0ef1b003f3ec4941c011ee7f79889e44416754de2f091d"

// LoadWasmFile 은 wasm 파일을 byte array로 return 해주는 함수.
func LoadWasmFile(path string) []byte {
	wasmCode, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return wasmCode
}

// EncodeToHexString 는 byte array를 hex string으로 변경해주는 함수.
func EncodeToHexString(src []byte) string {
	return hex.EncodeToString(src)
}

// DecodeHexString 는 hex string을 byte array로 변경해주는 함수.
func DecodeHexString(str string) []byte {
	res, err := hex.DecodeString(str)
	if err != nil {
		panic(err)
	}

	return res
}

// Blake2b256 는 blake2b 256 hash 결과 값을 return 해주는 함수.
func Blake2b256(ob []byte) []byte {
	hash, err := blake2b.New256(nil)
	if err != nil {
		panic(err)
	}

	hash.Write(ob)
	return hash.Sum(nil)
}

// MakeProtocolVersion 은 major, minor, patch의 값을 받아 ProtocolVersion 을 만들어주는 함수
func MakeProtocolVersion(major uint32, minor uint32, patch uint32) *state.ProtocolVersion {
	return &state.ProtocolVersion{Major: uint32(major), Minor: uint32(minor), Patch: uint32(patch)}
}

func AbiDeployArgsTobytes(src []*consensus.Deploy_Arg) ([]byte, error) {
	res := make([]byte, 4)
	binary.LittleEndian.PutUint32(res, uint32(len(src)))

	for _, deployArg := range src {
		var clValue storedvalue.CLValue
		clValue, err := clValue.FromDeployArgValue(deployArg.GetValue())
		if err != nil {
			return nil, err
		}
		res = append(res, clValue.ToBytes()...)
	}

	return res, nil
}

func JsonStringToDeployArgs(str string) (deployArgs []*consensus.Deploy_Arg, err error) {
	if str == "" {
		return []*consensus.Deploy_Arg{}, nil
	}

	jsonDecoder := json.NewDecoder(strings.NewReader(str))
	_, err = jsonDecoder.Token()
	if err != nil {
		return nil, err
	}

	for jsonDecoder.More() {
		arg := consensus.Deploy_Arg{}
		err := jsonpb.UnmarshalNext(jsonDecoder, &arg)
		if err != nil {
			return nil, err
		}
		deployArgs = append(deployArgs, &arg)
	}

	return deployArgs, nil
}

func DeployArgsToJsonString(args []*consensus.Deploy_Arg) (string, error) {
	m := &jsonpb.Marshaler{}
	str := "["
	for idx, arg := range args {
		if idx != 0 {
			str += ","
		}
		s, err := m.MarshalToString(arg)
		if err != nil {
			return "", err
		}
		str += s
	}
	str += "]"

	return str, nil
}

// MakeDeploy 는 address, sessionCode, sessionArgs, paymentCode, paymentArgs, gasPrice, timestamp를 받아 DeployItem을 만들어주는 함수.
//
// Seesion, Payment 데이터로 Deploy Body를 만들고 이를 Marshal한 값을 Blake2b256 Hash를 하여 Deploy Body Hash 값을 만든다.
// Deploy Body Hash 값을 포함한 Deploy Header 값을 만들고 이를 Marshal한 값을 Blake2b256 Hash하여 Deploy Header Hash 값을 만든다.
// Deploy Header Hash 값을 Deploy Item을 만들고 return 해준다.
func MakeDeploy(
	fromAddress []byte,
	sessionType ContractType,
	sessionData []byte,
	sessionArgsStr string,
	paymentType ContractType,
	paymentData []byte,
	paymentArgsStr string,
	gasPrice uint64,
	int64Timestamp int64,
	chainName string) (deploy *ipc.DeployItem, err error) {
	timestamp := uint64(int64Timestamp)

	sessionArgs, err := JsonStringToDeployArgs(sessionArgsStr)
	if err != nil {
		return nil, err
	}
	paymentArgs, err := JsonStringToDeployArgs(paymentArgsStr)
	if err != nil {
		return nil, err
	}
	deployBody := &consensus.Deploy_Body{
		Session: MakeDeployCode(sessionType, sessionData, sessionArgs),
		Payment: MakeDeployCode(paymentType, paymentData, paymentArgs)}

	marshalDeployBody, _ := proto.Marshal(deployBody)
	bodyHash := Blake2b256(marshalDeployBody)

	deployHeader := &consensus.Deploy_Header{
		AccountPublicKey: fromAddress,
		Timestamp:        timestamp,
		GasPrice:         gasPrice,
		BodyHash:         bodyHash,
		ChainName:        chainName}

	marshalDeployHeader, _ := proto.Marshal(deployHeader)
	headerHash := Blake2b256(marshalDeployHeader)

	sessionAbi, err := AbiDeployArgsTobytes(sessionArgs)
	if err != nil {
		return nil, err
	}
	paymentAbi, err := AbiDeployArgsTobytes(paymentArgs)
	if err != nil {
		return nil, err
	}

	deploy = &ipc.DeployItem{
		Address:           fromAddress,
		Session:           MakeDeployPayload(sessionType, sessionData, sessionAbi),
		Payment:           MakeDeployPayload(paymentType, paymentData, paymentAbi),
		GasPrice:          gasPrice,
		AuthorizationKeys: [][]byte{fromAddress},
		DeployHash:        headerHash}

	return deploy, nil
}

type ContractType int

const (
	WASM = iota
	HASH
	UREF
	LOCAL
	NAME
)

func MakeDeployCode(contractType ContractType, data []byte, args []*consensus.Deploy_Arg) *consensus.Deploy_Code {
	deployCode := &consensus.Deploy_Code{Args: args}
	switch contractType {
	case WASM:
		deployCode.Contract = &consensus.Deploy_Code_Wasm{Wasm: data}
	case UREF:
		deployCode.Contract = &consensus.Deploy_Code_Uref{Uref: data}
	case HASH:
		deployCode.Contract = &consensus.Deploy_Code_Hash{Hash: data}
	case NAME:
		deployCode.Contract = &consensus.Deploy_Code_Name{Name: string(data)}
	default:
		deployCode = nil
	}

	return deployCode
}

func MakeDeployPayload(contractType ContractType, data []byte, args []byte) *ipc.DeployPayload {
	deployPayload := &ipc.DeployPayload{}
	switch contractType {
	case WASM:
		deployPayload.Payload = &ipc.DeployPayload_DeployCode{DeployCode: &ipc.DeployCode{Code: data, Args: args}}
	case UREF:
		deployPayload.Payload = &ipc.DeployPayload_StoredContractUref{StoredContractUref: &ipc.StoredContractURef{Uref: data, Args: args}}
	case HASH:
		deployPayload.Payload = &ipc.DeployPayload_StoredContractHash{StoredContractHash: &ipc.StoredContractHash{Hash: data, Args: args}}
	case NAME:
		deployPayload.Payload = &ipc.DeployPayload_StoredContractName{StoredContractName: &ipc.StoredContractName{StoredContractName: string(data), Args: args}}
	default:
		deployPayload = nil
	}

	return deployPayload
}

// MakeInitDeploys 은 Deploy Item array를 할당 받기위한 함수.
func MakeInitDeploys() []*ipc.DeployItem {
	return []*ipc.DeployItem{}
}

// AddDeploy 는 deploy를 deploy list에 추가하기위한 함수.
func AddDeploy(deploys []*ipc.DeployItem, deploy *ipc.DeployItem) []*ipc.DeployItem {
	return append(deploys, deploy)
}

func GenesisConfigMock(
	chainName string, address []byte, balance string, bondedAmount string, protocolVersion *state.ProtocolVersion,
	mapCosts map[string]uint32, mintInstallWasmPath string, posInstallWasmPath string) (
	*ipc.ChainSpec_GenesisConfig, error) {
	genesisConfig := ipc.ChainSpec_GenesisConfig{}
	genesisConfig.Name = chainName
	genesisConfig.Timestamp = 0
	genesisConfig.ProtocolVersion = protocolVersion

	// load mint_install.wasm, pos_install.wasm
	genesisConfig.MintInstaller = LoadWasmFile(mintInstallWasmPath)
	genesisConfig.PosInstaller = LoadWasmFile(posInstallWasmPath)

	// GenesisAccount
	accounts := make([]*ipc.ChainSpec_GenesisAccount, 1)
	accounts[0] = &ipc.ChainSpec_GenesisAccount{}
	accounts[0].PublicKey = address
	accounts[0].Balance = &state.BigInt{Value: balance, BitWidth: 512}
	accounts[0].BondedAmount = &state.BigInt{Value: bondedAmount, BitWidth: 512}
	genesisConfig.Accounts = accounts

	// CostTable
	genesisConfig.Costs = &ipc.ChainSpec_CostTable{
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

	genesisConfig.DeployConfig = &ipc.ChainSpec_DeployConfig{
		MaxTtlMillis:    86400000,
		MaxDependencies: 10,
	}

	return &genesisConfig, nil
}
