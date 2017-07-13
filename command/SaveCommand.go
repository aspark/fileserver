package command

import (
	"errors"
	"io/ioutil"
	"os"

	"path/filepath"

	"github.com/aspark/fileserver/config"
	"github.com/aspark/fileserver/protocal"
	"github.com/aspark/fileserver/utils"
)

func handleSaveCommand(ctx *utils.CommandContext, args []protocal.IBlock) (protocal.IBlockResult, error) {
	fileFullName, ok := args[0].GetContent().(string)
	if ok {
		err := Save(ctx, fileFullName, args[1].GetContent().([]byte))
		if err != nil {
			return nil, err
		}

		return protocal.SimpleStringBlock{Content: "OK"}, nil // protocal.ConvertToIntegerBytes(1)
	}

	return nil, errors.New("invalid parameter")
}

//
func Save(ctx *utils.CommandContext, fileName string, bytes []byte) error {
	fileFullName, err := config.GetConfigPath(ctx, fileName)
	if err != nil {
		return err
	}

	var dir = filepath.Dir(fileFullName)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, os.ModePerm)
	}

	return ioutil.WriteFile(fileFullName, bytes, os.ModePerm)
}
