package storedvalue

import (
	"encoding/binary"
	"errors"
	"math/big"

	"github.com/hdac-io/casperlabs-ee-grpc-go-util/protobuf/io/casperlabs/casper/consensus/state"
)

const (
	TAG_LENGTH = 1
	TAG_INDEX  = 0

	UINT32_LENGTH      = 4
	INT32_LENGTH       = 4
	LONG_LENGTH        = 8
	BIGINT_SIZE_LENGTH = 1
	OPTION_SIZE_LENGTH = 1
)

type CL_TYPE_TAG int

const (
	TAG_BOOL CL_TYPE_TAG = iota
	TAG_I32
	TAG_I64
	TAG_U8
	TAG_U32
	TAG_U64
	TAG_U128
	TAG_U256
	TAG_U512
	TAG_UNIT
	TAG_STRING
	TAG_KEY
	TAG_UREF
	TAG_OPTION
	TAG_LIST
	TAG_FIXED_LIST
	TAG_RESULT
	TAG_MAP
	TAG_TUPLE1
	TAG_TUPLE2
	TAG_TUPLE3
	TAG_ANY
)

type CLValue struct {
	Bytes []byte        `json:"bytes"`
	Tags  []CL_TYPE_TAG `json:"tags"`
}

func NewClValue(bytes []byte, tags []CL_TYPE_TAG) CLValue {
	return CLValue{
		Bytes: bytes,
		Tags:  tags,
	}
}

func (c CLValue) FromBytes(src []byte) (clValue CLValue, err error, pos int) {
	valueLength := int(binary.LittleEndian.Uint32(src[:SIZE_LENGTH]))
	pos = SIZE_LENGTH

	serializedValue := src[pos : pos+valueLength]
	pos += valueLength
	clValue.Bytes = serializedValue

	switch CL_TYPE_TAG(src[pos]) {
	case TAG_BOOL:
		clValue.Tags = []CL_TYPE_TAG{CL_TYPE_TAG(src[pos])}
		pos += TAG_LENGTH
	case TAG_I32:
		clValue.Tags = []CL_TYPE_TAG{CL_TYPE_TAG(src[pos])}
		pos += TAG_LENGTH
	case TAG_I64:
		clValue.Tags = []CL_TYPE_TAG{CL_TYPE_TAG(src[pos])}
		pos += TAG_LENGTH
	case TAG_U8:
		clValue.Tags = []CL_TYPE_TAG{CL_TYPE_TAG(src[pos])}
		pos += TAG_LENGTH
	case TAG_U32:
		clValue.Tags = []CL_TYPE_TAG{CL_TYPE_TAG(src[pos])}
		pos += TAG_LENGTH
	case TAG_U64:
		clValue.Tags = []CL_TYPE_TAG{CL_TYPE_TAG(src[pos])}
		pos += TAG_LENGTH
	case TAG_U128:
		clValue.Tags = []CL_TYPE_TAG{CL_TYPE_TAG(src[pos])}
		pos += TAG_LENGTH
	case TAG_U256:
		clValue.Tags = []CL_TYPE_TAG{CL_TYPE_TAG(src[pos])}
		pos += TAG_LENGTH
	case TAG_U512:
		clValue.Tags = []CL_TYPE_TAG{CL_TYPE_TAG(src[pos])}
		pos += TAG_LENGTH
	case TAG_UNIT:
		clValue.Tags = []CL_TYPE_TAG{CL_TYPE_TAG(src[pos])}
		pos += TAG_LENGTH
	case TAG_STRING:
		clValue.Tags = []CL_TYPE_TAG{CL_TYPE_TAG(src[pos])}
		pos += TAG_LENGTH
	case TAG_KEY:
		clValue.Tags = []CL_TYPE_TAG{CL_TYPE_TAG(src[pos])}
		pos += TAG_LENGTH
	case TAG_UREF:
		clValue.Tags = []CL_TYPE_TAG{CL_TYPE_TAG(src[pos])}
		pos += TAG_LENGTH
	case TAG_OPTION:
		clValue.Tags = []CL_TYPE_TAG{CL_TYPE_TAG(src[pos]), CL_TYPE_TAG(src[pos+TAG_LENGTH])}
		pos += TAG_LENGTH + TAG_LENGTH
	case TAG_LIST:
		clValue.Tags = []CL_TYPE_TAG{CL_TYPE_TAG(src[pos]), CL_TYPE_TAG(src[pos+TAG_LENGTH])}
		pos += TAG_LENGTH + TAG_LENGTH
	case TAG_FIXED_LIST:
		clValue.Tags = []CL_TYPE_TAG{CL_TYPE_TAG(src[pos]), CL_TYPE_TAG(src[pos+TAG_LENGTH])}
		pos += TAG_LENGTH + TAG_LENGTH
		pos += UINT32_LENGTH
	case TAG_RESULT:
		clValue.Tags = []CL_TYPE_TAG{CL_TYPE_TAG(src[pos]), CL_TYPE_TAG(src[pos+TAG_LENGTH]), CL_TYPE_TAG(src[pos+TAG_LENGTH+TAG_LENGTH])}
		pos += TAG_LENGTH + TAG_LENGTH + TAG_LENGTH
	case TAG_MAP:
		clValue.Tags = []CL_TYPE_TAG{CL_TYPE_TAG(src[pos]), CL_TYPE_TAG(src[pos+TAG_LENGTH]), CL_TYPE_TAG(src[pos+TAG_LENGTH+TAG_LENGTH])}
		pos += TAG_LENGTH + TAG_LENGTH + TAG_LENGTH
	case TAG_TUPLE1:
		clValue.Tags = []CL_TYPE_TAG{CL_TYPE_TAG(src[pos]), CL_TYPE_TAG(src[pos+TAG_LENGTH])}
		pos += TAG_LENGTH + TAG_LENGTH
	case TAG_TUPLE2:
		clValue.Tags = []CL_TYPE_TAG{CL_TYPE_TAG(src[pos]), CL_TYPE_TAG(src[pos+TAG_LENGTH]), CL_TYPE_TAG(src[pos+TAG_LENGTH+TAG_LENGTH])}
		pos += TAG_LENGTH + TAG_LENGTH + TAG_LENGTH + TAG_LENGTH
	case TAG_TUPLE3:
		clValue.Tags = []CL_TYPE_TAG{CL_TYPE_TAG(src[pos]), CL_TYPE_TAG(src[pos+TAG_LENGTH]), CL_TYPE_TAG(src[pos+TAG_LENGTH+TAG_LENGTH]), CL_TYPE_TAG(src[pos+TAG_LENGTH+TAG_LENGTH+TAG_LENGTH])}
		pos += TAG_LENGTH + TAG_LENGTH + TAG_LENGTH + TAG_LENGTH + TAG_LENGTH
	case TAG_ANY:
		clValue.Tags = []CL_TYPE_TAG{CL_TYPE_TAG(src[pos])}
		pos += TAG_LENGTH
	}

	return clValue, nil, pos
}

