package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/Pallinder/go-randomdata"
	"github.com/thanhpk/randstr"
)

/*
	TODO make this a configurable command.
	support options for:
		- removing created files
		- listing disk space used (similar to `du -hs dat`)
		- configurable destination folder
		- partition? X files into Y folders (Z = X/Y) Z files per folder
		- delete:
			find ./dat/ -name "*.txt" -delete
*/

func main() {
	n := os.Args[1:]
	nInt, err := strconv.Atoi(n[0])
	if err != nil {
		log.Panic(err)
	}
	doNfiles(nInt)
}

func doNfiles(n int) {
	var wg sync.WaitGroup
	maxFd := make(chan (struct{}), 20)
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(wg *sync.WaitGroup, c *chan (struct{})) {
			*c <- struct{}{}
			defer func() { <-*c }()
			defer wg.Done()

			fd, err := os.Create(fmt.Sprintf("dat/%v.txt", randstr.Hex(10)))
			if err != nil {
				log.Panic(err)
			}
			defer fd.Close()

			writeGibberish(fd)
		}(&wg, &maxFd)
	}
	wg.Wait()
	fmt.Printf("created %v files", n)
}

func writeGibberish(w io.WriteCloser) {
	p := randomdata.Alphanumeric(1000)
	w.Write([]byte(p))
}
