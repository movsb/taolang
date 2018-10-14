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

Value* BinaryExpression::Evaluate(Context* ctx) {
    Value *lv, *rv;
    if(op != ttLogicalAnd && op != ttLogicalOr){
        lv = left->Evaluate(ctx);
        rv = right->Evaluate(ctx);
    }

    auto lt=lv->type, rt = rv->type;

    if(lt == ValueType::Nil && rt == ValueType::Nil) {
        switch(op) {
        case ttEqual:
            return Value::fromBoolean(true);
        case ttNotEqual:
            return Value::fromBoolean(false);
        }
    }

    if(lt == ValueType::Boolean && rt == ValueType::Boolean) {
        switch(op) {
        case ttEqual:
            return Value::fromBoolean(lv->boolean() == rv->boolean());
        case ttNotEqual:
            return Value::fromBoolean(lv->boolean() != rv->boolean());
        }
    }

    if(lt == ValueType::Number && rt == ValueType::Number){
        switch(op){
        case ttAddition:
            return Value::fromNumber(lv->number() + rv->number());
        case ttSubtraction:
            return Value::fromNumber(lv->number() - rv->number());
        case ttMultiply:
            return Value::fromNumber(lv->number() * rv->number());
        case ttDivision:
            if(rv->number() == 0) {
                throw TypeError("divide by zero");
            }
            return Value::fromNumber(lv->number() / rv->number());
        case ttGreaterThan:
            return Value::fromBoolean(lv->number() > rv->number());
        case ttGreaterThanOrEqual:
            return Value::fromBoolean(lv->number() >= rv->number());
        case ttLessThan:
            return Value::fromBoolean(lv->number() < rv->number());
        case ttLessThanOrEqual:
            return Value::fromBoolean(lv->number() <= rv->number());
        case ttEqual:
            return Value::fromBoolean(lv->number() == rv->number());
        case ttNotEqual:
            return Value::fromBoolean(lv->number() != rv->number());
        case ttPercent:
            return Value::fromNumber(lv->number() % rv->number());
        case ttStarStar:
            // TODO precision lost
            // val := math.Pow(float64(lv->number()), float64(rv->number()))
            //return Value::fromNumber(int(val))
        case ttLeftShift:
            return Value::fromNumber(lv->number() << uint(rv->number()));
        case ttRightShift:
            return Value::fromNumber(lv->number() >> uint(rv->number()));
        case ttBitAnd:
            return Value::fromNumber(lv->number() & rv->number());
        case ttBitOr:
            return Value::fromNumber(lv->number() | rv->number());
        case ttBitXor:
            return Value::fromNumber(lv->number() ^ rv->number());
        case ttBitAndNot:
            return Value::fromNumber(lv->number() &~ rv->number());
        }
    }

    if(lt == ValueType::String && rt == ValueType::String){
        switch(op) {
        case ttAddition:
            return Value::fromString(lv->str.p + rv->str.p);
        case ttEqual:
            return Value::fromBoolean(lv->str.p == rv->str.p);
        case ttNotEqual:
            return Value::fromBoolean(lv->str.p != rv->str.p);
        default:
            throw SyntaxError("not supported operator on two strings");
        }
    }

    if(op == ttLogicalAnd) {
        return Value::fromBoolean(
            left->Evaluate(ctx)->truth(ctx) &&
                right->Evaluate(ctx)->truth(ctx)
        );
    } else if(op == ttLogicalOr) {
        lv = left->Evaluate(ctx);
        if(lv->truth(ctx)) {
            return lv;
        }
        return right->Evaluate(ctx);
    }

    // TODO
    if(lt == ValueType::Builtin && rt == ValueType::Builtin) {
        p1 := reflect.ValueOf(lv->builtin().fn).Pointer()
        p2 := reflect.ValueOf(rv->builtin().fn).Pointer()
        switch op {
        case ttEqual:
            return Value::fromBoolean(p1 == p2)
        case ttNotEqual:
            return Value::fromBoolean(p1 != p2)
        default:
            throw SyntaxError("not supported operator on two builtins");
        }
    }

    throw SyntaxError("unknown binary operator and operands");
}

Value* TernaryExpression::Evaluate(Context* ctx) {
    return cond->Evaluate(ctx)->truth(ctx)
        ? left->Evaluate(ctx)
        : right->Evaluate(ctx)
        ;
}

}
