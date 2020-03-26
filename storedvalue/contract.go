package storedvalue

import (
	"encoding/binary"
	"fmt"

	"github.com/hdac-io/casperlabs-ee-grpc-go-util/protobuf/io/casperlabs/casper/consensus/state"
)

const (
	PROTOCOL_VERSION_MAJOR_LENGTH = 4
	PROTOCOL_VERSION_MINOR_LENGTH = 4
	PROTOCOL_VERSION_PATCH_LENGTH = 4
	PROTOCOL_VERSION_LENGTH       = PROTOCOL_VERSION_MAJOR_LENGTH + PROTOCOL_VERSION_MINOR_LENGTH + PROTOCOL_VERSION_PATCH_LENGTH
)

type Contract struct {
	Body            []byte          `json:"body"`
	NamedKeys       []NamedKey      `json:"named_key"`
	ProtocolVersion ProtocolVersion `json:"protocol_version"`
}

func NewContract(body []byte, namedKeys []NamedKey, protocolVersion ProtocolVersion) Contract {
	return Contract{
		Body:            body,
		NamedKeys:       namedKeys,
		ProtocolVersion: protocolVersion,
	}
}

func (c Contract) FromBytes(src []byte) (contract Contract, err error, pos int) {
	pos = 0
	bodySizeBytes := src[pos : pos+SIZE_LENGTH]
	pos += SIZE_LENGTH
	bodySize := int(binary.LittleEndian.Uint32(bodySizeBytes))
	body := src[pos : pos+bodySize]
	pos += bodySize

	// NamedKeys
	namedKeys := []NamedKey{}
	namedKeysSizeBytes := src[pos : pos+SIZE_LENGTH]
	pos += SIZE_LENGTH
	namedKeysSize := int(binary.LittleEndian.Uint32(namedKeysSizeBytes))
	for i := 0; i < namedKeysSize; i++ {
		var namedKey NamedKey
		namedKey, err, length := namedKey.FromBytes(src[pos:])
		if err != nil {
			return Contract{}, err, pos
		}
		pos += length

		namedKeys = append(namedKeys, namedKey)
	}

	var protocolVersion ProtocolVersion
	protocolVersion, err, length := protocolVersion.FromBytes(src[pos:])
	pos += length
	if err != nil {
		return Contract{}, err, pos
	}

	c = NewContract(body, namedKeys, protocolVersion)
	return c, nil, pos
}

func (c Contract) ToBytes() []byte {
	res := make([]byte, SIZE_LENGTH)
	binary.LittleEndian.PutUint32(res, uint32(len(c.Body)))

	namedKeysLengthBytes := make([]byte, SIZE_LENGTH)
	binary.LittleEndian.PutUint32(namedKeysLengthBytes, uint32(len(c.NamedKeys)))
	res = append(res, namedKeysLengthBytes...)
	for _, namedKey := range c.NamedKeys {
		res = append(res, namedKey.ToBytes()...)
	}

	res = append(res, c.ProtocolVersion.ToBytes()...)

	return res
}

func (c Contract) ToStateValue() *state.Contract {
	stateNamedKeys := []*state.NamedKey{}
	for _, namedKey := range c.NamedKeys {
		stateNamedKeys = append(stateNamedKeys, namedKey.ToStateValue())
	}

	return &state.Contract{
		Body:            c.Body,
		NamedKeys:       stateNamedKeys,
		ProtocolVersion: c.ProtocolVersion.ToStateValue(),
	}
}

func (c Contract) FromStateValue(state *state.Contract) (Contract, error) {
	namedKeys := []NamedKey{}
	for _, stateNamedKey := range state.GetNamedKeys() {
		var namedKey NamedKey
		namedKey, err := namedKey.FromStateValue(stateNamedKey)
		if err != nil {
			return Contract{}, nil
		}
		namedKeys = append(namedKeys, namedKey)
	}

	return NewContract(
		state.Body,
		namedKeys,
		NewProtocolVersion(state.ProtocolVersion.GetMajor(), state.ProtocolVersion.GetMinor(), state.ProtocolVersion.GetPatch()),
	), nil
}

type ProtocolVersion struct {
	Major uint32 `json:"major,omitempty"`
	Minor uint32 `json:"minor,omitempty"`
	Patch uint32 `json:"patch,omitempty"`
}

func NewProtocolVersion(major uint32, minor uint32, patch uint32) ProtocolVersion {
	return ProtocolVersion{
		Major: major,
		Minor: minor,
		Patch: patch,
	}
}

func (p ProtocolVersion) FromBytes(src []byte) (protocolVersion ProtocolVersion, err error, pos int) {
	pos = 0
	if len(src) != PROTOCOL_VERSION_LENGTH {
		return ProtocolVersion{}, fmt.Errorf("ActionThresholds must be %d, but %d", PROTOCOL_VERSION_LENGTH, len(src)), pos
	}

	major := binary.LittleEndian.Uint32(src[pos : pos+PROTOCOL_VERSION_MAJOR_LENGTH])
	pos += PROTOCOL_VERSION_MAJOR_LENGTH
	minor := binary.LittleEndian.Uint32(src[pos : pos+PROTOCOL_VERSION_MINOR_LENGTH])
	pos += PROTOCOL_VERSION_MINOR_LENGTH
	patch := binary.LittleEndian.Uint32(src[pos : pos+PROTOCOL_VERSION_PATCH_LENGTH])
	pos += PROTOCOL_VERSION_PATCH_LENGTH

	p = NewProtocolVersion(major, minor, patch)

	return p, nil, pos
}

func (p ProtocolVersion) ToBytes() []byte {
	pos := 0
	res := make([]byte, PROTOCOL_VERSION_LENGTH)
	binary.LittleEndian.PutUint32(res[pos:pos+PROTOCOL_VERSION_MAJOR_LENGTH], p.Major)
	pos += PROTOCOL_VERSION_MAJOR_LENGTH
	binary.LittleEndian.PutUint32(res[pos:pos+PROTOCOL_VERSION_MINOR_LENGTH], p.Minor)
	pos += PROTOCOL_VERSION_PATCH_LENGTH
	binary.LittleEndian.PutUint32(res[pos:pos+PROTOCOL_VERSION_PATCH_LENGTH], p.Patch)
	return res
}

func (p ProtocolVersion) ToStateValue() *state.ProtocolVersion {
	return &state.ProtocolVersion{
		Major: p.Major,
		Minor: p.Minor,
		Patch: p.Patch,
	}
}
