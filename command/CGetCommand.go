package command

import (
	"errors"

	"os"

	"strconv"

	"github.com/aspark/fileserver/config"
	"github.com/aspark/fileserver/protocal"
	"github.com/aspark/fileserver/utils"
)

// get continue
func handleCGetCommand(ctx *utils.CommandContext, args []protocal.IBlock) (protocal.IBlockResult, error) {
	fileName, ok := args[0].GetContent().(string)
	if ok {
		var offsize, length int64 = 0, -1

		if len(args) > 1 {
			str, ok := args[1].GetContent().(string)
			if ok {
				i, err := strconv.ParseInt(str, 10, 64)
				if err != nil {
					return nil, err
				}
				offsize = i
			} else {
				offsize = int64(args[1].GetContent().(int))
			}
		}

		if len(args) > 2 {
			str, ok := args[2].GetContent().(string)
			if ok {
				i, err := strconv.ParseInt(str, 10, 64)
				if err != nil {
					return nil, err
				}
				length = i
			} else {
				length = int64(args[2].GetContent().(int))
			}
		}

		var size, start, real, content, err = CGet(ctx, fileName, offsize, length)
		if err != nil {
			return nil, err
		}

		rstBlocks := []protocal.IBlock{
			protocal.Integer64Block{Content: size},
			protocal.Integer64Block{Content: start},
			protocal.Integer64Block{Content: real},
			protocal.BytesBlock{Content: content}, // file length; offsize; real read length; file content bytes
		}
		result := protocal.ArrayBlock{Content: rstBlocks}

		return result, nil
	}

	return nil, errors.New("invalid parameter")
}

//
func CGet(ctx *utils.CommandContext, fileName string, args ...int64) (int64, int64, int64, []byte, error) {
	fileFullName, err := config.GetConfigPath(ctx, fileName)
	if err != nil {
		return 0, 0, 0, nil, err
	}

	if r, e := ExistFile(ctx, fileFullName); !r || e != nil {
		return 0, 0, 0, nil, errors.New("file is not exists")
	}

	fi, _ := os.Stat(fileFullName)
	start, length := int64(0), int64(-1)

	if args != nil {
		if len(args) > 0 {
			start = args[0]
		}
		if len(args) > 1 {
			length = args[1]
		}
	}

	size := fi.Size()
	if start < 0 {
		return 0, 0, 0, nil, errors.New("start must great than zero")
	}

	if length < 0 || start+length > size {
		length = size - start
	}

	f, err := os.OpenFile(fileFullName, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return 0, 0, 0, nil, err
	}

	f.Seek(start, 0)

	bytes := make([]byte, length)
	l, err := f.Read(bytes)
	if err != nil {
		return 0, 0, 0, nil, err
	}

	if int64(l) < length {
		bytes = bytes[:l]
		length = int64(l) //real read bytes length
	}

	return size, start, length, bytes, err
}
