package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"errors"

	"net"

	"github.com/aspark/fileserver/utils"
)

//
func CreateContext(conn *net.Conn) *utils.CommandContext {
	ctx := new(utils.CommandContext) //utils.CreateContext()
	ctx.SessionID = utils.Random(18)

	for k := range GetConfig().Path {
		ctx.SelectedPath = k //set to first key
		break
	}

	if conn != nil {
		var remoteAddr = (*conn).RemoteAddr()
		ip, _, _ := net.SplitHostPort(remoteAddr.String())

		ctx.Conn = conn
		ctx.ClientIP = ip

		log.Println("connnect: ", remoteAddr)
	}

	return ctx
}

//
func GetConfigPath(ctx *utils.CommandContext, fileName string) (string, error) {

	var selectedPath = GetConfig().Path[ctx.SelectedPath]
	path := utils.ResolvePath(selectedPath)

	if utils.IsNullOrWhiteSpace(fileName) {
		return path, nil
	}

	if utils.IsNullOrWhiteSpace(fileName) || ((strings.Contains(fileName, "..") || filepath.IsAbs(fileName)) && !utils.IsTesting()) {
		return "", errors.New("filename contains invalid chars")
	}

	return filepath.Join(path, fileName), nil
}

//
type Config struct {
	lock      sync.RWMutex
	hasLoaded bool
	Port      int               // default: 19860
	Path      map[string]string //default: ./attachments
	User      string
	Pwd       string
	AllowedIP string //IP白名单
}

var cfg = new(Config)

//
func GetConfig() *Config {
	// log.Println("config", cfg, cfg.lock)
	if cfg.hasLoaded {
		return cfg
	}

	cfg.lock.Lock()
	defer cfg.lock.Unlock()

	if cfg.hasLoaded {
		return cfg
	}

	cfg.hasLoaded = false
	defer (func() { cfg.hasLoaded = true })()

	file, err := os.Open(utils.ResolvePath("./config.json"))
	if err != nil {
		// log.Fatal(err)
		log.Println("cannot find config file, use default")
	} else {
		decoder := json.NewDecoder(file)
		err = decoder.Decode(cfg)
		if err != nil {
			// log.Fatal(err)
			log.Println("decode config failed: ", err.Error())
		}
	}

	ensureDefaultValue(cfg)

	return cfg
}

//ReloadConfig Reload Config
func ReloadConfig() {
	cfg.lock.Lock()
	defer cfg.lock.Unlock()

	cfg.hasLoaded = false
}

func ensureDefaultValue(config *Config) {
	if config.Path == nil || len(config.Path) == 0 {
		config.Path = make(map[string]string)
		config.Path["default"] = "./attachments"
	}

	if config.Port <= 0 {
		config.Port = 19860
	}
}
