package util

import (
	"encoding/binary"

	"github.com/hdac-io/casperlabs-ee-grpc-go-util/protobuf/io/casperlabs/casper/consensus/state"
)

const (
	StoredValue_Type_ClValue = iota
	StoredValue_Type_Account
	StoredValue_Type_Contract
)

const (
	TYPE = 0
)

func MarshalStoreValue(src *state.StoredValue) []byte {
	return []byte{}
}

func UnmarshalStoreValue(src []byte) *state.StoredValue {
	var res *state.StoredValue
	switch src[TYPE] {
	case StoredValue_Type_ClValue:
		res = &state.StoredValue{Variants: &state.StoredValue_ClValue{ClValue: UnmarsahlClValue(src[1:])}}
	case StoredValue_Type_Account:
		res = &state.StoredValue{Variants: &state.StoredValue_Account{Account: UnmarsahlAccount(src[1:])}}
	case StoredValue_Type_Contract:
		res = &state.StoredValue{Variants: &state.StoredValue_Contract{Contract: UnmarsahlContract(src[1:])}}
	default:

	}
	return res
}

func UnmarsahlClValue(src []byte) *state.CLValue {
	valueLength := binary.LittleEndian.Uint32(src[:UINT32_LEN])

	serializedValue := src[UINT32_LEN : UINT32_LEN+valueLength]
	clType := &state.CLType{}
	switch state.CLType_Simple(src[UINT32_LEN+valueLength]) {
	case state.CLType_BOOL:

	case state.CLType_I32:
		clType = &state.CLType{Variants: &state.CLType_SimpleType{SimpleType: state.CLType_U32}}
	case state.CLType_I64:
		clType = &state.CLType{Variants: &state.CLType_SimpleType{SimpleType: state.CLType_U64}}
	case state.CLType_U8:

	case state.CLType_U32:
		clType = &state.CLType{Variants: &state.CLType_SimpleType{SimpleType: state.CLType_U32}}
	case state.CLType_U64:
		clType = &state.CLType{Variants: &state.CLType_SimpleType{SimpleType: state.CLType_U64}}
	case state.CLType_U128:

	case state.CLType_U256:

	case state.CLType_U512:
		clType = &state.CLType{Variants: &state.CLType_SimpleType{SimpleType: state.CLType_U512}}
	case state.CLType_UNIT:

	case state.CLType_STRING:

	case state.CLType_KEY:
		clType = &state.CLType{Variants: &state.CLType_SimpleType{SimpleType: state.CLType_KEY}}
	case state.CLType_UREF:
	default:
		switch int32(src[UINT32_LEN+valueLength-1]) {
		case Option:

		case List:

		case Fixed_list:

		case Result:

		case Map:

		case Tuple1:

		case Tuple2:

		case Tuple3:

		case Any:
		}
	}
	return &state.CLValue{
		SerializedValue: serializedValue,
		ClType:          clType,
	}
}

const (
	ACCOUNT_LEN = 32

	UINT32_LEN = 4
)

const (
	KEY_ADDRESS = iota
	KEY_HASH
	KEY_UREF
	KEY_LOCAL
)

func UnmarsahlAccount(src []byte) *state.Account {
	pos := ACCOUNT_LEN

	// NamedKeys
	namedKeys := []*state.NamedKey{}
	namedKeysSizeBytes := src[pos : pos+UINT32_LEN]
	pos += UINT32_LEN
	namedKeysSize := int(binary.LittleEndian.Uint32(namedKeysSizeBytes))
	for i := 0; i < namedKeysSize; i++ {
		nameSize := int(binary.LittleEndian.Uint32(src[pos : pos+UINT32_LEN]))
		pos += UINT32_LEN

		name := string(src[pos : pos+nameSize])
		pos += nameSize

		key := fromByteToKey(src[pos:])
		if src[pos] == KEY_UREF {
			pos += 2
		}
		pos += (ACCOUNT_LEN + 1)

		namedKey := &state.NamedKey{
			Name: name,
			Key:  key,
		}

		namedKeys = append(namedKeys, namedKey)
	}

	// Purse ID
	purseId := &state.Key_URef{Uref: src[pos : pos+ACCOUNT_LEN], AccessRights: state.Key_URef_AccessRights(src[pos+ACCOUNT_LEN+1])}
	pos += (ACCOUNT_LEN + 2)

	// Associate Key
	associatedKeys := []*state.Account_AssociatedKey{}
	associateKeySize := int(src[pos])
	pos++

	for i := 0; i < associateKeySize; i++ {
		weight := binary.BigEndian.Uint32(src[pos+ACCOUNT_LEN : pos+ACCOUNT_LEN+UINT32_LEN])
		associatedKey := &state.Account_AssociatedKey{PublicKey: src[pos : pos+ACCOUNT_LEN], Weight: weight}
		pos += (ACCOUNT_LEN + UINT32_LEN)

		associatedKeys = append(associatedKeys, associatedKey)
	}

	// ActionThresholds
	actionThresholds := &state.Account_ActionThresholds{DeploymentThreshold: uint32(src[pos]), KeyManagementThreshold: uint32(src[pos+1])}

	return &state.Account{
		PublicKey:        src[:ACCOUNT_LEN],
		NamedKeys:        namedKeys,
		PurseId:          purseId,
		AssociatedKeys:   associatedKeys,
		ActionThresholds: actionThresholds,
	}
}

func UnmarsahlContract(src []byte) *state.Contract {
	return &state.Contract{}
}
