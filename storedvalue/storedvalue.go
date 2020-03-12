package storedvalue

import (
	"errors"
)

type STORED_VALUE_TYPE = int

const (
	TYPE_CL_VALUE = iota
	TYPE_ACCOUNT
	TYPE_CONTRACT
)

const (
	STORED_VALUE_TYPE_POS    = 0
	STORED_VALUE_TYPE_LENGTH = 1

	ADDRESS_LENGTH = 32
	SIZE_LENGTH    = 4
)

type StoredValue struct {
	Type     STORED_VALUE_TYPE `json:"type"`
	ClValue  CLValue           `json:"cl_value"`
	Account  Account           `json:"account"`
	Contract Contract          `json:"contract"`
}

func NewStoredValueFromClValue(clValue CLValue) StoredValue {
	return StoredValue{
		Type:    TYPE_CL_VALUE,
		ClValue: clValue,
	}
}

func NewStoredValueFromAccount(account Account) StoredValue {
	return StoredValue{
		Type:    TYPE_ACCOUNT,
		Account: account,
	}
}

func NewStoredValueFromContract(contract Contract) StoredValue {
	return StoredValue{
		Type:     TYPE_CONTRACT,
		Contract: contract,
	}
}

func (s StoredValue) FromBytes(src []byte) (storedvalue StoredValue, err error, pos int) {
	pos = STORED_VALUE_TYPE_POS
	pos += STORED_VALUE_TYPE_LENGTH

	switch STORED_VALUE_TYPE(src[STORED_VALUE_TYPE_POS]) {
	case TYPE_CL_VALUE:
		var clValue CLValue
		clValue, err, length := clValue.FromBytes(src[pos:])
		if err != nil {
			return StoredValue{}, err, pos
		}

		s.Type = TYPE_CL_VALUE
		s.ClValue = clValue
		pos += length
	case TYPE_ACCOUNT:
		var account Account
		account, err, length := account.FromBytes(src[pos:])
		if err != nil {
			return StoredValue{}, err, pos
		}

		s.Type = TYPE_ACCOUNT
		s.Account = account
		pos += length
	case TYPE_CONTRACT:
		var contract Contract
		contract, err, length := contract.FromBytes(src[pos:])
		if err != nil {
			return StoredValue{}, err, pos
		}

		s.Type = TYPE_CONTRACT
		s.Contract = contract
		pos += length
	default:
		return StoredValue{}, errors.New(`StoredValue data is invalid.`), pos
	}

	return s, nil, pos
}
