package command

import (
	"github.com/aspark/fileserver/config"
	"github.com/aspark/fileserver/protocal"
	"github.com/aspark/fileserver/utils"
)

func handleReloadCommand(ctx *utils.CommandContext, args []protocal.IBlock) (protocal.IBlockResult, error) {
	config.ReloadConfig()

	return protocal.SimpleStringBlock{Content: "OK"}, nil
}
