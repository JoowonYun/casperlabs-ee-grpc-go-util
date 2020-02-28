package util

import (
	"encoding/binary"
	"math/big"

	"github.com/hdac-io/casperlabs-ee-grpc-go-util/protobuf/io/casperlabs/casper/consensus/state"
)

const (
	Option     = 13
	List       = 14
	Fixed_list = 15
	Result     = 16
	Map        = 17
	Tuple1     = 18
	Tuple2     = 19
	Tuple3     = 20
	Any        = 21
)

func ToBytes(c *state.CLValue) []byte {
	res := AbiBytesToBytes(c.SerializedValue)

	tags := []byte{}
	switch c.ClType.GetVariants().(type) {
	case *state.CLType_SimpleType:
		tags = append(tags, byte(c.GetClType().GetSimpleType()))
	case *state.CLType_OptionType:
		tags = append(tags, []byte{Option, byte(c.GetClType().GetOptionType().GetInner().GetSimpleType())}...)
	case *state.CLType_ListType:
		tags = append(tags, []byte{List, byte(c.GetClType().GetListType().GetInner().GetSimpleType())}...)
	case *state.CLType_FixedListType:
		tags = append(tags, []byte{Fixed_list, byte(c.GetClType().GetFixedListType().GetInner().GetSimpleType())}...)
		tags = append(tags, c.SerializedValue[:4]...)
	case *state.CLType_ResultType:
	case *state.CLType_MapType:
	case *state.CLType_Tuple1Type:
	case *state.CLType_Tuple2Type:
	case *state.CLType_Tuple3Type:
	case *state.CLType_AnyType:
		tags = append(tags, byte(c.GetClType().GetSimpleType()))
	}
	res = append(res, tags...)

	return res
}

func ToValues(c *state.CLValue) *state.Value {
	value := &state.Value{}
	switch c.ClType.GetVariants().(type) {
	case *state.CLType_SimpleType:
		switch c.GetClType().GetSimpleType() {
		case state.CLType_BOOL:
			// value = &state.Value{Value: &state.}
		case state.CLType_I32:
			value = &state.Value{Value: &state.Value_IntValue{IntValue: int32(binary.LittleEndian.Uint32(c.SerializedValue))}}
		case state.CLType_I64:
			value = &state.Value{Value: &state.Value_LongValue{LongValue: binary.LittleEndian.Uint64(c.SerializedValue)}}
		case state.CLType_U8:
			// value = &state.Value{Value: &state.Value_}
		case state.CLType_U32:
			value = &state.Value{Value: &state.Value_IntValue{IntValue: int32(binary.LittleEndian.Uint32(c.SerializedValue))}}
		case state.CLType_U64:
			value = &state.Value{Value: &state.Value_LongValue{LongValue: binary.LittleEndian.Uint64(c.SerializedValue)}}
		case state.CLType_U128:
			value = &state.Value{Value: &state.Value_BigInt{BigInt: &state.BigInt{Value: fromByteToBigInt(c.SerializedValue).String(), BitWidth: 128}}}
		case state.CLType_U256:
			value = &state.Value{Value: &state.Value_BigInt{BigInt: &state.BigInt{Value: fromByteToBigInt(c.SerializedValue).String(), BitWidth: 256}}}
		case state.CLType_U512:
			value = &state.Value{Value: &state.Value_BigInt{BigInt: &state.BigInt{Value: fromByteToBigInt(c.SerializedValue).String(), BitWidth: 512}}}
		case state.CLType_UNIT:
			// value = &state.Value{}
		case state.CLType_STRING:
			value = &state.Value{Value: &state.Value_StringValue{StringValue: string(c.SerializedValue)}}
		case state.CLType_KEY:
			value = &state.Value{Value: &state.Value_Key{Key: fromByteToKey(c.SerializedValue)}}
		case state.CLType_UREF:
			// value = &state.Value{}
		}
	case *state.CLType_OptionType:
	case *state.CLType_ListType:
	case *state.CLType_FixedListType:
		value = &state.Value{Value: &state.Value_BytesValue{BytesValue: c.GetSerializedValue()}}
	case *state.CLType_ResultType:
	case *state.CLType_MapType:
	case *state.CLType_Tuple1Type:
	case *state.CLType_Tuple2Type:
	case *state.CLType_Tuple3Type:
	case *state.CLType_AnyType:
	}

	return value
}

