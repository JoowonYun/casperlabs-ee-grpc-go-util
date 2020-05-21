package integration

import (
	"testing"

	"github.com/hdac-io/casperlabs-ee-grpc-go-util/grpc"
	"github.com/hdac-io/casperlabs-ee-grpc-go-util/protobuf/io/casperlabs/casper/consensus/state"
	"github.com/hdac-io/casperlabs-ee-grpc-go-util/protobuf/io/casperlabs/ipc"
	"github.com/hdac-io/casperlabs-ee-grpc-go-util/util"
	"github.com/stretchr/testify/assert"
)

func TestCustomContractCounter(t *testing.T) {
	client, rootStateHash, proxyHash, protocolVersion := InitalRunGenensis(DEFAULT_GENESIS_ACCOUNT, DEFAULT_GENESIS_STATE_INFO)

	// counterDefine
	rootStateHash, _ = RunCounterDefine(client, rootStateHash, GENESIS_ADDRESS, proxyHash, protocolVersion)

	// query
	storedValue := RunQuery(client, rootStateHash, "address", GENESIS_ADDRESS, []string{"counter", "count"}, protocolVersion)
	assert.Equal(t, int32(0), storedValue.ClValue.ToStateValues().GetIntValue())

	// First counter call
	rootStateHash, _ = RunCounterCall(client, rootStateHash, GENESIS_ADDRESS, proxyHash, protocolVersion)

	// query
	storedValue = RunQuery(client, rootStateHash, "address", GENESIS_ADDRESS, []string{"counter", "count"}, protocolVersion)
	assert.Equal(t, int32(1), storedValue.ClValue.ToStateValues().GetIntValue())

	// Second counter call
	rootStateHash, _ = RunCounterCall(client, rootStateHash, GENESIS_ADDRESS, proxyHash, protocolVersion)

	// query
	storedValue = RunQuery(client, rootStateHash, "address", GENESIS_ADDRESS, []string{"counter", "count"}, protocolVersion)
	assert.Equal(t, int32(2), storedValue.ClValue.ToStateValues().GetIntValue())
}

func TestTransferToAccount(t *testing.T) {
	client, rootStateHash, proxyHash, protocolVersion := InitalRunGenensis(DEFAULT_GENESIS_ACCOUNT, DEFAULT_GENESIS_STATE_INFO)
	amount := "900000000000000000"

	rootStateHash, _ = RunTransferToAccount(client, rootStateHash, GENESIS_ADDRESS, ADDRESS1, amount, proxyHash, protocolVersion)

	queryResult, errMessage := grpc.QueryBalance(client, rootStateHash, ADDRESS1, protocolVersion)
	assert.Equal(t, amount, queryResult)
	assert.Equal(t, "", errMessage)
}

func TestBond(t *testing.T) {
	client, rootStateHash, proxyHash, protocolVersion := InitalRunGenensis(DEFAULT_GENESIS_ACCOUNT, DEFAULT_GENESIS_STATE_INFO)
	bondAmount := "10000000000000000"

	rootStateHash, bonds := RunBond(client, rootStateHash, GENESIS_ADDRESS, bondAmount, proxyHash, protocolVersion)
	assert.Equal(t, "1010000000000000000", bonds[0].GetStake().GetValue())

	storedValue := RunQuery(client, rootStateHash, "address", GENESIS_ADDRESS, []string{"pos"}, protocolVersion)
	allValidator := storedValue.Contract.NamedKeys.GetAllValidators()

	assert.Equal(t, 1, len(allValidator))
	assert.Equal(t, "1010000000000000000", allValidator[GENESIS_ADDRESS_HEX])
}

func TestUnbond(t *testing.T) {
	client, rootStateHash, proxyHash, protocolVersion := InitalRunGenensis(DEFAULT_GENESIS_ACCOUNT, DEFAULT_GENESIS_STATE_INFO)
	amount := "10000000000000000"

	rootStateHash, bonds := RunUnbond(client, rootStateHash, GENESIS_ADDRESS, amount, proxyHash, protocolVersion)
	// step
	rootStateHash, bonds = RunStep(client, rootStateHash, SYSTEM_ACCOUNT, proxyHash, protocolVersion)
	assert.Equal(t, "990000000000000000", bonds[0].GetStake().GetValue())

	storedValue := RunQuery(client, rootStateHash, "address", GENESIS_ADDRESS, []string{"pos"}, protocolVersion)
	allValidator := storedValue.Contract.NamedKeys.GetAllValidators()

	assert.Equal(t, 1, len(allValidator))
	assert.Equal(t, "990000000000000000", allValidator[GENESIS_ADDRESS_HEX])
}

