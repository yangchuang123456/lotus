package lotus_common

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/lotus/chain/wallet"
	"github.com/filecoin-project/lotus/lib/sigs"
)

func GetBatchSignedMessages(initNonce, number uint64, privateKey wallet.Key,toAddr address.Address,value types.BigInt) ([]*types.SignedMessage, error) {
	msg := &types.Message{
		To:       toAddr,
		From:     privateKey.Address,
		Value:    value,
		Nonce:    initNonce,
		GasLimit: 1,
		GasPrice: types.NewInt(0),
	}

	sigMessages := make([]*types.SignedMessage, 0)
	i := uint64(0)
	for ; i < number; i++ {
		//generate dust tx
		sig, err := sigs.Sign(wallet.ActSigType(privateKey.Type), privateKey.PrivateKey, msg.Cid().Bytes())
		if err != nil {
			return nil, err
		}
		sigMessages = append(sigMessages, &types.SignedMessage{
			Message:   *msg,
			Signature: *sig,
		})

		msg.Nonce++
	}
	return sigMessages, nil
}
