package command

import (
	"errors"
	"os"

	"log"

	"github.com/aspark/fileserver/config"
	"github.com/aspark/fileserver/protocal"
	"github.com/aspark/fileserver/utils"
)

func handleLenCommand(ctx *utils.CommandContext, args []protocal.IBlock) (protocal.IBlockResult, error) {
	fileName, ok := args[0].GetContent().(string)
	if ok {
		len, err := Len(ctx, fileName)
		if err != nil {
			return nil, err
		}

		return protocal.Integer64Block{Content: len}, nil
	}

	return nil, errors.New("invalid parameter")
}

//
func Len(ctx *utils.CommandContext, fileName string) (int64, error) {
	fileFullName, err := config.GetConfigPath(ctx, fileName)
	if err != nil {
		return -1, err
	}

	if r, e := ExistFile(ctx, fileFullName); !r || e != nil {
		return -1, nil //errors.New("the file is not exists")
	}

	fi, err := os.Stat(fileFullName)

	if err != nil {
		log.Println(err)
		return -1, err
	}

	return fi.Size(), nil
}
