package storedvalue

import (
	"encoding/binary"
	"fmt"

	"github.com/hdac-io/casperlabs-ee-grpc-go-util/protobuf/io/casperlabs/casper/consensus/state"
)

const (
	ASSOCIATED_KEY_SIZE_LENGTH       = 4
	ASSOCIATED_KEY_LENGTH            = ADDRESS_LENGTH
	ASSOCIATED_KEY_WEIGHT_LENGTH     = 1
	ASSOCIATED_KEY_SERIALIZED_LENGTH = ASSOCIATED_KEY_LENGTH + ASSOCIATED_KEY_WEIGHT_LENGTH

	ACTION_THRESHOLD_DEPLOYMENT_LENGTH     = 1
	ACTION_THRESHOLD_KEY_MANAGEMENT_LENGTH = 1
)

type Account struct {
	PublicKey        []byte
	NamedKeys        NamedKeys
	MainPurse        URef
	AssociatedKeys   []AssociatedKey
	ActionThresholds ActionThresholds
}

func NewAccount(publicKey []byte, namedKeys NamedKeys, purseId URef, associatedKeys []AssociatedKey, actionThresholds ActionThresholds) Account {
	return Account{
		PublicKey:        publicKey,
		NamedKeys:        namedKeys,
		MainPurse:        purseId,
		AssociatedKeys:   associatedKeys,
		ActionThresholds: actionThresholds,
	}
}

func (a Account) FromBytes(src []byte) (account Account, err error, pos int) {
	pos = 0
	publicKey := src[pos:ADDRESS_LENGTH]
	pos += ADDRESS_LENGTH

	// NamedKeys
	namedKeys := NamedKeys{}
	namedKeysSizeBytes := src[pos : pos+SIZE_LENGTH]
	pos += SIZE_LENGTH
	namedKeysSize := int(binary.LittleEndian.Uint32(namedKeysSizeBytes))
	for i := 0; i < namedKeysSize; i++ {
		var namedKey NamedKey
		namedKey, err, length := namedKey.FromBytes(src[pos:])
		if err != nil {
			return Account{}, err, pos
		}
		pos += length

		namedKeys = append(namedKeys, namedKey)
	}

	// Purse ID
	var purseID URef
	purseID, err, length := purseID.FromBytes(src[pos:])
	if err != nil {
		return Account{}, err, pos
	}
	pos += length

	// Associate Key
	associatedKeys := []AssociatedKey{}
	associatedKeysSizeBytes := src[pos : pos+SIZE_LENGTH]
	pos += SIZE_LENGTH
	associateKeySize := int(binary.LittleEndian.Uint32(associatedKeysSizeBytes))
	for i := 0; i < associateKeySize; i++ {
		var associatedKey AssociatedKey
		associatedKey, err, length := associatedKey.FromBytes(src[pos:])
		if err != nil {
			return Account{}, err, pos
		}
		pos += length

		associatedKeys = append(associatedKeys, associatedKey)
	}

	// ActionThresholds
	var actionThresholds ActionThresholds
	actionThresholds, err, length = actionThresholds.FromBytes(src[pos:])
	if err != nil {
		return Account{}, err, pos
	}
	pos += length

	return NewAccount(publicKey, namedKeys, purseID, associatedKeys, actionThresholds), nil, pos
}

func (a Account) ToBytes() []byte {
	res := a.PublicKey

	namedKeysLengthBytes := make([]byte, SIZE_LENGTH)
	binary.LittleEndian.PutUint32(namedKeysLengthBytes, uint32(len(a.NamedKeys)))
	res = append(res, namedKeysLengthBytes...)
	for _, namedKey := range a.NamedKeys {
		res = append(res, namedKey.ToBytes()...)
	}

	res = append(res, a.MainPurse.ToBytes()...)

	associatedKeysLengthBytes := make([]byte, SIZE_LENGTH)
	binary.LittleEndian.PutUint32(associatedKeysLengthBytes, uint32(len(a.AssociatedKeys)))
	for _, associatedKey := range a.AssociatedKeys {
		res = append(res, associatedKey.ToBytes()...)
	}

	res = append(res, a.ActionThresholds.ToBytes()...)

	return res
}

func (a Account) ToStateValue() *state.Account {
	stateNamedKeys := []*state.NamedKey{}
	for _, namedKey := range a.NamedKeys {
		stateNamedKeys = append(stateNamedKeys, namedKey.ToStateValue())
	}

	stateAssociatedKeys := []*state.Account_AssociatedKey{}
	for _, associatedKey := range a.AssociatedKeys {
		stateAssociatedKeys = append(stateAssociatedKeys, associatedKey.ToStateValue())
	}

	return &state.Account{
		PublicKey:        a.PublicKey,
		MainPurse:        a.MainPurse.ToStateValue(),
		NamedKeys:        stateNamedKeys,
		AssociatedKeys:   stateAssociatedKeys,
		ActionThresholds: a.ActionThresholds.ToStateValue(),
	}
}

