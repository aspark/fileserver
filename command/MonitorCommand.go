package command

import (
	"errors"

	"github.com/aspark/fileserver/protocal"
	"github.com/aspark/fileserver/utils"
)

func handleMonitorCommand(ctx *utils.CommandContext, args []protocal.IBlock) (protocal.IBlockResult, error) {

	var monitor = make(chan string, 10)
	key := registerMonitor(monitor)
	defer unregisterMonitor(key)

	for line := range monitor {
		_, err := (*ctx.Conn).Write(protocal.BulkStringBlock{Content: line + "\r\n"}.GetBytes())
		if err != nil {
			break
		}
	}

	// return protocal.SimpleStringBlock{Content: "EOF"}, nil
	return nil, errors.New("stop monitor")
}
