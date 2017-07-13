package command

import (
	"errors"
	"strings"

	"log"

	"fmt"

	"sync"

	"github.com/aspark/fileserver/config"
	"github.com/aspark/fileserver/protocal"
	"github.com/aspark/fileserver/utils"
)

//CommandConfig struct
type commandConfig struct {
	// Size     int
	Handler  func(ctx *utils.CommandContext, blocks []protocal.IBlock) (protocal.IBlockResult, error)
	SkipAuth bool
}

var cmdMapper = map[string]commandConfig{
	"AUTH":    commandConfig{Handler: handleAuthCommand, SkipAuth: true},
	"RELOAD":  commandConfig{Handler: handleReloadCommand},
	"INFO":    commandConfig{Handler: handleInfoCommand},
	"MONITOR": commandConfig{Handler: handleMonitorCommand},
	"SELECT":  commandConfig{Handler: handleSelectCommand},
	"EXISTS":  commandConfig{Handler: handleExistsCommand},
	"LEN":     commandConfig{Handler: handleLenCommand},
	"GET":     commandConfig{Handler: handleGetCommand},
	"CGET":    commandConfig{Handler: handleCGetCommand}, // continue get/load
	"SAVE":    commandConfig{Handler: handleSaveCommand},
	"DEL":     commandConfig{Handler: handleDelCommand},
	"DIR":     commandConfig{Handler: handleDirCommand},
	"PING":    commandConfig{Handler: handlePingCommand, SkipAuth: true},
}

type monitorCollection struct {
	locker sync.Mutex
	items  map[string]chan<- string
}

var _monitors = new(monitorCollection)

func registerMonitor(monitor chan<- string) string {
	_monitors.locker.Lock()
	defer _monitors.locker.Unlock()
	var key = utils.Random(30)
	if _monitors.items == nil {
		_monitors.items = make(map[string]chan<- string)
	}
	_monitors.items[key] = monitor

	return key
}

func unregisterMonitor(key string) bool {
	_monitors.locker.Lock()
	defer _monitors.locker.Unlock()

	if _, ok := _monitors.items[key]; ok {
		delete(_monitors.items, key)
		return true
	}

	return false
}

func notifyMonitor(msg string) {
	if _monitors.items == nil {
		return
	}

	_monitors.locker.Lock()
	defer _monitors.locker.Unlock()

	for _, c := range _monitors.items {
		c <- msg
	}
}

func handleNotImplementCommand(ctx *utils.CommandContext, blocks []protocal.IBlock) (protocal.IBlockResult, error) {
	log.Println("NotImplement")
	return protocal.ErrorBlock{Content: "NotImplement"}, nil
}

//ExecuteBlocks method
func ExecuteBlocks(ctx *utils.CommandContext, blocks []protocal.IBlock) ([]byte, error) {
	var err error
	var bytes = make([]byte, 0, 8)

	notifyMonitor(fmt.Sprintf("%s", blocks))
	log.Println(ctx.SessionID, " executing: ", blocks)

	if blocks == nil || len(blocks) < 1 {
		err = errors.New("need command block")
	}

	var blk = blocks[0]
	cmd, ok := blk.GetContent().(string)
	if ok {
		cmd = strings.ToUpper(cmd)
		cmdCfg, ok := cmdMapper[cmd]
		if ok {
			if !cmdCfg.SkipAuth && !ctx.HasAuth && !utils.IsNullOrWhiteSpace(config.GetConfig().Pwd) {
				err = errors.New("No Auth")
				goto end
			}

			// log.Println("Reveive command: ", cmd)
			var result protocal.IBlockResult
			result, err = cmdCfg.Handler(ctx, blocks[1:])
			if err != nil {
				goto end
			}

			if result != nil {
				log.Println(ctx.SessionID, " executed: ", result)
				notifyMonitor(fmt.Sprintf("%s", result))

				var resultBytes = result.GetBytes()
				// log.Println("result len: ", len(resultBytes), resultBytes)
				bytes = append(bytes, resultBytes...)
			}
		} else {
			err = errors.New("not support command: " + cmd)
		}
	} else {
		err = errors.New("cannot get the command string")
	}

end:
	if err != nil {
		bytes = protocal.ConvertToErrorBytes_Error(err)
	}

	return bytes, err
}
