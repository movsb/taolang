package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Sync calls callback within program thread.
func Sync(callback func()) {
	queue <- callback
}

var queue = make(chan func(), 16)

func exec(input io.ReadCloser) {
	defer input.Close()
	tokenizer := NewTokenizer(input)
	parser := NewParser(tokenizer)
	program, err := parser.Parse()
	if err != nil {
		fmt.Println(err)
		return
	}
	ret, err := program.Execute()
	if err != nil {
		fmt.Println(err)
		return
	}
	_ = ret
}

func main() {
	if len(os.Args) == 1 {
		queue <- func() {
			exec(os.Stdin)
		}
	} else {
		file, err := os.Open(os.Args[1])
		if err != nil {
			panic(err)
		}
		queue <- func() {
			exec(file)
		}
	}

	for {
		select {
		case fn := <-queue:
			fn()
		case <-time.After(time.Hour):
		}
	}
}
