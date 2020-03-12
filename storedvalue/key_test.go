package storedvalue

import (
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
			1, 7,
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
		1, 7,
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
