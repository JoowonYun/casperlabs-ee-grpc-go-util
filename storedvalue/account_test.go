package storedvalue

import (
	"encoding/hex"
	"testing"

	"github.com/hdac-io/casperlabs-ee-grpc-go-util/protobuf/io/casperlabs/casper/consensus/state"
	"github.com/stretchr/testify/assert"
)

func TestAccountAbi(t *testing.T) {
	bytes, err := hex.DecodeString("000000000000000000000000000000000000000000000000000000000000000002000000040000006d696e74026cc261631cd46c959857de59ee0a5f61099457300012267bbde569820625c7f8010703000000706f7302bb0d91b8604970a269bf96ac55de5fa416135e2837d88a0bac938e2eca2d0fe20107c64bb588f000ca1b265648b53dae77559545fe93c24fd5f940602febdcdc91680107010000000000000000000000000000000000000000000000000000000000000000000000010101")
	assert.NoError(t, err)
	var account Account
	account, err, pos := account.FromBytes(bytes)

	assert.NoError(t, err)
	assert.Equal(t, len(bytes), pos)
	assert.Equal(t,
		[]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		account.PublicKey)

	mintAddressBytes, err := hex.DecodeString("6cc261631cd46c959857de59ee0a5f61099457300012267bbde569820625c7f8")
	assert.NoError(t, err)
	assert.Equal(t,
		"mint",
		account.NamedKeys[0].Name)
	assert.Equal(t,
		mintAddressBytes,
		account.NamedKeys[0].Key.Uref.Address)
	assert.Equal(t,
		state.Key_URef_READ_ADD_WRITE,
		account.NamedKeys[0].Key.Uref.AccessRights)

	posAddressBytes, err := hex.DecodeString("bb0d91b8604970a269bf96ac55de5fa416135e2837d88a0bac938e2eca2d0fe2")
	assert.NoError(t, err)
	assert.Equal(t,
		"pos",
		account.NamedKeys[1].Name)
	assert.Equal(t,
		posAddressBytes,
		account.NamedKeys[1].Key.Uref.Address)
	assert.Equal(t,
		state.Key_URef_READ_ADD_WRITE,
		account.NamedKeys[1].Key.Uref.AccessRights)

	purseIdAddress, err := hex.DecodeString("c64bb588f000ca1b265648b53dae77559545fe93c24fd5f940602febdcdc9168")
	assert.NoError(t, err)
	assert.Equal(t,
		purseIdAddress,
		account.PurseId.Address)
	assert.Equal(t,
		state.Key_URef_READ_ADD_WRITE,
		account.PurseId.AccessRights)

	assert.Equal(t,
		make([]byte, 32),
		account.AssociatedKeys[0].PublicKey)
	assert.Equal(t,
		uint32(1),
		account.AssociatedKeys[0].Weight)
	assert.Equal(t,
		uint32(1),
		account.ActionThresholds.DeploymentThreshold)
	assert.Equal(t,
		uint32(1),
		account.ActionThresholds.KeyManagementThreshold)
}
