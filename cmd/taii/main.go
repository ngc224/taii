package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/nissy/taii"
)

const (
	version = "0.0.1"
)

var (
	isHelp    = flag.Bool("h", false, "this help")
	isVersion = flag.Bool("v", false, "show version and exit")
)

func main() {
	os.Exit(exitcode(run()))
}

func exitcode(err error) int {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return 1
	}

	return 0
}

func run() error {
	flag.Parse()

	if *isHelp {
		fmt.Fprintf(os.Stderr, "Usage: %s [-h | -v] file ...\n", os.Args[0])
		flag.PrintDefaults()
		return nil
	}

	if *isVersion {
		fmt.Println("v" + version)
		return nil
	}

	var tails []*tail.Tail

	for _, v := range flag.Args() {
		t, err := tail.NewTail(v)

		if err != nil {
			return err
		}

		tails = append(tails, t)
	}

	defer func() {
		for _, v := range tails {
			v.Close()
		}
	}()

	done := make(chan bool)

	for _, v := range tails {
		go func(t *tail.Tail) {
			for {
				select {
				case e := <-t.Watcher.Events:
					if e.Op&fsnotify.Write == fsnotify.Write {
						fmt.Print(t.AddString())
					}
				case err := <-t.Watcher.Errors:
					fmt.Fprintf(os.Stderr, "Error: %s\n", err)
					done <- true
				}
			}
		}(v)
	}

	<-done

	return nil
}
