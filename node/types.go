package node

import (
	"github.com/barkisnet/barkis/codec"
	sdk "github.com/barkisnet/barkis/types"
	"github.com/barkisnet/barkis/x/auth"
	authcutils "github.com/barkisnet/barkis/x/auth/client/utils"
	"github.com/barkisnet/barkis/x/auth/exported"
	"github.com/barkisnet/monitor/key"
	"github.com/tendermint/tendermint/crypto"
	cmn "github.com/tendermint/tendermint/libs/common"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

type Node struct {
	ChainID    string
	Rpc        rpcclient.Client
	Cdc        *codec.Codec
	keyManager key.KeyManager
	TxEncoder  sdk.TxEncoder
}

func NewNode(chainID, url string, keyManager key.KeyManager, cdc *codec.Codec) *Node {
	rpc := rpcclient.NewHTTP(url, "/websocket")
	return &Node{
		ChainID:    chainID,
		Rpc:        rpc,
		Cdc:        cdc,
		keyManager: keyManager,
		TxEncoder:  authcutils.GetTxEncoder(cdc),
	}
}

func (node *Node) BuildAndSign(account exported.Account, memo string, msgs []sdk.Msg, fees auth.StdFee) ([]byte, error) {
	stdSignMsg := auth.StdSignMsg{
		ChainID:       node.ChainID,
		AccountNumber: account.GetAccountNumber(),
		Sequence:      account.GetSequence(),
		Memo:          memo,
		Msgs:          msgs,
		Fee:           fees,
	}

	sigBytes, err := node.keyManager.Sign(stdSignMsg)
	if err != nil {
		return nil, err
	}
	stdSignature := auth.StdSignature{
		PubKey:    node.keyManager.GetPrivKey().PubKey(),
		Signature: sigBytes,
	}
	return node.TxEncoder(auth.NewStdTx(msgs, fees, []auth.StdSignature{stdSignature}, memo))
}

type Validator struct {
	OperatorAddr    sdk.ValAddress `json:"operator_addr"`
	ConsensusAddr   cmn.HexBytes   `json:"consensus_addr"`
	ConsensusPubKey crypto.PubKey  `json:"consensus_pub_key"`
}
