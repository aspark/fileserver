package command

import (
	"github.com/aspark/fileserver/config"
	"github.com/aspark/fileserver/protocal"
	"github.com/aspark/fileserver/utils"
)

func handleInfoCommand(ctx *utils.CommandContext, args []protocal.IBlock) (protocal.IBlockResult, error) {
	var items = make([]protocal.IBlock, 0)

	var cfg = config.GetConfig()
	var paths string
	for key := range cfg.Path {
		paths += (key + ";")
	}
	items = append(items, protocal.BulkStringBlock{Content: "Path:" + paths})

	return protocal.ArrayBlock{Content: items}, nil
}
