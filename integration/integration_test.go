package integration

import (
	"encoding/hex"
	"testing"

	"github.com/hdac-io/casperlabs-ee-grpc-go-util/grpc"
	"github.com/hdac-io/casperlabs-ee-grpc-go-util/util"
	"github.com/stretchr/testify/assert"
)

func TestCustomContractCounter(t *testing.T) {
	client, rootStateHash, proxyHash, protocolVersion := InitalRunGenensis()

	// counterDefine
	rootStateHash, _ = RunCounterDefine(client, rootStateHash, genesisAddress, proxyHash, protocolVersion)

	// query
	storedValue := RunQuery(client, rootStateHash, "address", genesisAddress, []string{"counter", "count"}, protocolVersion)
	assert.Equal(t, int32(0), storedValue.ClValue.ToStateValues().GetIntValue())

	// First counter call
	rootStateHash, _ = RunCounterCall(client, rootStateHash, genesisAddress, proxyHash, protocolVersion)

	// query
	storedValue = RunQuery(client, rootStateHash, "address", genesisAddress, []string{"counter", "count"}, protocolVersion)
	assert.Equal(t, int32(1), storedValue.ClValue.ToStateValues().GetIntValue())

	// Second counter call
	rootStateHash, _ = RunCounterCall(client, rootStateHash, genesisAddress, proxyHash, protocolVersion)

	// query
	storedValue = RunQuery(client, rootStateHash, "address", genesisAddress, []string{"counter", "count"}, protocolVersion)
	assert.Equal(t, int32(2), storedValue.ClValue.ToStateValues().GetIntValue())
}

func TestTransferToAccount(t *testing.T) {
	client, rootStateHash, proxyHash, protocolVersion := InitalRunGenensis()
	amount := "900000000000000000"

	rootStateHash, _ = RunTransferToAccount(client, rootStateHash, genesisAddress, address1, amount, proxyHash, protocolVersion)

	queryResult, errMessage := grpc.QueryBalance(client, rootStateHash, address1, protocolVersion)
	assert.Equal(t, amount, queryResult)
	assert.Equal(t, "", errMessage)
}

func TestBond(t *testing.T) {
	client, rootStateHash, proxyHash, protocolVersion := InitalRunGenensis()
	bondAmount := "10000000000000000"

	rootStateHash, bonds := RunBond(client, rootStateHash, genesisAddress, bondAmount, proxyHash, protocolVersion)
	assert.Equal(t, "1010000000000000000", bonds[0].GetStake().GetValue())

	storedValue := RunQuery(client, rootStateHash, "address", genesisAddress, []string{"pos"}, protocolVersion)
	allValidator := storedValue.Contract.NamedKeys.GetAllValidators()

	assert.Equal(t, 1, len(allValidator))
	assert.Equal(t, "1010000000000000000", allValidator["d70243dd9d0d646fd6df282a8f7a8fa05a6629bec01d8024c3611eb1c1fb9f84"])
}

func TestUnbond(t *testing.T) {
	client, rootStateHash, proxyHash, protocolVersion := InitalRunGenensis()
	amount := "10000000000000000"

	rootStateHash, bonds := RunUnbond(client, rootStateHash, genesisAddress, amount, proxyHash, protocolVersion)
	assert.Equal(t, "990000000000000000", bonds[0].GetStake().GetValue())

	storedValue := RunQuery(client, rootStateHash, "address", genesisAddress, []string{"pos"}, protocolVersion)
	allValidator := storedValue.Contract.NamedKeys.GetAllValidators()

	assert.Equal(t, 1, len(allValidator))
	assert.Equal(t, "990000000000000000", allValidator["d70243dd9d0d646fd6df282a8f7a8fa05a6629bec01d8024c3611eb1c1fb9f84"])
}

func TestDelegate(t *testing.T) {
	client, rootStateHash, proxyHash, protocolVersion := InitalRunGenensis()
	amount := "100"

	rootStateHash, bonds := RunDelegate(client, rootStateHash, genesisAddress, genesisAddress, amount, proxyHash, protocolVersion)
	assert.Equal(t, "1000000000000000100", bonds[0].GetStake().GetValue())

	storedValue := RunQuery(client, rootStateHash, "address", SYSTEM_ACCOUNT, []string{"pos"}, protocolVersion)
	delegators := storedValue.Contract.NamedKeys.GetDelegateFromValidator(genesisAddress)
	assert.Equal(t, 1, len(delegators))
	assert.Equal(t, "1000000000000000100", delegators[hex.EncodeToString(genesisAddress)])
}

