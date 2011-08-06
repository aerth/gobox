package main

import (
	"flag"
	"os"
	"common"
	"path/filepath"
)

var (
	flagSet     = flag.NewFlagSet("gobox", flag.ExitOnError)
	helpFlag    = flagSet.Bool("help", false, "Show help")
	listFlag    = flagSet.Bool("list", false, "List applets")
	installFlag = flagSet.String("install", "", "Create symlinks for applets in given path")
)

func Gobox(call []string) (e os.Error) {
	e = flagSet.Parse(call[1:])
	if e != nil {
		return
	}

	if *listFlag {
		list()
	} else if *installFlag != "" {
		e = install(*installFlag)
	} else {
		help()
	}
	return
}

func help() {
	flagSet.PrintDefaults()
	println()
	list()
}

func list() {
	println("List of compiled applets:\n")
	for name, _ := range Applets {
		print(name, ", ")
	}
	println("")
}

func install(path string) os.Error {
	goboxpath, e := common.GetGoboxBinaryPath()
	if e != nil {
		return e
	}
	for name, _ := range Applets {
		// Don't overwrite the executable
		if name == "gobox" {
			continue
		}
		newpath := filepath.Join(path, name)
		e = common.ForcedSymlink(goboxpath, newpath)
		if e != nil {
			common.DumpError(e)
		}
	}
	return nil
}