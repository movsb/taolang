package main

import (
	"math"
	"reflect"
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
	op   TokenType
	expr Expression
}

// NewUnaryExpression news a UnaryExpression.
func NewUnaryExpression(op TokenType, expr Expression) *UnaryExpression {
	return &UnaryExpression{
		op:   op,
		expr: expr,
	}
}

// Evaluate implements Expression.
func (u *UnaryExpression) Evaluate(ctx *Context) Value {
	value := u.expr.Evaluate(ctx)
	switch u.op {
	case ttAddition:
		if value.Type != vtNumber {
			panic("+value is invalid")
		}
		return ValueFromNumber(+value.number())
	case ttSubtraction:
		if value.Type != vtNumber {
			panic("-value is invalid")
		}
		return ValueFromNumber(-value.number())
	case ttNot:
		return ValueFromBoolean(!value.Truth(ctx))
	case ttBitXor:
		if value.Type != vtNumber {
			panic("^value is invalid")
		}
		return ValueFromNumber(^value.number())
	}
	panicf("unknown unary operator: %v", u.op)
	return ValueFromNil()
}

// IncrementDecrementExpression is an a++ / a-- / ++a / --a expressions.
type IncrementDecrementExpression struct {
	prefix bool
	op     TokenType
	expr   Expression
}

// NewIncrementDecrementExpression news an IncrementDecrementExpression.
func NewIncrementDecrementExpression(op TokenType, prefix bool, expr Expression) *IncrementDecrementExpression {
	return &IncrementDecrementExpression{
		prefix: prefix,
		op:     op,
		expr:   expr,
	}
}

// Evaluate implements Expression.
func (i *IncrementDecrementExpression) Evaluate(ctx *Context) Value {
	oldval := i.expr.Evaluate(ctx)
	if oldval.isNumber() {
		assigner, ok := i.expr.(Assigner)
		if !ok {
			panic(NewNotAssignableError(oldval))
		}
		newnum := 0
		switch i.op {
		case ttIncrement:
			newnum = oldval.number() + 1
		case ttDecrement:
			newnum = oldval.number() - 1
		default:
			panic("won't go here")
		}
		newval := ValueFromNumber(newnum)
		assigner.Assign(ctx, newval)
		if i.prefix {
			return newval
		}
		return oldval
	}
	panic(NewNotAssignableError(oldval))
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

	// Logical values are evaluated "short-circuit"-ly.
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
		case ttSubtraction:
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
		case ttBitAndNot:
			return ValueFromNumber(lv.number() &^ rv.number())
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

	if lt == vtBuiltin && rt == vtBuiltin {
		p1 := reflect.ValueOf(lv.builtin().fn).Pointer()
		p2 := reflect.ValueOf(rv.builtin().fn).Pointer()
		switch op {
		case ttEqual:
			return ValueFromBoolean(p1 == p2)
		case ttNotEqual:
			return ValueFromBoolean(p1 != p2)
		default:
			panic("not supported operator on two builtins")
		}
	}

	panic("unknown binary operator and operands")
}

// TernaryExpression is the conditional(`?:`) expression.
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
		ctx.AddSymbol(name, arg)
	}
}

// EvaluatedFunctionExpression is the result of a FunctionExpression.
// The result is the closure and the expr itself.
//  Evaluate(FunctionExpression) -> EvaluatedFunctionExpression
//  Execute(EvaluatedFunctionExpression) -> Execute(FunctionExpression, this)
type EvaluatedFunctionExpression struct {
	this *Context // this is the scope where the function expression is defined
	fn   *FunctionExpression
}

// Execute evaluates the function expression within closure.
// This is not a statement interface implementation.
func (e *EvaluatedFunctionExpression) Execute(ctx *Context) Value {
	ctx.SetParent(e.this) // this is how closure works
	return e.fn.Execute(ctx)
}

// BindArguments binds actual arguments from call expression.
func (e *EvaluatedFunctionExpression) BindArguments(ctx *Context, args ...Value) {
	e.fn.params.BindArguments(ctx, args...)
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
	// a lambda function or an anonymous function doesn't have a name.
	if f.name != "" {
		ctx.AddSymbol(f.name, value)
	}
	return value
}

// Execute executes function statements.
// This is not a statement interface implementation.
func (f *FunctionExpression) Execute(ctx *Context) Value {
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
	key := i.key.Evaluate(ctx)
	indexable := i.indexable.Evaluate(ctx)

	// both obj.key, obj[0] are correct.
	// so, we need to query both interfaces.
	keyable, _ := indexable.value.(KeyIndexer)
	elemable, _ := indexable.value.(ElemIndexer)

	// convert from primitives to object when needed
	if keyable == nil && elemable == nil {
		switch indexable.Type {
		case vtString:
			keyable = KeyIndexer(NewString(indexable.str()))
		}
	}

	// get property
	if key.Type == vtString && keyable != nil {
		return keyable.Key(key.str())
	}

	// get element
	if key.Type == vtNumber && elemable != nil {
		return elemable.Elem(key.number())
	}

	if keyable == nil && elemable == nil {
		panic(NewNotIndexableError(indexable))
	}

	if keyable != nil && key.Type != vtString {
		panic(NewKeyTypeError(key))
	}
	if elemable != nil && key.Type != vtNumber {
		panic(NewKeyTypeError(key))
	}

	panic("won't go here")
}

// Assign implements Assigner.
func (i *IndexExpression) Assign(ctx *Context, val Value) {
	value := i.indexable.Evaluate(ctx)
	keyable, ok1 := value.value.(KeyAssigner)
	elemable, ok2 := value.value.(ElemAssigner)
	if !ok1 && !ok2 {
		panic(NewNotIndexableError(value))
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
	panic(NewKeyTypeError(key))
}

// CallExpression wraps a call.
type CallExpression struct {
	Callable Expression
	Args     *Arguments
}

// NewCallExpression news a call expression.
func NewCallExpression(callable Expression, args *Arguments) *CallExpression {
	c := &CallExpression{
		Callable: callable,
		Args:     args,
	}
	if c.Args == nil {
		c.Args = &Arguments{}
	}
	return c
}

// CallFunc calls user function.
func CallFunc(ctx *Context, callable Expression, args ...Expression) Value {
	c := NewCallExpression(callable, &Arguments{args})
	return c.Evaluate(ctx)
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
		fn := callable.function()
		newCtx := NewContext(fn.fn.name, nil)
		args := f.Args.EvaluateAll(ctx)
		fn.BindArguments(newCtx, args.values...)
		return fn.Execute(newCtx)
	case vtBuiltin:
		fn := callable.builtin()
		newCtx := NewContext(fn.name, nil)
		args := f.Args.EvaluateAll(ctx)
		return fn.Execute(newCtx, &args)
	default:
		panic(NewNotCallableError(callable))
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