func (c CLValue) ToBytes() []byte {
	res := make([]byte, 4)
	binary.LittleEndian.PutUint32(res, uint32(len(c.Bytes)))
	res = append(res, c.Bytes...)

	for _, tag := range c.Tags {
		res = append(res, byte(tag))
	}

	if c.Tags[TAG_INDEX] == TAG_FIXED_LIST {
		res = append(res, res[:SIZE_LENGTH]...)
	}

	return res
}

func (c CLValue) ToStateValues() *state.Value {
	value := &state.Value{}
	switch c.Tags[TAG_INDEX] {
	case TAG_BOOL:
	case TAG_I32:
		value = &state.Value{Value: &state.Value_IntValue{IntValue: int32(binary.LittleEndian.Uint32(c.Bytes))}}
	case TAG_I64:
		value = &state.Value{Value: &state.Value_LongValue{LongValue: binary.LittleEndian.Uint64(c.Bytes)}}
	case TAG_U8:
	case TAG_U32:
		value = &state.Value{Value: &state.Value_IntValue{IntValue: int32(binary.LittleEndian.Uint32(c.Bytes))}}
	case TAG_U64:
		value = &state.Value{Value: &state.Value_LongValue{LongValue: binary.LittleEndian.Uint64(c.Bytes)}}
	case TAG_U128:
		value = &state.Value{Value: &state.Value_BigInt{BigInt: &state.BigInt{Value: fromByteToBigInt(c.Bytes).String(), BitWidth: 128}}}
	case TAG_U256:
		value = &state.Value{Value: &state.Value_BigInt{BigInt: &state.BigInt{Value: fromByteToBigInt(c.Bytes).String(), BitWidth: 256}}}
	case TAG_U512:
		value = &state.Value{Value: &state.Value_BigInt{BigInt: &state.BigInt{Value: fromByteToBigInt(c.Bytes).String(), BitWidth: 512}}}
	case TAG_STRING:
		value = &state.Value{Value: &state.Value_StringValue{StringValue: string(c.Bytes)}}
	case TAG_KEY:
		var key Key
		key, _, _ = key.FromBytes(c.Bytes)
		value = &state.Value{Value: &state.Value_Key{Key: key.ToStateValue()}}
	case TAG_UREF:
		var uref URef
		uref, _, _ = uref.FromBytes(c.Bytes)
		value = &state.Value{Value: &state.Value_Key{Key: &state.Key{Value: &state.Key_Uref{Uref: uref.ToStateValue()}}}}
	case TAG_LIST:
		length := int(binary.LittleEndian.Uint32(c.Bytes[:SIZE_LENGTH]))
		pos := length
		switch c.Tags[TAG_LENGTH] {
		case TAG_I32, TAG_U32:
			var values []int32
			for i := 0; i < length; i++ {
				val := int32(binary.LittleEndian.Uint32(c.Bytes[pos : pos+INT32_LENGTH]))
				values = append(values, val)
				pos += INT32_LENGTH
			}

			value = &state.Value{Value: &state.Value_IntList{
				IntList: &state.IntList{Values: values},
			}}
		case TAG_STRING:
			var values []string
			for i := 0; i < length; i++ {
				strLength := binary.LittleEndian.Uint32(c.Bytes[pos : pos+SIZE_LENGTH])
				pos += SIZE_LENGTH
				values = append(values, string(c.Bytes[pos:strLength]))
				pos += int(strLength)
			}

			value = &state.Value{Value: &state.Value_StringList{
				StringList: &state.StringList{Values: values},
			}}
		}
	case TAG_FIXED_LIST:
		value = &state.Value{Value: &state.Value_BytesValue{BytesValue: c.Bytes}}
	// TODO
	case TAG_UNIT:

	case TAG_OPTION:

	case TAG_RESULT:

	case TAG_MAP:

	case TAG_TUPLE1:

	case TAG_TUPLE2:

	case TAG_TUPLE3:

	case TAG_ANY:

	}
	return value
}

