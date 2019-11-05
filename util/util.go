package util

import (
	"encoding/binary"
	"encoding/hex"
	"io/ioutil"
	"math/big"

	"github.com/golang/protobuf/proto"
	"github.com/hdac-io/casperlabs-ee-grpc-go-util/protobuf/io/casperlabs/casper/consensus"
	"github.com/hdac-io/casperlabs-ee-grpc-go-util/protobuf/io/casperlabs/casper/consensus/state"
	"github.com/hdac-io/casperlabs-ee-grpc-go-util/protobuf/io/casperlabs/ipc"
	"golang.org/x/crypto/blake2b"
)

const StrEmptyStateHash = "3307a54ca6d5bfbafc0ef1b003f3ec4941c011ee7f79889e44416754de2f091d"

func LoadWasmFile(path string) []byte {
	wasmCode, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return wasmCode
}

func EncodeToHexString(src []byte) string {
	return hex.EncodeToString(src)
}

func DecodeHexString(str string) []byte {
	res, err := hex.DecodeString(str)
	if err != nil {
		panic(err)
	}

	return res
}

func Blake2b256(ob []byte) []byte {
	hash, err := blake2b.New256(nil)
	if err != nil {
		panic(err)
	}

	hash.Write(ob)
	return hash.Sum(nil)
}

func MakeProtocolVersion(major uint32, minor uint32, patch uint32) *state.ProtocolVersion {
	return &state.ProtocolVersion{Major: uint32(major), Minor: uint32(minor), Patch: uint32(patch)}
}

func AbiUint32ToBytes(src uint32) []byte {
	res := make([]byte, 4)
	binary.LittleEndian.PutUint32(res, src)
	return res
}

func AbiUint64ToBytes(src uint64) []byte {
	res := make([]byte, 8)
	binary.LittleEndian.PutUint64(res, src)
	return res
}

func AbiBytesToBytes(src []byte) []byte {
	res := make([]byte, 4)
	binary.LittleEndian.PutUint32(res, uint32(len(src)))
	res = append(res, src...)
	return res
}

func AbiOptionToBytes(src []byte) []byte {
	res := make([]byte, 1)
	if len(src) > 0 {
		res[0] = 1
		res = append(res, src...)
	}

	return res
}

func AbiStringToBytes(src string) []byte {
	res := make([]byte, 4)
	binary.LittleEndian.PutUint32(res, uint32(len(src)))
	res = append(res, []byte(src)...)
	return res
}

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

func AbiBigIntTobytes(src *big.Int) []byte {
	bytes := reverseBytes(src.Bytes())
	res := []byte{byte(len(bytes))}
	res = append(res, bytes...)

	return res
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

func MakeArgsTransferToAccount(address string, amount uint64) []byte {
	addressBytes := DecodeHexString(address)
	addressAbi := AbiBytesToBytes(addressBytes)
	amountAbi := AbiUint64ToBytes(amount)
	sessionAbiList := [][]byte{addressAbi, amountAbi}
	return AbiMakeArgs(sessionAbiList)
}

func MakeArgsStandardPayment(amount *big.Int) []byte {
	paymentAbiList := [][]byte{AbiBigIntTobytes(amount)}
	paymentAbi := AbiMakeArgs(paymentAbiList)
	return paymentAbi
}

func MakeArgsBonding(amount uint64) []byte {
	abiList := [][]byte{AbiUint64ToBytes(amount)}
	abi := AbiMakeArgs(abiList)
	return abi
}

func MakeArgsUnBonding(amount uint64) []byte {
	abiList := [][]byte{AbiOptionToBytes(AbiUint64ToBytes(amount))}
	abi := AbiMakeArgs(abiList)
	return abi
}

func MakeDeploy(
	strFromAddress string,
	sessionCode []byte,
	sessionArgs []byte,
	paymentCode []byte,
	paymentArgs []byte,
	gasPrice uint64,
	int64Timestamp int64) (deploy *ipc.DeployItem) {
	fromAddress := DecodeHexString(strFromAddress)
	timestamp := uint64(int64Timestamp)

	deployBody := &consensus.Deploy_Body{
		Session: &consensus.Deploy_Code{Contract: &consensus.Deploy_Code_Wasm{Wasm: sessionCode}, AbiArgs: sessionArgs},
		Payment: &consensus.Deploy_Code{Contract: &consensus.Deploy_Code_Wasm{Wasm: paymentCode}, AbiArgs: paymentArgs}}

	marshalDeployBody, _ := proto.Marshal(deployBody)
	bodyHash := Blake2b256(marshalDeployBody)

	deployHeader := &consensus.Deploy_Header{
		AccountPublicKey: fromAddress,
		Timestamp:        timestamp,
		GasPrice:         gasPrice,
		BodyHash:         bodyHash}

	marshalDeployHeader, _ := proto.Marshal(deployHeader)
	headerHash := Blake2b256(marshalDeployHeader)

	deploy = &ipc.DeployItem{
		Address:           fromAddress,
		Session:           &ipc.DeployPayload{Payload: &ipc.DeployPayload_DeployCode{DeployCode: &ipc.DeployCode{Code: sessionCode, Args: sessionArgs}}},
		Payment:           &ipc.DeployPayload{Payload: &ipc.DeployPayload_DeployCode{DeployCode: &ipc.DeployCode{Code: paymentCode, Args: paymentArgs}}},
		GasPrice:          gasPrice,
		AuthorizationKeys: [][]byte{fromAddress},
		DeployHash:        headerHash}

	return deploy
}

func MakeInitDeploys() []*ipc.DeployItem {
	return []*ipc.DeployItem{}
}

func AddDeploy(deploys []*ipc.DeployItem, deploy *ipc.DeployItem) []*ipc.DeployItem {
	return append(deploys, deploy)
}
