package tests

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

var sectorTestInput = struct {
	commonRpcConfig
}{
	commonRpcConfig: rpcConfig,
}

func Test_GetSectorExpiration(t *testing.T) {
	op := getTestUsualApiOp(t,sectorTestInput.commonRpcConfig)

	minerActorAddr, err := op.StorageMiner.ActorAddress(op.Ctx)
	assert.NoError(t, err)
	log.Println("the miner actor addr is:", minerActorAddr)

	curTipSets, err := op.FullNode.ChainHead(op.Ctx)
	assert.NoError(t, err)

	_, err = op.FullNode.ChainGetTipSet(op.Ctx, curTipSets.Key())
	assert.NoError(t, err)

	log.Println("the tipSet cid and height is:", curTipSets.Cids(), curTipSets.Height())
	return

	sectors, err := op.FullNode.StateMinerSectors(op.Ctx, minerActorAddr, nil, true, curTipSets.Key())
	assert.NoError(t, err)

	log.Println("the sectors is:", sectors)
}


