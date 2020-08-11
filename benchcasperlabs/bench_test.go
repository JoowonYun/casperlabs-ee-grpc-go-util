package benchcasperlabs

import (
	"testing"

	"github.com/hdac-io/casperlabs-ee-grpc-go-util/bench"
	"github.com/hdac-io/casperlabs-ee-grpc-go-util/grpc"
	"github.com/hdac-io/casperlabs-ee-grpc-go-util/integration"
	"github.com/stretchr/testify/assert"
)

func TestCasperTransferToAccount100000(t *testing.T) {
	client, rootStateHash, _, protocolVersion := integration.InitalRunGenensis("../integration/contracts/mint_install.wasm", "../integration/contracts/pos_install.wasm", "../integration/contracts/standard_payment_install.wasm", integration.DEFAULT_GENESIS_ACCOUNT)
	amount := "1"

	for i := 0; i < 100000; i++ {
		rootStateHash, _ = bench.RunTransferToAccountWithWASM(client, rootStateHash, integration.GENESIS_ADDRESS, integration.ADDRESS1, amount, integration.BASIC_FEE, protocolVersion)
	}

	queryResult, errMessage := grpc.QueryBalance(client, rootStateHash, integration.ADDRESS1, protocolVersion)
	assert.Equal(t, "100000", queryResult)
	assert.Equal(t, "", errMessage)
}