func (c CLValue) FromStateValue(value *state.Value) (CLValue, error) {
	switch value.GetValue().(type) {
	case *state.Value_IntValue:
		res := make([]byte, INT32_LENGTH)
		binary.LittleEndian.PutUint32(res, uint32(value.GetIntValue()))

		c.Bytes = res
		c.Tags = []CL_TYPE_TAG{TAG_I32}
	case *state.Value_BytesValue:
		c.Bytes = value.GetBytesValue()
		c.Tags = []CL_TYPE_TAG{TAG_FIXED_LIST, TAG_U8}
	case *state.Value_IntList:
		res := make([]byte, SIZE_LENGTH)
		binary.LittleEndian.PutUint32(res, uint32(len(value.GetIntList().GetValues())))

		for _, intValue := range value.GetIntList().GetValues() {
			intBytes := make([]byte, INT32_LENGTH)
			binary.LittleEndian.PutUint32(intBytes, uint32(intValue))
			res = append(res, intBytes...)
		}

		c.Bytes = res
		c.Tags = []CL_TYPE_TAG{TAG_LIST, TAG_I32}
	case *state.Value_StringValue:
		res := make([]byte, SIZE_LENGTH)
		binary.LittleEndian.PutUint32(res, uint32(len(value.GetStringValue())))
		res = append(res, []byte(value.GetStringValue())...)

		c.Bytes = res
		c.Tags = []CL_TYPE_TAG{TAG_STRING}
	case *state.Value_Account:
		// stored value...
	case *state.Value_Contract:
		// stored value...
	case *state.Value_StringList:
		res := make([]byte, SIZE_LENGTH)
		binary.LittleEndian.PutUint32(res, uint32(len(value.GetStringList().GetValues())))

		for _, stringValue := range value.GetStringList().GetValues() {
			stringLengthBytes := make([]byte, SIZE_LENGTH)
			binary.LittleEndian.PutUint32(stringLengthBytes, uint32(len(stringValue)))

			res = append(res, stringLengthBytes...)
			res = append(res, []byte(stringValue)...)
		}

		c.Bytes = res
		c.Tags = []CL_TYPE_TAG{TAG_LIST, TAG_STRING}
	case *state.Value_NamedKey:
		var namedKey NamedKey
		namedKey, err := namedKey.FromStateValue(value.GetNamedKey())
		if err != nil {
			return CLValue{}, err
		}

		c.Bytes = namedKey.ToBytes()
		c.Tags = []CL_TYPE_TAG{TAG_TUPLE2, TAG_STRING, TAG_KEY}
	case *state.Value_BigInt:
		bigIntValue, ok := new(big.Int).SetString(value.GetBigInt().GetValue(), 10)
		if !ok {
			return CLValue{}, errors.New("Bigint data is invalid.")
		}
		bytes := reverseBytes(bigIntValue.Bytes())
		res := []byte{byte(len(bytes))}

		switch value.GetBigInt().GetBitWidth() {
		case 128:
			c.Tags = []CL_TYPE_TAG{TAG_U128}
		case 256:
			c.Tags = []CL_TYPE_TAG{TAG_U256}
		case 512:
			c.Tags = []CL_TYPE_TAG{TAG_U512}
		default:
			return CLValue{}, errors.New("Bigint data is invalid.")
		}

		c.Bytes = append(res, bytes...)
	case *state.Value_Key:
		var key Key
		key, err := key.FromStateValue(value.GetKey())
		if err != nil {
			return CLValue{}, err
		}

		c.Bytes = key.ToBytes()
		c.Tags = []CL_TYPE_TAG{TAG_KEY}
	case *state.Value_Unit:
		c.Bytes = []byte{}
		c.Tags = []CL_TYPE_TAG{TAG_UNIT}
	case *state.Value_LongValue:
		res := make([]byte, LONG_LENGTH)
		binary.LittleEndian.PutUint64(res, value.GetLongValue())
		c.Bytes = res

		c.Tags = []CL_TYPE_TAG{TAG_I64}
	default:
		return CLValue{}, errors.New("ClValue data is invalid.")
	}

	return c, nil
}

