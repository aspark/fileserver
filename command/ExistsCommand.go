package command

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/aspark/fileserver/config"
	"github.com/aspark/fileserver/protocal"
	"github.com/aspark/fileserver/utils"
)

func handleExistsCommand(ctx *utils.CommandContext, args []protocal.IBlock) (protocal.IBlockResult, error) {
	fileName, ok := args[0].GetContent().(string)
	if ok {
		var result bool
		result, err := ExistFile(ctx, fileName)
		if err != nil {
			return nil, err
		}

		if result {
			return protocal.IntegerBlock{Content: 1}, nil
		}

		return protocal.IntegerBlock{Content: 0}, nil

	}

	return nil, errors.New("invalid parameter")
}

//
func ExistFile(ctx *utils.CommandContext, fileName string) (result bool, err error) {
	fileFullName := fileName
	if filepath.IsAbs(fileFullName) == false {
		fileFullName, err = config.GetConfigPath(ctx, fileFullName)
		if err != nil {
			return false, err
		}
	}
	info, err := os.Stat(fileFullName)
	log.Println("result", info, err, fileFullName)
	if err != nil && os.IsNotExist(err) {
		result = false
		err = nil
	} else if err == nil || os.IsExist(err) {
		result = info.IsDir() == false
		err = nil
	}

	return result, err
}
