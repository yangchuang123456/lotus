package lotus_common

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/lotus/chain/wallet"
)

func GetPrivateFromString(sk string) (*wallet.Key, error) {
	key, err := hex.DecodeString(sk)
	if err != nil {
		return nil, fmt.Errorf("GetPrivateFromString decode sk error:%s", err.Error())
	}

	keyInfo := types.KeyInfo{}
	err = json.Unmarshal(key, &keyInfo)
	if err != nil {
		return nil, fmt.Errorf("GetPrivateFromString json unmarshal keyinfo error:%s", err.Error())
	}

	return wallet.NewKey(keyInfo)
}