func (a Account) FromStateValue(state *state.Account) (Account, error) {
	namedKeys := []NamedKey{}
	for _, stateNamedKey := range state.GetNamedKeys() {
		var namedKey NamedKey
		namedKey, err := namedKey.FromStateValue(stateNamedKey)
		if err != nil {
			return Account{}, nil
		}
		namedKeys = append(namedKeys, namedKey)
	}

	associatedKeys := []AssociatedKey{}
	for _, stateAssociatedKey := range state.GetAssociatedKeys() {
		associatedKey := NewAssociatedKey(stateAssociatedKey.GetPublicKey(), stateAssociatedKey.GetWeight())
		associatedKeys = append(associatedKeys, associatedKey)
	}

	return NewAccount(
		state.GetPublicKey(),
		namedKeys,
		NewURef(state.GetMainPurse().GetUref(), state.GetMainPurse().GetAccessRights()),
		associatedKeys,
		NewActionThresholds(state.ActionThresholds.GetDeploymentThreshold(), state.ActionThresholds.GetKeyManagementThreshold())), nil
}

type AssociatedKey struct {
	PublicKey []byte `json:"public_key"`
	Weight    uint32 `json:"weight"`
}

func NewAssociatedKey(publicKey []byte, weight uint32) AssociatedKey {
	return AssociatedKey{
		PublicKey: publicKey,
		Weight:    weight,
	}
}

func (a AssociatedKey) FromBytes(src []byte) (associatedKey AssociatedKey, err error, pos int) {
	pos = 0

	if len(src) < ASSOCIATED_KEY_LENGTH+ASSOCIATED_KEY_WEIGHT_LENGTH {
		return AssociatedKey{}, fmt.Errorf("AssociatedKey more than %d, but %d", ASSOCIATED_KEY_LENGTH+ASSOCIATED_KEY_WEIGHT_LENGTH, len(src)), pos
	}

	publicKey := src[pos:ASSOCIATED_KEY_LENGTH]
	pos += ASSOCIATED_KEY_LENGTH
	weight := uint32(src[pos])
	pos += ASSOCIATED_KEY_WEIGHT_LENGTH

	return NewAssociatedKey(publicKey, weight), nil, pos
}

func (a AssociatedKey) ToBytes() []byte {
	res := a.PublicKey

	weightBytes := make([]byte, SIZE_LENGTH)
	binary.BigEndian.PutUint32(weightBytes, a.Weight)

	res = append(res, weightBytes...)
	return res
}

func (a AssociatedKey) ToStateValue() *state.Account_AssociatedKey {
	return &state.Account_AssociatedKey{
		PublicKey: a.PublicKey,
		Weight:    a.Weight,
	}
}

type ActionThresholds struct {
	DeploymentThreshold    uint32 `json:"deployment_threshold`
	KeyManagementThreshold uint32 `json:"key_management_threshold"`
}

func NewActionThresholds(deploymentThreshold uint32, keyManagementThreshold uint32) ActionThresholds {
	return ActionThresholds{
		DeploymentThreshold:    deploymentThreshold,
		KeyManagementThreshold: keyManagementThreshold,
	}
}

func (a ActionThresholds) FromBytes(src []byte) (actionThresholds ActionThresholds, err error, pos int) {
	pos = 0
	if len(src) < ACTION_THRESHOLD_DEPLOYMENT_LENGTH+ACTION_THRESHOLD_KEY_MANAGEMENT_LENGTH {
		return ActionThresholds{}, fmt.Errorf("ActionThresholds more than %d, but %d", ACTION_THRESHOLD_DEPLOYMENT_LENGTH+ACTION_THRESHOLD_KEY_MANAGEMENT_LENGTH, len(src)), pos
	}

	deployment := uint32(src[pos])
	pos += ACTION_THRESHOLD_DEPLOYMENT_LENGTH
	keyManagement := uint32(src[pos])
	pos += ACTION_THRESHOLD_KEY_MANAGEMENT_LENGTH
	return NewActionThresholds(deployment, keyManagement), nil, pos
}

func (a ActionThresholds) ToBytes() []byte {
	return []byte{byte(a.DeploymentThreshold), byte(a.KeyManagementThreshold)}
}

func (a ActionThresholds) ToStateValue() *state.Account_ActionThresholds {
	return &state.Account_ActionThresholds{
		DeploymentThreshold:    a.DeploymentThreshold,
		KeyManagementThreshold: a.KeyManagementThreshold,
	}
}
