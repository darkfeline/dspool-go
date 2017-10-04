package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/fsnotify/fsnotify"
)

type watchConfig struct {
	watchDir string
	spoolDir string
}

func main() {
	configureLog()
	c := parseArgs()
	err := spoolDownloads(c)
	if err != nil {
		log.Fatal(err)
	}
}

func configureLog() {
	log.SetPrefix(fmt.Sprintf("%s: ", os.Args[0]))
	log.SetFlags(log.Flags() | log.Lshortfile)
}

func parseArgs() *watchConfig {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] WATCH_DIR SPOOL_DIR\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() != 2 {
		flag.Usage()
		os.Exit(1)
	}
	args := flag.Args()
	return &watchConfig{
		watchDir: args[0],
		spoolDir: args[1],
	}
}

func spoolDownloads(c *watchConfig) error {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer w.Close()

	err = w.Add(c.watchDir)
	if err != nil {
		return err
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for {
			select {
			case e := <-w.Events:
				log.Printf("Got event: %+v", e)
			case err := <-w.Errors:
				log.Printf("Error while watching: %s", err)
			case <-sc:
				wg.Done()
				return
			}
		}
	}()
	wg.Wait()
	return nil
}
