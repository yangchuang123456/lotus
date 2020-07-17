package tests

import (
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/specs-actors/actors/abi"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

var messageTestInput = struct {
	commonRpcConfig
}{
	commonRpcConfig: rpcConfig,
}

func Test_ChainMessage(t *testing.T) {
	op := getTestUsualApiOp(t, minerTestInput.commonRpcConfig)
/*	chainNotify, err := op.FullNode.ChainNotify(op.Ctx)
	assert.NoError(t, err)

	for {
		select {
		case change := <-chainNotify:
			for _, v := range change {
				message, err := op.FullNode.ChainGetBlockMessages(op.Ctx, v.Val.Blocks()[0].Cid())
				assert.NoError(t,err)
				log.Println("the block message is:",message.BlsMessages)
			}
		default:
			time.Sleep(time.Second)
		}
	}*/

	height := 6222
	for i:=0;i<=height;i++{
		tipSet,err:=op.FullNode.ChainGetTipSetByHeight(op.Ctx,abi.ChainEpoch(i),types.TipSetKey{})
		assert.NoError(t,err)
/*		log.Println("the tipSet height is:",tipSet.Height())
		dp2plog.L.Info("the return tipSet key is:",zap.Any("height",tipSet.Height()),zap.Any("tipSetKey",tipSet.Key()))
		log.Println("the tipSet cids is:",zap.Any("cids",tipSet.Cids()))*/
/*		for _,b:=range tipSet.Blocks(){
			log.Println("the tipSet block  cid is:",b.Cid())
		}*/

		m,err:=op.FullNode.ChainGetBlockMessages(op.Ctx,tipSet.Blocks()[0].Cid())
		assert.NoError(t,err)
		if len(m.BlsMessages) != 0{
			log.Println("the block height is:",tipSet.Height())
			log.Println("the message is:",*m.BlsMessages[0])
		}

		if len(m.SecpkMessages) !=0{
			log.Println("the block height is:",tipSet.Height())
			log.Println("the message is:",*m.SecpkMessages[0])
		}

	}
}
