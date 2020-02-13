// Package util 은 Casperlabs의 Execution Engine과 연동시 필요한 모듈을 정의한 모듈이다.
package util

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"math/big"
	"strconv"
	"strings"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/hdac-io/casperlabs-ee-grpc-go-util/protobuf/io/casperlabs/casper/consensus"
	"github.com/hdac-io/casperlabs-ee-grpc-go-util/protobuf/io/casperlabs/casper/consensus/state"
	"github.com/hdac-io/casperlabs-ee-grpc-go-util/protobuf/io/casperlabs/ipc"
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

// AbiUint32ToBytes 는 uint32 형식을 Abi 형태인 byte array로 변경해주는 함수.
//
// little endian의 uint32형태로 넣는다.
func AbiUint32ToBytes(src uint32) []byte {
	res := make([]byte, 4)
	binary.LittleEndian.PutUint32(res, src)
	return res
}

// AbiUint64ToBytes 는 uint64 형식을 Abi 형태인 byte array로 변경해주는 함수.
//
// little endian의 uint64형태로 넣는다.
func AbiUint64ToBytes(src uint64) []byte {
	res := make([]byte, 8)
	binary.LittleEndian.PutUint64(res, src)
	return res
}

// AbiBytesToBytes 는 byte array 형식을 Abi 형태인 byte array로 변경해주는 함수.
//
// byte array의 길이를 little endian의 uint32형태로 넣고, 그 뒤 src 내용을 붙인다.
func AbiBytesToBytes(src []byte) []byte {
	res := make([]byte, 4)
	binary.LittleEndian.PutUint32(res, uint32(len(src)))
	res = append(res, src...)
	return res
}

// AbiOptionToBytes 는 byte array 형식에 Abi 형태인 Option을 추가한 byte array로 변경해주는 함수.
//
// byte array에서 값이 있으면 앞에 1을 추가하고 없으면 0을 추가한다.
func AbiOptionToBytes(src []byte) []byte {
	res := make([]byte, 1)
	if len(src) > 0 {
		res[0] = 1
		res = append(res, src...)
	}

	return res
}

// AbiStringToBytes 는 string 형식에 Abi 형태인 byte array로 변경해주는 함수.
//
// string의 length를 little endian의 uint32형태로 넣고, src를 utf8인코딩의 byte array로 변환하여 붙인다.
func AbiStringToBytes(src string) []byte {
	res := make([]byte, 4)
	binary.LittleEndian.PutUint32(res, uint32(len(src)))
	res = append(res, []byte(src)...)
	return res
}

// AbiMakeArgs 는 deploy에 사용할 Args를 abi 형태인 byte array로 변경해주는 함수.
//
// Args 개 수를 little endian의 uint32형태로 넣고, 각 Arg를 순차적으로 붙인다.
// 이 때 Arg의 length를 little endian의 uint32형태로 넣고, data 내용을 붙인다.
func AbiMakeArgs(src [][]byte) []byte {
	res := make([]byte, 4)
	binary.LittleEndian.PutUint32(res, uint32(len(src)))

	for _, data := range src {
		bytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(bytes, uint32(len(data)))
		res = append(res, bytes...)
		res = append(res, data...)
	}

	return res
}

// AbiBigIntTobytes 는 big.Int 형식을 Abi 형태인 byte array로 변경해주는 함수.
//
// big.Int를 byte array로 변환 후 reverse하고, 해당 length를 맨 앞에 추가해준다.
func AbiBigIntTobytes(src *big.Int) []byte {
	bytes := reverseBytes(src.Bytes())
	res := []byte{byte(len(bytes))}
	res = append(res, bytes...)

	return res
}

func AbiDeployArgsTobytes(src []*consensus.Deploy_Arg) []byte {
	res := make([]byte, 4)
	binary.LittleEndian.PutUint32(res, uint32(len(src)))

	for _, deployArg := range src {
		bytes := AbiDeployArgTobytes(deployArg.GetValue())
		lenBytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(lenBytes, uint32(len(bytes)))
		res = append(res, lenBytes...)
		res = append(res, bytes...)
	}

	return res
}

func AbiDeployArgTobytes(src *consensus.Deploy_Arg_Value) (res []byte) {
	var data []byte
	switch src.GetValue().(type) {
	case *consensus.Deploy_Arg_Value_OptionalValue:
		data = AbiOptionToBytes(AbiDeployArgTobytes(src.GetOptionalValue()))
	case *consensus.Deploy_Arg_Value_BytesValue:
		data = src.GetBytesValue()
	case *consensus.Deploy_Arg_Value_IntValue:
		data = AbiUint32ToBytes(uint32(src.GetIntValue()))
	case *consensus.Deploy_Arg_Value_IntList:
		// TODO : Need to test
		intValues := src.GetIntList().GetValues()
		for _, value := range intValues {
			data = append(data, AbiUint32ToBytes(uint32(value))...)
		}
	case *consensus.Deploy_Arg_Value_StringValue:
		data = AbiStringToBytes(src.GetStringValue())
	case *consensus.Deploy_Arg_Value_StringList:
		// TODO : Need to test
		strValues := src.GetStringList().GetValues()
		for _, value := range strValues {
			data = append(data, AbiStringToBytes(value)...)
		}
	case *consensus.Deploy_Arg_Value_LongValue:
		data = AbiUint64ToBytes(uint64(src.GetLongValue()))
	case *consensus.Deploy_Arg_Value_BigInt:
		// TODO : ParseUint didn't support 512.. need to change more.
		bitWidth := uint32(64)
		if src.GetBigInt().GetBitWidth() < 64 {
			bitWidth = src.GetBigInt().GetBitWidth()
		}
		val, err := strconv.ParseUint(src.GetBigInt().GetValue(), 10, int(bitWidth))
		if err != nil {
			return []byte{}
		}
		data = AbiBigIntTobytes(new(big.Int).SetUint64(val))
	case *consensus.Deploy_Arg_Value_Key:
		switch src.GetKey().GetValue().(type) {
		case *state.Key_Address_:
			data = append([]byte{WASM}, src.GetKey().GetAddress().GetAccount()...)
		case *state.Key_Hash_:
			data = append([]byte{HASH}, src.GetKey().GetHash().GetHash()...)
		case *state.Key_Uref:
			data = append([]byte{UREF}, src.GetKey().GetUref().GetUref()...)
			data = append(data, []byte{byte(src.GetKey().GetUref().GetAccessRights())}...)
		case *state.Key_Local_:
			data = append([]byte{LOCAL}, src.GetKey().GetLocal().GetHash()...)
		}
	default:
		return []byte{}
	}
	res = append(res, data...)

	return res
}

