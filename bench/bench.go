package bench

import (
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/hdac-io/casperlabs-ee-grpc-go-util/grpc"
	"github.com/hdac-io/casperlabs-ee-grpc-go-util/integration"
	"github.com/hdac-io/casperlabs-ee-grpc-go-util/protobuf/io/casperlabs/casper/consensus"
	"github.com/hdac-io/casperlabs-ee-grpc-go-util/protobuf/io/casperlabs/casper/consensus/state"
	"github.com/hdac-io/casperlabs-ee-grpc-go-util/protobuf/io/casperlabs/ipc"
	"github.com/hdac-io/casperlabs-ee-grpc-go-util/util"
)

func RunTransferToAccountWithWASM(client ipc.ExecutionEngineServiceClient, stateHash []byte, runAddress []byte,
	toAddress []byte, amount string,
	fee string, protocolVersion *state.ProtocolVersion) (resultStateHash []byte, bonds []*ipc.Bond) {

	transferToAccountWASM := util.LoadWasmFile("../integration/contracts/transfer_to_account_u512.wasm")

	sessionArgs := []*consensus.Deploy_Arg{
		&consensus.Deploy_Arg{
			Name: "address",
			Value: &state.CLValueInstance{
				ClType: &state.CLType{Variants: &state.CLType_ListType{ListType: &state.CLType_List{Inner: &state.CLType{Variants: &state.CLType_SimpleType{SimpleType: state.CLType_U8}}}}},
				Value: &state.CLValueInstance_Value{
					Value: &state.CLValueInstance_Value_BytesValue{
						BytesValue: toAddress}}}},
		&consensus.Deploy_Arg{
			Name: "amount",
			Value: &state.CLValueInstance{
				ClType: &state.CLType{Variants: &state.CLType_SimpleType{SimpleType: state.CLType_U512}},
				Value: &state.CLValueInstance_Value{
					Value: &state.CLValueInstance_Value_U512{
						U512: &state.CLValueInstance_U512{
							Value: amount}}}}}}

	standardpaymentWASM := util.LoadWasmFile("../integration/contracts/standard_payment.wasm")
	paymentArgs := []*consensus.Deploy_Arg{
		&consensus.Deploy_Arg{
			Name: "amount",
			Value: &state.CLValueInstance{
				ClType: &state.CLType{Variants: &state.CLType_SimpleType{SimpleType: state.CLType_U512}},
				Value: &state.CLValueInstance_Value{
					Value: &state.CLValueInstance_Value_U512{
						U512: &state.CLValueInstance_U512{
							Value: fee}}}}}}

	timestamp := time.Now().Unix()

	deployBody := &consensus.Deploy_Body{
		Session: util.MakeDeployCode(util.WASM, transferToAccountWASM, sessionArgs),
		Payment: util.MakeDeployCode(util.WASM, standardpaymentWASM, paymentArgs)}

	marshalDeployBody, _ := proto.Marshal(deployBody)
	bodyHash := util.Blake2b256(marshalDeployBody)

	deployHeader := &consensus.Deploy_Header{
		AccountPublicKey: runAddress,
		Timestamp:        uint64(timestamp),
		GasPrice:         uint64(10),
		BodyHash:         bodyHash,
		ChainName:        integration.CHAIN_NAME}

	marshalDeployHeader, _ := proto.Marshal(deployHeader)
	headerHash := util.Blake2b256(marshalDeployHeader)

	sessionAbi, err := util.AbiDeployArgsTobytes(sessionArgs)
	if err != nil {
		panic(err)
	}
	paymentAbi, err := util.AbiDeployArgsTobytes(paymentArgs)
	if err != nil {
		panic(err)
	}

	deploy := &ipc.DeployItem{
		Address:           runAddress,
		Session:           util.MakeDeployPayload(util.WASM, transferToAccountWASM, sessionAbi),
		Payment:           util.MakeDeployPayload(util.WASM, standardpaymentWASM, paymentAbi),
		GasPrice:          uint64(10),
		AuthorizationKeys: [][]byte{runAddress},
		DeployHash:        headerHash}

	deploys := []*ipc.DeployItem{deploy}

	res, err := grpc.Execute(client, stateHash, timestamp, deploys, protocolVersion)
	effect, err := integration.ExecuteErrorHandler(res)
	if err != nil {
		panic(err)
	}

	stateHash, bonds, errMessage := grpc.Commit(client, stateHash, effect, protocolVersion)
	if errMessage != "" {
		panic(errMessage)
	}

	return stateHash, bonds
}
