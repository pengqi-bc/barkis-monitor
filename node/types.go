package node

import (
	"github.com/barkisnet/barkis/codec"
	sdk "github.com/barkisnet/barkis/types"
	authcutils "github.com/barkisnet/barkis/x/auth/client/utils"
	"github.com/tendermint/tendermint/crypto"
	cmn "github.com/tendermint/tendermint/libs/common"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

type Node struct {
	ChainID    string
	Rpc        rpcclient.Client
	Cdc        *codec.Codec
	TxEncoder  sdk.TxEncoder
}

func NewNode(chainID, url string, cdc *codec.Codec) *Node {
	rpc := rpcclient.NewHTTP(url, "/websocket")
	return &Node{
		ChainID:    chainID,
		Rpc:        rpc,
		Cdc:        cdc,
		TxEncoder:  authcutils.GetTxEncoder(cdc),
	}
}

type Validator struct {
	Counter         int64          `json:"counter"`
	OperatorAddr    sdk.ValAddress `json:"operator_addr"`
	ConsensusAddr   cmn.HexBytes   `json:"consensus_addr"`
	ConsensusPubKey crypto.PubKey  `json:"consensus_pub_key"`
}
