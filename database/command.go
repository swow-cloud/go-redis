package database

import "strings"

var cmdTable = make(map[string]*Command)

type Command struct {
	exector ExecFunc
	arity   int
}

func RegisterCommand(name string, exector ExecFunc, arity int) {
	name = strings.ToLower(name)
	cmdTable[name] = &Command{
		exector: exector,
		arity:   arity,
	}
}
