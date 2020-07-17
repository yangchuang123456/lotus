package tests

import (
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

var minerTestInput = struct {
	commonRpcConfig
}{
	commonRpcConfig: rpcConfig,
}

func Test_getMinerInfo(t *testing.T) {
	op := getTestUsualApiOp(t, minerTestInput.commonRpcConfig)
	minerAddr, err := op.StorageMiner.ActorAddress(op.Ctx)
	assert.NoError(t, err)
	log.Println("the minerAddr is:", minerAddr)

	minerInfo, err := op.FullNode.StateMinerInfo(op.Ctx, minerAddr, types.TipSetKey{})
	assert.NoError(t, err)
	log.Println("the miner info is:", minerInfo)

	ownerBalance, err := op.FullNode.WalletBalance(op.Ctx, minerInfo.Owner)
	log.Println("the owner balance is:", ownerBalance)
}

func Test_MinerRemoteSignWithdraw(t *testing.T) {
/*	op := getTestUsualApiOp(t, minerTestInput.commonRpcConfig)
	minerAddr, err := op.StorageMiner.ActorAddress(op.Ctx)



	message := &types.Message{
		To:       minerAddr,
		From:     mi.Owner,
		Value:    types.NewInt(0),
		GasPrice: types.NewInt(1),
		GasLimit: gasLimit,
		Method:   builtin.MethodsMiner.WithdrawBalance,
		Params:   params,
	}*/
}
