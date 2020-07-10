package storedvalue

import (
	"fmt"

	"github.com/hdac-io/casperlabs-ee-grpc-go-util/protobuf/io/casperlabs/casper/consensus/state"
)

const (
	UREF_ACCESS_RIGHTS_SERIALIZED_LENGTH = 1
)

type URef struct {
	Address      []byte                      `json:"address"`
	AccessRights state.Key_URef_AccessRights `json:"access_rights"`
}

func NewURef(address []byte, accessRights state.Key_URef_AccessRights) URef {
	return URef{
		Address:      address,
		AccessRights: accessRights}
}

func (u URef) GetAddress() []byte {
	return u.Address
}

func (u URef) GetAccessRights() state.Key_URef_AccessRights {
	return u.AccessRights
}

func (u URef) ToBytes() []byte {
	res := u.Address

	res = append(res, byte(u.AccessRights))
	return res
}

func (u URef) FromBytes(src []byte) (uref URef, err error, pos int) {
	if len(src) < ADDRESS_LENGTH+UREF_ACCESS_RIGHTS_SERIALIZED_LENGTH {
		return URef{}, fmt.Errorf("URef bytes more than %d, but %d", ADDRESS_LENGTH+UREF_ACCESS_RIGHTS_SERIALIZED_LENGTH, len(src)), pos
	}

	u.Address = src[:ADDRESS_LENGTH]
	pos = ADDRESS_LENGTH
	u.AccessRights = state.Key_URef_AccessRights(src[pos])
	pos += UREF_ACCESS_RIGHTS_SERIALIZED_LENGTH

	return u, nil, pos
}

func (u URef) ToStateValue() *state.Key_URef {
	return &state.Key_URef{
		Uref:         u.Address,
		AccessRights: u.AccessRights,
	}
}

func (u URef) FromStateValue(uref *state.Key_URef) (URef, error) {
	return NewURef(uref.GetUref(), u.AccessRights), nil
}

func (u URef) ToCLInstanceValue() *state.CLValueInstance_Value {
	return &state.CLValueInstance_Value{
		Value: &state.CLValueInstance_Value_Uref{
			Uref: u.ToStateValue(),
		},
	}
}
