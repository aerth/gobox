package grep

import (
	"flag"
	"fmt"
	"github.com/surma/gobox/pkg/common"
	"io"
	"os"
	"path/filepath"
	"regexp"
)

var (
	flagSet       = flag.NewFlagSet("grep", flag.PanicOnError)
	helpFlag      = flagSet.Bool("help", false, "Show this help")
	invertFlag    = flagSet.Bool("v", false, "select non-matching lines")
	linenumFlag   = flagSet.Bool("n", false, "show line number")
	recursiveFlag = flagSet.Bool("r", false, "search for string in all files recursively")
)

func Grep(call []string) error {
	e := flagSet.Parse(call[1:])
	if e != nil {
		return e
	}

	if flagSet.NArg() == 0 || *helpFlag {
		println("`grep` <pattern> [<file>...]")
		flagSet.PrintDefaults()
		return nil
	}

	pattern, err := regexp.Compile(flagSet.Arg(0))
	if err != nil {
		return err
	}
	if *recursiveFlag {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		walkFn := func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			if fh, err := os.Open(path); err == nil {
				func() {
					doGrep(pattern, fh, path, true, *linenumFlag)
				}()
			}

			return nil
		}
		return filepath.Walk(wd, walkFn)

	}
	if flagSet.NArg() == 1 {
		doGrep(pattern, os.Stdin, "<stdin>", false, *linenumFlag)
	} else {
		for _, fn := range flagSet.Args()[1:] {
			if fh, err := os.Open(fn); err == nil {
				func() {
					doGrep(pattern, fh, fn, (flagSet.NArg() > 1), *linenumFlag)
					fh.Close()
				}()
			} else {
				fmt.Fprintf(os.Stderr, "grep: %s: %v\n", fn, err)
			}
		}
	}

	return nil
}

func doGrep(pattern *regexp.Regexp, fh io.Reader, fn string, printFilename, print_ln bool) {
	buf := common.NewBufferedReader(fh)
	ln := 0
	for {
		ln++
		line, err := buf.ReadWholeLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while reading from %s: %v\n", fn, err)
			return
		}

		if line == "" {
			continue
		}
		if pattern.MatchString(line) {
			if printFilename {
				fmt.Printf("%v:", fn)
			}
			if print_ln {
				fmt.Printf("%v:", ln)
			}
			fmt.Printf("%s\n", line)
		}
	}
}
