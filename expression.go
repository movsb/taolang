package main

type Expression interface {
	Evaluate(ctx *Context) *Value
}

type UnaryExpression struct {
	tt   TokenType
	expr Expression
}

func NewUnaryExpression(tt TokenType, expr Expression) *UnaryExpression {
	return &UnaryExpression{
		tt:   tt,
		expr: expr,
	}
}

func (u *UnaryExpression) Evaluate(ctx *Context) *Value {
	value := u.expr.Evaluate(ctx)
	switch u.tt {
	case ttSubstraction:
		if value.Type != vtNumber {
			panic("-value is invalid")
		}
		return ValueFromNumber(-value.Number)
	case ttNot:
		switch value.Type {
		case vtNil:
			return ValueFromBoolean(true)
		case vtBoolean:
			return ValueFromBoolean(!value.Bool)
		case vtNumber:
			return ValueFromBoolean(!(value.Number != 0))
		case vtString:
			return ValueFromBoolean(!(len(value.Str) != 0))
		default:
			panic("!value is invalid")
		}
	}
	panicf("unknown unary operator: %v", u.tt) // TODO
	return nil
}

type BinaryExpression struct {
	left  Expression
	op    TokenType
	right Expression
}

func NewBinaryExpression(left Expression, op TokenType, right Expression) *BinaryExpression {
	return &BinaryExpression{
		left:  left,
		op:    op,
		right: right,
	}
}

func (b *BinaryExpression) Evaluate(ctx *Context) *Value {
	lv := b.left.Evaluate(ctx)
	rv := b.right.Evaluate(ctx)
	lt, rt := lv.Type, rv.Type
	op := b.op

	if lt == vtNil && rt == vtNil {
		if op == ttEqual {
			return ValueFromBoolean(true)
		} else if op == ttNotEqual {
			return ValueFromBoolean(false)
		} else {
			panic("not supported operator on two nils")
		}
	}

	if lt == vtBoolean && rt == vtBoolean {
		switch op {
		case ttEqual:
			return ValueFromBoolean(lv.Bool == rv.Bool)
		case ttNotEqual:
			return ValueFromBoolean(lv.Bool != rv.Bool)
		default:
			panic("not supported operator on two booleans")
		}
	}

	if lt == vtNumber && rt == vtNumber {
		switch op {
		case ttAddition:
			return ValueFromNumber(lv.Number + rv.Number)
		case ttSubstraction:
			return ValueFromNumber(lv.Number - rv.Number)
		case ttMultiply:
			return ValueFromNumber(lv.Number * rv.Number)
		case ttDivision:
			if rv.Number == 0 {
				panic("divide by zero")
			}
			return ValueFromNumber(lv.Number / rv.Number)
		case ttGreaterThan:
			return ValueFromBoolean(lv.Number > rv.Number)
		case ttGreaterThanOrEqual:
			return ValueFromBoolean(lv.Number >= rv.Number)
		case ttLessThan:
			return ValueFromBoolean(lv.Number < rv.Number)
		case ttLessThanOrEqual:
			return ValueFromBoolean(lv.Number <= rv.Number)
		case ttEqual:
			return ValueFromBoolean(lv.Number == rv.Number)
		case ttNotEqual:
			return ValueFromBoolean(lv.Number != rv.Number)
		default:
			panic("not supported operator on two numbers")
		}
	}

	if lt == vtString && rt == vtString {
		switch op {
		case ttAddition:
			return ValueFromString(lv.Str + rv.Str)
		default:
			panic("not supported operator on two strings")
		}
	}

	panic("unknown binary operator and operands")
}

type Parameters struct {
	names []string
}

func (p *Parameters) Len() int {
	return len(p.names)
}

func (p *Parameters) GetParam(name string) string {
	for _, param := range p.names {
		if param == name {
			return param
		}
	}
	return ""
}

func (p *Parameters) GetAt(index int) string {
	if index > len(p.names)-1 {
		panic("parameter index out of range")
	}
	return p.names[index]
}

func (p *Parameters) PutParam(name string) {
	p.names = append(p.names, name)
}

type FunctionExpression struct {
	name   string
	params *Parameters
	block  *BlockStatement
}

func NewFunctionExpression(name string, params *Parameters, block *BlockStatement) *FunctionExpression {
	return &FunctionExpression{
		name:   name,
		params: params,
		block:  block,
	}
}

func (f *FunctionExpression) Evaluate(ctx *Context) *Value {
	value := ValueFromFunction(f.name, f)
	if f.name != "" {
		ctx.AddValue(f.name, value)
	}
	return value
}

type Arguments struct {
	exprs []Expression
}

func (a *Arguments) Len() int {
	return len(a.exprs)
}

func (a *Arguments) PutArgument(expr Expression) {
	a.exprs = append(a.exprs, expr)
}

func (a *Arguments) EvaluateAll(ctx *Context) Values {
	args := []*Value{}
	for _, expr := range a.exprs {
		args = append(args, expr.Evaluate(ctx))
	}
	return args
}

// IndexExpression is
// obj.key    -> key: identifier whose name is "key"
// obj[key]   -> key: expression that returns string
// formally: obj should be `indexable', which supports
// syntaxes like: "str".len(), or: 123.str()
type IndexExpression struct {
	indexable Expression
	key       Expression
}

func (i *IndexExpression) Evaluate(ctx *Context) *Value {
	value := i.indexable.Evaluate(ctx)
	indexer, ok := value.Interface().(Indexer)
	if !ok {
		panic("value of expr is not indexable")
	}
	key := i.key.Evaluate(ctx)
	if key.Type != vtString {
		panic("key is not string")
	}
	return indexer.Index(key.Str)
}

type CallExpression struct {
	Callable Expression
	Args     *Arguments
}

func (f *CallExpression) Evaluate(ctx *Context) *Value {
	callable := f.Callable.Evaluate(ctx)
	if callable.Type == vtVariable {
		callable = callable.Evaluate(ctx)
	}

	switch callable.Type {
	case vtFunction:
		break
	case vtBuiltin:
		break
	case vtNil:
		panic("cannot call on nil value")
	case vtBoolean:
		panic("cannot call on boolean value")
	case vtNumber:
		panic("cannot call on number value")
	case vtString:
		panic("cannot call on string value")
	case vtObject:
		panic("cannot call on object literal")
	default:
		panic("cannot call on unknown expr")
	}

	switch callable.Type {
	case vtFunction:
		fn := callable.Func.(*FunctionExpression)
		if len(f.Args.exprs) != fn.params.Len() {
			panic("parameters and arguments don't match")
		}
		newCtx := NewContext(ctx)
		for i := 0; i < f.Args.Len(); i++ {
			newCtx.AddValue(
				fn.params.GetAt(i),
				f.Args.exprs[i].Evaluate(ctx),
			)
		}
		fn.block.Execute(newCtx)
		if ret, ok := fn.block.Return(); ok {
			return ret
		} else {
			return ValueFromNil()
		}
	case vtBuiltin:
		newCtx := NewContext(ctx)
		args := f.Args.EvaluateAll(ctx)
		return callable.Builtin(newCtx, args)
	default:
		panic("bad call")
	}
	return ValueFromNil()
}

type ObjectExpression struct {
	props map[string]Expression
}

func NewObjectExpression() *ObjectExpression {
	return &ObjectExpression{
		props: make(map[string]Expression),
	}
}

func (o *ObjectExpression) Evaluate(ctx *Context) *Value {
	obj := NewObject()
	for k, v := range o.props {
		obj.props[k] = *v.Evaluate(ctx)
	}
	return ValueFromObject(obj)
}
