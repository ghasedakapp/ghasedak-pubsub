package cmd

import (
	"ghasedak-pubsub/api/rpc"
	"ghasedak-pubsub/pkg"
)

func init() {
	pkg.InitConfig("")
	pkg.InitLog()
}

func Main() {
	rpc.InitGrpc(":5050")
	pkg.RunAwaitSignal()
}
