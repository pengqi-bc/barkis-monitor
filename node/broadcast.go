package node

import (
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

func (node *Node) Broadcast(signedTxBytes []byte) (*ctypes.ResultBroadcastTxCommit, error) {
	return node.Rpc.BroadcastTxCommit(signedTxBytes)
}
