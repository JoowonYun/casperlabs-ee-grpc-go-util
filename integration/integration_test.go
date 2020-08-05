package integration

import (
	"testing"
)

// func TestCustomContractCounter(t *testing.T) {
// 	client, rootStateHash, proxyHash, protocolVersion := InitalRunGenensis(DEFAULT_GENESIS_ACCOUNT)

// 	// counterDefine
// 	rootStateHash, _ = RunCounterDefine(client, rootStateHash, GENESIS_ADDRESS, proxyHash, protocolVersion)

// 	// query
// 	storedValue := RunQuery(client, rootStateHash, "address", GENESIS_ADDRESS, []string{"counter", "count"}, protocolVersion)
// 	assert.Equal(t, int32(0), storedValue.ClValue.ToStateValues().GetIntValue())

// 	// First counter call
// 	rootStateHash, _ = RunCounterCall(client, rootStateHash, GENESIS_ADDRESS, proxyHash, protocolVersion)

// 	// query
// 	storedValue = RunQuery(client, rootStateHash, "address", GENESIS_ADDRESS, []string{"counter", "count"}, protocolVersion)
// 	assert.Equal(t, int32(1), storedValue.ClValue.ToStateValues().GetIntValue())

// 	// Second counter call
// 	rootStateHash, _ = RunCounterCall(client, rootStateHash, GENESIS_ADDRESS, proxyHash, protocolVersion)

// 	// query
// 	storedValue = RunQuery(client, rootStateHash, "address", GENESIS_ADDRESS, []string{"counter", "count"}, protocolVersion)
// 	assert.Equal(t, int32(2), storedValue.ClValue.ToStateValues().GetIntValue())
// }

func TestTransferToAccount(t *testing.T) {
	client, rootStateHash, _, protocolVersion := InitalRunGenensis(DEFAULT_GENESIS_ACCOUNT)
	amount := "1"

	for i := 0; i < 10000; i++ {
		rootStateHash, _ = RunTransferToAccount(client, rootStateHash, GENESIS_ADDRESS, ADDRESS1, amount, proxyHash, protocolVersion)
	}

	// queryResult, errMessage := grpc.QueryBalance(client, rootStateHash, GENESIS_ADDRESS, protocolVersion)
	// assert.Equal(t, amount, queryResult)
	// assert.Equal(t, "", errMessage)
}

// func TestBondAndUnbond(t *testing.T) {
// 	client, rootStateHash, proxyHash, protocolVersion := InitalRunGenensis(DEFAULT_GENESIS_ACCOUNT)
// 	bondAmount := "10000000000000000"

// 	// bond
// 	rootStateHash, bonds := RunBond(client, rootStateHash, GENESIS_ADDRESS, bondAmount, proxyHash, protocolVersion)
// 	assert.Equal(t, "1000000000000000000", bonds[0].GetStake().GetValue())

// 	stakeAmount, errMsg := grpc.QueryStake(client, rootStateHash, GENESIS_ADDRESS, protocolVersion)
// 	if errMsg != "" {
// 		panic(errMsg)
// 	}

// 	assert.Equal(t, "1010000000000000000", stakeAmount)

// 	// unbond
// 	unbondAmount := "10000000000000000"
// 	rootStateHash, bonds = RunUnbond(client, rootStateHash, GENESIS_ADDRESS, unbondAmount, proxyHash, protocolVersion)
// 	rootStateHash, bonds = RunStep(client, rootStateHash, SYSTEM_ACCOUNT, proxyHash, protocolVersion)
// 	assert.Equal(t, "1000000000000000000", bonds[0].GetStake().GetValue())

// 	stakeAmount, errMsg = grpc.QueryStake(client, rootStateHash, GENESIS_ADDRESS, protocolVersion)
// 	if errMsg != "" {
// 		panic(errMsg)
// 	}
// }

// func TestDelegate(t *testing.T) {
// 	client, rootStateHash, proxyHash, protocolVersion := InitalRunGenensis(DEFAULT_GENESIS_ACCOUNT)
// 	amount := "100"

