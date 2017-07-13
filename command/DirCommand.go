package command

import (
	"os"

	"errors"

	"io/ioutil"

	"github.com/aspark/fileserver/config"
	"github.com/aspark/fileserver/protocal"
	"github.com/aspark/fileserver/utils"
)

func handleDirCommand(ctx *utils.CommandContext, args []protocal.IBlock) (protocal.IBlockResult, error) {

	dir, err := config.GetConfigPath(ctx, args[0].GetContent().(string))
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(dir)
	if os.IsNotExist(err) || !info.IsDir() {
		return nil, errors.New(dir + " is not directory or not exists")
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var items = make([]protocal.IBlock, 0, len(files))

	// var dirBase, _ = config.GetConfigPath(ctx, "./")
	for _, file := range files {
		var name = file.Name()
		if file.IsDir() {
			// continue
			name += "/"
		}

		// rel, err := filepath.Rel(dirBase, file.Name())
		if err != nil {
			return nil, err
		}

		items = append(items, protocal.BulkStringBlock{Content: name}) //rel
	}

	var result = protocal.ArrayBlock{Content: items}

	return result, nil
}
