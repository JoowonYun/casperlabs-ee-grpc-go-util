package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/hdac-io/casperlabs-ee-grpc-go-util/grpc"
	"github.com/hdac-io/casperlabs-ee-grpc-go-util/protobuf/io/casperlabs/casper/consensus"
	"github.com/hdac-io/casperlabs-ee-grpc-go-util/protobuf/io/casperlabs/casper/consensus/state"
	"github.com/hdac-io/casperlabs-ee-grpc-go-util/protobuf/io/casperlabs/ipc"
	"github.com/hdac-io/casperlabs-ee-grpc-go-util/protobuf/io/casperlabs/ipc/transforms"
	"github.com/hdac-io/casperlabs-ee-grpc-go-util/storedvalue"
	"github.com/hdac-io/casperlabs-ee-grpc-go-util/util"
)

func main() {
	// Init variable
	emptyStateHash := util.DecodeHexString(util.StrEmptyStateHash)
	rootStateHash := emptyStateHash
	genesisAddress := util.DecodeHexString("d70243dd9d0d646fd6df282a8f7a8fa05a6629bec01d8024c3611eb1c1fb9f84")
	systemContract := make([]byte, 32)

	chainName := "hdac"
	costs := map[string]uint32{
		"regular":            1,
		"div-multiplier":     16,
		"mul-multiplier":     4,
		"mem-multiplier":     2,
		"mem-initial-pages":  4096,
		"mem-grow-per-page":  8192,
		"mem-copy-per-byte":  1,
		"max-stack-height":   65536,
		"opcodes-multiplier": 3,
		"opcodes-divisor":    8}

	protocolVersion := util.MakeProtocolVersion(1, 0, 0)

	// Connect to ee sock.
	socketPath := os.Getenv("HOME") + `/.casperlabs/.casper-node.sock`
	client := grpc.Connect(socketPath)

	// laod wasm code
	cntDefCode, cntCallCode := loadWasmCode()

	// run genesis
	println(`RunGenesis`)
	genesisConfig, err := util.GenesisConfigMock(
		chainName, genesisAddress, "5000000000000000000", "1000000000000000000", protocolVersion, costs,
		"./example/contracts/hdac_mint_install.wasm", "./example/contracts/pop_install.wasm")
	if err != nil {
		fmt.Printf("Bad GenesisConfigMock err : %v", err)
		return
	}

	response, err := grpc.RunGenesis(client, genesisConfig)
	if err != nil {
		panic(err)
	}

	var parentStateHash []byte
	var effects []*transforms.TransformEntry

	switch response.GetResult().(type) {
	case *ipc.GenesisResponse_Success:
		parentStateHash = response.GetSuccess().GetPoststateHash()
		effects = response.GetSuccess().GetEffect().GetTransformMap()
	case *ipc.GenesisResponse_FailedDeploy:
		panic(response.GetFailedDeploy().GetMessage())
	}

	postStateHash, bonds, errMessage := grpc.Commit(client, rootStateHash, effects, protocolVersion)
	if bytes.Equal(postStateHash, parentStateHash) {
		rootStateHash = postStateHash
	}
	printCommitResult(rootStateHash, bonds)

	queryResult10, errMessage := grpc.Query(client, rootStateHash, "address", systemContract, []string{}, protocolVersion)
	var storedValue storedvalue.StoredValue
	storedValue.FromBytes(queryResult10)
	storedValue, err, _ = storedValue.FromBytes(queryResult10)
	if err != nil {
		panic(err)
	}

	proxyHash := storedValue.Account.NamedKeys[0].Key.Hash
	println("Proxy hash : " + util.EncodeToHexString(proxyHash))
	println(errMessage)

	// Run "Counter Define contract"
	println(`Counter Define`)
	timestamp := time.Now().Unix()
	paymentArgs := []*consensus.Deploy_Arg{
		&consensus.Deploy_Arg{
			Name: "method",
			Value: &consensus.Deploy_Arg_Value{
				Value: &consensus.Deploy_Arg_Value_StringValue{
					StringValue: "standard_payment"}}},
		&consensus.Deploy_Arg{
			Name: "fee",
			Value: &consensus.Deploy_Arg_Value{
				Value: &consensus.Deploy_Arg_Value_BigInt{
					BigInt: &state.BigInt{Value: "10000000000000000", BitWidth: 512}}}}}
	paymentArgsStr, err := util.DeployArgsToJsonString(paymentArgs)
	if err != nil {
		panic(err)
	}
	deploy, _ := util.MakeDeploy(genesisAddress, util.WASM, cntDefCode, "", util.HASH, proxyHash, paymentArgsStr, uint64(10), timestamp, chainName)
	deploys := util.MakeInitDeploys()
	deploys = util.AddDeploy(deploys, deploy)

	res, err := grpc.Execute(client, rootStateHash, timestamp, deploys, protocolVersion)
	effects2, err := executeErrorHandler(res)

	postStateHash2, bonds2, errMessage := grpc.Commit(client, rootStateHash, effects2, protocolVersion)
	if errMessage != "" {
		panic(errMessage)
	}
	rootStateHash = postStateHash2
	printCommitResult(rootStateHash, bonds2)

	queryResult, errMessage := grpc.QueryBalance(client, rootStateHash, genesisAddress, protocolVersion)
	println(util.EncodeToHexString(genesisAddress), ": ", queryResult)
	println(errMessage)

	// Run "Counter Call contract"
	println(`Counter Call`)
	timestamp = time.Now().Unix()
	deploy, _ = util.MakeDeploy(genesisAddress, util.WASM, cntCallCode, "", util.HASH, proxyHash, paymentArgsStr, uint64(10), timestamp, chainName)
	deploys = util.MakeInitDeploys()
	deploys = util.AddDeploy(deploys, deploy)
	res, err = grpc.Execute(client, rootStateHash, timestamp, deploys, protocolVersion)
	effects3, err := executeErrorHandler(res)

	postStateHash3, bonds3, errMessage := grpc.Commit(client, rootStateHash, effects3, protocolVersion)
	if errMessage != "" {
		panic(errMessage)
	}
	rootStateHash = postStateHash3
	printCommitResult(rootStateHash, bonds3)

	// Query counter contract.
	println(`Counter Query`)
	path := []string{"counter", "count"}
	queryResult1, errMessage := grpc.Query(client, rootStateHash, "address", genesisAddress, path, protocolVersion)
	storedValue, err, _ = storedValue.FromBytes(queryResult1)
	if err != nil {
		panic(err)
	}
	println(storedValue.ClValue.ToStateValues().GetIntValue())
	println(errMessage)

	queryResult, errMessage = grpc.QueryBalance(client, rootStateHash, genesisAddress, protocolVersion)
	println(util.EncodeToHexString(genesisAddress), ": ", queryResult)
	println(errMessage)

	timestamp = time.Now().Unix()
	deploy, _ = util.MakeDeploy(genesisAddress, util.WASM, cntCallCode, "", util.HASH, proxyHash, paymentArgsStr, uint64(10), timestamp, chainName)
	deploys = util.MakeInitDeploys()
	deploys = util.AddDeploy(deploys, deploy)
	res, err = grpc.Execute(client, rootStateHash, timestamp, deploys, protocolVersion)
	effects4, err := executeErrorHandler(res)

	postStateHash4, bonds4, errMessage := grpc.Commit(client, rootStateHash, effects4, protocolVersion)
	rootStateHash = postStateHash4
	if errMessage != "" {
		panic(errMessage)
	}
	printCommitResult(rootStateHash, bonds4)

	queryResult2, errMessage := grpc.Query(client, rootStateHash, "address", genesisAddress, path, protocolVersion)
	storedValue, err, _ = storedValue.FromBytes(queryResult2)
	println(storedValue.ClValue.ToStateValues().GetIntValue())
	println(errMessage)

	queryResult3, errMessage := grpc.QueryBalance(client, rootStateHash, genesisAddress, protocolVersion)
	println(util.EncodeToHexString(genesisAddress), ": ", queryResult3)
	println(errMessage)

	// Run "Send transaction"
	println(`Send Transaction`)
	timestamp = time.Now().Unix()
	address1 := util.DecodeHexString("93236a9263d2ac6198c5ed211774c745d5dc62a910cb84276f8a7c4959208915")

	sessionArgs := []*consensus.Deploy_Arg{
		&consensus.Deploy_Arg{
			Name: "method",
			Value: &consensus.Deploy_Arg_Value{
				Value: &consensus.Deploy_Arg_Value_StringValue{
					StringValue: "transfer_to_account"}}},
		&consensus.Deploy_Arg{
			Name: "address",
			Value: &consensus.Deploy_Arg_Value{
				Value: &consensus.Deploy_Arg_Value_BytesValue{
					BytesValue: address1}}},
		&consensus.Deploy_Arg{
			Name: "amount",
			Value: &consensus.Deploy_Arg_Value{
				Value: &consensus.Deploy_Arg_Value_BigInt{
					BigInt: &state.BigInt{Value: "1000000000000000000", BitWidth: 512}}}},
	}
	sessionArgsStr, err := util.DeployArgsToJsonString(sessionArgs)
	if err != nil {
		panic(err)
	}
	deploy, _ = util.MakeDeploy(genesisAddress, util.HASH, proxyHash, sessionArgsStr, util.HASH, proxyHash, paymentArgsStr, uint64(10), timestamp, chainName)
	deploys = util.MakeInitDeploys()
	deploys = util.AddDeploy(deploys, deploy)
	res, err = grpc.Execute(client, rootStateHash, timestamp, deploys, protocolVersion)
	effects5, err := executeErrorHandler(res)

	postStateHash5, bonds5, errMessage := grpc.Commit(client, rootStateHash, effects5, protocolVersion)
	if errMessage != "" {
		panic(errMessage)
	}
	rootStateHash = postStateHash5
	printCommitResult(rootStateHash, bonds5)

	queryResult4, errMessage := grpc.QueryBalance(client, rootStateHash, genesisAddress, protocolVersion)
	println(util.EncodeToHexString(genesisAddress), ": ", queryResult4)
	println(errMessage)

	queryResult5, errMessage := grpc.QueryBalance(client, rootStateHash, address1, protocolVersion)
	println(util.EncodeToHexString(address1), ": ", queryResult5)
	println(errMessage)

	// bonding
	println(`Bonding`)
	timestamp = time.Now().Unix()
	bondingArgs := []*consensus.Deploy_Arg{
		&consensus.Deploy_Arg{
			Name: "method",
			Value: &consensus.Deploy_Arg_Value{
				Value: &consensus.Deploy_Arg_Value_StringValue{
					StringValue: "bond"}}},
		&consensus.Deploy_Arg{
			Name: "amount",
			Value: &consensus.Deploy_Arg_Value{
				Value: &consensus.Deploy_Arg_Value_BigInt{
					BigInt: &state.BigInt{Value: "10000000000000000", BitWidth: 512}}}},
	}
	bondingArgsStr, err := util.DeployArgsToJsonString(bondingArgs)
	if err != nil {
		panic(err)
	}
	deploy, _ = util.MakeDeploy(address1, util.HASH, proxyHash, bondingArgsStr, util.HASH, proxyHash, paymentArgsStr, uint64(10), timestamp, chainName)
	deploys = util.MakeInitDeploys()
	deploys = util.AddDeploy(deploys, deploy)
	res, err = grpc.Execute(client, rootStateHash, timestamp, deploys, protocolVersion)
	effects6, err := executeErrorHandler(res)
	postStateHash6, bonds6, errMessage := grpc.Commit(client, rootStateHash, effects6, protocolVersion)
	if errMessage != "" {
		panic(errMessage)
	}
	rootStateHash = postStateHash6
	printCommitResult(rootStateHash, bonds6)

	println(`delegateion`)
	delegationArgs := []*consensus.Deploy_Arg{
		&consensus.Deploy_Arg{
			Name: "method",
			Value: &consensus.Deploy_Arg_Value{
				Value: &consensus.Deploy_Arg_Value_StringValue{
					StringValue: "delegate"}}},
		&consensus.Deploy_Arg{
			Name: "validator",
			Value: &consensus.Deploy_Arg_Value{
				Value: &consensus.Deploy_Arg_Value_BytesValue{
					BytesValue: genesisAddress,
				},
			},
		},
		&consensus.Deploy_Arg{
			Name: "amount",
			Value: &consensus.Deploy_Arg_Value{
				Value: &consensus.Deploy_Arg_Value_BigInt{
					BigInt: &state.BigInt{
						Value:    "100",
						BitWidth: 512,
					},
				},
			},
		},
	}
	delegationArgsStr, err := util.DeployArgsToJsonString(delegationArgs)
	if err != nil {
		panic(err)
	}
	deploy, _ = util.MakeDeploy(address1, util.HASH, proxyHash, delegationArgsStr, util.HASH, proxyHash, paymentArgsStr, uint64(10), timestamp, chainName)
	deploys = util.MakeInitDeploys()
	deploys = util.AddDeploy(deploys, deploy)
	res, err = grpc.Execute(client, rootStateHash, timestamp, deploys, protocolVersion)
	effects, err = executeErrorHandler(res)
	if err != nil {
		println(err.Error())
	}
	postStateHash, bonds, errMessage = grpc.Commit(client, rootStateHash, effects, protocolVersion)
	if errMessage != "" {
		panic(errMessage)
	}
	rootStateHash = postStateHash
	printCommitResult(rootStateHash, bonds)

	queryResult11, errMessage := grpc.Query(client, rootStateHash, "address", address1, []string{"pos"}, protocolVersion)
	var storedValue1 storedvalue.StoredValue
	storedValue1.FromBytes(queryResult11)
	storedValue1, err, _ = storedValue1.FromBytes(queryResult11)
	if err != nil {
		panic(err)
	}

	println(`redelegateion`)
	redelegationArgs := []*consensus.Deploy_Arg{
		&consensus.Deploy_Arg{
			Name: "method",
			Value: &consensus.Deploy_Arg_Value{
				Value: &consensus.Deploy_Arg_Value_StringValue{
					StringValue: "redelegate"}}},
		&consensus.Deploy_Arg{
			Name: "src",
			Value: &consensus.Deploy_Arg_Value{
				Value: &consensus.Deploy_Arg_Value_BytesValue{
					BytesValue: genesisAddress,
				},
			},
		},
		&consensus.Deploy_Arg{
			Name: "dest",
			Value: &consensus.Deploy_Arg_Value{
				Value: &consensus.Deploy_Arg_Value_BytesValue{
					BytesValue: address1,
				},
			},
		},
		&consensus.Deploy_Arg{
			Name: "amount",
			Value: &consensus.Deploy_Arg_Value{
				Value: &consensus.Deploy_Arg_Value_BigInt{
					BigInt: &state.BigInt{
						Value:    "50",
						BitWidth: 512,
					},
				},
			},
		},
	}

	redelegationArgsStr, err := util.DeployArgsToJsonString(redelegationArgs)
	if err != nil {
		panic(err)
	}
	deploy, _ = util.MakeDeploy(address1, util.HASH, proxyHash, redelegationArgsStr, util.HASH, proxyHash, paymentArgsStr, uint64(10), timestamp, chainName)
	deploys = util.MakeInitDeploys()
	deploys = util.AddDeploy(deploys, deploy)
	res, err = grpc.Execute(client, rootStateHash, timestamp, deploys, protocolVersion)
	effects, err = executeErrorHandler(res)
	if err != nil {
		println(err.Error())
	}
	postStateHash, bonds, errMessage = grpc.Commit(client, rootStateHash, effects, protocolVersion)
	if errMessage != "" {
		panic(errMessage)
	}
	rootStateHash = postStateHash
	printCommitResult(rootStateHash, bonds)

	println(`undelegateion`)
	undelegationArgs := []*consensus.Deploy_Arg{
		&consensus.Deploy_Arg{
			Name: "method",
			Value: &consensus.Deploy_Arg_Value{
				Value: &consensus.Deploy_Arg_Value_StringValue{
					StringValue: "undelegate"}}},
		&consensus.Deploy_Arg{
			Name: "validator",
			Value: &consensus.Deploy_Arg_Value{
				Value: &consensus.Deploy_Arg_Value_BytesValue{
					BytesValue: address1,
				},
			},
		},
		&consensus.Deploy_Arg{
			Name: "amount",
			Value: &consensus.Deploy_Arg_Value{
				Value: &consensus.Deploy_Arg_Value_OptionalValue{
					OptionalValue: &consensus.Deploy_Arg_Value{
						Value: &consensus.Deploy_Arg_Value_BigInt{
							BigInt: &state.BigInt{
								Value:    "50",
								BitWidth: 512,
							},
						},
					},
				},
			},
		},
	}
	undelegationArgsStr, err := util.DeployArgsToJsonString(undelegationArgs)
	if err != nil {
		println(err)
	}
	deploy, _ = util.MakeDeploy(address1, util.HASH, proxyHash, undelegationArgsStr, util.HASH, proxyHash, paymentArgsStr, uint64(10), timestamp, chainName)
	deploys = util.MakeInitDeploys()
	deploys = util.AddDeploy(deploys, deploy)
	res, err = grpc.Execute(client, rootStateHash, timestamp, deploys, protocolVersion)
	effects, err = executeErrorHandler(res)
	if err != nil {
		println(err.Error())
	}
	if err != nil {
		println(err.Error())
	}
	postStateHash, bonds, errMessage = grpc.Commit(client, rootStateHash, effects, protocolVersion)
	if errMessage != "" {
		panic(errMessage)
	}
	rootStateHash = postStateHash
	printCommitResult(rootStateHash, bonds)

	// unbonding
	println(`Unbonding`)
	timestamp = time.Now().Unix()
	unbondingArgs := []*consensus.Deploy_Arg{
		&consensus.Deploy_Arg{
			Name: "method",
			Value: &consensus.Deploy_Arg_Value{
				Value: &consensus.Deploy_Arg_Value_StringValue{
					StringValue: "unbond"}}},
		&consensus.Deploy_Arg{
			Name: "amount",
			Value: &consensus.Deploy_Arg_Value{
				Value: &consensus.Deploy_Arg_Value_OptionalValue{
					OptionalValue: &consensus.Deploy_Arg_Value{
						Value: &consensus.Deploy_Arg_Value_BigInt{
							BigInt: &state.BigInt{Value: "100", BitWidth: 512}}}}}},
	}
	unbondingArgsStr, err := util.DeployArgsToJsonString(unbondingArgs)
	if err != nil {
		panic(err)
	}
	deploy, _ = util.MakeDeploy(address1, util.HASH, proxyHash, unbondingArgsStr, util.HASH, proxyHash, paymentArgsStr, uint64(10), timestamp, chainName)
	deploys = util.MakeInitDeploys()
	deploys = util.AddDeploy(deploys, deploy)
	res, err = grpc.Execute(client, rootStateHash, timestamp, deploys, protocolVersion)
	effects7, err := executeErrorHandler(res)
	if err != nil {
		println(err.Error())
	}
	postStateHash7, bonds7, errMessage := grpc.Commit(client, rootStateHash, effects7, protocolVersion)
	if errMessage != "" {
		panic(errMessage)
	}
	rootStateHash = postStateHash7
	printCommitResult(rootStateHash, bonds7)

	// Voting
	println(`Vote`)
	timestamp = time.Now().Unix()
	voteArgs := []*consensus.Deploy_Arg{
		&consensus.Deploy_Arg{
			Name: "method",
			Value: &consensus.Deploy_Arg_Value{
				Value: &consensus.Deploy_Arg_Value_StringValue{
					StringValue: "vote"}}},
		&consensus.Deploy_Arg{
			Name: "hash",
			Value: &consensus.Deploy_Arg_Value{
				Value: &consensus.Deploy_Arg_Value_Key{
					Key: &state.Key{Value: &state.Key_Hash_{
						Hash: &state.Key_Hash{
							Hash: address1,
						},
					}},
				}}},
		&consensus.Deploy_Arg{
			Name: "amount",
			Value: &consensus.Deploy_Arg_Value{
				Value: &consensus.Deploy_Arg_Value_BigInt{
					BigInt: &state.BigInt{Value: "9999999999999800", BitWidth: 512}}}},
	}
	voteArgsStr, err := util.DeployArgsToJsonString(voteArgs)
	if err != nil {
		panic(err)
	}
	deploy, _ = util.MakeDeploy(address1, util.HASH, proxyHash, voteArgsStr, util.HASH, proxyHash, paymentArgsStr, uint64(10), timestamp, chainName)
	deploys = util.MakeInitDeploys()
	deploys = util.AddDeploy(deploys, deploy)
	res, err = grpc.Execute(client, rootStateHash, timestamp, deploys, protocolVersion)
	effects7, err = executeErrorHandler(res)
	if err != nil {
		println(err.Error())
	}
	postStateHash7, bonds7, errMessage = grpc.Commit(client, rootStateHash, effects7, protocolVersion)
	if errMessage != "" {
		panic(errMessage)
	}
	rootStateHash = postStateHash7
	printCommitResult(rootStateHash, bonds7)

	// Voting res
	queryResult10, errMessage = grpc.Query(client, rootStateHash, "address", systemContract, []string{"pos"}, protocolVersion)
	var storedValue20 storedvalue.StoredValue
	storedValue20.FromBytes(queryResult10)
	storedValue20, err, _ = storedValue20.FromBytes(queryResult10)
	if err != nil {
		panic(err)
	}
	users := storedValue20.Contract.NamedKeys.GetVotingDappFromUser(address1)
	for user, value := range users {
		println(user + " : " + value)
	}

	// Unvoting
	println(`Unvote`)
	timestamp = time.Now().Unix()
	voteArgs = []*consensus.Deploy_Arg{
		&consensus.Deploy_Arg{
			Name: "method",
			Value: &consensus.Deploy_Arg_Value{
				Value: &consensus.Deploy_Arg_Value_StringValue{
					StringValue: "unvote"}}},
		&consensus.Deploy_Arg{
			Name: "hash",
			Value: &consensus.Deploy_Arg_Value{
				Value: &consensus.Deploy_Arg_Value_Key{
					Key: &state.Key{Value: &state.Key_Hash_{
						Hash: &state.Key_Hash{
							Hash: address1,
						},
					}},
				}}},
		&consensus.Deploy_Arg{
			Name: "amount",
			Value: &consensus.Deploy_Arg_Value{
				Value: &consensus.Deploy_Arg_Value_OptionalValue{
					OptionalValue: &consensus.Deploy_Arg_Value{
						Value: &consensus.Deploy_Arg_Value_BigInt{
							BigInt: &state.BigInt{Value: "90", BitWidth: 512}}}}}}}
	voteArgsStr, err = util.DeployArgsToJsonString(voteArgs)
	if err != nil {
		panic(err)
	}
	deploy, _ = util.MakeDeploy(address1, util.HASH, proxyHash, voteArgsStr, util.HASH, proxyHash, paymentArgsStr, uint64(10), timestamp, chainName)
	deploys = util.MakeInitDeploys()
	deploys = util.AddDeploy(deploys, deploy)
	res, err = grpc.Execute(client, rootStateHash, timestamp, deploys, protocolVersion)
	effects7, err = executeErrorHandler(res)
	if err != nil {
		println(err.Error())
	}
	postStateHash7, bonds7, errMessage = grpc.Commit(client, rootStateHash, effects7, protocolVersion)
	if errMessage != "" {
		panic(errMessage)
	}
	rootStateHash = postStateHash7
	printCommitResult(rootStateHash, bonds7)

	// Unvoting res
	queryResult10, errMessage = grpc.Query(client, rootStateHash, "address", systemContract, []string{"pos"}, protocolVersion)
	var storedValue10 storedvalue.StoredValue
	storedValue10.FromBytes(queryResult10)
	storedValue10, err, _ = storedValue10.FromBytes(queryResult10)
	if err != nil {
		panic(err)
	}
	users = storedValue10.Contract.NamedKeys.GetVotingDappFromUser(address1)
	for user, value := range users {
		println(user + " : " + value)
	}

	// Upgrade costs data..
	println(`Upgrade`)
	costs["regular"] = 2
	nextProtocolVersion := util.MakeProtocolVersion(2, 0, 0)
	postStateHash8, effects8, errMessage := grpc.Upgrade(client, rootStateHash, cntDefCode, costs, protocolVersion, nextProtocolVersion)
	postStateHash9, bonds8, errMessage := grpc.Commit(client, rootStateHash, effects8, nextProtocolVersion)
	if bytes.Equal(postStateHash8, postStateHash9) {
		rootStateHash = postStateHash8
		protocolVersion = nextProtocolVersion
	}
	printCommitResult(rootStateHash, bonds8)
}

