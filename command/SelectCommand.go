package command

import (
	"errors"

	"github.com/aspark/fileserver/config"
	"github.com/aspark/fileserver/protocal"
	"github.com/aspark/fileserver/utils"
)

func handleSelectCommand(ctx *utils.CommandContext, args []protocal.IBlock) (protocal.IBlockResult, error) {

	name := args[0].GetContent().(string)
	if utils.IsNullOrWhiteSpace(name) {
		return nil, errors.New("name cannot be null or empty")
	}

	var cfg = config.GetConfig()
	if _, ok := cfg.Path[name]; !ok {
		return nil, errors.New("not exist the 'path name' in config")
	}

	ctx.SelectedPath = name

	return protocal.SimpleStringBlock{Content: "OK"}, nil
}
