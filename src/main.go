package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"
)

// the message queue.
var queue = make(chan func(), 16)

// Async queues and calls the callback within main thread.
func Async(callback func()) {
	queue <- callback
}

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
	var wait bool
	flag.BoolVar(&wait, "wait", false, "wait to exit")
	flag.Parse()

	var err error
	var file io.ReadCloser

	if flag.NArg() == 0 || flag.NArg() == 1 && flag.Arg(0) == "-" {
		file = ioutil.NopCloser(os.Stdin)
	} else {
		file, err = os.Open(flag.Arg(0))
		if err != nil {
			panic(err)
		}
	}

	if wait {
		queue <- func() {
			exec(file)
		}
		for {
			select {
			case fn := <-queue:
				fn()
			case <-time.After(time.Hour):
			}
		}
	} else {
		exec(file)
	}
}