func FromString(s string) *state.CLValue {
	return &state.CLValue{
		ClType:          &state.CLType{Variants: &state.CLType_SimpleType{SimpleType: state.CLType_STRING}},
		SerializedValue: AbiStringToBytes(s)}
}

func FromU512(value *big.Int) *state.CLValue {
	return &state.CLValue{
		ClType:          &state.CLType{Variants: &state.CLType_SimpleType{SimpleType: state.CLType_U512}},
		SerializedValue: AbiBigIntTobytes(value),
	}
}

func FromU32(value uint32) *state.CLValue {
	return &state.CLValue{
		ClType:          &state.CLType{Variants: &state.CLType_SimpleType{SimpleType: state.CLType_U32}},
		SerializedValue: AbiUint32ToBytes(value),
	}
}

func FromU64(value uint64) *state.CLValue {
	return &state.CLValue{
		ClType:          &state.CLType{Variants: &state.CLType_SimpleType{SimpleType: state.CLType_I64}},
		SerializedValue: AbiUint64ToBytes(value),
	}
}

func FromStringList(values []string) *state.CLValue {
	res := make([]byte, 4)
	binary.LittleEndian.PutUint32(res, uint32(len(values)))

	for _, value := range values {
		res = append(res, AbiStringToBytes(value)...)
	}

	return &state.CLValue{
		ClType:          &state.CLType{Variants: &state.CLType_ListType{ListType: &state.CLType_List{Inner: &state.CLType{Variants: &state.CLType_SimpleType{SimpleType: state.CLType_STRING}}}}},
		SerializedValue: res,
	}
}

func FromOption(value []byte, tag state.CLType_Simple) *state.CLValue {
	return &state.CLValue{
		ClType:          &state.CLType{Variants: &state.CLType_OptionType{OptionType: &state.CLType_Option{Inner: &state.CLType{Variants: &state.CLType_SimpleType{SimpleType: tag}}}}},
		SerializedValue: AbiOptionToBytes(value),
	}
}

func FromFixedList(value []byte, tag state.CLType_Simple) *state.CLValue {
	return &state.CLValue{
		ClType:          &state.CLType{Variants: &state.CLType_FixedListType{FixedListType: &state.CLType_FixedList{Inner: &state.CLType{Variants: &state.CLType_SimpleType{SimpleType: tag}}, Len: uint32(len(value))}}},
		SerializedValue: value,
	}
}

/**	TODO

func FromI32(value int32) ClValue {
	return NewClValue([]byte{}, []byte(byte(state.CLType_I32)))
}

func FromI64(value int64) ClValue {
	return NewClValue([]byte{}, []byte(byte(state.CLType_I64)))
}

func FromKey(key) ClValue
func FromURef(key) ClValue
**/

func fromByteToKey(src []byte) *state.Key {
	pos := 0
	var key *state.Key

	keyType := src[pos]
	pos++

	switch keyType {
	case KEY_ADDRESS:
		key = &state.Key{Value: &state.Key_Address_{Address: &state.Key_Address{Account: src[pos : pos+ACCOUNT_LEN]}}}
	case KEY_HASH:
		key = &state.Key{Value: &state.Key_Hash_{Hash: &state.Key_Hash{Hash: src[pos : pos+ACCOUNT_LEN]}}}
	case KEY_UREF:
		key = &state.Key{Value: &state.Key_Uref{Uref: &state.Key_URef{Uref: src[pos : pos+ACCOUNT_LEN], AccessRights: state.Key_URef_AccessRights(src[pos+ACCOUNT_LEN+1])}}}
		pos += 2
	case KEY_LOCAL:
		key = &state.Key{Value: &state.Key_Local_{Local: &state.Key_Local{Hash: src[pos : pos+ACCOUNT_LEN]}}}
	}

	return key
}

func fromByteToBigInt(src []byte) *big.Int {
	res := reverseBytes(src[1:])
	return new(big.Int).SetBytes(res)
}
