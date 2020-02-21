package main

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/hdac-io/casperlabs-ee-grpc-go-util/grpc"
	"github.com/hdac-io/casperlabs-ee-grpc-go-util/protobuf/io/casperlabs/casper/consensus"
	"github.com/hdac-io/casperlabs-ee-grpc-go-util/protobuf/io/casperlabs/ipc"
	"github.com/hdac-io/casperlabs-ee-grpc-go-util/protobuf/io/casperlabs/ipc/transforms"
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

	genesisConfig, err := util.GenesisConfigMock(
		chainName, genesisAddress, "500000000", "1000000", protocolVersion, costs,
		"./example/contracts/mint_install.wasm", "./example/contracts/pos_install.wasm")
	if err != nil {
		fmt.Printf("Bad GenesisConfigMock err : %v", err)
		return
	}

	//var parentStateHash , effects
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
	println(util.EncodeToHexString(rootStateHash))
	println(bonds[0].String())

	queryResult10, errMessage := grpc.Query(client, rootStateHash, "address", systemContract, []string{}, protocolVersion)
	proxyHash := queryResult10.GetAccount().GetNamedKeys()[0].GetKey().GetHash().GetHash()
	println(util.EncodeToHexString(genesisAddress), ": ", proxyHash)
	println(errMessage)

	// Run "Counter Define contract"
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
				Value: &consensus.Deploy_Arg_Value_IntValue{
					IntValue: 100000000}}}}
	deploy, _ := util.MakeDeploy(genesisAddress, util.WASM, cntDefCode, []*consensus.Deploy_Arg{}, util.HASH, proxyHash, paymentArgs, uint64(10), timestamp, chainName)
	deploys := util.MakeInitDeploys()
	deploys = util.AddDeploy(deploys, deploy)

	res, err := grpc.Execute(client, rootStateHash, timestamp, deploys, protocolVersion)
	effects2, err := executeErrorHandler(res)

	postStateHash2, bonds2, errMessage := grpc.Commit(client, rootStateHash, effects2, protocolVersion)
	rootStateHash = postStateHash2
	println(util.EncodeToHexString(postStateHash2))
	println(bonds2[0].String())

	queryResult, errMessage := grpc.QueryBalance(client, rootStateHash, genesisAddress, protocolVersion)
	println(util.EncodeToHexString(genesisAddress), ": ", queryResult)
	println(errMessage)

	// Run "Counter Call contract"
	timestamp = time.Now().Unix()
	deploy, _ = util.MakeDeploy(genesisAddress, util.WASM, cntCallCode, []*consensus.Deploy_Arg{}, util.HASH, proxyHash, paymentArgs, uint64(10), timestamp, chainName)
	deploys = util.MakeInitDeploys()
	deploys = util.AddDeploy(deploys, deploy)
	res, err = grpc.Execute(client, rootStateHash, timestamp, deploys, protocolVersion)
	effects3, err := executeErrorHandler(res)

	postStateHash3, bonds3, errMessage := grpc.Commit(client, rootStateHash, effects3, protocolVersion)
	rootStateHash = postStateHash3
	println(util.EncodeToHexString(postStateHash3))
	println(bonds3[0].String())

	// Query counter contract.
	path := []string{"counter", "count"}
	queryResult1, errMessage := grpc.Query(client, rootStateHash, "address", genesisAddress, path, protocolVersion)
	println(queryResult1.GetIntValue())
	println(errMessage)

	queryResult, errMessage = grpc.QueryBalance(client, rootStateHash, genesisAddress, protocolVersion)
	println(util.EncodeToHexString(genesisAddress), ": ", queryResult)
	println(errMessage)

	timestamp = time.Now().Unix()
	deploy, _ = util.MakeDeploy(genesisAddress, util.WASM, cntCallCode, []*consensus.Deploy_Arg{}, util.HASH, proxyHash, paymentArgs, uint64(10), timestamp, chainName)
	deploys = util.MakeInitDeploys()
	deploys = util.AddDeploy(deploys, deploy)
	res, err = grpc.Execute(client, rootStateHash, timestamp, deploys, protocolVersion)
	effects4, err := executeErrorHandler(res)

	postStateHash4, bonds4, errMessage := grpc.Commit(client, rootStateHash, effects4, protocolVersion)
	rootStateHash = postStateHash4
	println(util.EncodeToHexString(postStateHash4))
	println(bonds4[0].String())

	queryResult2, errMessage := grpc.Query(client, rootStateHash, "address", genesisAddress, path, protocolVersion)
	println(queryResult2.GetIntValue())
	println(errMessage)

	queryResult3, errMessage := grpc.QueryBalance(client, rootStateHash, genesisAddress, protocolVersion)
	println(util.EncodeToHexString(genesisAddress), ": ", queryResult3)
	println(errMessage)

	// Run "Send transaction"
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
				Value: &consensus.Deploy_Arg_Value_LongValue{
					LongValue: int64(10)}}},
	}

	deploy, _ = util.MakeDeploy(genesisAddress, util.HASH, proxyHash, sessionArgs, util.HASH, proxyHash, paymentArgs, uint64(10), timestamp, chainName)
	deploys = util.MakeInitDeploys()
	deploys = util.AddDeploy(deploys, deploy)
	res, err = grpc.Execute(client, rootStateHash, timestamp, deploys, protocolVersion)
	effects5, err := executeErrorHandler(res)

	postStateHash5, bonds5, errMessage := grpc.Commit(client, rootStateHash, effects5, protocolVersion)
	rootStateHash = postStateHash5
	println(util.EncodeToHexString(postStateHash5))
	println(bonds5[0].String())

	queryResult4, errMessage := grpc.QueryBalance(client, rootStateHash, genesisAddress, protocolVersion)
	println(util.EncodeToHexString(genesisAddress), ": ", queryResult4)
	println(errMessage)

	queryResult5, errMessage := grpc.QueryBalance(client, rootStateHash, address1, protocolVersion)
	println(util.EncodeToHexString(address1), ": ", queryResult5)
	println(errMessage)

	// bonding
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
				Value: &consensus.Deploy_Arg_Value_LongValue{
					LongValue: int64(10)}}},
	}
	deploy, _ = util.MakeDeploy(genesisAddress, util.HASH, proxyHash, bondingArgs, util.HASH, proxyHash, paymentArgs, uint64(10), timestamp, chainName)
	deploys = util.MakeInitDeploys()
	deploys = util.AddDeploy(deploys, deploy)
	res, err = grpc.Execute(client, rootStateHash, timestamp, deploys, protocolVersion)
	effects6, err := executeErrorHandler(res)
	postStateHash6, bonds6, errMessage := grpc.Commit(client, rootStateHash, effects6, protocolVersion)
	rootStateHash = postStateHash6
	println(util.EncodeToHexString(rootStateHash))
	println(bonds6[0].String())

	// unbonding
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
						Value: &consensus.Deploy_Arg_Value_LongValue{
							LongValue: int64(100)}}}}},
	}
	deploy, _ = util.MakeDeploy(genesisAddress, util.HASH, proxyHash, unbondingArgs, util.HASH, proxyHash, paymentArgs, uint64(10), timestamp, chainName)
	deploys = util.MakeInitDeploys()
	deploys = util.AddDeploy(deploys, deploy)
	res, err = grpc.Execute(client, rootStateHash, timestamp, deploys, protocolVersion)
	effects7, err := executeErrorHandler(res)
	postStateHash7, bonds7, errMessage := grpc.Commit(client, rootStateHash, effects7, protocolVersion)
	rootStateHash = postStateHash7
	println(util.EncodeToHexString(rootStateHash))
	println(bonds7[0].String())

	// Upgrade costs data..
	costs["regular"] = 2
	nextProtocolVersion := util.MakeProtocolVersion(2, 0, 0)
	postStateHash8, effects8, errMessage := grpc.Upgrade(client, rootStateHash, cntDefCode, costs, protocolVersion, nextProtocolVersion)
	postStateHash9, bonds8, errMessage := grpc.Commit(client, rootStateHash, effects8, nextProtocolVersion)
	if bytes.Equal(postStateHash8, postStateHash9) {
		rootStateHash = postStateHash8
		protocolVersion = nextProtocolVersion
	}
	println(util.EncodeToHexString(rootStateHash))
	println(bonds8[0].String())
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
