package lotus_common

import (
	"context"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/ipfs/go-cid"
)

var (
	DefaultGasPrice = types.NewInt(0)
	DefaultGasLimit = int64(10000)
)

type ApiConfig struct {
	FullNode       api.FullNode
	FullNodeCloser jsonrpc.ClientCloser
	StorageMiner   api.StorageMiner
	StorageMinerCloser jsonrpc.ClientCloser
}

type UsualApiOperation struct {
	ApiConfig
	Ctx context.Context
}

func NewUsualApiOperation(conf ApiConfig) *UsualApiOperation {
	return &UsualApiOperation{
		ApiConfig: conf,
		Ctx:       context.Background(),
	}
}

func (o *UsualApiOperation) GetDefaultAddr() (address.Address, error) {
	return o.FullNode.WalletDefaultAddress(o.Ctx)
}

func (o *UsualApiOperation) GetAddrBalance(addr string) (types.BigInt, error) {
	a, err := address.NewFromString(addr)
	if err != nil {
		return types.BigInt{}, err
	}
	return o.FullNode.WalletBalance(o.Ctx, a)
}

func (o *UsualApiOperation) GetAddrNonce(addr string) (uint64, error) {
	a, err := address.NewFromString(addr)
	if err != nil {
		return 0, err
	}
	return o.FullNode.MpoolGetNonce(o.Ctx, a)
}

func (o *UsualApiOperation) PushBatchSignedMessage(fromAddr, toAddr, fromSk string, number uint64, value types.BigInt) ([]cid.Cid, error) {
	nonce, err := o.GetAddrNonce(fromAddr)
	if err != nil {
		return nil, err
	}
	privateKey, err := GetPrivateFromString(fromSk)
	if err != nil {
		return nil, err
	}

	addr, err := address.NewFromString(toAddr)
	if err != nil {
		return nil, err
	}

	sigMes, err := GetBatchSignedMessages(nonce, number, *privateKey, addr, value)
	if err != nil {
		return nil, err
	}
	cidS := make([]cid.Cid, 0)
	for _, m := range sigMes {
		cid, err := o.FullNode.MpoolPush(o.Ctx, m)
		if err != nil {
			return nil, err
		}
		cidS = append(cidS, cid)
	}

	return cidS, nil
}

func (o *UsualApiOperation) GetPendingMessage() ([]*types.SignedMessage, error) {
	//todo: find reason
	return o.FullNode.MpoolPending(o.Ctx, types.TipSetKey{})
/*	tipSet, err := o.FullNode.ChainHead(o.Ctx)
	if err != nil {
		return nil, err
	}

	waitFunc := func(ctx context.Context, baseTime uint64) (func(bool), error) {
		// Wait around for half the block time in case other parents come in
		deadline := baseTime + build.PropagationDelaySecs
		time.Sleep(time.Until(time.Unix(int64(deadline), 0)))

		return func(bool) {}, nil
	}

	waitFunc(o.Ctx,tipSet.MinTimestamp())

	return o.FullNode.MpoolPending(o.Ctx, tipSet.Key())*/
}

func (o *UsualApiOperation) TransferCoinToAddr(fromAddr, toAddr string, value types.BigInt) (*types.SignedMessage, error) {
	from, err := address.NewFromString(fromAddr)
	if err != nil {
		return nil, err
	}

	to, err := address.NewFromString(toAddr)
	if err != nil {
		return nil, err
	}

	msg := &types.Message{
		From:     from,
		To:       to,
		Value:    value,
		GasLimit: DefaultGasLimit,
		GasPrice: DefaultGasPrice,
	}
	return o.FullNode.MpoolPushMessage(o.Ctx, msg)
}
