package storedvalue

import (
	"testing"

	"github.com/hdac-io/casperlabs-ee-grpc-go-util/protobuf/io/casperlabs/casper/consensus/state"
	"github.com/stretchr/testify/assert"
)

func TestURefToByteUnknown(t *testing.T) {
	uref := NewURef(make([]byte, 32), state.Key_URef_NONE)

	res := uref.ToBytes()
	assert.Equal(
		t,
		[]byte{
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0},
		res)
}

func TestURefToByteReadAddWrite(t *testing.T) {
	uref := NewURef(make([]byte, 32), state.Key_URef_READ_ADD_WRITE)

	res := uref.ToBytes()
	assert.Equal(
		t,
		[]byte{
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			7},
		res)
}

func TestURefFromByteUnknown(t *testing.T) {
	src := []byte{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0}

	var uref URef
	uref, err, pos := uref.FromBytes(src)

	assert.NoError(t, err)
	assert.Equal(
		t,
		make([]byte, 32),
		uref.GetAddress())
	assert.Equal(
		t,
		state.Key_URef_NONE,
		uref.GetAccessRights())
	assert.Equal(t, len(src), pos)
}

func TestURefFromByteReadAddWrite(t *testing.T) {
	src := []byte{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		7}

	var uref URef
	uref, err, pos := uref.FromBytes(src)

	assert.NoError(t, err)
	assert.Equal(
		t,
		make([]byte, 32),
		uref.GetAddress())
	assert.Equal(
		t,
		state.Key_URef_READ_ADD_WRITE,
		uref.GetAccessRights())
	assert.Equal(t, len(src), pos)
}

func TestURefFromByteError(t *testing.T) {
	src := []byte{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	var uref URef
	uref, err, pos := uref.FromBytes(src)

	assert.Error(t, err)
	assert.Equal(t, 0, pos)
}
