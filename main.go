package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/LazerSharp/go-git/gogit"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	args := os.Args
	//fmt.Println("Hello World!", args)
	switch args[1] {
	case "init":
		cmd_init(args[2:])
	case "_":
		panic("Bad command!")
	}
}

func cmd_init(args []string) {

	var path string
	if len(args) < 1 {
		path = "."
	} else {
		path = args[0]
	}
	gogit.Check(gogit.RepoCreate(path))
	fmt.Println("Git repo initialized!")
}