// 	rootStateHash, bonds := RunBond(client, rootStateHash, GENESIS_ADDRESS, amount, proxyHash, protocolVersion)
// 	rootStateHash, bonds = RunDelegate(client, rootStateHash, GENESIS_ADDRESS, GENESIS_ADDRESS, amount, proxyHash, protocolVersion)
// 	assert.Equal(t, "1000000000000000100", bonds[0].GetStake().GetValue())

// 	storedValue := RunQuery(client, rootStateHash, "address", SYSTEM_ACCOUNT, []string{"pos"}, protocolVersion)
// 	delegators := storedValue.Contract.NamedKeys.GetDelegateFromValidator(GENESIS_ADDRESS)
// 	assert.Equal(t, 1, len(delegators))
// 	assert.Equal(t, "1000000000000000100", delegators[GENESIS_ADDRESS_HEX])
// }

// func TestDelegateFromAnotherAddress(t *testing.T) {
// 	client, rootStateHash, proxyHash, protocolVersion := InitalRunGenensis(DEFAULT_GENESIS_ACCOUNT)
// 	amount := "10000000000000000000"
// 	delegateAmount := "1000000000000000"

// 	rootStateHash, _ = RunTransferToAccount(client, rootStateHash, GENESIS_ADDRESS, ADDRESS1, amount, proxyHash, protocolVersion)
// 	balance, errMessage := grpc.QueryBalance(client, rootStateHash, ADDRESS1, protocolVersion)
// 	assert.Equal(t, amount, balance)
// 	assert.Equal(t, "", errMessage)

// 	rootStateHash, _ = RunBond(client, rootStateHash, ADDRESS1, delegateAmount, proxyHash, protocolVersion)
// 	rootStateHash, _ = RunDelegate(client, rootStateHash, ADDRESS1, GENESIS_ADDRESS, delegateAmount, proxyHash, protocolVersion)

// 	storedValue := RunQuery(client, rootStateHash, "address", SYSTEM_ACCOUNT, []string{"pos"}, protocolVersion)
// 	delegators := storedValue.Contract.NamedKeys.GetDelegateFromValidator(GENESIS_ADDRESS)
// 	assert.Equal(t, 2, len(delegators))
// 	assert.Equal(t, "1000000000000000000", delegators[GENESIS_ADDRESS_HEX])
// 	assert.Equal(t, "1000000000000000", delegators[ADDRESS1_HEX])
// }

// func TestUndelegation(t *testing.T) {
// 	client, rootStateHash, proxyHash, protocolVersion := InitalRunGenensis(DEFAULT_GENESIS_ACCOUNT)
// 	amount := "100"

// 	rootStateHash, bonds := RunUndelegate(client, rootStateHash, GENESIS_ADDRESS, GENESIS_ADDRESS, amount, proxyHash, protocolVersion)
// 	rootStateHash, bonds = RunStep(client, rootStateHash, SYSTEM_ACCOUNT, proxyHash, protocolVersion)
// 	assert.Equal(t, 1, len(bonds))
// }

// func TestRedelegation(t *testing.T) {
// 	client, rootStateHash, proxyHash, protocolVersion := InitalRunGenensis(DEFAULT_GENESIS_ACCOUNT)
// 	amount := "100"

// 	rootStateHash, bonds := RunRedelegate(client, rootStateHash, GENESIS_ADDRESS, GENESIS_ADDRESS, ADDRESS1, amount, proxyHash, protocolVersion)
// 	rootStateHash, bonds = RunStep(client, rootStateHash, SYSTEM_ACCOUNT, proxyHash, protocolVersion)
// 	assert.Equal(t, 2, len(bonds))
// }

// func TestVoteAndUnvote(t *testing.T) {
// 	client, rootStateHash, proxyHash, protocolVersion := InitalRunGenensis(DEFAULT_GENESIS_ACCOUNT)
// 	voteAmount := "123"

// 	rootStateHash, _ = RunVote(client, rootStateHash, GENESIS_ADDRESS, ADDRESS1, voteAmount, proxyHash, protocolVersion)

// 	votedAmount, errMsg := grpc.QueryVoted(client, rootStateHash, ADDRESS1_DAPP, protocolVersion)
// 	if errMsg != "" {
// 		panic(errMsg)
// 	}
// 	assert.Equal(t, voteAmount, votedAmount)

