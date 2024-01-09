package main

import (
	"github.com/alecthomas/kong"
	"todo-simple/pkg/cmd"
)

var CLI struct {
	Add cmd.AddCmd `cmd:"" help:"Add task"`
	Delete cmd.DeleteCmd `cmd:"" help:"Delete task"`
	List cmd.ListCmd `cmd:"" help:"List task"`
}

func main() {
	ctx := kong.Parse(&CLI)
	err := ctx.Run(&kong.Context{})
	ctx.FatalIfErrorf(err)
}
