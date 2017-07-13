package core

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"reflect"

	"errors"

	"path/filepath"

	"github.com/aspark/fileserver/command"
	"github.com/aspark/fileserver/config"
	"github.com/aspark/fileserver/protocal"
	"github.com/aspark/fileserver/utils"
)

var _flag = true
var _tcp net.Listener

//
func StartFileServer() error {
	port := config.GetConfig().Port
	if port <= 0 {
		msg := "tcp port must be great than zero"
		log.Println(msg)

		return errors.New(msg)
	}

	var err error
	_tcp, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		msg := "listen error:" + err.Error()
		log.Println(msg)

		return errors.New(msg)
	}

	log.Println("start fileserver, listening on:", port)

	for {
		c, err := _tcp.Accept()
		if err != nil {
			log.Println("accept error:", err)
			break
		}

		go handleConn(c)
	}

	return nil
}

//
func StopFileServer() {
	_flag = false
	_tcp.Close()
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	//create cmd ctx
	ctx := config.CreateContext(&conn)

	var buf = make([]byte, 1024)
	for {
		//ip check
		var cfg = config.GetConfig()
		if len(cfg.AllowedIP) > 0 {
			if strings.Index(cfg.AllowedIP, ctx.ClientIP) < 0 {
				log.Println("not allowed client:", ctx.ClientIP)
				return
			}
		}

		//read from conn
		n, err := conn.Read(buf)
		if err == io.EOF {
			log.Println("disconnect", conn.RemoteAddr())
			break
		} else if err != nil {
			log.Println("conn error:", err)
			break
		} else if n > 0 {
			//begin parse to block
			blocks, err := protocal.ParseAll(buf[:n], conn) //data

			if err != nil {
				log.Println(err)
				break
			}

			log.Println("receive blocks: ", blocks)

			//execute the block
			for _, block := range blocks {
				if reflect.TypeOf(block) == reflect.TypeOf(protocal.ArrayBlock{}) {
					bytes, err := command.ExecuteBlocks(ctx, block.GetContent().([]protocal.IBlock))
					if bytes != nil {
						_, err := conn.Write(bytes) //write the response
						if err != nil {
							log.Println(err)
							break
						}
					}

					if err != nil {
						log.Println(err)
						break
					}
				} else {
					log.Println("only exec array block, skip: ", block.GetContent())
				}
			}
		}
	}
}

//
func LogToFile() {
	var now = time.Now()
	var fileName = utils.ResolvePath(fmt.Sprintf("log/fileserver.%d-%02d-%02d.log", now.Year(), now.Month(), now.Day()))
	os.Mkdir(filepath.Dir(fileName), os.ModePerm)
	var logFile, err = os.OpenFile(fileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModePerm)
	if err != nil {
		fmt.Println("server start failed: ", err)
		return
	}

	log.SetOutput(logFile)

}
