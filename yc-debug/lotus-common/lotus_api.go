package lotus_common

import (
	"bytes"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/api/client"
	"github.com/filecoin-project/lotus/node/repo"
	log "github.com/filecoin-project/lotus/tools/dlog/dp2plog"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"os"
)

func GetHttpRequestHeaders(tokenPath string) (http.Header, error) {
	f, err := os.Open(tokenPath)
	if os.IsNotExist(err) {
		return nil, repo.ErrNoAPIEndpoint
	} else if err != nil {
		return nil, err
	}
	defer f.Close() //nolint: errcheck // Read only op

	tb, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	token := string(bytes.TrimSpace(tb))
	log.L.Debug("connect node", zap.String("token", token))
	headers := http.Header{}
	headers.Add("Authorization", "Bearer "+token)

	return headers, nil
}

func GetLotusFullNodeRpc(tokenPath string, addr string) (api.FullNode, jsonrpc.ClientCloser, error) {
	if addr == "" {
		addr = "ws://127.0.0.1:1234/rpc/v0"
	}

	headers, err := GetHttpRequestHeaders(tokenPath)
	if err != nil {
		return nil, nil, err
	}

	return client.NewFullNodeRPC(addr, headers)
}

func GetStorageMinerRpc(tokenPath string, addr string) (api.StorageMiner, jsonrpc.ClientCloser, error) {
	if addr == "" {
		addr = "ws://127.0.0.1:2345/rpc/v0"
	}

	headers, err := GetHttpRequestHeaders(tokenPath)
	if err != nil {
		return nil, nil, err
	}

	return client.NewStorageMinerRPC(addr, headers)
}
