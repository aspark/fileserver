package main

import (
	"crypto/md5"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/aspark/fileserver/command"
	"github.com/aspark/fileserver/config"
	"github.com/aspark/fileserver/protocal"
)

func Test_AuthCommand(t *testing.T) {
	stamp := strconv.FormatInt(time.Now().Unix(), 10)
	var bytes = protocal.ConvertToBulkStringBytesEx("AUTH", "user", "pwd", stamp)

	blocks, _ := protocal.ParseAll(bytes, nil)
	ctx := config.CreateContext(nil)
	_, err := command.ExecuteBlocks(ctx, blocks)

	if err == nil {
		t.Fatal("need auth failed")
	}

	bytes = protocal.ConvertToBulkStringBytesEx("AUTH", "user", fmt.Sprintf("%x", md5.Sum([]byte("pwd"+stamp))), stamp)

	blocks, _ = protocal.ParseAll(bytes, nil)
	_, err = command.ExecuteBlocks(ctx, blocks)

	if err != nil || ctx.HasAuth == false {
		t.Fatal("auth failed")
	}
}
