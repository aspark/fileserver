package command

import (
	"errors"

	"strconv"

	"time"

	"crypto/md5"

	"fmt"

	"io"

	"strings"

	"log"

	"math"

	"github.com/aspark/fileserver/config"
	"github.com/aspark/fileserver/protocal"
	"github.com/aspark/fileserver/utils"
)

func handleAuthCommand(ctx *utils.CommandContext, args []protocal.IBlock) (protocal.IBlockResult, error) {
	if ctx.HasAuth {
		return protocal.SimpleStringBlock{Content: "OK"}, nil
	}

	user, ok := args[0].GetContent().(string)
	if ok {
		pwd, _ := args[1].GetContent().(string)
		timestamp, _ := args[2].GetContent().(string)
		success, err := Auth(user, pwd, timestamp)
		if err != nil {
			return nil, err
		}

		if success {
			log.Println("Auth success: ", user)

			ctx.HasAuth = true
			ctx.User = user
			return protocal.SimpleStringBlock{Content: "OK"}, err
		}
	} else {
		return nil, errors.New("invalid parameter")
	}

	return nil, errors.New("UNKnown FAIL")
}

//
func Auth(user string, hash string, timestamp string) (bool, error) {

	stamp, err := strconv.ParseInt(timestamp, 10, 64)

	if err != nil {
		return false, errors.New("invalid timestamp: " + timestamp)
	}

	var sub = time.Now().Sub(time.Unix(stamp, 0)).Minutes()
	if math.Abs(sub) > 3 { //3 mins
		return false, errors.New("timestamp expired: " + timestamp)
	}

	cfg := config.GetConfig()
	enc := md5.New()
	io.WriteString(enc, cfg.Pwd)
	io.WriteString(enc, timestamp)
	pwd := fmt.Sprintf("%x", enc.Sum(nil))

	if cfg.User != user || strings.ToLower(hash) != strings.ToLower(pwd) {
		return false, errors.New("user or pwd is wrong")
	}

	return true, nil
}
