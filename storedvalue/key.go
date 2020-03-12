package storedvalue

import (
	"errors"
	"fmt"

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
