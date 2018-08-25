package main

type Program struct {
	stmts []Statement
}

func (p *Program) Execute() {
	global := NewContext(nil)
	InitBuiltins(global)
	for _, stmt := range p.stmts {
		stmt.Execute(global)
	}
	main := &CallExpression{
		Callable: ValueFromVariable("main"),
		Args:     &Arguments{},
	}
	main.Evaluate(global)
}
