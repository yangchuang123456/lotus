package dp2plog

import (
	"github.com/filecoin-project/lotus/tools/util"
	"go.uber.org/zap"
	"os"
)

var L *zap.Logger

func init() {
	var err error
	if os.Getenv("XLotusLogOn") != "" {
		L, err = util.LogToWorkDir(util.CurExecPath(), "p2p", zap.DebugLevel).Build()
		if err != nil {
			panic(err)
		}
		return
	}

	L, err = util.LogNothing().Build()
	if err != nil {
		panic(err)
	}
}
