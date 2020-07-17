package tests

import (
	lotus_common "github.com/filecoin-project/lotus/yc-debug/lotus-common"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

type commonRpcConfig struct {
	LotusHost      string
	TokenPath      string
	LotusTokenFile string

	StorageMinerHost string
	StorageTokenFile string
}

var rpcConfig = commonRpcConfig{
	LotusHost:        "ws://127.0.0.1:1234/rpc/v0",
	TokenPath:        ".",
	StorageMinerHost: "ws://127.0.0.1:2345/rpc/v0",

	LotusTokenFile:   "lotus-token",
	StorageTokenFile: "storage-token",
}

func getTestUsualApiOp(t *testing.T,config commonRpcConfig) *lotus_common.UsualApiOperation {
	fNodeApi, cCloser, err := lotus_common.GetLotusFullNodeRpc(filepath.Join(config.TokenPath, config.LotusTokenFile), config.LotusHost)
	assert.NoError(t, err)

	storageMinerApi, storageMinerCloser, err := lotus_common.GetStorageMinerRpc(filepath.Join(config.TokenPath, config.StorageTokenFile), config.StorageMinerHost)
	assert.NoError(t, err)

	return lotus_common.NewUsualApiOperation(lotus_common.ApiConfig{
		FullNode:           fNodeApi,
		FullNodeCloser:     cCloser,
		StorageMiner:       storageMinerApi,
		StorageMinerCloser: storageMinerCloser,
	})
}