// 	votingAmount, errMsg := grpc.QueryVoting(client, rootStateHash, GENESIS_ADDRESS, protocolVersion)
// 	if errMsg != "" {
// 		panic(errMsg)
// 	}
// 	assert.Equal(t, voteAmount, votingAmount)

// 	unvoteAmount := "23"
// 	rootStateHash, _ = RunUnvote(client, rootStateHash, GENESIS_ADDRESS, ADDRESS1, unvoteAmount, proxyHash, protocolVersion)
// 	votedAmount, errMsg = grpc.QueryVoted(client, rootStateHash, ADDRESS1_DAPP, protocolVersion)
// 	if errMsg != "" {
// 		panic(errMsg)
// 	}
// 	assert.Equal(t, "100", votedAmount)

// 	votingAmount, errMsg = grpc.QueryVoting(client, rootStateHash, GENESIS_ADDRESS, protocolVersion)
// 	if errMsg != "" {
// 		panic(errMsg)
// 	}
// 	assert.Equal(t, "100", votingAmount)
// }

// func TestVoteMoreAccount(t *testing.T) {
// 	client, rootStateHash, proxyHash, protocolVersion := InitalRunGenensis(DEFAULT_GENESIS_ACCOUNT)
// 	address2 := util.DecodeHexString("03170a2e7597b7b7e3d84c05391d139a62b157e78786d8c082f29dcf4c111314")
// 	address3 := util.DecodeHexString("f0f84944e0ccfa9e67383e6a448291787d208c8e46adc849f714078663d1dd36")

// 	address2_dapp := util.DecodeHexString("0103170a2e7597b7b7e3d84c05391d139a62b157e78786d8c082f29dcf4c111314")
// 	address3_dapp := util.DecodeHexString("01f0f84944e0ccfa9e67383e6a448291787d208c8e46adc849f714078663d1dd36")

// 	amount1 := "100"
// 	amount2 := "200"
// 	amount3 := "300"

// 	rootStateHash, _ = RunVote(client, rootStateHash, GENESIS_ADDRESS, ADDRESS1, amount1, proxyHash, protocolVersion)
// 	rootStateHash, _ = RunVote(client, rootStateHash, GENESIS_ADDRESS, address2, amount2, proxyHash, protocolVersion)
// 	rootStateHash, _ = RunVote(client, rootStateHash, GENESIS_ADDRESS, address3, amount3, proxyHash, protocolVersion)

// 	address1DappVoterAmount, errMsg := grpc.QueryVoted(client, rootStateHash, ADDRESS1_DAPP, protocolVersion)
// 	if errMsg != "" {
// 		panic(errMsg)
// 	}
// 	assert.Equal(t, "100", address1DappVoterAmount)

// 	address2DappVoterAmount, errMsg := grpc.QueryVoted(client, rootStateHash, address2_dapp, protocolVersion)
// 	if errMsg != "" {
// 		panic(errMsg)
// 	}
// 	assert.Equal(t, "200", address2DappVoterAmount)
// 	address3DappVoterAmount, errMsg := grpc.QueryVoted(client, rootStateHash, address3_dapp, protocolVersion)
// 	if errMsg != "" {
// 		panic(errMsg)
// 	}
// 	assert.Equal(t, "300", address3DappVoterAmount)

// 	genensisVoteAmount, errMsg := grpc.QueryVoting(client, rootStateHash, GENESIS_ADDRESS, protocolVersion)
// 	if errMsg != "" {
// 		panic(errMsg)
// 	}
// 	assert.Equal(t, "600", genensisVoteAmount)
// }

// func TestStepAndClaim(t *testing.T) {
// 	client, rootStateHash, proxyHash, protocolVersion := InitalRunGenensis(DEFAULT_GENESIS_ACCOUNT)
// 	amount := "1000000000000000000"
// 	stakeAmount := "100000000000000000"

// 	// delegate from ADDRESS1 to GENESIS_ADDRESS
// 	rootStateHash, _ = RunTransferToAccount(client, rootStateHash, GENESIS_ADDRESS, ADDRESS1, amount, proxyHash, protocolVersion)
// 	rootStateHash, _ = RunBond(client, rootStateHash, ADDRESS1, stakeAmount, proxyHash, protocolVersion)
// 	rootStateHash, _ = RunDelegate(client, rootStateHash, ADDRESS1, GENESIS_ADDRESS, stakeAmount, proxyHash, protocolVersion)
// 	beforeAddress1Amount, errMessage := grpc.QueryBalance(client, rootStateHash, ADDRESS1, protocolVersion)
// 	assert.Equal(t, "", errMessage)
// 	beforeGenesisAmount, errMessage := grpc.QueryBalance(client, rootStateHash, GENESIS_ADDRESS, protocolVersion)
// 	assert.Equal(t, "", errMessage)

