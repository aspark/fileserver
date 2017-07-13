package command

import (
	"errors"

	"io/ioutil"

	"github.com/aspark/fileserver/config"
	"github.com/aspark/fileserver/protocal"
	"github.com/aspark/fileserver/utils"
)

func handleGetCommand(ctx *utils.CommandContext, args []protocal.IBlock) (protocal.IBlockResult, error) {
	fileName, ok := args[0].GetContent().(string)
	if ok {
		result, err := Get(ctx, fileName)
		if err != nil {
			return nil, err
		}

		return protocal.BytesBlock{Content: result}, nil
	}

	return nil, errors.New("invalid parameter")
}

//
func Get(ctx *utils.CommandContext, fileName string) ([]byte, error) {
	fileFullName, err := config.GetConfigPath(ctx, fileName)
	if err != nil {
		return nil, err
	}

	if r, e := ExistFile(ctx, fileFullName); !r || e != nil {
		return nil, errors.New("file is not exists")
	}

	bytes, err := ioutil.ReadFile(fileFullName)

	return bytes, err
}