func (c CLValue) FromDeployArgValue(value *state.CLValueInstance_Value) (CLValue, error) {
	switch value.GetValue().(type) {
	case *state.CLValueInstance_Value_I32:
		res := make([]byte, INT32_LENGTH)
		binary.LittleEndian.PutUint32(res, uint32(value.GetI32()))
		c.Bytes = res
		c.Tags = []CL_TYPE_TAG{TAG_I32}
	case *state.CLValueInstance_Value_I64:
		res := make([]byte, LONG_LENGTH)
		binary.LittleEndian.PutUint64(res, uint64(value.GetI64()))
		c.Bytes = res
		c.Tags = []CL_TYPE_TAG{TAG_I64}
	case *state.CLValueInstance_Value_U8:
		c.Bytes = []byte{byte(value.GetU8())}
		c.Tags = []CL_TYPE_TAG{TAG_U8}
	case *state.CLValueInstance_Value_U32:
		res := make([]byte, INT32_LENGTH)
		binary.LittleEndian.PutUint32(res, uint32(value.GetU32()))
		c.Bytes = res
		c.Tags = []CL_TYPE_TAG{TAG_U32}
	case *state.CLValueInstance_Value_U64:
		res := make([]byte, LONG_LENGTH)
		binary.LittleEndian.PutUint64(res, uint64(value.GetU64()))
		c.Bytes = res
		c.Tags = []CL_TYPE_TAG{TAG_U64}
	case *state.CLValueInstance_Value_U128:
		bigIntValue, ok := new(big.Int).SetString(value.GetU512().GetValue(), 10)
		if !ok {
			return CLValue{}, errors.New("Bigint data is invalid.")
		}
		bytes := reverseBytes(bigIntValue.Bytes())
		res := []byte{byte(len(bytes))}
		c.Tags = []CL_TYPE_TAG{TAG_U128}
		c.Bytes = append(res, bytes...)
	case *state.CLValueInstance_Value_U256:
		bigIntValue, ok := new(big.Int).SetString(value.GetU512().GetValue(), 10)
		if !ok {
			return CLValue{}, errors.New("Bigint data is invalid.")
		}
		bytes := reverseBytes(bigIntValue.Bytes())
		res := []byte{byte(len(bytes))}
		c.Tags = []CL_TYPE_TAG{TAG_U256}
		c.Bytes = append(res, bytes...)
	case *state.CLValueInstance_Value_U512:
		bigIntValue, ok := new(big.Int).SetString(value.GetU512().GetValue(), 10)
		if !ok {
			return CLValue{}, errors.New("Bigint data is invalid.")
		}
		bytes := reverseBytes(bigIntValue.Bytes())
		res := []byte{byte(len(bytes))}
		c.Tags = []CL_TYPE_TAG{TAG_U512}
		c.Bytes = append(res, bytes...)
	case *state.CLValueInstance_Value_StrValue:
		res := make([]byte, SIZE_LENGTH)
		binary.LittleEndian.PutUint32(res, uint32(len(value.GetStrValue())))
		res = append(res, []byte(value.GetStrValue())...)
		c.Bytes = res
		c.Tags = []CL_TYPE_TAG{TAG_STRING}
	case *state.CLValueInstance_Value_Key:
		var key Key
		key, err := key.FromStateValue(value.GetKey())
		if err != nil {
			return CLValue{}, err
		}
		c.Bytes = key.ToBytes()
		c.Tags = []CL_TYPE_TAG{TAG_KEY}
	case *state.CLValueInstance_Value_Uref:
		var uref URef
		uref, err := uref.FromStateValue(value.GetUref())
		if err != nil {
			return CLValue{}, err
		}
		c.Bytes = uref.ToBytes()
		c.Tags = []CL_TYPE_TAG{TAG_UREF}
	case *state.CLValueInstance_Value_OptionValue:
		clValue, err := c.FromDeployArgValue(value.GetOptionValue().GetValue())
		if err != nil {
			return CLValue{}, err
		}

		res := make([]byte, OPTION_SIZE_LENGTH)
		if len(clValue.Bytes) > 0 {
			res[0] = 1
			res = append(res, clValue.Bytes...)
		}
		c.Bytes = res
		c.Tags = append([]CL_TYPE_TAG{TAG_OPTION}, clValue.Tags...)
	case *state.CLValueInstance_Value_ListValue:
		c.Tags = []CL_TYPE_TAG{TAG_LIST}
		res := make([]byte, SIZE_LENGTH)
		binary.LittleEndian.PutUint32(res, uint32(len(value.GetListValue().GetValues())))
		for idx, val := range value.GetListValue().GetValues() {
			clValue, err := c.FromDeployArgValue(val)
			if err != nil {
				return CLValue{}, err
			}

			c.Bytes = append(res, clValue.Bytes...)
			if idx == 0 {
				c.Tags = append(c.Tags, clValue.Tags...)
			}
		}
	case *state.CLValueInstance_Value_FixedListValue:
		c.Tags = []CL_TYPE_TAG{TAG_FIXED_LIST}
		res := make([]byte, SIZE_LENGTH)
		binary.LittleEndian.PutUint32(res, uint32(len(value.GetListValue().GetValues())))
		for idx, val := range value.GetFixedListValue().GetValues() {
			clValue, err := c.FromDeployArgValue(val)
			if err != nil {
				return CLValue{}, err
			}

			c.Bytes = append(res, clValue.Bytes...)
			if idx == 0 {
				c.Tags = append(c.Tags, clValue.Tags...)
			}
		}
	case *state.CLValueInstance_Value_MapValue:
		c.Tags = []CL_TYPE_TAG{TAG_MAP}
		res := make([]byte, SIZE_LENGTH)
		binary.LittleEndian.PutUint32(res, uint32(len(value.GetMapValue().GetValues())))
		for idx, val := range value.GetMapValue().GetValues() {
			clValueKey, err := c.FromDeployArgValue(val.GetKey())
			if err != nil {
				return CLValue{}, err
			}
			c.Bytes = append(res, clValueKey.Bytes...)

			clValueValue, err := c.FromDeployArgValue(val.GetValue())
			if err != nil {
				return CLValue{}, err
			}
			c.Bytes = append(res, clValueValue.Bytes...)

			if idx == 0 {
				c.Tags = append(c.Tags, clValueKey.Tags...)
				c.Tags = append(c.Tags, clValueValue.Tags...)
			}
		}
	// TODO
	case *state.CLValueInstance_Value_Unit:
	case *state.CLValueInstance_Value_ResultValue:
	case *state.CLValueInstance_Value_Tuple1Value:
	case *state.CLValueInstance_Value_Tuple2Value:
	case *state.CLValueInstance_Value_Tuple3Value:
	case *state.CLValueInstance_Value_BytesValue:
		c.Bytes = value.GetBytesValue()
		c.Tags = []CL_TYPE_TAG{TAG_FIXED_LIST, TAG_U8}
	default:
		return CLValue{}, errors.New("ClValue data is invalid.")
	}

	return c, nil
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

func fromByteToBigInt(src []byte) *big.Int {
	res := reverseBytes(src[BIGINT_SIZE_LENGTH:])
	return new(big.Int).SetBytes(res)
}
