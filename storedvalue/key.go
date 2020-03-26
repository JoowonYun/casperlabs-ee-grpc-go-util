package storedvalue

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/hdac-io/casperlabs-ee-grpc-go-util/protobuf/io/casperlabs/casper/consensus/state"
)

const (
	KEY_ID_POS    = 0
	KEY_ID_LENGTH = 1
)

type KEY_ID int

const (
	KEY_ID_ACCOUNT KEY_ID = iota
	KEY_ID_HASH
	KEY_ID_UREF
	KEY_ID_LOCAL
)

type Key struct {
	KeyID   KEY_ID  `json:"key_id"`
	Account Account `json:"account"`
	Uref    URef    `json:"uref"`
	Hash    []byte  `json:"hash"`
	Local   []byte  `json:"local"`
}

func NewKeyFromURef(uref URef) Key {
	return Key{
		KeyID: KEY_ID_UREF,
		Uref:  uref,
	}
}

func NewKeyFromAccount(account Account) Key {
	return Key{
		KeyID:   KEY_ID_ACCOUNT,
		Account: account,
	}
}

func NewKeyFromHash(hash []byte) Key {
	return Key{
		KeyID: KEY_ID_HASH,
		Hash:  hash,
	}
}

func NewKeyFromLocal(local []byte) Key {
	return Key{
		KeyID: KEY_ID_LOCAL,
		Local: local,
	}
}

func (k Key) FromBytes(src []byte) (key Key, err error, pos int) {
	pos = KEY_ID_POS
	k.KeyID = KEY_ID(src[pos])
	pos += KEY_ID_LENGTH

	if len(src) < KEY_ID_LENGTH+ADDRESS_LENGTH {
		return Key{}, fmt.Errorf("Key length more than %d, but %d", KEY_ID_LENGTH+ADDRESS_LENGTH, len(src)), pos
	}

	switch k.KeyID {
	case KEY_ID_ACCOUNT:
		var account Account
		account, err, length := account.FromBytes(src[pos:])
		if err != nil {
			return Key{}, err, pos
		}
		k.Account = account
		pos += length
	case KEY_ID_HASH:
		k.Hash = src[pos : pos+ADDRESS_LENGTH]
		pos += ADDRESS_LENGTH
	case KEY_ID_UREF:
		var uref URef
		uref, err, length := uref.FromBytes(src[pos:])
		if err != nil {
			return Key{}, err, pos
		}
		k.Uref = uref
		pos += length
	case KEY_ID_LOCAL:
		k.Local = src[pos : pos+ADDRESS_LENGTH]
		pos += ADDRESS_LENGTH
	}

	return k, nil, pos
}

func (k Key) ToBytes() []byte {
	res := []byte{byte(k.KeyID)}

	switch k.KeyID {
	case KEY_ID_ACCOUNT:
		res = append(res, k.ToBytes()...)
	case KEY_ID_HASH:
		res = append(res, k.Hash...)
	case KEY_ID_UREF:
		res = append(res, k.Uref.ToBytes()...)
	case KEY_ID_LOCAL:
		res = append(res, k.Local...)
	}

	return res
}

func (k Key) ToStateValue() *state.Key {
	var value *state.Key
	switch k.KeyID {
	case KEY_ID_ACCOUNT:
		value = &state.Key{Value: &state.Key_Address_{Address: &state.Key_Address{Account: k.Account.ToBytes()}}}
	case KEY_ID_HASH:
		value = &state.Key{Value: &state.Key_Hash_{Hash: &state.Key_Hash{Hash: k.Hash}}}
	case KEY_ID_UREF:
		value = &state.Key{Value: &state.Key_Uref{Uref: k.Uref.ToStateValue()}}
	case KEY_ID_LOCAL:
		value = &state.Key{Value: &state.Key_Local_{Local: &state.Key_Local{Hash: k.Local}}}
	}

	return value
}

func (k Key) FromStateValue(key *state.Key) (Key, error) {
	switch key.GetValue().(type) {
	case *state.Key_Address_:
		var account Account
		account, err, _ := account.FromBytes(key.GetAddress().GetAccount())
		if err != nil {
			return Key{}, err
		}
		k = NewKeyFromAccount(account)
	case *state.Key_Hash_:
		k = NewKeyFromHash(key.GetHash().GetHash())
	case *state.Key_Uref:
		uref := NewURef(key.GetUref().GetUref(), key.GetUref().GetAccessRights())
		k = NewKeyFromURef(uref)
	case *state.Key_Local_:
		k = NewKeyFromLocal(key.GetLocal().GetHash())
	default:
		errors.New("Key data is invalid.")
	}

	return k, nil
}

