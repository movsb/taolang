package main

import (
	"fmt"
	"os"
)

func main() {
	input, err := os.Open("main.tao")
	if err != nil {
		panic(err)
	}
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

	fmt.Println("main exit:", ret)
}
