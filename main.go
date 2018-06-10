package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/ktr0731/ecnv/lib/ecnv"
)

var tokens chan struct{}

func main() {
	var out, format string
	var x, y, worker uint
	flag.StringVar(&out, "out", "thumb", "output dir")
	flag.UintVar(&x, "x", 120, "x size")
	flag.UintVar(&y, "y", 0, "y size")
	flag.UintVar(&worker, "worker", 100, "worker num")
	flag.StringVar(&format, "format", ".jpg", "ext format")
	flag.Parse()

	tokens = make(chan struct{}, worker)

	if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "please pass jpg images")
		os.Exit(1)
	}

	if _, err := os.Stat(out); os.IsNotExist(err) {
		err := os.MkdirAll(out, 0755)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to create dir: %s\n", err)
			os.Exit(1)
		}
	}

	var wg sync.WaitGroup
	i := 0
	for _, name := range flag.Args() {
		wg.Add(1)
		i++
		log.Printf("[%f] %s goroutine %d waiting...\n", float64(i)/float64(len(flag.Args())), name, i)
		go func(name string, i int) {
			defer wg.Done()
			if err := resizeImage(i, out, x, y, name, format); err != nil {
				fmt.Fprintln(os.Stderr, color.RedString(err.Error()))
				os.Exit(1)
			}
		}(name, i)
	}

	log.Println("waiting for finish all goroutines")
	wg.Wait()
}

func resizeImage(idx int, out string, x, y uint, name, format string) error {
	tokens <- struct{}{}

	log.Printf("%s goroutine %d started\n", name, idx)

	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()

	n := filepath.Join(out, strings.Split(filepath.Base(name), ".")[0]+".thumb"+format)
	fw, err := os.Create(n)
	if err != nil {
		return err
	}
	defer fw.Close()

	err = ecnv.ResizeImage(fw, f, idx, x, y, name, format)
	if err != nil {
		return err
	}

	<-tokens
	log.Printf("done %d\n", idx)
	return nil
}