func TestDelegate(t *testing.T) {
	client, rootStateHash, proxyHash, protocolVersion := InitalRunGenensis(DEFAULT_GENESIS_ACCOUNT, DEFAULT_GENESIS_STATE_INFO)
	amount := "100"

	rootStateHash, bonds := RunDelegate(client, rootStateHash, GENESIS_ADDRESS, GENESIS_ADDRESS, amount, proxyHash, protocolVersion)
	assert.Equal(t, "1000000000000000100", bonds[0].GetStake().GetValue())

	storedValue := RunQuery(client, rootStateHash, "address", SYSTEM_ACCOUNT, []string{"pos"}, protocolVersion)
	delegators := storedValue.Contract.NamedKeys.GetDelegateFromValidator(GENESIS_ADDRESS)
	assert.Equal(t, 1, len(delegators))
	assert.Equal(t, "1000000000000000100", delegators[GENESIS_ADDRESS_HEX])
}

func TestDelegateFromAnotherAddress(t *testing.T) {
	client, rootStateHash, proxyHash, protocolVersion := InitalRunGenensis(DEFAULT_GENESIS_ACCOUNT, DEFAULT_GENESIS_STATE_INFO)
	amount := "100000000000000000"
	delegateAmount := "1000000000000000"

	rootStateHash, _ = RunTransferToAccount(client, rootStateHash, GENESIS_ADDRESS, ADDRESS1, amount, proxyHash, protocolVersion)
	balance, errMessage := grpc.QueryBalance(client, rootStateHash, ADDRESS1, protocolVersion)
	assert.Equal(t, "100000000000000000", balance)
	assert.Equal(t, "", errMessage)

	rootStateHash, _ = RunDelegate(client, rootStateHash, ADDRESS1, GENESIS_ADDRESS, delegateAmount, proxyHash, protocolVersion)

	storedValue := RunQuery(client, rootStateHash, "address", SYSTEM_ACCOUNT, []string{"pos"}, protocolVersion)
	delegators := storedValue.Contract.NamedKeys.GetDelegateFromValidator(GENESIS_ADDRESS)
	assert.Equal(t, 2, len(delegators))
	assert.Equal(t, "1000000000000000000", delegators[GENESIS_ADDRESS_HEX])
	assert.Equal(t, "1000000000000000", delegators[ADDRESS1_HEX])
}

func TestUndelegation(t *testing.T) {
	client, rootStateHash, proxyHash, protocolVersion := InitalRunGenensis(DEFAULT_GENESIS_ACCOUNT, DEFAULT_GENESIS_STATE_INFO)
	amount := "100"

	rootStateHash, bonds := RunUndelegate(client, rootStateHash, GENESIS_ADDRESS, GENESIS_ADDRESS, amount, proxyHash, protocolVersion)
	rootStateHash, bonds = RunStep(client, rootStateHash, SYSTEM_ACCOUNT, proxyHash, protocolVersion)
	assert.Equal(t, "999999999999999900", bonds[0].GetStake().GetValue())

	storedValue := RunQuery(client, rootStateHash, "address", SYSTEM_ACCOUNT, []string{"pos"}, protocolVersion)
	delegators := storedValue.Contract.NamedKeys.GetDelegateFromValidator(GENESIS_ADDRESS)
	assert.Equal(t, 1, len(delegators))
	assert.Equal(t, "999999999999999900", delegators[GENESIS_ADDRESS_HEX])
}

