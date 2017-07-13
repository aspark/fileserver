package main

import (
	"crypto/md5"
	"flag"
	"io"
	"time"

	"net"

	"fmt"

	"log"

	"bufio"
	"os"

	"strings"

	"strconv"

	"github.com/aspark/fileserver/protocal"
	"github.com/aspark/fileserver/utils"
)

var _ip = flag.String("h", "127.0.0.1", "specified host")
var _port = flag.Int("p", 19860, "port")
var _user = flag.String("u", "", "user")
var _pwd = flag.String("a", "", "password")
var _selected = "default"

// var help = flag.Bool("help")

var printer = make(chan string, 1)

func main() {
	flag.Parse()

	if len(os.Args) > 1 && strings.ToLower(os.Args[1]) == "help" {
		flag.Usage()
		return
	}

	conn := createConn()

	reader := bufio.NewReader(os.Stdin)
	for {
		os.Stdout.WriteString(fmt.Sprintf("fileserver %s:%d[%s]> ", *_ip, *_port, _selected))
		line, err := reader.ReadString('\r')
		if err != nil {
			log.Println(err.Error())
			continue
		}

		line = string(line[:len(line)-1])

		var cmd string
		var args = make([]string, 0, 2)
		segs := strings.Split(line, " ")
		for i, seg := range segs {
			if utils.IsNullOrWhiteSpace(seg) {
				continue
			}

			if i == 0 {
				cmd = strings.TrimSpace(seg)
			} else {
				args = append(args, strings.Trim(strings.TrimSpace(seg), "\""))
			}
		}

		if strings.ToLower(cmd) == "quit" {
			releaseConn()
			return
		}

		err = sendCommand(conn, cmd, args)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		if strings.ToLower(cmd) == "monitor" {
			for {
				blocks, err := protocal.ParseAll(nil, conn)
				if err != nil {
					log.Println(err.Error())
					return
				}

				printBlocks(blocks)
			}
		}

		blocks, err := protocal.ParseAll(nil, conn)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		if strings.ToLower(cmd) == "select" && blocks[0].GetContent() == "OK" {
			_selected = args[0]
		}

		printBlocks(blocks)
	}
}

func execCommnad(cmd string, args []string) ([]protocal.IBlock, error) {
	var conn = createConn()
	err := sendCommand(conn, cmd, args)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return readBlocks(conn)
}

func sendCommand(conn net.Conn, cmd string, args []string) error {
	var items = make([]protocal.IBlock, 0, 1)
	items = append(items, protocal.SimpleStringBlock{Content: cmd})
	if args != nil && len(args) > 0 {
		for _, arg := range args {
			items = append(items, protocal.SimpleStringBlock{Content: arg})
		}
	}

	// log.Println("cmd: ", []byte(cmd))
	// log.Println("args: ", args)

	var bytes = protocal.ArrayBlock{Content: items}.GetBytes()
	_, err := conn.Write(bytes)
	if err != nil {
		log.Println(err.Error())
		releaseConn()

		return err
	}

	return nil
}

func readBlocks(conn net.Conn) ([]protocal.IBlock, error) {
	return protocal.ParseAll(nil, conn)
}

func printBlocks(blocks []protocal.IBlock) {
	if blocks != nil {
		for _, blk := range blocks {
			os.Stdout.Write([]byte(fmt.Sprintf("%s\r\n", blk.RawString())))
		}
	}
}

var __conn net.Conn

func createConn() net.Conn {
	if __conn == nil {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", *_ip, *_port))

		if err != nil {
			log.Println(err.Error())
		}

		__conn = conn

		if !utils.IsNullOrWhiteSpace(*_user) {
			enc := md5.New()
			io.WriteString(enc, *_pwd)
			timestamp := strconv.FormatInt(time.Now().Unix(), 10)
			io.WriteString(enc, timestamp)
			pwd := fmt.Sprintf("%x", enc.Sum(nil))
			blocks, err := execCommnad("auth", []string{*_user, pwd, timestamp})
			if err == nil {
				printBlocks(blocks)
			}
		}
	}

	return __conn
}

func releaseConn() {
	err := __conn.Close()
	if err != nil {
		log.Println(err.Error())
	}
	__conn = nil
}
