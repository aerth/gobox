package rm

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	flagSet       = flag.NewFlagSet("rm", flag.PanicOnError)
	recursiveFlag = flagSet.Bool("r", false, "remove directories and their contents recursively")
	forceFlag     = flagSet.Bool("f", false, "ignore nonexistent files and arguments, never prompt")
	helpFlag      = flagSet.Bool("help", false, "Show this help")
)

func Rm(call []string) error {
	e := flagSet.Parse(call[1:])
	if e != nil {
		return e
	}

	if flagSet.NArg() <= 0 || *helpFlag {
		println("`rm` [options] <files...>")
		flagSet.PrintDefaults()
		return nil
	}

	for _, file := range flagSet.Args() {
		e = delete(file)
		if e != nil {
			if *forceFlag {
				continue
			}
			return e
		}
	}
	return e
}

func delete(file string) error {
	fi, e := os.Stat(file)
	if e != nil {
		if *forceFlag {
			return nil
		}
		return e
	}
	if fi.IsDir() && *recursiveFlag {
		e := deleteDir(file)
		if e != nil && !*forceFlag {
			return e
		}

	}
	return os.Remove(file)
}

func deleteDir(dir string) error {
	files, e := ioutil.ReadDir(dir)
	if e != nil {
		return e
	}
	for _, file := range files {
		e = delete(filepath.Join(dir, file.Name()))
		if e != nil {
			if *forceFlag {
				continue
			}
			return e
		}
	}
	return nil
}