func TestDelegateFromAnotherAddress(t *testing.T) {
	client, rootStateHash, proxyHash, protocolVersion := InitalRunGenensis()
	amount := "100000000000000000"
	delegateAmount := "1000000000000000"

	rootStateHash, _ = RunTransferToAccount(client, rootStateHash, genesisAddress, address1, amount, proxyHash, protocolVersion)
	balance, errMessage := grpc.QueryBalance(client, rootStateHash, address1, protocolVersion)
	assert.Equal(t, "100000000000000000", balance)
	assert.Equal(t, "", errMessage)

	rootStateHash, _ = RunDelegate(client, rootStateHash, address1, genesisAddress, delegateAmount, proxyHash, protocolVersion)

	storedValue := RunQuery(client, rootStateHash, "address", SYSTEM_ACCOUNT, []string{"pos"}, protocolVersion)
	delegators := storedValue.Contract.NamedKeys.GetDelegateFromValidator(genesisAddress)
	assert.Equal(t, 2, len(delegators))
	assert.Equal(t, "1000000000000000000", delegators[hex.EncodeToString(genesisAddress)])
	assert.Equal(t, "1000000000000000", delegators[hex.EncodeToString(address1)])
}

func TestUndelegation(t *testing.T) {
	client, rootStateHash, proxyHash, protocolVersion := InitalRunGenensis()
	amount := "100"

	rootStateHash, bonds := RunUndelegate(client, rootStateHash, genesisAddress, genesisAddress, amount, proxyHash, protocolVersion)
	assert.Equal(t, "999999999999999900", bonds[0].GetStake().GetValue())

	storedValue := RunQuery(client, rootStateHash, "address", SYSTEM_ACCOUNT, []string{"pos"}, protocolVersion)
	delegators := storedValue.Contract.NamedKeys.GetDelegateFromValidator(genesisAddress)
	assert.Equal(t, 1, len(delegators))
	assert.Equal(t, "999999999999999900", delegators[hex.EncodeToString(genesisAddress)])
}

func TestRedelegation(t *testing.T) {
	client, rootStateHash, proxyHash, protocolVersion := InitalRunGenensis()
	amount := "100"

	rootStateHash, bonds := RunRedelegate(client, rootStateHash, genesisAddress, genesisAddress, address1, amount, proxyHash, protocolVersion)
	assert.Equal(t, 2, len(bonds))

	storedValue := RunQuery(client, rootStateHash, "address", SYSTEM_ACCOUNT, []string{"pos"}, protocolVersion)
	delegators := storedValue.Contract.NamedKeys.GetDelegateFromDelegator(genesisAddress)
	assert.Equal(t, 2, len(delegators))
	assert.Equal(t, "999999999999999900", delegators[hex.EncodeToString(genesisAddress)])
	assert.Equal(t, "100", delegators[hex.EncodeToString(address1)])
}

func TestVoteAndUnvote(t *testing.T) {
	client, rootStateHash, proxyHash, protocolVersion := InitalRunGenensis()
	amount := "123"

	rootStateHash, _ = RunVote(client, rootStateHash, genesisAddress, address1, amount, proxyHash, protocolVersion)
	storedValue := RunQuery(client, rootStateHash, "address", SYSTEM_ACCOUNT, []string{"pos"}, protocolVersion)
	voter := storedValue.Contract.NamedKeys.GetVotingDappFromUser(genesisAddress)

	assert.Equal(t, 1, len(voter))
	assert.Equal(t, amount, voter[hex.EncodeToString(address1)])

	unvoteAmount := "23"
	rootStateHash, _ = RunUnvote(client, rootStateHash, genesisAddress, address1, unvoteAmount, proxyHash, protocolVersion)
	storedValue = RunQuery(client, rootStateHash, "address", SYSTEM_ACCOUNT, []string{"pos"}, protocolVersion)
	voter = storedValue.Contract.NamedKeys.GetVotingDappFromUser(genesisAddress)

	assert.Equal(t, 1, len(voter))
	assert.Equal(t, "100", voter[hex.EncodeToString(address1)])

	voter = storedValue.Contract.NamedKeys.GetVotingUserFromDapp(address1)
	assert.Equal(t, "100", voter[hex.EncodeToString(genesisAddress)])
}