func TestRedelegation(t *testing.T) {
	client, rootStateHash, proxyHash, protocolVersion := InitalRunGenensis(DEFAULT_GENESIS_ACCOUNT, DEFAULT_GENESIS_STATE_INFO)
	amount := "100"

	rootStateHash, bonds := RunRedelegate(client, rootStateHash, GENESIS_ADDRESS, GENESIS_ADDRESS, ADDRESS1, amount, proxyHash, protocolVersion)
	rootStateHash, bonds = RunStep(client, rootStateHash, SYSTEM_ACCOUNT, proxyHash, protocolVersion)
	assert.Equal(t, 2, len(bonds))

	storedValue := RunQuery(client, rootStateHash, "address", SYSTEM_ACCOUNT, []string{"pos"}, protocolVersion)
	delegators := storedValue.Contract.NamedKeys.GetDelegateFromDelegator(GENESIS_ADDRESS)
	assert.Equal(t, 2, len(delegators))
	assert.Equal(t, "999999999999999900", delegators[GENESIS_ADDRESS_HEX])
	assert.Equal(t, "100", delegators[ADDRESS1_HEX])
}

func TestVoteAndUnvote(t *testing.T) {
	client, rootStateHash, proxyHash, protocolVersion := InitalRunGenensis(DEFAULT_GENESIS_ACCOUNT, DEFAULT_GENESIS_STATE_INFO)
	amount := "123"

	rootStateHash, _ = RunVote(client, rootStateHash, GENESIS_ADDRESS, ADDRESS1, amount, proxyHash, protocolVersion)
	storedValue := RunQuery(client, rootStateHash, "address", SYSTEM_ACCOUNT, []string{"pos"}, protocolVersion)
	voter := storedValue.Contract.NamedKeys.GetVotingDappFromUser(GENESIS_ADDRESS)

	assert.Equal(t, 1, len(voter))
	assert.Equal(t, amount, voter[ADDRESS1_DAPP_HEX])

	unvoteAmount := "23"
	rootStateHash, _ = RunUnvote(client, rootStateHash, GENESIS_ADDRESS, ADDRESS1, unvoteAmount, proxyHash, protocolVersion)
	storedValue = RunQuery(client, rootStateHash, "address", SYSTEM_ACCOUNT, []string{"pos"}, protocolVersion)
	voter = storedValue.Contract.NamedKeys.GetVotingDappFromUser(GENESIS_ADDRESS)

	assert.Equal(t, 1, len(voter))
	assert.Equal(t, "100", voter[ADDRESS1_DAPP_HEX])

	voter = storedValue.Contract.NamedKeys.GetVotingUserFromDapp(ADDRESS1_DAPP)
	assert.Equal(t, "100", voter[GENESIS_ADDRESS_HEX])
}

func TestVoteMoreAccount(t *testing.T) {
	client, rootStateHash, proxyHash, protocolVersion := InitalRunGenensis(DEFAULT_GENESIS_ACCOUNT, DEFAULT_GENESIS_STATE_INFO)
	address2 := util.DecodeHexString("03170a2e7597b7b7e3d84c05391d139a62b157e78786d8c082f29dcf4c111314")
	address3 := util.DecodeHexString("f0f84944e0ccfa9e67383e6a448291787d208c8e46adc849f714078663d1dd36")

	address2_dapp_hex := "0103170a2e7597b7b7e3d84c05391d139a62b157e78786d8c082f29dcf4c111314"
	address3_dapp_hex := "01f0f84944e0ccfa9e67383e6a448291787d208c8e46adc849f714078663d1dd36"
	address2_dapp := util.DecodeHexString("0103170a2e7597b7b7e3d84c05391d139a62b157e78786d8c082f29dcf4c111314")
	address3_dapp := util.DecodeHexString("01f0f84944e0ccfa9e67383e6a448291787d208c8e46adc849f714078663d1dd36")

	amount1 := "100"
	amount2 := "200"
	amount3 := "300"

	rootStateHash, _ = RunVote(client, rootStateHash, GENESIS_ADDRESS, ADDRESS1, amount1, proxyHash, protocolVersion)
	rootStateHash, _ = RunVote(client, rootStateHash, GENESIS_ADDRESS, address2, amount2, proxyHash, protocolVersion)
	rootStateHash, _ = RunVote(client, rootStateHash, GENESIS_ADDRESS, address3, amount3, proxyHash, protocolVersion)

	storedValue := RunQuery(client, rootStateHash, "address", SYSTEM_ACCOUNT, []string{"pos"}, protocolVersion)
	values := storedValue.Contract.NamedKeys.GetVotingDappFromUser(GENESIS_ADDRESS)

	assert.Equal(t, 3, len(values))
	assert.Equal(t, "100", values[ADDRESS1_DAPP_HEX])
	assert.Equal(t, "200", values[address2_dapp_hex])
	assert.Equal(t, "300", values[address3_dapp_hex])

	values = storedValue.Contract.NamedKeys.GetVotingUserFromDapp(ADDRESS1_DAPP)
	assert.Equal(t, "100", values[GENESIS_ADDRESS_HEX])

	values = storedValue.Contract.NamedKeys.GetVotingUserFromDapp(address2_dapp)
	assert.Equal(t, "200", values[GENESIS_ADDRESS_HEX])

	values = storedValue.Contract.NamedKeys.GetVotingUserFromDapp(address3_dapp)
	assert.Equal(t, "300", values[GENESIS_ADDRESS_HEX])
}

