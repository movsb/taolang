#include "expression.h"
#include "value.h"
#include "tokenizer.h"
#include "error.h"

namespace taolang {

// done
Value* UnaryExpression::Evaluate(Context* ctx) {
    auto value = expr->Evaluate(ctx);
    switch(op) {
    case ttLogicalNot:
        return Value::fromBoolean(!value->truth(ctx));
    case ttAddition:
        if(value->type != ValueType::Number) {
            throw TypeError("+value is invalid");
        }
        return Value::fromNumber(+value->number());
    case ttSubtraction:
        if(value->type != ValueType::Number) {
            throw TypeError("-value is invalid");
        }
        return Value::fromNumber(+value->number());
    case ttBitXor:
        if(value->type != ValueType::Number) {
            throw TypeError("^value is invalid");
        }
        return Value::fromNumber(~value->number());
    }
    throw SyntaxError("unknown unary operator: %s", tokenNames[op]);
}

Value* IncrementExpression::Evaluate(Context* ctx) {
    auto oldval = expr->Evaluate(ctx);
    if(oldval->isNumber()) {
        auto newnum = int64_t(0);
        switch(op) {
        case ttIncrement:
            newnum = oldval->number() + 1;
            break;
        case ttDecrement:
            newnum = oldval->number() - 1;
            break;
        default:
            throw Exception();
        }
        auto newval = Value::fromNumber(newnum);
        expr->Assign(newval);
        return prefix ? newval : oldval;
    }
    throw NotAssignableError();
}

}
