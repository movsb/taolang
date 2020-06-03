package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/movsb/taolang/taolang"
)

func main() {
	var err error
	var file io.ReadCloser

	var doMain bool
	flag.BoolVar(&doMain, "m", false, "run main")
	flag.BoolVar(&doMain, "main", false, "run main")

	flag.Parse()

	if flag.NArg() == 0 || flag.NArg() == 1 && flag.Arg(0) == "-" {
		file = ioutil.NopCloser(os.Stdin)
	} else {
		file, err = os.Open(flag.Arg(0))
		if err != nil {
			panic(err)
		}
	}

	defer file.Close()
	program := taolang.NewProgram()
	err = program.LoadInput(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if doMain {
		_, err := program.CallFunc("main")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
