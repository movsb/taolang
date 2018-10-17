#include "expression.h"
#include "value.h"
#include "tokenizer.h"
#include "error.h"

namespace taolang {

// done
Value* UnaryExpression::Evaluate(Context* ctx) {
    auto value = _expr->Evaluate(ctx);
    switch(_op) {
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
    default:
        break;
    }
    throw SyntaxError("unknown unary operator: %s", tokenNames[_op]);
}

Value* IncrementExpression::Evaluate(Context* ctx) {
    auto oldval = _expr->Evaluate(ctx);
    if(oldval->isNumber()) {
        auto newnum = int64_t(0);
        switch(_op) {
        case ttIncrement:
            newnum = oldval->number() + 1;
            break;
        case ttDecrement:
            newnum = oldval->number() - 1;
            break;
        default:
            throw Error();
        }
        auto newval = Value::fromNumber(newnum);
        _expr->Assign(ctx, newval);
        return _prefix ? newval : oldval;
    }
    throw NotAssignableError();
}

Value* BinaryExpression::Evaluate(Context* ctx) {
    Value *lv, *rv;
    if(_op != ttLogicalAnd && _op != ttLogicalOr){
        lv = _left->Evaluate(ctx);
        rv = _right->Evaluate(ctx);
    }

    auto lt=lv->type, rt = rv->type;

    if(lt == ValueType::Nil && rt == ValueType::Nil) {
        switch(_op) {
        case ttEqual:
            return Value::fromBoolean(true);
        case ttNotEqual:
            return Value::fromBoolean(false);
        default:
            break;
        }
    }

    if(lt == ValueType::Boolean && rt == ValueType::Boolean) {
        switch(_op) {
        case ttEqual:
            return Value::fromBoolean(lv->boolean() == rv->boolean());
        case ttNotEqual:
            return Value::fromBoolean(lv->boolean() != rv->boolean());
        default:
            break;
        }
    }

    if(lt == ValueType::Number && rt == ValueType::Number){
        switch(_op){
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
        default:
            break;
        }
    }

    if(lt == ValueType::String && rt == ValueType::String){
        switch(_op) {
        case ttAddition:
            return Value::fromString(lv->str + rv->str);
        case ttEqual:
            return Value::fromBoolean(lv->str == rv->str);
        case ttNotEqual:
            return Value::fromBoolean(lv->str != rv->str);
        default:
            throw SyntaxError("not supported operator on two strings");
        }
    }

    if(_op == ttLogicalAnd) {
        return Value::fromBoolean(
            _left->Evaluate(ctx)->truth(ctx) &&
                _right->Evaluate(ctx)->truth(ctx)
        );
    } else if(_op == ttLogicalOr) {
        lv = _left->Evaluate(ctx);
        if(lv->truth(ctx)) {
            return lv;
        }
        return _right->Evaluate(ctx);
    }

    // TODO
    if(lt == ValueType::Builtin && rt == ValueType::Builtin) {
        /*
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
        */
    }

    throw SyntaxError("unknown binary operator and operands");
}

Value* TernaryExpression::Evaluate(Context* ctx) {
    return cond->Evaluate(ctx)->truth(ctx)
        ? left->Evaluate(ctx)
        : right->Evaluate(ctx)
        ;
}

Value* NewExpression::Evaluate(Context* ctx) {
    // TODO
}

Value* AssignmentExpression::Evaluate(Context* ctx) {
    auto val = _expr->Evaluate(ctx);
    _left->Assign(ctx, val);
    return val;
}

Value* EvaluatedFunctionExpression::Execute(Context* ctx, Values* args) {

}

Value* FunctionExpression::Evaluate(Context* ctx) {
    auto val = Value::fromFunction(this, ctx);
    if(!_name.empty()) {
        ctx->AddSymbol(_name, val);
    }
    return val;
}

Value* FunctionExpression::Execute(Context* ctx, Values* args) {
    // TODO
}


}
