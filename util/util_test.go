package util

import (
	"testing"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/stretchr/testify/assert"
)

func TestEncodeToHexString(t *testing.T) {
	res := EncodeToHexString([]byte{147, 37, 61, 17})
	assert.Equal(t, res, "93253d11", "they should be equal")
}

func TestDecodeHexString(t *testing.T) {
	res := DecodeHexString("1d526a")
	assert.Equal(t, res, []byte{29, 82, 106}, "they should be equal")
}

func TestBlake2b256_0(t *testing.T) {
	res := Blake2b256([]byte{0})
	assert.Equal(t, EncodeToHexString(res), "03170a2e7597b7b7e3d84c05391d139a62b157e78786d8c082f29dcf4c111314", "they should be equal")
}

func TestBlake2b256_LEN_8(t *testing.T) {
	res := Blake2b256([]byte{147, 68, 51, 37, 29, 6, 244, 90})
	assert.Equal(t, EncodeToHexString(res), "f0f84944e0ccfa9e67383e6a448291787d208c8e46adc849f714078663d1dd36", "they should be equal")
}

func TestJsonStringToDeployArgs(t *testing.T) {
	inputStr := `[{"name":"amount","value":{"value":{"i32":123456}}},{"name":"fee","value":{"value":{"i32":54321}}}]`
	args, err := JsonStringToDeployArgs(inputStr)
	assert.NoError(t, err)

	assert.Equal(t, "amount", args[0].GetName())
	assert.Equal(t, int32(123456), args[0].GetValue().GetValue().GetI32())
	assert.Equal(t, "fee", args[1].GetName())
	assert.Equal(t, int32(54321), args[1].GetValue().GetValue().GetI32())

	m := &jsonpb.Marshaler{}
	str := "["
	for idx, arg := range args {
		if idx != 0 {
			str += ","
		}
		s, _ := m.MarshalToString(arg)
		str += s
	}
	str += "]"

	args2, err := JsonStringToDeployArgs(str)
	assert.NoError(t, err)
	assert.Equal(t, args[0].GetName(), args2[0].GetName())
	assert.Equal(t, args[0].GetValue().GetValue().GetI32(), args2[0].GetValue().GetValue().GetI32())
	assert.Equal(t, args[1].GetName(), args2[1].GetName())
	assert.Equal(t, args[1].GetValue().GetValue().GetI32(), args2[1].GetValue().GetValue().GetI32())
}
