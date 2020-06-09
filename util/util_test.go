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

func TestJson(t *testing.T) {
	js := `[
		{
		   "name":"method",
		   "value":{
			  "cl_type":{
				 "simple_type":"STRING"
			  },
			  "value":{
				 "str_value":"get_token"
			  }
		   }
		},
		{
		   "name":"ver1_pubkey",
		   "value":{
			  "cl_type":{
				 "list_type":{
					"inner":{
					   "simple_type":"STRING"
					}
				 }
			  },
			  "value":{
				 "list_value":{
					"values":[
					   {
						  "str_value":"02c4ef70543e18889167ca67c8aa28c1d4c259e89cb34483a8ed6cfd3a03e8246b"
					   }
					]
				 }
			  }
		   }
		},
		{
		   "name":"message",
		   "value":{
			  "cl_type":{
				 "list_type":{
					"inner":{
					   "simple_type":"STRING"
					}
				 }
			  },
			  "value":{
				 "list_value":{
					"values":[
					   {
						  "str_value":"69046d44e3d75d48436377626372a44a5066966b5d72c00b67769c1cc6a8619a"
					   }
					]
				 }
			  }
		   }
		},
		{
		   "name":"signature",
		   "value":{
			  "cl_type":{
				 "list_type":{
					"inner":{
					   "simple_type":"STRING"
					}
				 }
			  },
			  "value":{
				 "list_value":{
					"values":[
					   {
						  "str_value":"24899366fd3d5dfe6740df1e5f467a53f1a3aaafce26d8df1497a925c55b5c266339a95fe6507bd611b0e3b6e74e3bb7f19eeb1165615e5cebe7f40e5765bc41"
					   }
					]
				 }
			  }
		   }
		}
	 ]`

	deployArgs, err := JsonStringToDeployArgs(js)
	assert.NoError(t, err)

	assert.Equal(t, "get_token", deployArgs[0].GetValue().GetValue().GetStrValue())
	assert.Equal(t, "02c4ef70543e18889167ca67c8aa28c1d4c259e89cb34483a8ed6cfd3a03e8246b", deployArgs[1].GetValue().GetValue().GetListValue().GetValues()[0].GetStrValue())
	assert.Equal(t, "69046d44e3d75d48436377626372a44a5066966b5d72c00b67769c1cc6a8619a", deployArgs[2].GetValue().GetValue().GetListValue().GetValues()[0].GetStrValue())
	assert.Equal(t, "24899366fd3d5dfe6740df1e5f467a53f1a3aaafce26d8df1497a925c55b5c266339a95fe6507bd611b0e3b6e74e3bb7f19eeb1165615e5cebe7f40e5765bc41", deployArgs[3].GetValue().GetValue().GetListValue().GetValues()[0].GetStrValue())
}
