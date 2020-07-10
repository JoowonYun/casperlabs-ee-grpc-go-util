package util

import (
	"encoding/hex"
	"testing"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/hdac-io/casperlabs-ee-grpc-go-util/protobuf/io/casperlabs/casper/consensus/state"
	"github.com/hdac-io/casperlabs-ee-grpc-go-util/storedvalue"
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

func TestListAbi(t *testing.T) {
	js := `[{
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
					   "str_value":"123"
					},
					{
					   "str_value":"456"
					}
				 ]
			  }
		   }
		}
	 }]`

	deployArgs, err := JsonStringToDeployArgs(js)
	assert.NoError(t, err)

	abi, err := AbiDeployArgsTobytes(deployArgs)
	assert.NoError(t, err)

	assert.Equal(t, []byte{
		1, 0, 0, 0,
		18, 0, 0, 0,
		2, 0, 0, 0,
		3, 0, 0, 0,
		49, 50, 51,
		3, 0, 0, 0,
		52, 53, 54,
		14, 10,
	}, abi)
}

func TestMapAbi(t *testing.T) {

	src := []byte{
		0,
		45, 0, 0, 0,
		2, 0, 0, 0,
		9, 0, 0, 0,
		107, 121, 99, 95, 108, 101, 118, 101, 108,
		1, 0, 0, 0,
		49,
		14, 0, 0, 0,
		115, 119, 97, 112, 112, 101, 100, 95, 97, 109, 111, 117, 110, 116,
		1, 0, 0, 0,
		48,
		17, 10, 10}

	var storedValue storedvalue.StoredValue
	storedValue, err, _ := storedValue.FromBytes(src)
	assert.NoError(t, err)

	value := storedValue.ClValue.ToCLInstanceValue()

	assert.Equal(t, len(value.GetMapValue().GetValues()), 2)

	assert.Equal(t, value.GetMapValue().GetValues()[0].GetKey().GetStrValue(), "kyc_level")
	assert.Equal(t, value.GetMapValue().GetValues()[0].GetValue().GetStrValue(), "1")

	assert.Equal(t, value.GetMapValue().GetValues()[1].GetKey().GetStrValue(), "swapped_amount")
	assert.Equal(t, value.GetMapValue().GetValues()[1].GetValue().GetStrValue(), "0")
}

func TestJsonAllType(t *testing.T) {
	jsonStr := `
    [
		{"name" : "bool", "value" : {"cl_type" : { "simple_type" : "BOOL" }, "value" : { "bool_value" : true }}},
		{"name" : "i32", "value" : {"cl_type" : { "simple_type" : "I32" }, "value" : { "i32" : 314 }}},
		{"name" : "i64", "value" : {"cl_type" : { "simple_type" : "I64" }, "value" : { "i64" : 2342 }}},
		{"name" : "u8", "value" : {"cl_type" : { "simple_type" : "U8" }, "value" : { "u8" : 7 }}},
		{"name" : "u32", "value" : {"cl_type" : { "simple_type" : "U32" }, "value" : { "u32" : 314 }}},		
		{"name" : "u64", "value" : {"cl_type" : { "simple_type" : "U64" }, "value" : { "u64" : 2342 }}},		
		{"name" : "u128", "value" : {"cl_type" : { "simple_type" : "U128" }, "value" : { "u128" :  {"value" : "123456789101112131415161718"}}}},		
		{"name" : "u256", "value" : {"cl_type" : { "simple_type" : "U256" }, "value" : { "u256" :  {"value" : "123456789101112131415161718"}}}},		
		{"name" : "u512", "value" : {"cl_type" : { "simple_type" : "U512" }, "value" : { "u512" :  {"value" : "123456789101112131415161718"}}}},		
		{"name" : "unit", "value" : {"cl_type" : { "simple_type" : "UNIT" }, "value" : { "unit" : {} }}},		
		{"name" : "string", "value" : {"cl_type" : { "simple_type" : "STRING" }, "value" : { "str_value" : "Hello, world!" }}},		
		{"name" : "accountKey", "value" : {"cl_type" : { "simple_type" : "KEY" }, "value" : {"key": {"address": {"account": "1wJD3Z0NZG/W3ygqj3qPoFpmKb7AHYAkw2EescH7n4Q="}}}}},		
		{"name" : "hashKey", "value" : {"cl_type" : { "simple_type" : "KEY" }, "value" : {"key": {"hash": {"hash": "1wJD3Z0NZG/W3ygqj3qPoFpmKb7AHYAkw2EescH7n4Q="}}}}},		
		{"name" : "urefKey", "value" : {"cl_type" : { "simple_type" : "KEY" }, "value" : {"key": {"uref": {"uref": "1wJD3Z0NZG/W3ygqj3qPoFpmKb7AHYAkw2EescH7n4Q=", "access_rights": 5}}}}},		
		{"name" : "uref", "value" : {"cl_type" : { "simple_type" : "UREF" }, "value" : {"uref": {"uref": "1wJD3Z0NZG/W3ygqj3qPoFpmKb7AHYAkw2EescH7n4Q=", "access_rights": 5}}}},		
		{
			"name" : "maybe_u64",
			"value" : {
				"cl_type" : {"option_type" : {"inner" : {"simple_type" : "U64"}}},
				"value" : {"option_value" : {}}
			}
		},		
		{
			"name" : "maybe_u64",
			"value" : {
				"cl_type" : {"option_type" : {"inner" : {"simple_type" : "U64"}}},
				"value" : {"option_value" : {"value" : {"u64" : 2342}}}
			}
		},		
		{
			"name" : "list_i32",
			"value" : {
				"cl_type" : {"list_type" : {"inner" : {"simple_type" : "I32"}}},
				"value" : {"list_value" : {"values" : [{"i32" : 0}, {"i32" : 1}, {"i32" : 2}, {"i32" : 3}]}}
			}
		},		
		{
			"name" : "fixed_list_str",
			"value" : {
				"cl_type" : {"fixed_list_type" : {"inner" : {"simple_type" : "STRING"}, "len" : 3}},
				"value" : {"fixed_list_value" : {"length" : 3, "values" : [{"str_value" : "A"}, {"str_value" : "B"}, {"str_value" : "C"}]}}
			}
		},		
		{
			"name" : "err_string",
			"value" : {
				"cl_type" : {"result_type" : {"ok" : {"simple_type" : "BOOL"}, "err" : {"simple_type" : "STRING"}}},
				"value" : {"result_value" : {"err" : {"str_value" : "Hello, world!"}}}
			}
		},		
		{
			"name" : "ok_bool",
			"value" : {
				"cl_type" : {"result_type" : {"ok" : {"simple_type" : "BOOL"}, "err" : {"simple_type" : "STRING"}}},
				"value" : {"result_value" : {"ok" : {"bool_value" : true}}}
			}
		},		
		{
			"name" : "map_string_i32",
			"value" : {
				"cl_type" : {"map_type" : {"key" : {"simple_type" : "STRING"}, "value" : {"simple_type" : "I32"}}},
				"value" : {
					"map_value" : {
						"values" : [
							{"key" : {"str_value" : "A"}, "value" : {"i32" : 0}},
							{"key" : {"str_value" : "B"}, "value" : {"i32" : 1}},
							{"key" : {"str_value" : "C"}, "value" : {"i32" : 2}}
						]
					}
				}
			}
		},		
		{
			"name" : "tuple1",
			"value" : {
				"cl_type" : {"tuple1_type" : {"type0" : {"simple_type" : "U8"}}},
				"value" : {"tuple1_value" : {"value_1" : {"u8" : 8}}}
			}
		},		
		{
			"name" : "tuple2",
			"value" : {
				"cl_type" : {"tuple2_type" : {"type0" : {"simple_type" : "U8"}, "type1" : {"simple_type" : "U32"}}},
				"value" : {"tuple2_value" : {"value_1" : {"u8" : 8}, "value_2" : {"u32" : 314}}}
			}
		},		
		{
			"name" : "tuple3",
			"value" : {
			"cl_type" : {"tuple3_type" : {
				"type0" : {"simple_type" : "U8"},
				"type1" : {"simple_type" : "U32"},
				"type2" : {"simple_type" : "U64"}
			}},
			"value" : {"tuple3_value" : {
				"value_1" : {"u8" : 8},
				"value_2" : {"u32" : 314},
				"value_3" : {"u64" : 2342}
			}}
			}
		},		
		{
			"name" : "raw_bytes",
			"value" : {
				"cl_type" : {"list_type" : {"inner" : {"simple_type" : "U8"}}},
				"value" : {"bytes_value" : "1wJD3Z0NZG/W3ygqj3qPoFpmKb7AHYAkw2EescH7n4Q="}
			}
		},
		{
			"name" : "raw_bytes_fixed",
			"value" : {
				"cl_type" : {"fixed_list_type" : {"inner" : {"simple_type" : "U8"}, "len" : 32}},
				"value" : {"bytes_value" : "1wJD3Z0NZG/W3ygqj3qPoFpmKb7AHYAkw2EescH7n4Q="}
			}
		}
	]`

	res, err := JsonStringToDeployArgs(jsonStr)
	assert.NoError(t, err)

	assert.Equal(t, 27, len(res))

	assert.Equal(t, state.CLType_BOOL, res[0].Value.GetClType().GetSimpleType())
	assert.Equal(t, true, res[0].GetValue().GetValue().GetBoolValue())

	assert.Equal(t, state.CLType_I32, res[1].Value.GetClType().GetSimpleType())
	assert.Equal(t, int32(314), res[1].GetValue().GetValue().GetI32())

	assert.Equal(t, state.CLType_I64, res[2].Value.GetClType().GetSimpleType())
	assert.Equal(t, int64(2342), res[2].GetValue().GetValue().GetI64())

	assert.Equal(t, state.CLType_U8, res[3].Value.GetClType().GetSimpleType())
	assert.Equal(t, int32(7), res[3].GetValue().GetValue().GetU8())

	assert.Equal(t, state.CLType_U32, res[4].Value.GetClType().GetSimpleType())
	assert.Equal(t, uint32(314), res[4].GetValue().GetValue().GetU32())

	assert.Equal(t, state.CLType_U64, res[5].Value.GetClType().GetSimpleType())
	assert.Equal(t, uint64(2342), res[5].GetValue().GetValue().GetU64())

	assert.Equal(t, state.CLType_U128, res[6].Value.GetClType().GetSimpleType())
	assert.Equal(t, "123456789101112131415161718", res[6].GetValue().GetValue().GetU128().GetValue())

	assert.Equal(t, state.CLType_U256, res[7].Value.GetClType().GetSimpleType())
	assert.Equal(t, "123456789101112131415161718", res[7].GetValue().GetValue().GetU256().GetValue())

	assert.Equal(t, state.CLType_U512, res[8].Value.GetClType().GetSimpleType())
	assert.Equal(t, "123456789101112131415161718", res[8].GetValue().GetValue().GetU512().GetValue())

	assert.Equal(t, state.CLType_UNIT, res[9].Value.GetClType().GetSimpleType())
	assert.Equal(t, "", res[9].GetValue().GetValue().GetUnit().String())

	assert.Equal(t, state.CLType_STRING, res[10].Value.GetClType().GetSimpleType())
	assert.Equal(t, "Hello, world!", res[10].GetValue().GetValue().GetStrValue())

	assert.Equal(t, state.CLType_KEY, res[11].Value.GetClType().GetSimpleType())
	assert.Equal(t, "d70243dd9d0d646fd6df282a8f7a8fa05a6629bec01d8024c3611eb1c1fb9f84", hex.EncodeToString(res[11].GetValue().GetValue().GetKey().GetAddress().GetAccount()))

	assert.Equal(t, state.CLType_KEY, res[12].Value.GetClType().GetSimpleType())
	assert.Equal(t, "d70243dd9d0d646fd6df282a8f7a8fa05a6629bec01d8024c3611eb1c1fb9f84", hex.EncodeToString(res[12].GetValue().GetValue().GetKey().GetHash().GetHash()))

	assert.Equal(t, state.CLType_KEY, res[13].Value.GetClType().GetSimpleType())
	assert.Equal(t, "d70243dd9d0d646fd6df282a8f7a8fa05a6629bec01d8024c3611eb1c1fb9f84", hex.EncodeToString(res[13].GetValue().GetValue().GetKey().GetUref().GetUref()))
	assert.Equal(t, state.Key_URef_READ_ADD, res[13].GetValue().GetValue().GetKey().GetUref().GetAccessRights())

	assert.Equal(t, state.CLType_UREF, res[14].Value.GetClType().GetSimpleType())
	assert.Equal(t, "d70243dd9d0d646fd6df282a8f7a8fa05a6629bec01d8024c3611eb1c1fb9f84", hex.EncodeToString(res[14].GetValue().GetValue().GetUref().GetUref()))
	assert.Equal(t, state.Key_URef_READ_ADD, res[14].GetValue().GetValue().GetUref().GetAccessRights())

	assert.Equal(t, state.CLType_U64, res[15].Value.GetClType().GetOptionType().GetInner().GetSimpleType())
	assert.Equal(t, uint64(0), res[15].GetValue().GetValue().GetOptionValue().GetValue().GetU64())

	assert.Equal(t, state.CLType_U64, res[16].Value.GetClType().GetOptionType().GetInner().GetSimpleType())
	assert.Equal(t, uint64(2342), res[16].GetValue().GetValue().GetOptionValue().GetValue().GetU64())

	assert.Equal(t, state.CLType_I32, res[17].Value.GetClType().GetListType().GetInner().GetSimpleType())
	assert.Equal(t, int32(0), res[17].GetValue().GetValue().GetListValue().GetValues()[0].GetI32())
	assert.Equal(t, int32(1), res[17].GetValue().GetValue().GetListValue().GetValues()[1].GetI32())
	assert.Equal(t, int32(2), res[17].GetValue().GetValue().GetListValue().GetValues()[2].GetI32())
	assert.Equal(t, int32(3), res[17].GetValue().GetValue().GetListValue().GetValues()[3].GetI32())

	assert.Equal(t, state.CLType_STRING, res[18].Value.GetClType().GetFixedListType().GetInner().GetSimpleType())
	assert.Equal(t, uint32(3), res[18].GetValue().GetValue().GetFixedListValue().GetLength())
	assert.Equal(t, "A", res[18].GetValue().GetValue().GetFixedListValue().GetValues()[0].GetStrValue())
	assert.Equal(t, "B", res[18].GetValue().GetValue().GetFixedListValue().GetValues()[1].GetStrValue())
	assert.Equal(t, "C", res[18].GetValue().GetValue().GetFixedListValue().GetValues()[2].GetStrValue())

	assert.Equal(t, state.CLType_BOOL, res[19].Value.GetClType().GetResultType().GetOk().GetSimpleType())
	assert.Equal(t, state.CLType_STRING, res[19].Value.GetClType().GetResultType().GetErr().GetSimpleType())
	assert.Equal(t, "Hello, world!", res[19].GetValue().GetValue().GetResultValue().GetErr().GetStrValue())

	assert.Equal(t, state.CLType_BOOL, res[20].Value.GetClType().GetResultType().GetOk().GetSimpleType())
	assert.Equal(t, state.CLType_STRING, res[20].Value.GetClType().GetResultType().GetErr().GetSimpleType())
	assert.Equal(t, true, res[20].GetValue().GetValue().GetResultValue().GetOk().GetBoolValue())

	assert.Equal(t, state.CLType_STRING, res[21].Value.GetClType().GetMapType().GetKey().GetSimpleType())
	assert.Equal(t, state.CLType_I32, res[21].Value.GetClType().GetMapType().GetValue().GetSimpleType())
	assert.Equal(t, "A", res[21].GetValue().GetValue().GetMapValue().GetValues()[0].GetKey().GetStrValue())
	assert.Equal(t, int32(0), res[21].GetValue().GetValue().GetMapValue().GetValues()[0].GetValue().GetI32())
	assert.Equal(t, "B", res[21].GetValue().GetValue().GetMapValue().GetValues()[1].GetKey().GetStrValue())
	assert.Equal(t, int32(1), res[21].GetValue().GetValue().GetMapValue().GetValues()[1].GetValue().GetI32())
	assert.Equal(t, "C", res[21].GetValue().GetValue().GetMapValue().GetValues()[2].GetKey().GetStrValue())
	assert.Equal(t, int32(2), res[21].GetValue().GetValue().GetMapValue().GetValues()[2].GetValue().GetI32())

	assert.Equal(t, state.CLType_U8, res[22].Value.GetClType().GetTuple1Type().GetType0().GetSimpleType())
	assert.Equal(t, int32(8), res[22].GetValue().GetValue().GetTuple1Value().GetValue_1().GetU8())

	assert.Equal(t, state.CLType_U8, res[23].Value.GetClType().GetTuple2Type().GetType0().GetSimpleType())
	assert.Equal(t, state.CLType_U32, res[23].Value.GetClType().GetTuple2Type().GetType1().GetSimpleType())
	assert.Equal(t, int32(8), res[23].GetValue().GetValue().GetTuple2Value().GetValue_1().GetU8())
	assert.Equal(t, uint32(314), res[23].GetValue().GetValue().GetTuple2Value().GetValue_2().GetU32())

	assert.Equal(t, state.CLType_U8, res[24].Value.GetClType().GetTuple3Type().GetType0().GetSimpleType())
	assert.Equal(t, state.CLType_U32, res[24].Value.GetClType().GetTuple3Type().GetType1().GetSimpleType())
	assert.Equal(t, state.CLType_U64, res[24].Value.GetClType().GetTuple3Type().GetType2().GetSimpleType())
	assert.Equal(t, int32(8), res[24].GetValue().GetValue().GetTuple3Value().GetValue_1().GetU8())
	assert.Equal(t, uint32(314), res[24].GetValue().GetValue().GetTuple3Value().GetValue_2().GetU32())
	assert.Equal(t, uint64(2342), res[24].GetValue().GetValue().GetTuple3Value().GetValue_3().GetU64())

	assert.Equal(t, state.CLType_U8, res[25].Value.GetClType().GetListType().GetInner().GetSimpleType())
	assert.Equal(t, "d70243dd9d0d646fd6df282a8f7a8fa05a6629bec01d8024c3611eb1c1fb9f84", EncodeToHexString(res[25].GetValue().GetValue().GetBytesValue()))

	assert.Equal(t, state.CLType_U8, res[26].Value.GetClType().GetFixedListType().GetInner().GetSimpleType())
	assert.Equal(t, "d70243dd9d0d646fd6df282a8f7a8fa05a6629bec01d8024c3611eb1c1fb9f84", EncodeToHexString(res[26].GetValue().GetValue().GetBytesValue()))
}
