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

}
