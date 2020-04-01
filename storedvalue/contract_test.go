package storedvalue

import (
	"encoding/hex"
	"testing"

	"github.com/hdac-io/casperlabs-ee-grpc-go-util/protobuf/io/casperlabs/casper/consensus/state"
	"github.com/stretchr/testify/assert"
)

func TestContractAbi(t *testing.T) {
	bytes, err := hex.DecodeString("01000000000500000097000000645f643730323433646439643064363436666436646632383261386637613866613035613636323962656330316438303234633336313165623163316662396638345f643730323433646439643064363436666436646632383261386637613866613035613636323962656330316438303234633336313165623163316662396638345f3130303030303030303030303030303030303001000000000000000000000000000000000000000000000000000000000000000011000000706f735f626f6e64696e675f7075727365027cdb081c47a129b41273a1d2830f7f8481eae8380978e17cec5b4e4f9e1d0b68010711000000706f735f7061796d656e745f70757273650251f1ddda0933696150cf78fe7a2141653e6a841d2f4ecaaa915a299cb7a4d19c010711000000706f735f726577617264735f707572736502c32d411249f72f9da9d61c8e0d115f3000ce00d6889b8195b94bc020ba522b1b010756000000765f643730323433646439643064363436666436646632383261386637613866613035613636323962656330316438303234633336313165623163316662396638345f31303030303030303030303030303030303030010000000000000000000000000000000000000000000000000000000000000000010000000000000000000000")
	assert.NoError(t, err)

	var contract Contract
	contract, err, pos := contract.FromBytes(bytes)

	assert.NoError(t, err)
	assert.Equal(t, len(bytes), pos)

	assert.Equal(t,
		[]byte{0},
		contract.Body)

	assert.Equal(t, 5, len(contract.NamedKeys))

	assert.Equal(t, "d_d70243dd9d0d646fd6df282a8f7a8fa05a6629bec01d8024c3611eb1c1fb9f84_d70243dd9d0d646fd6df282a8f7a8fa05a6629bec01d8024c3611eb1c1fb9f84_1000000000000000000", contract.NamedKeys[0].Name)
	assert.Equal(t, make([]byte, 32), contract.NamedKeys[4].Key.Hash)
	assert.Equal(t, "pos_bonding_purse", contract.NamedKeys[1].Name)
	assert.Equal(t, "7cdb081c47a129b41273a1d2830f7f8481eae8380978e17cec5b4e4f9e1d0b68", hex.EncodeToString(contract.NamedKeys[1].Key.Uref.GetAddress()))
	assert.Equal(t, state.Key_URef_READ_ADD_WRITE, contract.NamedKeys[1].Key.Uref.GetAccessRights())
	assert.Equal(t, "pos_payment_purse", contract.NamedKeys[2].Name)
	assert.Equal(t, "51f1ddda0933696150cf78fe7a2141653e6a841d2f4ecaaa915a299cb7a4d19c", hex.EncodeToString(contract.NamedKeys[2].Key.Uref.GetAddress()))
	assert.Equal(t, state.Key_URef_READ_ADD_WRITE, contract.NamedKeys[2].Key.Uref.GetAccessRights())
	assert.Equal(t, "pos_rewards_purse", contract.NamedKeys[3].Name)
	assert.Equal(t, "c32d411249f72f9da9d61c8e0d115f3000ce00d6889b8195b94bc020ba522b1b", hex.EncodeToString(contract.NamedKeys[3].Key.Uref.GetAddress()))
	assert.Equal(t, state.Key_URef_READ_ADD_WRITE, contract.NamedKeys[3].Key.Uref.GetAccessRights())
	assert.Equal(t, "v_d70243dd9d0d646fd6df282a8f7a8fa05a6629bec01d8024c3611eb1c1fb9f84_1000000000000000000", contract.NamedKeys[4].Name)
	assert.Equal(t, make([]byte, 32), contract.NamedKeys[4].Key.Hash)

	assert.Equal(t,
		NewProtocolVersion(1, 0, 0),
		contract.ProtocolVersion)
}
