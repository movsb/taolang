package main

import (
	"io"
	"strings"
)

// Chunk is a piece of code.
type Chunk struct {
	stmts []Statement
}

// Execute executes this chunk.
func (c *Chunk) Execute(ctx *Context) (err error) {
	defer catchAsError(&err)

	for _, stmt := range c.stmts {
		stmt.Execute(ctx)
	}
	return
}

// Program is a runtime state.
type Program struct {
	globalObject  *Global
	globalContext *Context
}

// NewProgram news a program.
func NewProgram() *Program {
	p := &Program{}
	p.globalContext = NewContext("--global--", nil)
	p.globalObject = NewGlobal()
	p.globalContext.AddObject("global", p.globalObject)
	return p
}

// Load compiles and executes the source code.
func (p *Program) Load(source string) error {
	return p.LoadInput(strings.NewReader(source))
}

// MustLoad must Load.
func (p *Program) MustLoad(source string) {
	err := p.LoadInput(strings.NewReader(source))
	if err != nil {
		panic(err)
	}
}

// LoadInput compiles and executes the source code.
func (p *Program) LoadInput(input io.Reader) error {
	tokenizer := NewTokenizer(input)
	parser := NewParser(tokenizer)
	chunk, err := parser.Parse()
	if err != nil {
		return err
	}
	return chunk.Execute(p.globalContext)
}

// MustLoadInput must LoadInput.
func (p *Program) MustLoadInput(input io.Reader) {
	err := p.LoadInput(input)
	if err != nil {
		panic(err)
	}
}

// CallFunc calls a named function with arguments.
func (p *Program) CallFunc(fn string, args ...Expression) (ret Value, err error) {
	defer catchAsError(&err)
	f := ValueFromVariable(fn)
	return CallFunc(p.globalContext, f, args...), nil
}
