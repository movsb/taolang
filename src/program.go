package main

// Program is a parsed source code, and it is executable.
type Program struct {
	stmts []Statement
}

// Execute executes the program.
func (p *Program) Execute() (ret Value, err interface{}) {
	defer func() {
		err = recover()
	}()

	globalObject := NewGlobal()
	globalContext := NewContext("--global--", nil)
	globalContext.AddObject("global", globalObject)

	globalContext.AddClass("Timer", &BuiltinConstructor{_NewTimer})

	for _, stmt := range p.stmts {
		stmt.Execute(globalContext)
	}

	main := ValueFromVariable("main")
	return CallFunc(globalContext, main), nil
}
