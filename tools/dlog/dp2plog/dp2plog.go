package dp2plog

import (
	"github.com/filecoin-project/lotus/tools/util"
	"go.uber.org/zap"
)

var L *zap.Logger

func init() {
	L = util.GetXDebugLog("p2p")
}
