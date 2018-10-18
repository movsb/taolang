#include "value.h"
#include "expression.h"

namespace taolang {

std::map<ValueType::Value, std::string> typeNames = {
    {ValueType::Nil, "nil"},
    {ValueType::Boolean, "boolean"},
    {ValueType::Number, "number"},
    {ValueType::String, "string"},
    {ValueType::Variable, "variable"},
    {ValueType::Object, "object"},
    {ValueType::Function, "function"},
    {ValueType::Builtin, "builtin"},
    {ValueType::Class, "class"},
};

Value* Builtin::Execute(Context* ctx, Values* args) {
    return (*_func)(_that, ctx, args);
}

Value* Value::fromFunction(FunctionExpression* func, Context* closure) {
    auto v = new Value();
    v->type = ValueType::Function;
    auto f = new EvaluatedFunctionExpression();
    f->_func = func;
    f->_closure = closure;
    v->func = f;
    return v;
}

Value* Value::Evaluate(Context* ctx) {

}

void Value::Assign(Context* ctx, Value* val) {

}

bool Value::truth(Context* ctx) {
    switch(type) {
    case ValueType::Nil:
        return false;
    case ValueType::Boolean:
        return boolean();
    case ValueType::Number:
        return number() != 0;
    case ValueType::String:
        return !str.empty();
    case ValueType::Function:
    case ValueType::Builtin:
        return true;
    case ValueType::Variable:
        return ctx->MustFind(var, true)->truth(ctx);
    default:
        break;
    }
    throw SyntaxError("unknown truth type");
}

}
