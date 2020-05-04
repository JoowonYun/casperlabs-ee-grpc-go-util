package storedvalue

import (
	"encoding/hex"
	"testing"

	"github.com/hdac-io/casperlabs-ee-grpc-go-util/protobuf/io/casperlabs/casper/consensus/state"
	"github.com/stretchr/testify/assert"
)

func TestKeyHashToBytes(t *testing.T) {
	key := NewKeyFromHash(make([]byte, 32))

	res := key.ToBytes()

	assert.Equal(
		t,
		[]byte{
			1,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		},
		res)
}

func TestKeyLocalToBytes(t *testing.T) {
	key := NewKeyFromLocal(make([]byte, 32))

	res := key.ToBytes()

	assert.Equal(
		t,
		[]byte{
			3,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		},
		res)
}

func TestKeyUrefToBytes(t *testing.T) {
	uref := NewURef(make([]byte, 32), state.Key_URef_READ_ADD_WRITE)
	key := NewKeyFromURef(uref)

	res := key.ToBytes()

	assert.Equal(
		t,
		[]byte{
			2,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			7,
		},
		res)
}

func TestKeyHashFromBytes(t *testing.T) {
	src := []byte{
		1,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}
	var k Key
	k, err, pos := k.FromBytes(src)

	assert.NoError(t, err)
	assert.Equal(
		t,
		KEY_ID_HASH,
		k.KeyID,
	)
	assert.Equal(
		t,
		make([]byte, 32),
		k.Hash,
	)
	assert.Equal(t, len(src), pos)
}

func TestKeyLocalFromBytes(t *testing.T) {
	src := []byte{
		3,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}
	var k Key
	k, err, pos := k.FromBytes(src)

	assert.NoError(t, err)
	assert.Equal(
		t,
		KEY_ID_LOCAL,
		k.KeyID,
	)
	assert.Equal(
		t,
		make([]byte, 32),
		k.Local,
	)
	assert.Equal(t, len(src), pos)
}

func TestKeyUrefFromBytes(t *testing.T) {
	src := []byte{
		2,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		7,
	}
	var k Key
	k, err, pos := k.FromBytes(src)

	assert.NoError(t, err)
	assert.Equal(
		t,
		KEY_ID_UREF,
		k.KeyID,
	)
	assert.Equal(
		t,
		make([]byte, 32),
		k.Uref.Address,
	)
	assert.Equal(
		t,
		state.Key_URef_READ_ADD_WRITE,
		k.Uref.AccessRights,
	)
	assert.Equal(t, len(src), pos)
}

func TestKeyFromBytesError(t *testing.T) {
	src := []byte{
		3,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}
	var k Key
	k, err, pos := k.FromBytes(src)

	assert.Error(t, err)
	assert.NotEqual(t, len(src), pos)
}

func TestNamedKeysGetAllValidators(t *testing.T) {
	namedkeys := NamedKeys{
		NamedKey{Name: "v_d70243dd9d0d646fd6df282a8f7a8fa05a6629bec01d8024c3611eb1c1fb9f84_1000000000000000000", Key: Key{}},
		NamedKey{Name: "pos_bonding_purse", Key: Key{}},
		NamedKey{Name: "v_51f1ddda0933696150cf78fe7a2141653e6a841d2f4ecaaa915a299cb7a4d19c_2000000000000000000", Key: Key{}},
		NamedKey{Name: "d_d70243dd9d0d646fd6df282a8f7a8fa05a6629bec01d8024c3611eb1c1fb9f84_d70243dd9d0d646fd6df282a8f7a8fa05a6629bec01d8024c3611eb1c1fb9f84_1000000000000000000", Key: Key{}},
	}

	validators := namedkeys.GetAllValidators()

	assert.Equal(t, 2, len(validators))
	assert.Equal(t, "1000000000000000000", validators["d70243dd9d0d646fd6df282a8f7a8fa05a6629bec01d8024c3611eb1c1fb9f84"])
	assert.Equal(t, "2000000000000000000", validators["51f1ddda0933696150cf78fe7a2141653e6a841d2f4ecaaa915a299cb7a4d19c"])
}

func TestNamedKeysGetValidatorStake(t *testing.T) {
	namedkeys := NamedKeys{
		NamedKey{Name: "v_d70243dd9d0d646fd6df282a8f7a8fa05a6629bec01d8024c3611eb1c1fb9f84_1000000000000000000", Key: Key{}},
		NamedKey{Name: "pos_bonding_purse", Key: Key{}},
		NamedKey{Name: "v_51f1ddda0933696150cf78fe7a2141653e6a841d2f4ecaaa915a299cb7a4d19c_2000000000000000000", Key: Key{}},
		NamedKey{Name: "d_d70243dd9d0d646fd6df282a8f7a8fa05a6629bec01d8024c3611eb1c1fb9f84_d70243dd9d0d646fd6df282a8f7a8fa05a6629bec01d8024c3611eb1c1fb9f84_1000000000000000000", Key: Key{}},
	}

	address, err := hex.DecodeString("d70243dd9d0d646fd6df282a8f7a8fa05a6629bec01d8024c3611eb1c1fb9f84")
	assert.NoError(t, err)

	stake := namedkeys.GetValidatorStake(address)
	assert.Equal(t, "1000000000000000000", stake)

	address, err = hex.DecodeString("51f1ddda0933696150cf78fe7a2141653e6a841d2f4ecaaa915a299cb7a4d19c")
	assert.NoError(t, err)

	stake = namedkeys.GetValidatorStake(address)
	assert.Equal(t, "2000000000000000000", stake)
}

