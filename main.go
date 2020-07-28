package main

import (
	"fmt"

	"github.com/barkisnet/barkis-monitor/constant"
	"github.com/barkisnet/barkis-monitor/daemon"
	"github.com/barkisnet/barkis-monitor/node"
	"github.com/barkisnet/barkis/app"
	sdk "github.com/barkisnet/barkis/types"
)

func main() {
	cdc := app.MakeCodec()

	rpcNode := node.NewNode(constant.ChainID, constant.NodeRpc, cdc)

	latestHeight, err := rpcNode.GetLatestHeight()
	if err != nil {
		fmt.Println(err.Error())
	}

	vals, err := rpcNode.QueryValiatorList(1, 100, sdk.BondStatusBonded)
	if err != nil {
		fmt.Println(err.Error())
	}

	go daemon.MonitorDaemon(latestHeight, rpcNode, vals)

	select {}
}
