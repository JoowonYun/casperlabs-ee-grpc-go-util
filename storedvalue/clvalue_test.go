package storedvalue

import (
	"testing"

	"github.com/hdac-io/casperlabs-ee-grpc-go-util/protobuf/io/casperlabs/casper/consensus/state"
	"github.com/stretchr/testify/assert"
)

func TestFromStringToBytes(t *testing.T) {
	stateValue := &state.Value{
		Value: &state.Value_StringValue{StringValue: "안녕하세요"},
	}

	var clValue CLValue
	clValue, err := clValue.FromStateValue(stateValue)

	assert.NoError(t, err)
	assert.Equal(t,
		[]byte{
			19, 0, 0, 0,
			15, 0, 0, 0,
			236, 149, 136, 235, 133, 149, 237, 149, 152, 236, 132, 184, 236, 154, 148,
			10},
		clValue.ToBytes())
}

func TestFromBigIntTobytes(t *testing.T) {
	stateValue := &state.Value{
		Value: &state.Value_BigInt{
			BigInt: &state.BigInt{Value: "256", BitWidth: 512},
		},
	}

	var clValue CLValue
	clValue, err := clValue.FromStateValue(stateValue)

	assert.NoError(t, err)
	assert.Equal(t, []byte{
		3, 0, 0, 0,
		2, 0, 1,
		8},
		clValue.ToBytes())
}

func TestFromByteToBigInt(t *testing.T) {
	src := []byte{2, 0, 1}
	res := fromByteToBigInt(src)

	assert.Equal(t, "256", res.String())
}

func TestFromU32ToBytes(t *testing.T) {

	stateValue := &state.Value{
		Value: &state.Value_IntValue{
			IntValue: 10,
		},
	}

	var clValue CLValue
	clValue, err := clValue.FromStateValue(stateValue)

	assert.NoError(t, err)
	assert.Equal(t, []byte{
		4, 0, 0, 0,
		10, 0, 0, 0,
		1},
		clValue.ToBytes())
}

func TestFromU64ToBytes(t *testing.T) {
	stateValue := &state.Value{
		Value: &state.Value_LongValue{
			LongValue: uint64(67305985),
		},
	}

	var clValue CLValue
	clValue, err := clValue.FromStateValue(stateValue)

	assert.NoError(t, err)
	assert.Equal(t, []byte{
		8, 0, 0, 0,
		1, 2, 3, 4, 0, 0, 0, 0,
		2},
		clValue.ToBytes())
}

func TestFromStringListToBytes(t *testing.T) {
	stateValue := &state.Value{
		Value: &state.Value_StringList{
			StringList: &state.StringList{Values: []string{"abc", "defgh", "ijklmnop"}},
		},
	}

	var clValue CLValue
	clValue, err := clValue.FromStateValue(stateValue)

	assert.NoError(t, err)
	assert.Equal(t, []byte{
		32, 0, 0, 0,
		3, 0, 0, 0,
		3, 0, 0, 0,
		97, 98, 99,
		5, 0, 0, 0,
		100, 101, 102, 103, 104,
		8, 0, 0, 0,
		105, 106, 107, 108, 109, 110, 111, 112,
		14, 10},
		clValue.ToBytes())
}
