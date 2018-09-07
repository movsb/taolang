package main

import (
	"math"
)

// Expression is the interface that is implemented by all expressions.
type Expression interface {
	Evaluate(ctx *Context) Value
}

// Assigner is implemented by those who can be assigned.
type Assigner interface {
	Assign(ctx *Context, value Value)
}

// UnaryExpression is a unary expression.
type UnaryExpression struct {
	tt   TokenType
	expr Expression
}

// NewUnaryExpression new a UnaryExpression.
func NewUnaryExpression(tt TokenType, expr Expression) *UnaryExpression {
	return &UnaryExpression{
		tt:   tt,
		expr: expr,
	}
}

// Evaluate implements
func (u *UnaryExpression) Evaluate(ctx *Context) Value {
	value := u.expr.Evaluate(ctx)
	switch u.tt {
	case ttAddition:
		if value.Type != vtNumber {
			panic("+value is invalid")
		}
		return ValueFromNumber(+value.number())
	case ttSubstraction:
		if value.Type != vtNumber {
			panic("-value is invalid")
		}
		return ValueFromNumber(-value.number())
	case ttNot:
		return ValueFromBoolean(!value.Truth(ctx))
	}
	panicf("unknown unary operator: %v", u.tt) // TODO
	return ValueFromNil()
}

// IncrementDecrementExpression is a++ / a-- / ++a / --a expressions.
type IncrementDecrementExpression struct {
	prefix bool
	op     Token
	expr   Expression
}

// NewIncrementDecrementExpression new an IncrementDecrementExpression.
func NewIncrementDecrementExpression(op Token, prefix bool, expr Expression) *IncrementDecrementExpression {
	return &IncrementDecrementExpression{
		prefix: prefix,
		op:     op,
		expr:   expr,
	}
}

// Evaluate implements
func (i *IncrementDecrementExpression) Evaluate(ctx *Context) Value {
	oldval := i.expr.Evaluate(ctx)
	if oldval.isNumber() {
		assigner, ok := i.expr.(Assigner)
		if !ok {
			panicf("not assignable: %v (type: %s)", oldval, oldval.TypeName())
		}
		newval := Value{}
		switch i.op.typ {
		case ttIncrement:
			newval = ValueFromNumber(oldval.number() + 1)
			assigner.Assign(ctx, newval)
		case ttDecrement:
			newval = ValueFromNumber(oldval.number() - 1)
			assigner.Assign(ctx, newval)
		default:
			panic("bad op")
		}
		if i.prefix {
			return newval
		}
		return oldval
	}
	panicf("not assignable: %v (type: %s)", oldval, oldval.TypeName())
	return Value{}
}

// BinaryExpression is a binary expression.
type BinaryExpression struct {
	left  Expression
	op    TokenType
	right Expression
}

// NewBinaryExpression news a BinaryExpression.
func NewBinaryExpression(left Expression, op TokenType, right Expression) *BinaryExpression {
	return &BinaryExpression{
		left:  left,
		op:    op,
		right: right,
	}
}

