package utils

import (
	"net"
	"sync"
)

type sessionKey string

const contextSessionKey = sessionKey("aspark")

//
type CommandContext struct {
	// context.Context
	SessionID    string
	User         string
	HasAuth      bool
	ClientIP     string
	Conn         *net.Conn
	sync         sync.Mutex
	SelectedPath string
}

//
// func CreateContext() *CommandContext {
// 	// return context.WithValue(context.Background(), utils.ContextSessionKey, utils.Session{})
// 	var ctx = new(CommandContext)

// 	return ctx
// }
