package main

import (
	"github.com/alecthomas/kong"
	"todo-simple/pkg/cmd"
)

var CLI struct {
	List cmd.ListCmd `cmd:"" help:"List task"`
}

func main() {
	ctx := kong.Parse(&CLI)
	err := ctx.Run(&kong.Context{})
	ctx.FatalIfErrorf(err)
}