// Evaluate implements Expression.
func (b *BinaryExpression) Evaluate(ctx *Context) Value {
	op := b.op
	lv, rv := Value{}, Value{}
	// Logical values are evaluated shortcutted
	if op != ttAndAnd && op != ttOrOr {
		lv = b.left.Evaluate(ctx)
		rv = b.right.Evaluate(ctx)
	}
	lt, rt := lv.Type, rv.Type

	if lt == vtNil && rt == vtNil {
		if op == ttEqual {
			return ValueFromBoolean(true)
		} else if op == ttNotEqual {
			return ValueFromBoolean(false)
		}
	}

	if lt == vtBoolean && rt == vtBoolean {
		switch op {
		case ttEqual:
			return ValueFromBoolean(lv.boolean() == rv.boolean())
		case ttNotEqual:
			return ValueFromBoolean(lv.boolean() != rv.boolean())
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
		case ttPercent:
			return ValueFromNumber(lv.number() % rv.number())
		case ttStarStar:
			// TODO precision lost
			val := math.Pow(float64(lv.number()), float64(rv.number()))
			return ValueFromNumber(int(val))
		case ttLeftShift:
			return ValueFromNumber(lv.number() << uint(rv.number()))
		case ttRightShift:
			return ValueFromNumber(lv.number() >> uint(rv.number()))
		case ttBitAnd:
			return ValueFromNumber(lv.number() & rv.number())
		case ttBitOr:
			return ValueFromNumber(lv.number() | rv.number())
		case ttBitXor:
			return ValueFromNumber(lv.number() ^ rv.number())
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

	if op == ttAndAnd {
		return ValueFromBoolean(
			b.left.Evaluate(ctx).Truth(ctx) &&
				b.right.Evaluate(ctx).Truth(ctx),
		)
	} else if op == ttOrOr {
		lv = b.left.Evaluate(ctx)
		if lv.Truth(ctx) {
			return lv
		}
		return b.right.Evaluate(ctx)
	}

	panic("unknown binary operator and operands")
}

// TernaryExpression is `?:` expression.
type TernaryExpression struct {
	cond  Expression
	left  Expression
	right Expression
}

// NewTernaryExpression news a ternary expression.
func NewTernaryExpression(cond, left, right Expression) *TernaryExpression {
	return &TernaryExpression{
		cond:  cond,
		left:  left,
		right: right,
	}
}

// Evaluate implements Expression.
func (t *TernaryExpression) Evaluate(ctx *Context) Value {
	if t.cond.Evaluate(ctx).Truth(ctx) {
		return t.left.Evaluate(ctx)
	}
	return t.right.Evaluate(ctx)
}

// Parameters is a collection of function parameters.
type Parameters struct {
	names []string
}

// NewParameters news
func NewParameters(names ...string) *Parameters {
	p := &Parameters{}
	for _, name := range names {
		p.names = append(p.names, name)
	}
	return p
}

// Len returns the count of parameters.
func (p *Parameters) Len() int {
	return len(p.names)
}

// GetAt gets n-th parameter.
func (p *Parameters) GetAt(index int) string {
	if index > len(p.names)-1 {
		panic("parameter index out of range")
	}
	return p.names[index]
}

// PutParam adds a parameter.
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

// EvaluatedFunctionExpression is the result of a FunctionExpression.
// The result is the closure and the expr itself.
//  Evaluate(FunctionExpression) -> EvaluatedFunctionExpression
//  Execute(EvaluatedFunctionExpression) -> Execute(FunctionExpression, this)
type EvaluatedFunctionExpression struct {
	this *Context // this is the scope where the function expression is defined
	expr *FunctionExpression
}

// Execute evaluates the function expression within closure.
// This is not a statement interface implementation.
func (e *EvaluatedFunctionExpression) Execute(ctx *Context) Value {
	return e.expr.Execute(e.this, ctx)
}

// BindArguments binds actual arguments from call expression.
func (e *EvaluatedFunctionExpression) BindArguments(ctx *Context, args ...Value) {
	e.expr.params.BindArguments(ctx, args...)
}

// FunctionExpression is
type FunctionExpression struct {
	name   string
	params *Parameters
	block  *BlockStatement
}

// Evaluate is
func (f *FunctionExpression) Evaluate(ctx *Context) Value {
	value := ValueFromFunction(f, ctx)
	if f.name != "" {
		ctx.AddValue(f.name, value)
	}
	return value
}

// Execute executes function statements.
// This is not a statement interface implementation.
func (f *FunctionExpression) Execute(this *Context, ctx *Context) Value {
	ctx.SetParent(this) // this is how closure works
	f.block.Execute(ctx)
	if ctx.hasret {
		return ctx.retval
	}
	return ValueFromNil()
}

// Arguments is the collection of arguments for function call.
type Arguments struct {
	exprs []Expression
}

// Len returns the length of arguments.
func (a *Arguments) Len() int {
	return len(a.exprs)
}

// PutArgument adds an argument.
func (a *Arguments) PutArgument(expr Expression) {
	a.exprs = append(a.exprs, expr)
}

// EvaluateAll evaluates all values of arguments.
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

// Evaluate implements Expression.
func (i *IndexExpression) Evaluate(ctx *Context) Value {
	value := i.indexable.Evaluate(ctx)
	keyable, ok1 := value.value.(KeyIndexer)
	elemable, ok2 := value.value.(ElemIndexer)
	if !ok1 && !ok2 {
		panicf("not indexable: %v (type: %s)", value, value.TypeName())
	}
	key := i.key.Evaluate(ctx)
	if key.Type == vtString && keyable != nil {
		return keyable.Key(key.str())
	}
	if key.Type == vtNumber && elemable != nil {
		return elemable.Elem(key.number())
	}
	panic("not indexable")
}

// Assign implements Assigner.
func (i *IndexExpression) Assign(ctx *Context, val Value) {
	value := i.indexable.Evaluate(ctx)
	keyable, ok1 := value.value.(KeyAssigner)
	elemable, ok2 := value.value.(ElemAssigner)
	if !ok1 && !ok2 {
		panicf("not assignable: %v (type: %s)", value, value.TypeName())
	}
	key := i.key.Evaluate(ctx)
	if key.isString() && keyable != nil {
		keyable.KeyAssign(key.str(), val)
		return
	}
	if key.isNumber() && elemable != nil {
		elemable.ElemAssign(key.number(), val)
		return
	}
	panic("not assignable")
}

// CallExpression wrap a method call.
type CallExpression struct {
	Callable Expression
	Args     *Arguments
}

// Evaluate implements Expression.
// It calls the callable.
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
	default:
		panicf("not callable: %v (type: %s)", callable, callable.TypeName())
	}

	switch callable.Type {
	case vtFunction:
		fn := callable.function()
		newCtx := NewContext(fn.expr.name, nil)
		args := f.Args.EvaluateAll(ctx)
		fn.BindArguments(newCtx, args.values...)
		return fn.Execute(newCtx)
	case vtBuiltin:
		fn := callable.builtin()
		newCtx := NewContext(fn.name, nil)
		args := f.Args.EvaluateAll(ctx)
		return fn.fn(newCtx, &args)
	default:
		panic("bad call")
	}
}

// ObjectExpression is the object literal expression.
type ObjectExpression struct {
	props map[string]Expression
}

// NewObjectExpression news an object literal expression.
func NewObjectExpression() *ObjectExpression {
	return &ObjectExpression{
		props: make(map[string]Expression),
	}
}

// Evaluate implements Expression.
func (o *ObjectExpression) Evaluate(ctx *Context) Value {
	obj := NewObject()
	for k, v := range o.props {
		obj.props[k] = v.Evaluate(ctx)
	}
	return ValueFromObject(obj)
}

// ArrayExpression is the array literal expression.
type ArrayExpression struct {
	elements []Expression
}

// NewArrayExpression news an array literal expression.
func NewArrayExpression() *ArrayExpression {
	return &ArrayExpression{}
}

// Evaluate implements Expression.
func (a *ArrayExpression) Evaluate(ctx *Context) Value {
	arr := NewArray()
	for _, element := range a.elements {
		arr.PushElem(element.Evaluate(ctx))
	}
	return ValueFromObject(arr)
}
