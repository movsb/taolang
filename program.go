package main

type Program struct {
	stmts []Statement
}

func (p *Program) Execute() (ret Value, err error) {
	defer func() {
		err = toErr(recover())
	}()

	global := NewContext(nil)

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
