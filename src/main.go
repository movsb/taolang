package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"
)

type MessageType uint

const (
	_ MessageType = iota
	mtExit
	mtExec
	mtTimer
)

type Message struct {
	mt   MessageType
	data interface{}
}

var queue = make(chan Message, 1)
var exit = make(chan int)

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
	go func() {
		for {
			select {
			case msg := <-queue:
				switch msg.mt {
				case mtExit:
					exit <- 0
				case mtExec:
					exec(msg.data.(io.ReadCloser))
				case mtTimer:
					t := msg.data.(*Timer)
					c := &CallExpression{
						Callable: t.callback,
						Args:     &Arguments{},
					}
					c.Evaluate(t.ctx)
				}
			case <-time.After(time.Hour):
			}
		}
	}()

	if len(os.Args) == 1 {
		queue <- Message{
			mt:   mtExec,
			data: ioutil.NopCloser(os.Stdin),
		}
	} else {
		file, err := os.Open(os.Args[1])
		if err != nil {
			panic(err)
		}
		queue <- Message{
			mt:   mtExec,
			data: file,
		}
	}

	<-exit
}
