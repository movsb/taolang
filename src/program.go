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

	global := NewContext("--global--", nil)

	InitBuiltins(global)

	for _, stmt := range p.stmts {
		stmt.Execute(global)
	}

	main := &CallExpression{
		Callable: ValueFromVariable("main"),
		Args:     &Arguments{},
	}

	return main.Evaluate(global), nil
}
