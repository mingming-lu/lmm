package cli

import (
	"context"
	"log"
)

var commands map[string]Command

// Command defines the cli command interface
type Command interface {
	Execute(context.Context) error
}

func Register(name string, cmd Command) {
	if name == "" {
		panic("command name is not allowed to be empty")
	}
	if _, exists := commands[name]; exists {
		panic("duplicate command name: " + name)
	}
	commands[name] = cmd
}

func Execute(c context.Context, name string) {
	cmd, ok := commands[name]
	if !ok {
		log.Print("[warn] not found command named:", name)
		return
	}
	if err := cmd.Execute(c); err != nil {
		log.Print("[error] failed to execute comamnd:", `'`+name+`'`, err.Error())
	}
}

func init() {
	commands = make(map[string]Command)
}
