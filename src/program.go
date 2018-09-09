package main

// Program is a parsed source code, and it is executable.
type Program struct {
	stmts []Statement
}

// Execute executes the program.
func (p *Program) Execute() (ret Value, err error) {
	defer func() {
		err = toErr(recover())
	}()

	globalObject := NewGlobal()
	globalContext := NewContext("--global--", nil)
	globalContext.AddObject("global", globalObject)

	for _, stmt := range p.stmts {
		stmt.Execute(globalContext)
	}

	main := &CallExpression{
		Callable: ValueFromVariable("main"),
		Args:     &Arguments{},
	}

	return main.Evaluate(globalContext), nil
}