// 	// run step 10 times
// 	for i := 0; i < 10; i++ {
// 		rootStateHash, _ = RunStep(client, rootStateHash, SYSTEM_ACCOUNT, proxyHash, protocolVersion)
// 	}
// 	afterStepAddress1Amount, errMessage := grpc.QueryBalance(client, rootStateHash, ADDRESS1, protocolVersion)
// 	assert.Equal(t, "", errMessage)
// 	assert.Equal(t, beforeAddress1Amount, afterStepAddress1Amount)
// 	afterStepGenesisAmount, errMessage := grpc.QueryBalance(client, rootStateHash, GENESIS_ADDRESS, protocolVersion)
// 	assert.Equal(t, "", errMessage)
// 	assert.Equal(t, beforeGenesisAmount, afterStepGenesisAmount)

// 	// check reward and commission is generated
// 	step10Address1Reward, errMessage := grpc.QueryReward(client, rootStateHash, ADDRESS1, protocolVersion)
// 	assert.Equal(t, "", errMessage)
// 	assert.NotEqual(t, "", step10Address1Reward)
// 	step10GenesisAddressReward, errMessage := grpc.QueryReward(client, rootStateHash, GENESIS_ADDRESS, protocolVersion)
// 	assert.Equal(t, "", errMessage)
// 	assert.NotEqual(t, "", step10GenesisAddressReward)
// 	step10GenesisAddressCommission, errMessage := grpc.QueryCommission(client, rootStateHash, GENESIS_ADDRESS, protocolVersion)
// 	assert.Equal(t, "", errMessage)
// 	assert.NotEqual(t, "", step10GenesisAddressCommission)

// 	// claim Commission
// 	rootStateHash, _ = RunClaimCommission(client, rootStateHash, GENESIS_ADDRESS, proxyHash, protocolVersion)
// 	afterClaimCommissionAddress1Amount, errMessage := grpc.QueryBalance(client, rootStateHash, ADDRESS1, protocolVersion)
// 	assert.Equal(t, "", errMessage)
// 	assert.Equal(t, beforeAddress1Amount, afterClaimCommissionAddress1Amount)
// 	afterClaimCommissionGenesisAmount, errMessage := grpc.QueryBalance(client, rootStateHash, GENESIS_ADDRESS, protocolVersion)
// 	assert.Equal(t, "", errMessage)
// 	assert.NotEqual(t, beforeGenesisAmount, afterClaimCommissionGenesisAmount)

// 	// claim Reward
// 	rootStateHash, _ = RunClaimReward(client, rootStateHash, GENESIS_ADDRESS, proxyHash, protocolVersion)
// 	rootStateHash, _ = RunClaimReward(client, rootStateHash, ADDRESS1, proxyHash, protocolVersion)
// 	afterClaimRewardAddress1Amount, errMessage := grpc.QueryBalance(client, rootStateHash, ADDRESS1, protocolVersion)
// 	assert.Equal(t, "", errMessage)
// 	assert.NotEqual(t, afterClaimCommissionAddress1Amount, afterClaimRewardAddress1Amount)
// 	afterClaimRewardGenesisAmount, errMessage := grpc.QueryBalance(client, rootStateHash, GENESIS_ADDRESS, protocolVersion)
// 	assert.Equal(t, "", errMessage)
// 	assert.NotEqual(t, afterClaimCommissionGenesisAmount, afterClaimRewardGenesisAmount)

