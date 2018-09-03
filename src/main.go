package main

import (
	"fmt"
	"io"
	"os"
)

func exec(input io.Reader) {
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
		exec(os.Stdin)
	} else {
		file, err := os.Open(os.Args[1])
		if err != nil {
			panic(err)
		}
		defer file.Close()
		exec(file)
	}
}