func loadWasmCode() (cntDefCode []byte, cntCallCode []byte) {
	cntDefCode = util.LoadWasmFile("./example/contracts/counter_define.wasm")

	cntCallCode = util.LoadWasmFile("./example/contracts/counter_call.wasm")

	return cntDefCode, cntCallCode
}

func executeErrorHandler(r *ipc.ExecuteResponse) (effects []*transforms.TransformEntry, err error) {
	switch r.GetResult().(type) {
	case *ipc.ExecuteResponse_Success:
		for _, res := range r.GetSuccess().GetDeployResults() {
			switch res.GetExecutionResult().GetError().GetValue().(type) {
			case *ipc.DeployError_GasError:
				err = fmt.Errorf("DeployError_GasError")
			case *ipc.DeployError_ExecError:
				err = fmt.Errorf("DeployError_ExecError : %s", res.GetExecutionResult().GetError().GetExecError().GetMessage())
			default:
				effects = append(effects, res.GetExecutionResult().GetEffects().GetTransformMap()...)
			}

		}
	case *ipc.ExecuteResponse_MissingParent:
		err = fmt.Errorf("Missing parentstate : %s", util.EncodeToHexString(r.GetMissingParent().GetHash()))
	default:
		err = fmt.Errorf("Unknown result : %s", r.String())
	}

	return effects, err
}

func printCommitResult(stateHash []byte, bonds []*ipc.Bond) {
	println("State hash : " + hex.EncodeToString(stateHash))
	for _, bond := range bonds {
		println(hex.EncodeToString(bond.ValidatorPublicKey) + " : " + bond.GetStake().GetValue())
	}
	println()
}
