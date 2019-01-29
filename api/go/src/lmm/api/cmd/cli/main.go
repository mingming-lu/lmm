package main

import (
	"context"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"lmm/api/cli"
	_ "lmm/api/cli/simple"
)

func main() {
	commands := os.Args[1:]

	c := context.TODO()

	for _, command := range commands {
		cli.Execute(c, command)
	}
}