func TestNamedKeysGetDelegateFromValidator(t *testing.T) {
	namedkeys := NamedKeys{
		NamedKey{Name: "d_d70243dd9d0d646fd6df282a8f7a8fa05a6629bec01d8024c3611eb1c1fb9f84_d70243dd9d0d646fd6df282a8f7a8fa05a6629bec01d8024c3611eb1c1fb9f84_1000000000000000000", Key: Key{}},
		NamedKey{Name: "v_d70243dd9d0d646fd6df282a8f7a8fa05a6629bec01d8024c3611eb1c1fb9f84_6000000000000000000", Key: Key{}},
		NamedKey{Name: "pos_bonding_purse", Key: Key{}},
		NamedKey{Name: "d_d70243dd9d0d646fd6df282a8f7a8fa05a6629bec01d8024c3611eb1c1fb9f84_51f1ddda0933696150cf78fe7a2141653e6a841d2f4ecaaa915a299cb7a4d19c_2000000000000000000", Key: Key{}},
		NamedKey{Name: "d_d70243dd9d0d646fd6df282a8f7a8fa05a6629bec01d8024c3611eb1c1fb9f84_93236a9263d2ac6198c5ed211774c745d5dc62a910cb84276f8a7c4959208915_3000000000000000000", Key: Key{}},
	}

	address, err := hex.DecodeString("d70243dd9d0d646fd6df282a8f7a8fa05a6629bec01d8024c3611eb1c1fb9f84")
	assert.NoError(t, err)

	delegators := namedkeys.GetDelegateFromValidator(address)

	assert.Equal(t, 1, len(delegators))
	assert.Equal(t, "1000000000000000000", delegators["d70243dd9d0d646fd6df282a8f7a8fa05a6629bec01d8024c3611eb1c1fb9f84"])
}

func TestNamedKeysGetDelegateFromDelegators(t *testing.T) {
	namedkeys := NamedKeys{
		NamedKey{Name: "d_d70243dd9d0d646fd6df282a8f7a8fa05a6629bec01d8024c3611eb1c1fb9f84_d70243dd9d0d646fd6df282a8f7a8fa05a6629bec01d8024c3611eb1c1fb9f84_1000000000000000000", Key: Key{}},
		NamedKey{Name: "v_d70243dd9d0d646fd6df282a8f7a8fa05a6629bec01d8024c3611eb1c1fb9f84_4000000000000000000", Key: Key{}},
		NamedKey{Name: "d_93236a9263d2ac6198c5ed211774c745d5dc62a910cb84276f8a7c4959208915_93236a9263d2ac6198c5ed211774c745d5dc62a910cb84276f8a7c4959208915_1000000000000000000", Key: Key{}},
		NamedKey{Name: "v_93236a9263d2ac6198c5ed211774c745d5dc62a910cb84276f8a7c4959208915_3000000000000000000", Key: Key{}},
		NamedKey{Name: "pos_bonding_purse", Key: Key{}},
		NamedKey{Name: "d_93236a9263d2ac6198c5ed211774c745d5dc62a910cb84276f8a7c4959208915_51f1ddda0933696150cf78fe7a2141653e6a841d2f4ecaaa915a299cb7a4d19c_2000000000000000000", Key: Key{}},
		NamedKey{Name: "d_d70243dd9d0d646fd6df282a8f7a8fa05a6629bec01d8024c3611eb1c1fb9f84_51f1ddda0933696150cf78fe7a2141653e6a841d2f4ecaaa915a299cb7a4d19c_3000000000000000000", Key: Key{}},
	}
	address, err := hex.DecodeString("93236a9263d2ac6198c5ed211774c745d5dc62a910cb84276f8a7c4959208915")
	assert.NoError(t, err)

	delegators := namedkeys.GetDelegateFromDelegator(address)

	assert.Equal(t, 2, len(delegators))
	assert.Equal(t, "1000000000000000000", delegators["93236a9263d2ac6198c5ed211774c745d5dc62a910cb84276f8a7c4959208915"])
	assert.Equal(t, "2000000000000000000", delegators["51f1ddda0933696150cf78fe7a2141653e6a841d2f4ecaaa915a299cb7a4d19c"])
}

