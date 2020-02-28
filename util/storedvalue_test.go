package util

import (
	"testing"

	"github.com/hdac-io/casperlabs-ee-grpc-go-util/protobuf/io/casperlabs/casper/consensus/state"
	"github.com/stretchr/testify/assert"
)

func TestAccountAbi(t *testing.T) {
	bytes := DecodeHexString("01000000000000000000000000000000000000000000000000000000000000000002000000040000006d696e74026cc261631cd46c959857de59ee0a5f61099457300012267bbde569820625c7f8010703000000706f7302bb0d91b8604970a269bf96ac55de5fa416135e2837d88a0bac938e2eca2d0fe20107c64bb588f000ca1b265648b53dae77559545fe93c24fd5f940602febdcdc91680107010000000000000000000000000000000000000000000000000000000000000000000000010101")
	res := UnmarshalStoreValue(bytes)

	accout := &state.StoredValue{Variants: &state.StoredValue_Account{
		Account: &state.Account{
			PublicKey: make([]byte, 32),
			PurseId:   &state.Key_URef{Uref: DecodeHexString("c64bb588f000ca1b265648b53dae77559545fe93c24fd5f940602febdcdc9168"), AccessRights: state.Key_URef_READ_ADD_WRITE},
			NamedKeys: []*state.NamedKey{
				&state.NamedKey{Name: "mint", Key: &state.Key{Value: &state.Key_Uref{Uref: &state.Key_URef{Uref: DecodeHexString("6cc261631cd46c959857de59ee0a5f61099457300012267bbde569820625c7f8"), AccessRights: state.Key_URef_READ_ADD_WRITE}}}},
				&state.NamedKey{Name: "pos", Key: &state.Key{Value: &state.Key_Uref{Uref: &state.Key_URef{Uref: DecodeHexString("bb0d91b8604970a269bf96ac55de5fa416135e2837d88a0bac938e2eca2d0fe2"), AccessRights: state.Key_URef_READ_ADD_WRITE}}}},
			},
			AssociatedKeys: []*state.Account_AssociatedKey{
				&state.Account_AssociatedKey{PublicKey: make([]byte, 32), Weight: 1},
			},
			ActionThresholds: &state.Account_ActionThresholds{DeploymentThreshold: uint32(1), KeyManagementThreshold: uint32(1)},
		}}}

	assert.Equal(t, res, accout)
}
