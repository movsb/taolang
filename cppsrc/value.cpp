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
    switch(type) {
    case ValueType::Variable:
        return ctx->MustFind(variable(), true);
    case ValueType::Class:
        throw SyntaxError("%s is a constructor", "TODO");
    default:
        return this;
    }
}

void Value::Assign(Context* ctx, Value* val) {
    throw NotAssignableError();
}

ICallable* Value::callable() {
    switch(type) {
    case ValueType::Function:
        return func;
    case ValueType::Builtin:
        return bi;
    default:
        throw NotCallableError();
    }
}

std::string Value::ToString() {
    switch(type) {
    case ValueType::Nil:
        return "nil";
    case ValueType::Boolean:
        return boolean() ? "true" : "false";
    case ValueType::Number:
        return std::to_string(number());
    case ValueType::String:
        return string();
    case ValueType::Function: {
            auto f = function()->_func->_name;
            return "function(" + (!f.empty() ? f : "\"\"") + ")";
        }
    case ValueType::Builtin: {
            auto o = builtin()->_that->TypeName();
            auto p = builtin()->_name;
            return "builtin(" + o + "." + p + ")";
        }
    case ValueType::Object:
        return static_cast<IString*>(object())->ToString();
    case ValueType::Array:
        return static_cast<IString*>(array())->ToString();
    case ValueType::Variable:
        return variable();
    case ValueType::Class:
        break;
    }

    throw Error("unknown value type to stringify");
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
        return ctx->MustFind(str, true)->truth(ctx);
    default:
        break;
    }
    throw SyntaxError("unknown truth type");
}

}
