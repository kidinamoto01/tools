package monitor

import (
	amino "github.com/tendermint/go-amino"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

var cdc = amino.NewCodec()

func init() {
	ctypes.RegisterAmino(cdc)
	cdc.RegisterConcrete(&NodeStatus{}, "tendermint/NodeStatus", nil)

}
