package node

import (
	"fmt"

	"github.com/barkisnet/barkis/codec"
	sdk "github.com/barkisnet/barkis/types"
	"github.com/barkisnet/barkis/x/auth/exported"
	authtypes "github.com/barkisnet/barkis/x/auth/types"
	"github.com/barkisnet/barkis/x/staking/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

func (node *Node) GetAccountWithHeight(cdc *codec.Codec, addr sdk.AccAddress) (exported.Account, error) {
	bs, err := cdc.MarshalJSON(authtypes.NewQueryAccountParams(addr))
	if err != nil {
		return nil, err
	}

	res, err := node.query(fmt.Sprintf("custom/%s/%s", "acc", "account"), bs)
	if err != nil {
		return nil, err
	}

	var account exported.Account
	if err := cdc.UnmarshalJSON(res, &account); err != nil {
		return nil, err
	}

	return account, nil
}

func (node *Node) query(path string, key cmn.HexBytes) (res []byte, err error) {
	opts := rpcclient.ABCIQueryOptions{
		Height: 0,
		Prove:  false,
	}

	result, err := node.Rpc.ABCIQueryWithOptions(path, key, opts)
	if err != nil {
		return res, err
	}

	resp := result.Response
	if !resp.IsOK() {
		return res, fmt.Errorf(resp.Log)
	}

	return resp.Value, nil
}

func (node *Node) QueryValiatorList(page, limit int, status string) ([]*Validator, error) {
	params := types.NewQueryValidatorsParams(page, limit, status)
	bz, err := node.Cdc.MarshalJSON(params)
	if err != nil {
		return nil, err
	}

	route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryValidators)
	res, err := node.query(route, bz)
	if err != nil {
		return nil, err
	}

	var validators []types.Validator
	err = node.Cdc.UnmarshalJSON(res, &validators)
	if err != nil {
		return nil, err
	}

	var vals []*Validator
	for _, val := range validators {
		vals = append(vals, &Validator{
			Counter:         0,
			OperatorAddr:    val.OperatorAddress,
			ConsensusAddr:   val.ConsPubKey.Address(),
			ConsensusPubKey: val.ConsPubKey,
		})
	}
	return vals, nil
}

func (node *Node) GetLatestHeight() (int64, error) {
	status, err := node.Rpc.ABCIInfo()
	if err != nil {
		return 0, nil
	}
	return status.Response.LastBlockHeight, nil
}