// 	// check reward and commission is claimed
// 	afterClaimAddress1Reward, errMessage := grpc.QueryReward(client, rootStateHash, ADDRESS1, protocolVersion)
// 	assert.Equal(t, "", errMessage)
// 	assert.Equal(t, "", errMessage)
// 	afterClaimGenesisAddressReward, errMessage := grpc.QueryReward(client, rootStateHash, GENESIS_ADDRESS, protocolVersion)
// 	assert.Equal(t, "0", afterClaimAddress1Reward)
// 	assert.Equal(t, "0", afterClaimGenesisAddressReward)
// 	afterClaimGenesisAddressCommission, errMessage := grpc.QueryCommission(client, rootStateHash, GENESIS_ADDRESS, protocolVersion)
// 	assert.Equal(t, "", errMessage)
// 	assert.Equal(t, "0", afterClaimGenesisAddressCommission)
// }

// // Do net support import & export
// func xTestImport(t *testing.T) {
// 	genesisAccounts := []*ipc.ChainSpec_GenesisAccount{
// 		{
// 			PublicKey:    GENESIS_ADDRESS,
// 			Balance:      &state.BigInt{Value: INITIAL_BALANCE, BitWidth: 512},
// 			BondedAmount: &state.BigInt{Value: INITIAL_BOND_AMOUNT, BitWidth: 512},
// 		},
// 		{
// 			PublicKey:    ADDRESS1,
// 			Balance:      &state.BigInt{Value: INITIAL_BALANCE, BitWidth: 512},
// 			BondedAmount: &state.BigInt{Value: INITIAL_BOND_AMOUNT, BitWidth: 512},
// 		},
// 	}

// 	DELEAGE_AMOUNT_ADDRESS1_FROM_GENESIS_ADDR := "100000000000000000"
// 	SELF_DELEAGE_AMOUNT_FROM_ADDRESS1 := "900000000000000000"
// 	VOTE_AMOUNT_FROM_GENESIS_ADDR := "20000000000000000"
// 	REWARD_FROM_GENESIS_ADDR := "123"
// 	COMMISSION_FROM_GENSIS_ADDR := "456"

// 	client, rootStateHash, _, protocolVersion := InitalRunGenensis(genesisAccounts)

// 	storedValue := RunQuery(client, rootStateHash, "address", SYSTEM_ACCOUNT, []string{"pos"}, protocolVersion)

// 	genesisAddressDelegateInfo := storedValue.Contract.NamedKeys.GetDelegateFromDelegator(GENESIS_ADDRESS)
// 	assert.Equal(t, 2, len(genesisAddressDelegateInfo))
// 	assert.Equal(t, INITIAL_BOND_AMOUNT, genesisAddressDelegateInfo[GENESIS_ADDRESS_HEX])
// 	assert.Equal(t, DELEAGE_AMOUNT_ADDRESS1_FROM_GENESIS_ADDR, genesisAddressDelegateInfo[ADDRESS1_HEX])

// 	address1DelegateInfo := storedValue.Contract.NamedKeys.GetDelegateFromDelegator(ADDRESS1)
// 	assert.Equal(t, SELF_DELEAGE_AMOUNT_FROM_ADDRESS1, address1DelegateInfo[ADDRESS1_HEX])

// 	genesisAddressVoteInfo := storedValue.Contract.NamedKeys.GetVotingDappFromUser(GENESIS_ADDRESS)
// 	assert.Equal(t, VOTE_AMOUNT_FROM_GENESIS_ADDR, genesisAddressVoteInfo[DAPP_HASH_HEX])

// 	assert.Equal(t, REWARD_FROM_GENESIS_ADDR, storedValue.Contract.NamedKeys.GetUserReward(GENESIS_ADDRESS))
// 	assert.Equal(t, COMMISSION_FROM_GENSIS_ADDR, storedValue.Contract.NamedKeys.GetValidatorCommission(GENESIS_ADDRESS))
// }

// func TestStandardPayment(t *testing.T) {
// 	client, rootStateHash, proxyHash, protocolVersion := InitalRunGenensis(DEFAULT_GENESIS_ACCOUNT)

// 	paymentStr := GetPaymentArgsJson(BASIC_FEE)

// 	rootStateHash, _ = RunExecute(client, rootStateHash, GENESIS_ADDRESS, util.HASH, proxyHash, paymentStr, proxyHash, "1000000000000000000", protocolVersion)

// 	queryResult, errMessage := grpc.QueryBalance(client, rootStateHash, GENESIS_ADDRESS, protocolVersion)
// 	assert.Equal(t, "49998900000000000000000", queryResult)
// 	assert.Equal(t, "", errMessage)
// }
