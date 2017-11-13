package main

import (
	"flag"
	"fmt"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/nfnt/resize"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Usage: ecnv image [image ...]")
		return
	}

	var out string
	var x, y uint
	flag.StringVar(&out, "out", "thumb", "output dir")
	flag.UintVar(&x, "x", 120, "x size")
	flag.UintVar(&x, "y", 0, "y size")
	flag.Parse()

	err := os.Mkdir(out, 0755)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to create dir: %s", err)
		os.Exit(1)
	}

	var wg sync.WaitGroup
	i := 0
	for _, name := range flag.Args() {
		wg.Add(1)
		i++
		log.Printf("goroutine %d started: %s\n", i, name)
		go func(name string) {
			defer wg.Done()
			if err := resizeImage(out, x, y, name); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			log.Println("done")
		}(name)
	}

	wg.Wait()
}

func resizeImage(out string, x, y uint, name string) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()

	i, err := jpeg.Decode(f)
	if err != nil {
		return err
	}

	i = resize.Resize(x, y, i, resize.Lanczos3)
	n := filepath.Join(out, strings.Split(name, ".")[0]+".thumb.jpg")
	fw, err := os.Create(n)
	if err != nil {
		return err
	}
	defer fw.Close()
	if err := jpeg.Encode(fw, i, nil); err != nil {
		return err
	}
	return nil
}
