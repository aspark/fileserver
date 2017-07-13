package command

import (
	"errors"

	"os"

	"github.com/aspark/fileserver/config"
	"github.com/aspark/fileserver/protocal"
	"github.com/aspark/fileserver/utils"
)

func handleDelCommand(ctx *utils.CommandContext, args []protocal.IBlock) (protocal.IBlockResult, error) {
	fileName, ok := args[0].GetContent().(string)
	if ok {
		result, err := Del(ctx, fileName)
		if err != nil {
			return nil, err
		}

		if result {
			return protocal.SimpleStringBlock{Content: "OK"}, nil
		}

		return protocal.SimpleStringBlock{Content: "OK"}, nil
	}

	return nil, errors.New("invalid parameter")
}

//
func Del(ctx *utils.CommandContext, fileName string) (bool, error) {
	fileFullName, err := config.GetConfigPath(ctx, fileName)
	if err != nil {
		return false, err
	}

	if r, e := ExistFile(ctx, fileFullName); !r || e != nil {
		return true, nil
	}

	err = os.Remove(fileFullName)

	return true, err
}
