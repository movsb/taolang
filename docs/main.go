//go:build wasm

package main

import (
	"bytes"
	"embed"
	"io/fs"
	"log"
	"syscall/js"

	"github.com/movsb/taolang/taolang"
)

//go:embed examples
var root embed.FS

func execute(source string) (string, error) {
	buf := bytes.NewBuffer(nil)
	taolang.Stdout = buf
	p := taolang.NewProgram()
	if err := p.Load(source); err != nil {
		log.Println(err)
		return ``, err
	}
	if _, err := p.CallFunc(`main`); err != nil {
		log.Println(err)
		return ``, err
	}
	return buf.String(), nil
}

func examples() map[string]any {
	dir, _ := fs.Sub(root, `examples`)
	entries, err := fs.Glob(dir, `*.tao`)
	if err != nil {
		panic(err)
	}

	contents := map[string]any{}
	for _, name := range entries {
		content, _ := fs.ReadFile(dir, name)
		contents[name] = string(content)
	}
	return contents
}

func main() {
	js.Global().Set(`execute`, js.FuncOf(func(this js.Value, args []js.Value) any {
		source := args[0].String()
		output, err := execute(source)
		if err != nil {
			return js.ValueOf(err.Error())
		}
		return js.ValueOf(output)
	}))
	js.Global().Set(`examples`, js.FuncOf(func(this js.Value, args []js.Value) any {
		return js.ValueOf(examples())
	}))
	select {}
}