func TestVoteMoreAccount(t *testing.T) {
	client, rootStateHash, proxyHash, protocolVersion := InitalRunGenensis()
	address2 := util.DecodeHexString("03170a2e7597b7b7e3d84c05391d139a62b157e78786d8c082f29dcf4c111314")
	address3 := util.DecodeHexString("f0f84944e0ccfa9e67383e6a448291787d208c8e46adc849f714078663d1dd36")
	amount1 := "100"
	amount2 := "200"
	amount3 := "300"

	rootStateHash, _ = RunVote(client, rootStateHash, genesisAddress, address1, amount1, proxyHash, protocolVersion)
	rootStateHash, _ = RunVote(client, rootStateHash, genesisAddress, address2, amount2, proxyHash, protocolVersion)
	rootStateHash, _ = RunVote(client, rootStateHash, genesisAddress, address3, amount3, proxyHash, protocolVersion)

	storedValue := RunQuery(client, rootStateHash, "address", SYSTEM_ACCOUNT, []string{"pos"}, protocolVersion)
	values := storedValue.Contract.NamedKeys.GetVotingDappFromUser(genesisAddress)

	assert.Equal(t, 3, len(values))
	assert.Equal(t, "100", values[hex.EncodeToString(address1)])
	assert.Equal(t, "200", values[hex.EncodeToString(address2)])
	assert.Equal(t, "300", values[hex.EncodeToString(address3)])

	values = storedValue.Contract.NamedKeys.GetVotingUserFromDapp(address1)
	assert.Equal(t, "100", values[hex.EncodeToString(genesisAddress)])

	values = storedValue.Contract.NamedKeys.GetVotingUserFromDapp(address2)
	assert.Equal(t, "200", values[hex.EncodeToString(genesisAddress)])

	values = storedValue.Contract.NamedKeys.GetVotingUserFromDapp(address3)
	assert.Equal(t, "300", values[hex.EncodeToString(genesisAddress)])
}

func TestStepAndCommission(t *testing.T) {
	client, rootStateHash, proxyHash, protocolVersion := InitalRunGenensis()
	amount := "1000000000000000000"

	// ready to create block
	rootStateHash, _ = RunTransferToAccount(client, rootStateHash, genesisAddress, SYSTEM_ACCOUNT, "1000000000000000000", proxyHash, protocolVersion)

	rootStateHash, _ = RunTransferToAccount(client, rootStateHash, genesisAddress, address1, amount, proxyHash, protocolVersion)
	rootStateHash, _ = RunDelegate(client, rootStateHash, address1, genesisAddress, "100000000000000000", proxyHash, protocolVersion)

	beforeAddress1Amount, errMessage := grpc.QueryBalance(client, rootStateHash, address1, protocolVersion)
	assert.Equal(t, "", errMessage)

	beforeGenesisAmount, errMessage := grpc.QueryBalance(client, rootStateHash, genesisAddress, protocolVersion)
	assert.Equal(t, "", errMessage)

	// step
	rootStateHash, _ = RunStep(client, rootStateHash, SYSTEM_ACCOUNT, proxyHash, protocolVersion)

	afterStepAddress1Amount, errMessage := grpc.QueryBalance(client, rootStateHash, address1, protocolVersion)
	assert.Equal(t, "", errMessage)
	assert.Equal(t, beforeAddress1Amount, afterStepAddress1Amount)

	afterStepGenesisAmount, errMessage := grpc.QueryBalance(client, rootStateHash, genesisAddress, protocolVersion)
	assert.Equal(t, beforeGenesisAmount, afterStepGenesisAmount)

	// claim Commission
	rootStateHash, _ = RunClaimCommission(client, rootStateHash, genesisAddress, proxyHash, protocolVersion)
	rootStateHash, _ = RunStep(client, rootStateHash, SYSTEM_ACCOUNT, proxyHash, protocolVersion)

	afterClaimCommissionAddress1Amount, errMessage := grpc.QueryBalance(client, rootStateHash, address1, protocolVersion)
	assert.Equal(t, "", errMessage)
	assert.Equal(t, beforeAddress1Amount, afterClaimCommissionAddress1Amount)

	afterClaimCommissionGenesisAmount, errMessage := grpc.QueryBalance(client, rootStateHash, genesisAddress, protocolVersion)
	assert.NotEqual(t, beforeGenesisAmount, afterClaimCommissionGenesisAmount)

	// claim Reward
	rootStateHash, _ = RunClaimReward(client, rootStateHash, genesisAddress, proxyHash, protocolVersion)
	rootStateHash, _ = RunClaimReward(client, rootStateHash, address1, proxyHash, protocolVersion)
	rootStateHash, _ = RunStep(client, rootStateHash, SYSTEM_ACCOUNT, proxyHash, protocolVersion)

	afterClaimRewardAddress1Amount, errMessage := grpc.QueryBalance(client, rootStateHash, address1, protocolVersion)
	assert.Equal(t, "", errMessage)
	assert.NotEqual(t, afterClaimCommissionAddress1Amount, afterClaimRewardAddress1Amount)

	afterClaimRewardGenesisAmount, errMessage := grpc.QueryBalance(client, rootStateHash, genesisAddress, protocolVersion)
	assert.Equal(t, "", errMessage)
	assert.NotEqual(t, afterClaimCommissionGenesisAmount, afterClaimRewardGenesisAmount)
}

