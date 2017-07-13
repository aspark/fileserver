package command

import (
	"github.com/aspark/fileserver/protocal"
	"github.com/aspark/fileserver/utils"
)

func handlePingCommand(ctx *utils.CommandContext, args []protocal.IBlock) (protocal.IBlockResult, error) {
	return protocal.SimpleStringBlock{Content: "PONG"}, nil
}
