package daemon

import (
	"fmt"
	"time"

	"github.com/barkisnet/barkis-monitor/node"
	sdk "github.com/barkisnet/barkis/types"
)

const (
	sleepSecond = 5
)

func MonitorDaemon(startHeight int64, rpcNode *node.Node, validatorList []*node.Validator) {
	height := startHeight + 1
	for {
		block, err := rpcNode.Rpc.Block(&height)
		if err != nil {
			time.Sleep(sleepSecond * time.Second)
			continue
		}
		height++

		fmt.Println(fmt.Sprintf("Get commit from height %d", height))
		commitsMap := make(map[string]bool)
		for _, commit := range block.Block.LastCommit.Precommits {
			if commit != nil {
				commitsMap[commit.ValidatorAddress.String()] = true
			}
		}

		var missedVals []*node.Validator
		for _, val := range validatorList {
			_, ok := commitsMap[val.ConsensusAddr.String()]
			if !ok {
				missedVals = append(missedVals, val)
			}
		}

		if len(missedVals) != 0 {
			fmt.Println(fmt.Println("Miss commit from validators:"))
			for _, val := range missedVals {
				consensusPubKey, _ := sdk.Bech32ifyConsPub(val.ConsensusPubKey)
				fmt.Println(fmt.Sprintf("operatorAddr: %s, consensusPubKey: %s", val.OperatorAddr.String(), consensusPubKey))
			}
		}
	}
}