func TestStepAndCommission(t *testing.T) {
	client, rootStateHash, proxyHash, protocolVersion := InitalRunGenensis(DEFAULT_GENESIS_ACCOUNT, DEFAULT_GENESIS_STATE_INFO)
	amount := "1000000000000000000"

	// ready to create block
	rootStateHash, _ = RunTransferToAccount(client, rootStateHash, GENESIS_ADDRESS, SYSTEM_ACCOUNT, "1000000000000000000", proxyHash, protocolVersion)

	rootStateHash, _ = RunTransferToAccount(client, rootStateHash, GENESIS_ADDRESS, ADDRESS1, amount, proxyHash, protocolVersion)
	rootStateHash, _ = RunDelegate(client, rootStateHash, ADDRESS1, GENESIS_ADDRESS, "100000000000000000", proxyHash, protocolVersion)

	beforeAddress1Amount, errMessage := grpc.QueryBalance(client, rootStateHash, ADDRESS1, protocolVersion)
	assert.Equal(t, "", errMessage)

	beforeGenesisAmount, errMessage := grpc.QueryBalance(client, rootStateHash, GENESIS_ADDRESS, protocolVersion)
	assert.Equal(t, "", errMessage)

	// step
	rootStateHash, _ = RunStep(client, rootStateHash, SYSTEM_ACCOUNT, proxyHash, protocolVersion)

	afterStepAddress1Amount, errMessage := grpc.QueryBalance(client, rootStateHash, ADDRESS1, protocolVersion)
	assert.Equal(t, "", errMessage)
	assert.Equal(t, beforeAddress1Amount, afterStepAddress1Amount)

	afterStepGenesisAmount, errMessage := grpc.QueryBalance(client, rootStateHash, GENESIS_ADDRESS, protocolVersion)
	assert.Equal(t, beforeGenesisAmount, afterStepGenesisAmount)

	// claim Commission
	rootStateHash, _ = RunClaimCommission(client, rootStateHash, GENESIS_ADDRESS, proxyHash, protocolVersion)
	rootStateHash, _ = RunStep(client, rootStateHash, SYSTEM_ACCOUNT, proxyHash, protocolVersion)

	afterClaimCommissionAddress1Amount, errMessage := grpc.QueryBalance(client, rootStateHash, ADDRESS1, protocolVersion)
	assert.Equal(t, "", errMessage)
	assert.Equal(t, beforeAddress1Amount, afterClaimCommissionAddress1Amount)

	afterClaimCommissionGenesisAmount, errMessage := grpc.QueryBalance(client, rootStateHash, GENESIS_ADDRESS, protocolVersion)
	assert.NotEqual(t, beforeGenesisAmount, afterClaimCommissionGenesisAmount)

	// claim Reward
	rootStateHash, _ = RunClaimReward(client, rootStateHash, GENESIS_ADDRESS, proxyHash, protocolVersion)
	rootStateHash, _ = RunClaimReward(client, rootStateHash, ADDRESS1, proxyHash, protocolVersion)
	rootStateHash, _ = RunStep(client, rootStateHash, SYSTEM_ACCOUNT, proxyHash, protocolVersion)

	afterClaimRewardAddress1Amount, errMessage := grpc.QueryBalance(client, rootStateHash, ADDRESS1, protocolVersion)
	assert.Equal(t, "", errMessage)
	assert.NotEqual(t, afterClaimCommissionAddress1Amount, afterClaimRewardAddress1Amount)

	afterClaimRewardGenesisAmount, errMessage := grpc.QueryBalance(client, rootStateHash, GENESIS_ADDRESS, protocolVersion)
	assert.Equal(t, "", errMessage)
	assert.NotEqual(t, afterClaimCommissionGenesisAmount, afterClaimRewardGenesisAmount)
}