func TestClaimAmount(t *testing.T) {
	client, rootStateHash, proxyHash, protocolVersion := InitalRunGenensis()
	amount := "1000000000000000000"

	// ready to create block
	rootStateHash, _ = RunTransferToAccount(client, rootStateHash, genesisAddress, SYSTEM_ACCOUNT, "1000000000000000000", proxyHash, protocolVersion)

	rootStateHash, _ = RunTransferToAccount(client, rootStateHash, genesisAddress, address1, amount, proxyHash, protocolVersion)
	rootStateHash, _ = RunDelegate(client, rootStateHash, address1, genesisAddress, "100000000000000000", proxyHash, protocolVersion)

	// 10 step
	for i := 0; i < 10; i++ {
		rootStateHash, _ = RunStep(client, rootStateHash, SYSTEM_ACCOUNT, proxyHash, protocolVersion)
	}

	storedValue := RunQuery(client, rootStateHash, "address", SYSTEM_ACCOUNT, []string{"pos"}, protocolVersion)
	step10Address1Reward := storedValue.Contract.NamedKeys.GetUserReward(address1)
	step10GenesisAddressReward := storedValue.Contract.NamedKeys.GetUserReward(genesisAddress)
	step10GenesisAddressCommission := storedValue.Contract.NamedKeys.GetValidatorCommission(genesisAddress)

	assert.NotEqual(t, "", step10Address1Reward)
	assert.NotEqual(t, "", step10GenesisAddressReward)
	assert.NotEqual(t, "", step10GenesisAddressCommission)

	rootStateHash, _ = RunClaimCommission(client, rootStateHash, genesisAddress, proxyHash, protocolVersion)
	rootStateHash, _ = RunClaimReward(client, rootStateHash, genesisAddress, proxyHash, protocolVersion)
	rootStateHash, _ = RunClaimReward(client, rootStateHash, address1, proxyHash, protocolVersion)
	rootStateHash, _ = RunStep(client, rootStateHash, SYSTEM_ACCOUNT, proxyHash, protocolVersion)

	storedValue = RunQuery(client, rootStateHash, "address", SYSTEM_ACCOUNT, []string{"pos"}, protocolVersion)
	stepAddress1Reward := storedValue.Contract.NamedKeys.GetUserReward(address1)
	stepGenesisAddressReward := storedValue.Contract.NamedKeys.GetUserReward(genesisAddress)
	stepGenesisAddressCommission := storedValue.Contract.NamedKeys.GetValidatorCommission(genesisAddress)

	assert.NotEqual(t, len(step10Address1Reward)+1, len(stepAddress1Reward))
	assert.NotEqual(t, len(step10GenesisAddressReward)+1, len(stepGenesisAddressReward))
	assert.NotEqual(t, len(step10GenesisAddressCommission)+1, len(stepGenesisAddressCommission))
}
