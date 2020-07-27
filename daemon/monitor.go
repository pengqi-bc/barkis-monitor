package daemon

import (
	"fmt"
	"time"

	"github.com/barkisnet/barkis-monitor/node"
	sdk "github.com/barkisnet/barkis/types"
)

const (
	sleepSecond           = 5
	threshholdMissedBlock = 100
	botId                 = "1277648385:AAEFE5fFExzmGvcIN5Z1Gyulzma-8nLE8os"
	chatId                = "@bks_monitor"
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

		for _, val := range validatorList {
			_, ok := commitsMap[val.ConsensusAddr.String()]
			if !ok {
				val.Counter++
				consensusPubKey, _ := sdk.Bech32ifyConsPub(val.ConsensusPubKey)
				fmt.Println(fmt.Sprintf("operatorAddr: %s, consensusPubKey: %s", val.OperatorAddr.String(), consensusPubKey))
				if val.Counter >= threshholdMissedBlock {
					alertToTg(fmt.Sprintf("validator missed more than %d blocks, operator addr: %s, consensus pubkey: %s", threshholdMissedBlock, val.OperatorAddr.String(), consensusPubKey))
				}
			} else {
				val.Counter = 0
			}
		}
	}
}

func alertToTg(msg string) {
	sendTelegramMessage(botId, chatId, msg)
}
