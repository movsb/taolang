#include "statement.h"

namespace taolang {

void EmptyStatement::Execute(Context* ctx) {
    // nop
}

void LetStatement::Execute(Context* ctx) {
    auto init = Value::fromNil();
    if(_expr != nullptr) {
        init = _expr->Evaluate(ctx);
        // TODO check builtin copy
    }
    ctx->AddSymbol(_name, init);
}

void FunctionStatement::Execute(Context* ctx) {
    _expr->Evaluate(ctx);
}

void ReturnStatement::Execute(Context* ctx) {
    auto value = _expr == nullptr 
        ? Value::fromNil()
        : _expr->Evaluate(ctx)
        ;
}

void BlockStatement::Execute(Context* ctx) {
    for(auto& stmt : _stmts) {
        Context* newCtx;
        switch(Type) {
        case StmtType::Block:
            newCtx = new Context(ctx);
            stmt->Execute(newCtx);
            break;
        default:
            newCtx = ctx;
            stmt->Execute(newCtx);
        }
        if(newCtx->_broke) {
            ctx->_broke = true;
            break;
        }
        if(newCtx->_hasRet) {
            ctx->SetReturn(newCtx->_retVal);
            return;
        }
    }
}

void ExpressionStatement::Execute(Context* ctx) {
    (void)_expr->Evaluate(ctx);
}

void ForStatement::Execute(Context* ctx) {
    if(_init != nullptr) {
        _init->Execute(ctx);
    }
    while((1)) {
        if(_test != nullptr) {
            if(!_test->Evaluate(ctx)->truth(ctx)) {
                break;
            }
        }
        auto newCtx = new Context(ctx);
        _block->Execute(newCtx);
        if(newCtx->_hasRet) {
            ctx->SetReturn(newCtx->_retVal);
            return;
        }
        if(newCtx->_broke) {
            ctx->SetBreak();
            break;
        }
        if(_incr != nullptr) {
            _incr->Evaluate(ctx);
        }
    }
}

void BreakStatement::Execute(Context* ctx) {
    ctx->SetBreak();
}

void IfStatement::Execute(Context* ctx) {
    auto cond = _cond->Evaluate(ctx);
    Context* newCtx;
    if(cond->truth(ctx)) {
        newCtx = new Context(ctx);
        _ifBlock->Execute(newCtx);
    } else {
        newCtx = new Context(ctx);
        _elseBlock->Execute(newCtx);
    }
    if(newCtx->_broke) {
        ctx->SetBreak();
        return;
    }
    if(newCtx->_hasRet) {
        ctx->SetReturn(newCtx->_retVal);
        return;
    }
}

void SwitchStatement::Execute(Context* ctx) {
    // TODO
}

}
