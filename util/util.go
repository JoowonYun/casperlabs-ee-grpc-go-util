package util

import (
	"encoding/hex"
	"io/ioutil"

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
