#pragma once

#include <unordered_map>

#include "error.h"
#include "context.h"
#include "value.h"
#include "tokenizer.h"

namespace taolang {

class BlockStatement;

enum class ExprType {
    Unary,
    Increment,
    Binary,
    Ternary,
    New,
    Assignment,
    Function,
    EvaluatedFunction,
    Index,
    Call,
    Object,
    Array,
};

class Arguments {
public:
    int Size() {
        return (int)_args.size();
    }
    IExpression* Get(int i) {
        return _args[i];
    }
    void Put(IExpression* arg) {
        _args.emplace_back(arg);
    }
    Values* EvaluateAll(Context* ctx) {
        auto values = new Values();
        for(auto& arg : _args) {
            values->Push(arg->Evaluate(ctx));
        }
        return values;
    }
private:
    std::vector<IExpression*> _args;
};

class Parameters {
public:
    int Size() {
        return (int)_params.size();
    }
    void Put(const std::string& name) {
        _params.emplace_back(name);
    }
    void BindArguments(Context* ctx, Values* args) {
        for(size_t i = 0; i < _params.size(); i++) {
            if(i < args->Size()) {
                ctx->AddSymbol(_params[i], args->Get(i));
            } else {
                ctx->AddSymbol(_params[i], Value::fromNil());
            }
        }
    }

private:
    std::vector<std::string> _params;
};

class BaseExpression : public IExpression {
public:
    BaseExpression(ExprType type)
        : type(type)
    {}
    ExprType type;
    virtual void Assign(Context* ctx, Value* value) override {
        throw NotAssignableError(
            "%s is not assignable",
            value->ToString().c_str()
        );
    }
};

class UnaryExpression : public BaseExpression {
public:
    UnaryExpression()
        : BaseExpression(ExprType::Unary)
    {}
    UnaryExpression(TokenType op, IExpression* expr)
        : BaseExpression(ExprType::Unary)
        , _op(op)
        , _expr(expr)
    {}
    TokenType _op;
    IExpression* _expr;
    virtual Value* Evaluate(Context* ctx) override;
};

class IncrementExpression : public BaseExpression {
public:
    IncrementExpression()
        : BaseExpression(ExprType::Increment)
    {}
    IncrementExpression(TokenType op, bool prefix, IExpression* expr)
        : BaseExpression(ExprType::Increment)
        , _op(op)
        , _prefix(prefix)
        , _expr(expr)
    {}
    bool _prefix;
    TokenType _op;
    IExpression* _expr;
    virtual Value* Evaluate(Context* ctx) override;
};

class BinaryExpression : public BaseExpression {
public:
    BinaryExpression()
        : BaseExpression(ExprType::Binary)
    {}
    BinaryExpression(IExpression* left, TokenType op, IExpression* right)
        : BaseExpression(ExprType::Binary)
        , _left(left)
        , _op(op)
        , _right(right)
    {}
    IExpression* _left;
    TokenType _op;
    IExpression* _right;
    virtual Value* Evaluate(Context* ctx) override;
};

class TernaryExpression : public BaseExpression {
public:
    TernaryExpression()
        : BaseExpression(ExprType::Ternary)
    {}
    IExpression* cond;
    IExpression* left;
    IExpression* right;
    virtual Value* Evaluate(Context* ctx) override;
};

class NewExpression : public BaseExpression {
public:
    NewExpression()
        : BaseExpression(ExprType::New)
    {}
    std::string _name;
    Arguments _args;
    virtual Value* Evaluate(Context* ctx) override;
};

class AssignmentExpression : public BaseExpression {
public:
    AssignmentExpression()
        : BaseExpression(ExprType::Assignment)
    {}
    IExpression* _left;
    IExpression* _expr;
    virtual Value* Evaluate(Context* ctx) override;
};

class FunctionExpression : public BaseExpression, public ICallable {
public:
    FunctionExpression()
        : BaseExpression(ExprType::Function)
    {}
    std::string _name;
    Parameters _params;
    BlockStatement* _body;
    virtual Value* Evaluate(Context* ctx) override;
    virtual Value* Execute(Context* ctx, Values* args) override;
};

class EvaluatedFunctionExpression : public ICallable {
public:
    EvaluatedFunctionExpression()
    {}
    Context* _closure;
    FunctionExpression* _func;
    virtual Value* Execute(Context* ctx, Values* args) override;
};

class IndexExpression : public BaseExpression {
public:
    IndexExpression()
        : BaseExpression(ExprType::Index)
    {}
    IExpression* _indexable;
    IExpression* _key;
    virtual Value* Evaluate(Context* ctx) override;
    virtual void Assign(Context* ctx, Value* value) override;
};

class CallExpression : public BaseExpression {
public:
    CallExpression()
        : BaseExpression(ExprType::Call)
    {}
    IExpression* _callable;
    Arguments    _args;
    virtual Value* Evaluate(Context* ctx) override;
};

Value* CallFunc(Context* ctx, IExpression* callable, Arguments* args);

class ObjectExpression : public BaseExpression {
public:
    ObjectExpression()
        : BaseExpression(ExprType::Object)
    {}
    std::unordered_map<std::string, IExpression*> _props;
    virtual Value* Evaluate(Context* ctx) override;
};

class ArrayExpression : public BaseExpression {
public:
    ArrayExpression()
        : BaseExpression(ExprType::Array)
    {}
    Arguments _elems;
    virtual Value* Evaluate(Context* ctx) override;
};

}
