package util

import (
	"encoding/binary"
	"encoding/hex"
	"io/ioutil"
	"math/big"

	state "github.com/hdac-io/casperlabs-ee-grpc-go-util/protobuf/io/casperlabs/casper/consensus/state"
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