type NamedKey struct {
	Name string `json:"name"`
	Key  Key    `json:"key"`
}

func NewNamedKey(name string, key Key) NamedKey {
	return NamedKey{
		Name: name,
		Key:  key,
	}
}

func (n NamedKey) FromBytes(src []byte) (namedKey NamedKey, err error, pos int) {
	pos = 0
	nameLength := int(binary.LittleEndian.Uint32(src[pos:SIZE_LENGTH]))
	pos += SIZE_LENGTH

	name := string(src[pos : pos+nameLength])
	pos += nameLength

	var key Key
	key, err, length := key.FromBytes(src[pos:])
	pos += length
	if err != nil {
		return NamedKey{}, err, pos
	}

	return NewNamedKey(name, key), nil, pos
}

func (n NamedKey) ToBytes() []byte {
	res := make([]byte, SIZE_LENGTH)
	binary.BigEndian.PutUint32(res, uint32(len(n.Name)))
	res = append(res, []byte(n.Name)...)

	res = append(res, n.Key.ToBytes()...)

	return res
}

func (n NamedKey) ToStateValue() *state.NamedKey {
	return &state.NamedKey{
		Name: n.Name,
		Key:  n.Key.ToStateValue(),
	}
}

func (n NamedKey) FromStateValue(state *state.NamedKey) (NamedKey, error) {
	var key Key
	key, err := key.FromStateValue(state.GetKey())
	if err != nil {
		return NamedKey{}, err
	}
	return NewNamedKey(state.GetName(), key), nil
}

type NamedKeys []*NamedKey

const (
	VALIDATOR_PREFIX_POS = iota
	VALIDATOR_ADDRESS_POS
	VALIDATOR_STAKE_POS
)

const (
	DELEGATOR_PREFIX_POS = iota
	DELEGATOR_VALIDATOR_POS
	DELEGATOR_DELEGATOR_POS
	DELEGATOR_STAKE_POS
)

const (
	VALIDATOR_PREFIX = "v"
	VALIDATOR_LENGTH = 3

	DELEGATE_PREFIX = "d"
	DELEGATE_LENGTH = 4
)

func (ns NamedKeys) GetAllValidators() map[string]string {
	validators := map[string]string{}

	for _, validator := range ns {
		values := strings.Split(validator.Name, "_")

		if values[VALIDATOR_PREFIX_POS] != VALIDATOR_PREFIX {
			continue
		}

		validators[values[VALIDATOR_ADDRESS_POS]] = values[VALIDATOR_STAKE_POS]
	}

	return validators
}

func (ns NamedKeys) GetValidatorStake(address []byte) string {
	validators := ns.GetAllValidators()
	addressStr := hex.EncodeToString(address)

	return validators[addressStr]
}

func (ns NamedKeys) GetDelegateFromValidator(address []byte) map[string]string {
	delegators := map[string]string{}
	addressStr := hex.EncodeToString(address)

	for _, delegator := range ns {
		values := strings.Split(delegator.Name, "_")

		if values[DELEGATOR_PREFIX_POS] != DELEGATE_PREFIX {
			continue
		}

		if values[DELEGATOR_VALIDATOR_POS] == addressStr {
			delegators[values[DELEGATOR_DELEGATOR_POS]] = values[DELEGATOR_STAKE_POS]
		}
	}

	return delegators
}

func (ns NamedKeys) GetDelegateFromDelegator(address []byte) map[string]string {
	delegators := map[string]string{}
	addressStr := hex.EncodeToString(address)

	for _, delegator := range ns {
		values := strings.Split(delegator.Name, "_")

		if values[DELEGATOR_PREFIX_POS] != DELEGATE_PREFIX {
			continue
		}

		if values[DELEGATOR_DELEGATOR_POS] == addressStr {
			delegators[values[DELEGATOR_VALIDATOR_POS]] = values[DELEGATOR_STAKE_POS]
		}
	}

	return delegators
}
