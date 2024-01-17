package main

import (
	"github.com/alecthomas/kong"
	"todo-simple/pkg/cmd"
)

var CLI struct {
	List   cmd.ListCmd   `cmd:"" help:"List task"`
	Add    cmd.AddCmd    `cmd:"" help:"Add task"`
	Delete cmd.DeleteCmd `cmd:"" help:"Add task"`
	Update cmd.UpdateCmd `cmd:"" help:"Update status of the specified task"`
}

func main() {
	ctx := kong.Parse(&CLI)
	err := ctx.Run(&kong.Context{})
	ctx.FatalIfErrorf(err)
}