func TestNamedKeysGetCommission(t *testing.T) {
	namedkeys := NamedKeys{
		NamedKey{Name: "d_d70243dd9d0d646fd6df282a8f7a8fa05a6629bec01d8024c3611eb1c1fb9f84_d70243dd9d0d646fd6df282a8f7a8fa05a6629bec01d8024c3611eb1c1fb9f84_1000000000000000000", Key: Key{}},
		NamedKey{Name: "v_d70243dd9d0d646fd6df282a8f7a8fa05a6629bec01d8024c3611eb1c1fb9f84_4000000000000000000", Key: Key{}},
		NamedKey{Name: "d_93236a9263d2ac6198c5ed211774c745d5dc62a910cb84276f8a7c4959208915_93236a9263d2ac6198c5ed211774c745d5dc62a910cb84276f8a7c4959208915_1000000000000000000", Key: Key{}},
		NamedKey{Name: "c_93236a9263d2ac6198c5ed211774c745d5dc62a910cb84276f8a7c4959208915_100000000000000000", Key: Key{}},
		NamedKey{Name: "v_93236a9263d2ac6198c5ed211774c745d5dc62a910cb84276f8a7c4959208915_3000000000000000000", Key: Key{}},
		NamedKey{Name: "pos_bonding_purse", Key: Key{}},
		NamedKey{Name: "d_93236a9263d2ac6198c5ed211774c745d5dc62a910cb84276f8a7c4959208915_51f1ddda0933696150cf78fe7a2141653e6a841d2f4ecaaa915a299cb7a4d19c_2000000000000000000", Key: Key{}},
		NamedKey{Name: "d_d70243dd9d0d646fd6df282a8f7a8fa05a6629bec01d8024c3611eb1c1fb9f84_51f1ddda0933696150cf78fe7a2141653e6a841d2f4ecaaa915a299cb7a4d19c_3000000000000000000", Key: Key{}},
	}
	address, err := hex.DecodeString("93236a9263d2ac6198c5ed211774c745d5dc62a910cb84276f8a7c4959208915")
	assert.NoError(t, err)

	values := namedkeys.GetValidatorCommission(address)

	assert.Equal(t, "100000000000000000", values)
}

func TestNamedKeysGetReward(t *testing.T) {
	namedkeys := NamedKeys{
		NamedKey{Name: "d_d70243dd9d0d646fd6df282a8f7a8fa05a6629bec01d8024c3611eb1c1fb9f84_d70243dd9d0d646fd6df282a8f7a8fa05a6629bec01d8024c3611eb1c1fb9f84_1000000000000000000", Key: Key{}},
		NamedKey{Name: "v_d70243dd9d0d646fd6df282a8f7a8fa05a6629bec01d8024c3611eb1c1fb9f84_4000000000000000000", Key: Key{}},
		NamedKey{Name: "d_93236a9263d2ac6198c5ed211774c745d5dc62a910cb84276f8a7c4959208915_93236a9263d2ac6198c5ed211774c745d5dc62a910cb84276f8a7c4959208915_1000000000000000000", Key: Key{}},
		NamedKey{Name: "c_93236a9263d2ac6198c5ed211774c745d5dc62a910cb84276f8a7c4959208915_100000000000000000", Key: Key{}},
		NamedKey{Name: "r_93236a9263d2ac6198c5ed211774c745d5dc62a910cb84276f8a7c4959208915_200000000000000000", Key: Key{}},
		NamedKey{Name: "r_d70243dd9d0d646fd6df282a8f7a8fa05a6629bec01d8024c3611eb1c1fb9f84_300000000000000000", Key: Key{}},
		NamedKey{Name: "v_93236a9263d2ac6198c5ed211774c745d5dc62a910cb84276f8a7c4959208915_3000000000000000000", Key: Key{}},
		NamedKey{Name: "pos_bonding_purse", Key: Key{}},
		NamedKey{Name: "d_93236a9263d2ac6198c5ed211774c745d5dc62a910cb84276f8a7c4959208915_51f1ddda0933696150cf78fe7a2141653e6a841d2f4ecaaa915a299cb7a4d19c_2000000000000000000", Key: Key{}},
		NamedKey{Name: "d_d70243dd9d0d646fd6df282a8f7a8fa05a6629bec01d8024c3611eb1c1fb9f84_51f1ddda0933696150cf78fe7a2141653e6a841d2f4ecaaa915a299cb7a4d19c_3000000000000000000", Key: Key{}},
	}
	address1, err := hex.DecodeString("93236a9263d2ac6198c5ed211774c745d5dc62a910cb84276f8a7c4959208915")
	address2, err := hex.DecodeString("d70243dd9d0d646fd6df282a8f7a8fa05a6629bec01d8024c3611eb1c1fb9f84")
	assert.NoError(t, err)

	address1Values := namedkeys.GetUserReward(address1)
	address2Values := namedkeys.GetUserReward(address2)

	assert.Equal(t, "200000000000000000", address1Values)
	assert.Equal(t, "300000000000000000", address2Values)
}
