package util

import (
	"testing"

	"math/big"

	"github.com/hdac-io/casperlabs-ee-grpc-go-util/protobuf/io/casperlabs/casper/consensus/state"
	"github.com/stretchr/testify/assert"
)

func TestFromStringToBytes(t *testing.T) {
	src := FromString("안녕하세요")
	res := ToBytes(src)
	assert.Equal(t, res, []byte{
		19, 0, 0, 0,
		15, 0, 0, 0,
		236, 149, 136, 235, 133, 149, 237, 149, 152, 236, 132, 184, 236, 154, 148,
		10})
}

func TestFromBigIntTobytes(t *testing.T) {
	src := FromU512(new(big.Int).SetUint64(256))
	res := ToBytes(src)
	assert.Equal(t, res, []byte{
		3, 0, 0, 0,
		2, 0, 1,
		8})
}

func TestFromByteToBigInt(t *testing.T) {
	src := []byte{2, 0, 1}
	res := fromByteToBigInt(src)

	assert.Equal(t, "256", res.String())
}

func TestFromU32ToBytes(t *testing.T) {
	src := FromU32(uint32(10))
	res := ToBytes(src)
	assert.Equal(t, res, []byte{
		4, 0, 0, 0,
		10, 0, 0, 0,
		4})
}

func TestFromU64ToBytes(t *testing.T) {
	src := FromU64(uint64(67305985))
	res := ToBytes(src)
	assert.Equal(t, res, []byte{
		8, 0, 0, 0,
		1, 2, 3, 4, 0, 0, 0, 0,
		2})
}

func TestFromStringListToBytes(t *testing.T) {
	src := FromStringList([]string{"abc", "defgh", "ijklmnop"})
	res := ToBytes(src)
	assert.Equal(t, res, []byte{
		32, 0, 0, 0,
		3, 0, 0, 0,
		3, 0, 0, 0,
		97, 98, 99,
		5, 0, 0, 0,
		100, 101, 102, 103, 104,
		8, 0, 0, 0,
		105, 106, 107, 108, 109, 110, 111, 112,
		14, 10})
}

func TestFromOptionToBytes(t *testing.T) {
	src := FromOption(AbiUint64ToBytes(uint64(10)), state.CLType_U64)
	res := ToBytes(src)
	assert.Equal(t, []byte{
		9, 0, 0, 0,
		1,
		10, 0, 0, 0, 0, 0, 0, 0,
		13, 5}, res)
}