func TestClaimAmount(t *testing.T) {
	client, rootStateHash, proxyHash, protocolVersion := InitalRunGenensis(DEFAULT_GENESIS_ACCOUNT, DEFAULT_GENESIS_STATE_INFO)
	amount := "1000000000000000000"

	// ready to create block
	rootStateHash, _ = RunTransferToAccount(client, rootStateHash, GENESIS_ADDRESS, SYSTEM_ACCOUNT, "1000000000000000000", proxyHash, protocolVersion)

	rootStateHash, _ = RunTransferToAccount(client, rootStateHash, GENESIS_ADDRESS, ADDRESS1, amount, proxyHash, protocolVersion)
	rootStateHash, _ = RunDelegate(client, rootStateHash, ADDRESS1, GENESIS_ADDRESS, "100000000000000000", proxyHash, protocolVersion)

	// 10 step
	for i := 0; i < 10; i++ {
		rootStateHash, _ = RunStep(client, rootStateHash, SYSTEM_ACCOUNT, proxyHash, protocolVersion)
	}

	storedValue := RunQuery(client, rootStateHash, "address", SYSTEM_ACCOUNT, []string{"pos"}, protocolVersion)
	step10Address1Reward := storedValue.Contract.NamedKeys.GetUserReward(ADDRESS1)
	step10GenesisAddressReward := storedValue.Contract.NamedKeys.GetUserReward(GENESIS_ADDRESS)
	step10GenesisAddressCommission := storedValue.Contract.NamedKeys.GetValidatorCommission(GENESIS_ADDRESS)

	assert.NotEqual(t, "", step10Address1Reward)
	assert.NotEqual(t, "", step10GenesisAddressReward)
	assert.NotEqual(t, "", step10GenesisAddressCommission)

	rootStateHash, _ = RunClaimCommission(client, rootStateHash, GENESIS_ADDRESS, proxyHash, protocolVersion)
	rootStateHash, _ = RunClaimReward(client, rootStateHash, GENESIS_ADDRESS, proxyHash, protocolVersion)
	rootStateHash, _ = RunClaimReward(client, rootStateHash, ADDRESS1, proxyHash, protocolVersion)
	rootStateHash, _ = RunStep(client, rootStateHash, SYSTEM_ACCOUNT, proxyHash, protocolVersion)

	storedValue = RunQuery(client, rootStateHash, "address", SYSTEM_ACCOUNT, []string{"pos"}, protocolVersion)
	stepAddress1Reward := storedValue.Contract.NamedKeys.GetUserReward(ADDRESS1)
	stepGenesisAddressReward := storedValue.Contract.NamedKeys.GetUserReward(GENESIS_ADDRESS)
	stepGenesisAddressCommission := storedValue.Contract.NamedKeys.GetValidatorCommission(GENESIS_ADDRESS)

	assert.NotEqual(t, len(step10Address1Reward)+1, len(stepAddress1Reward))
	assert.NotEqual(t, len(step10GenesisAddressReward)+1, len(stepGenesisAddressReward))
	assert.NotEqual(t, len(step10GenesisAddressCommission)+1, len(stepGenesisAddressCommission))
}