func JsonStringToDeployArgs(str string) (deployArgs []*consensus.Deploy_Arg, err error) {
	jsonDecoder := json.NewDecoder(strings.NewReader(str))
	_, err = jsonDecoder.Token()
	if err != nil {
		return deployArgs, err
	}

	for jsonDecoder.More() {
		arg := consensus.Deploy_Arg{}
		err := jsonpb.UnmarshalNext(jsonDecoder, &arg)
		if err != nil {
			return deployArgs, err
		}
		deployArgs = append(deployArgs, &arg)
	}

	return deployArgs, nil
}

func reverseBytes(src []byte) []byte {
	len := len(src)
	for i := 0; i < (len / 2); i++ {
		tmp := src[i]
		src[i] = src[len-i-1]
		src[len-i-1] = tmp
	}

	return src
}

// MakeArgsTransferToAccount 는 transfer_to_account.wasm을 사용할 때 Args를 만드는 함수.
//
// string의 수신자 address와 amount를 받아 amount를 abi 형태로 만든다.
// 이 후 2개의 값을 AbiMakeArgs를 통해 하나의 Abi args로 만들어 return 해준다.
func MakeArgsTransferToAccount(address []byte, amount uint64) []byte {
	amountAbi := AbiUint64ToBytes(amount)
	sessionAbiList := [][]byte{address, amountAbi}
	return AbiMakeArgs(sessionAbiList)
}

func MakeArgsStoreTransferToAccount(address []byte, amount uint64) []byte {
	sessionAbiList := [][]byte{AbiStringToBytes("send"), address, AbiUint64ToBytes(amount)}
	return AbiMakeArgs(sessionAbiList)
}

// MakeArgsStandardPayment 는 standard_payment.wasm을 사용할 때 Args를 만드는 함수.
//
// big.Int은 amount를 받아 Abi형태로 만들고, AbiMakeArgs를 통해 Abi args로 만들어 return 해준다.
func MakeArgsStandardPayment(amount *big.Int) []byte {
	paymentAbiList := [][]byte{AbiBigIntTobytes(amount)}
	paymentAbi := AbiMakeArgs(paymentAbiList)
	return paymentAbi
}

// MakeArgsBonding 은 bonding.wasm을 사용할 때 Args를 만드는 함수.
//
// uint64의 amount를 받아 Abi형태로 만들고, AbiMakeArgs를 통해 Abi args로 만들어 return 해준다.
func MakeArgsBonding(amount uint64) []byte {
	abiList := [][]byte{AbiUint64ToBytes(amount)}
	abi := AbiMakeArgs(abiList)
	return abi
}

func MakeArgsStoreBonding(method string, amount uint64, purseId []byte) []byte {
	abiList := [][]byte{AbiStringToBytes(method), AbiUint64ToBytes(amount), purseId}
	abi := AbiMakeArgs(abiList)
	return abi
}

// MakeArgsUnBonding 은 unbonding.wasm을 사용할 때 Args를 만드는 함수.
//
// uint64의 amount를 받아 Abi형태로 만들고, Option Abi를 추가한다. 이 후 AbiMakeArgs를 통해 Abi args로 만들어 return 해준다.
func MakeArgsUnBonding(amount uint64) []byte {
	abiList := [][]byte{AbiOptionToBytes(AbiUint64ToBytes(amount))}
	abi := AbiMakeArgs(abiList)
	return abi
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
	sessionArgs []byte,
	paymentType ContractType,
	paymentData []byte,
	paymentArgs []byte,
	gasPrice uint64,
	int64Timestamp int64,
	chainName string) (deploy *ipc.DeployItem) {
	timestamp := uint64(int64Timestamp)

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

	deploy = &ipc.DeployItem{
		Address:           fromAddress,
		Session:           MakeDeployPayload(sessionType, sessionData, sessionArgs),
		Payment:           MakeDeployPayload(paymentType, paymentData, paymentArgs),
		GasPrice:          gasPrice,
		AuthorizationKeys: [][]byte{fromAddress},
		DeployHash:        headerHash}

	return deploy
}

type ContractType int

const (
	WASM = iota
	HASH
	UREF
	LOCAL
	NAME
)

func MakeDeployCode(contractType ContractType, data []byte, args []byte) *consensus.Deploy_Code {
	deployCode := &consensus.Deploy_Code{AbiArgs: args}
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
