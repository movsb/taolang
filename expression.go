package main

// Expression is the interface that is implemented by all expressions.
type Expression interface {
	Evaluate(ctx *Context) Value
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

func (u *UnaryExpression) Evaluate(ctx *Context) Value {
	value := u.expr.Evaluate(ctx)
	switch u.tt {
	case ttSubstraction:
		if value.Type != vtNumber {
			panic("-value is invalid")
		}
		return ValueFromNumber(-value.number())
	case ttNot:
		switch value.Type {
		case vtNil:
			return ValueFromBoolean(true)
		case vtBoolean:
			return ValueFromBoolean(!value.boolean())
		case vtNumber:
			return ValueFromBoolean(!(value.number() != 0))
		case vtString:
			return ValueFromBoolean(!(len(value.str()) != 0))
		default:
			panic("!value is invalid")
		}
	}
	panicf("unknown unary operator: %v", u.tt) // TODO
	return ValueFromNil()
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

func (b *BinaryExpression) Evaluate(ctx *Context) Value {
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
			return ValueFromBoolean(lv.boolean() == rv.boolean())
		case ttNotEqual:
			return ValueFromBoolean(lv.boolean() != rv.boolean())
		default:
			panic("not supported operator on two booleans")
		}
	}

	if lt == vtNumber && rt == vtNumber {
		switch op {
		case ttAddition:
			return ValueFromNumber(lv.number() + rv.number())
		case ttSubstraction:
			return ValueFromNumber(lv.number() - rv.number())
		case ttMultiply:
			return ValueFromNumber(lv.number() * rv.number())
		case ttDivision:
			if rv.number() == 0 {
				panic("divide by zero")
			}
			return ValueFromNumber(lv.number() / rv.number())
		case ttGreaterThan:
			return ValueFromBoolean(lv.number() > rv.number())
		case ttGreaterThanOrEqual:
			return ValueFromBoolean(lv.number() >= rv.number())
		case ttLessThan:
			return ValueFromBoolean(lv.number() < rv.number())
		case ttLessThanOrEqual:
			return ValueFromBoolean(lv.number() <= rv.number())
		case ttEqual:
			return ValueFromBoolean(lv.number() == rv.number())
		case ttNotEqual:
			return ValueFromBoolean(lv.number() != rv.number())
		default:
			panic("not supported operator on two numbers")
		}
	}

	if lt == vtString && rt == vtString {
		switch op {
		case ttAddition:
			return ValueFromString(lv.str() + rv.str())
		default:
			panic("not supported operator on two strings")
		}
	}

	panic("unknown binary operator and operands")
}

type Parameters struct {
	names []string
}

func NewParameters(names ...string) *Parameters {
	p := &Parameters{}
	for _, name := range names {
		p.names = append(p.names, name)
	}
	return p
}

func (p *Parameters) Len() int {
	return len(p.names)
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

// BindArguments assigns actual arguments.
// un-aligned parameters and arguments are set to nil.
func (p *Parameters) BindArguments(ctx *Context, args ...Value) {
	for index, name := range p.names {
		var arg Value
		if index < len(args) {
			arg = args[index]
		}
		ctx.AddValue(name, arg)
	}
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

func (f *FunctionExpression) Evaluate(ctx *Context) Value {
	value := ValueFromFunction(f)
	if f.name != "" {
		ctx.AddValue(f.name, value)
	}
	return value
}

// Execute executes function statements.
// This is not a statement interface implementation.
func (f *FunctionExpression) Execute(ctx *Context) Value {
	f.block.Execute(ctx)
	if ret, ok := f.block.Return(); ok {
		return ret
	} else {
		return ValueFromNil()
	}
}

func (f *FunctionExpression) BindArguments(ctx *Context, args ...Value) {
	f.params.BindArguments(ctx, args...)
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
	args := Values{}
	for _, expr := range a.exprs {
		args.values = append(args.values, expr.Evaluate(ctx))
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

func (i *IndexExpression) Evaluate(ctx *Context) Value {
	value := i.indexable.Evaluate(ctx).value
	keyer, ok1 := value.(KeyIndexer)
	elemer, ok2 := value.(ElemIndexer)
	if !ok1 && !ok2 {
		panic("value of expr is not indexable")
	}
	key := i.key.Evaluate(ctx)
	if key.Type == vtString && keyer != nil {
		return keyer.Key(key.str())
	}
	if key.Type == vtNumber && elemer != nil {
		return elemer.Elem(key.number())
	}
	panic("not indexable")
}

type CallExpression struct {
	Callable Expression
	Args     *Arguments
}

func (f *CallExpression) Evaluate(ctx *Context) Value {
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
		fn := callable.function()
		newCtx := NewContext(ctx)
		for i := 0; i < fn.params.Len() && i < f.Args.Len(); i++ {
			newCtx.AddValue(
				fn.params.GetAt(i),
				f.Args.exprs[i].Evaluate(ctx),
			)
		}
		return fn.Execute(newCtx)
	case vtBuiltin:
		newCtx := NewContext(ctx)
		args := f.Args.EvaluateAll(ctx)
		return callable.builtin().fn(newCtx, &args)
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

func (o *ObjectExpression) Evaluate(ctx *Context) Value {
	obj := NewObject()
	for k, v := range o.props {
		obj.props[k] = v.Evaluate(ctx)
	}
	return ValueFromObject(obj)
}

type ArrayExpression struct {
	elements []Expression
}

func NewArrayExpression() *ArrayExpression {
	return &ArrayExpression{}
}

func (a *ArrayExpression) Evaluate(ctx *Context) Value {
	arr := NewArray()
	for _, element := range a.elements {
		arr.PushElem(element.Evaluate(ctx))
	}
	return ValueFromObject(arr)
}
