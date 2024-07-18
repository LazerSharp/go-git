package main

import (
	"fmt"
	"log"
	"os"

	"github.com/LazerSharp/go-git/gogit"
)

func main() {
	//slog.SetLogLoggerLevel(slog.LevelDebug)

	args := os.Args
	//fmt.Println("Hello World!", args)
	subcmd := args[1]
	restArgs := args[2:]
	switch subcmd {
	case "init":
		cmdInit(restArgs)
	case "cat-file":
		cmdCatFile(restArgs)
	case "hash-object":
		cmdHashObject(restArgs)
	default:
		panic("Bad command!")
	}
}

func cmdInit(args []string) {

	var path string
	if len(args) < 1 {
		path = "."
	} else {
		path = args[0]
	}
	gogit.Check(gogit.RepoCreate(path))
	fmt.Println("Git repo initialized!")
}

func cmdCatFile(args []string) {

	if len(args) != 2 {
		log.Fatal("cat-file: invalid args")
	}
	typ := args[0]
	obj := args[1]
	gogit.CatFile(typ, obj)

}

func cmdHashObject(args []string) {

	if len(args) < 1 {
		log.Fatal("hash-object: file name missing")
	}
	fpth := args[0]
	f := gogit.Must(os.Open(fpth))
	defer func() {
		gogit.Check(f.Close())
	}()
	obj := gogit.Must(gogit.NewBlob(f))
	sha := gogit.Must(gogit.WriteObject(obj, nil))
	fmt.Println(sha)
}