func TestImport(t *testing.T) {
	genesisAccounts := []*ipc.ChainSpec_GenesisAccount{
		{
			PublicKey:    GENESIS_ADDRESS,
			Balance:      &state.BigInt{Value: INITIAL_BALANCE, BitWidth: 512},
			BondedAmount: &state.BigInt{Value: INITIAL_BOND_AMOUNT, BitWidth: 512},
		},
		{
			PublicKey:    ADDRESS1,
			Balance:      &state.BigInt{Value: INITIAL_BALANCE, BitWidth: 512},
			BondedAmount: &state.BigInt{Value: INITIAL_BOND_AMOUNT, BitWidth: 512},
		},
	}

	DELEAGE_AMOUNT_ADDRESS1_FROM_GENESIS_ADDR := "100000000000000000"
	SELF_DELEAGE_AMOUNT_FROM_ADDRESS1 := "900000000000000000"
	VOTE_AMOUNT_FROM_GENESIS_ADDR := "20000000000000000"
	REWARD_FROM_GENESIS_ADDR := "123"
	COMMISSION_FROM_GENSIS_ADDR := "456"
	stateInfos := []string{
		"d_" + GENESIS_ADDRESS_HEX + "_" + GENESIS_ADDRESS_HEX + "_" + INITIAL_BOND_AMOUNT,
		"d_" + GENESIS_ADDRESS_HEX + "_" + ADDRESS1_HEX + "_" + DELEAGE_AMOUNT_ADDRESS1_FROM_GENESIS_ADDR,
		"d_" + ADDRESS1_HEX + "_" + ADDRESS1_HEX + "_" + SELF_DELEAGE_AMOUNT_FROM_ADDRESS1,
		"a_" + GENESIS_ADDRESS_HEX + "_" + DAPP_HASH_HEX + "_" + VOTE_AMOUNT_FROM_GENESIS_ADDR,
		"r_" + GENESIS_ADDRESS_HEX + "_" + REWARD_FROM_GENESIS_ADDR,
		"c_" + GENESIS_ADDRESS_HEX + "_" + COMMISSION_FROM_GENSIS_ADDR,
	}

	client, rootStateHash, _, protocolVersion := InitalRunGenensis(genesisAccounts, stateInfos)

	storedValue := RunQuery(client, rootStateHash, "address", SYSTEM_ACCOUNT, []string{"pos"}, protocolVersion)

	genesisAddressDelegateInfo := storedValue.Contract.NamedKeys.GetDelegateFromDelegator(GENESIS_ADDRESS)
	assert.Equal(t, 2, len(genesisAddressDelegateInfo))
	assert.Equal(t, INITIAL_BOND_AMOUNT, genesisAddressDelegateInfo[GENESIS_ADDRESS_HEX])
	assert.Equal(t, DELEAGE_AMOUNT_ADDRESS1_FROM_GENESIS_ADDR, genesisAddressDelegateInfo[ADDRESS1_HEX])

	address1DelegateInfo := storedValue.Contract.NamedKeys.GetDelegateFromDelegator(ADDRESS1)
	assert.Equal(t, SELF_DELEAGE_AMOUNT_FROM_ADDRESS1, address1DelegateInfo[ADDRESS1_HEX])

	genesisAddressVoteInfo := storedValue.Contract.NamedKeys.GetVotingDappFromUser(GENESIS_ADDRESS)
	assert.Equal(t, VOTE_AMOUNT_FROM_GENESIS_ADDR, genesisAddressVoteInfo[DAPP_HASH_HEX])

	assert.Equal(t, REWARD_FROM_GENESIS_ADDR, storedValue.Contract.NamedKeys.GetUserReward(GENESIS_ADDRESS))
	assert.Equal(t, COMMISSION_FROM_GENSIS_ADDR, storedValue.Contract.NamedKeys.GetValidatorCommission(GENESIS_ADDRESS))
}

func TestStandardPayment(t *testing.T) {
	client, rootStateHash, proxyHash, protocolVersion := InitalRunGenensis(DEFAULT_GENESIS_ACCOUNT, DEFAULT_GENESIS_STATE_INFO)

	paymentStr := GetPaymentArgsJson("1000000000000000000")

	rootStateHash, _ = RunExecute(client, rootStateHash, GENESIS_ADDRESS, util.HASH, proxyHash, paymentStr, proxyHash, "1000000000000000000", protocolVersion)

	queryResult, errMessage := grpc.QueryBalance(client, rootStateHash, GENESIS_ADDRESS, protocolVersion)
	assert.Equal(t, "3000000000000000000", queryResult)
	assert.Equal(t, "", errMessage)
}
