#pragma once

#include <unordered_map>

#include "error.h"
#include "context.h"
#include "value.h"
#include "tokenizer.h"

namespace taolang {

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
    void Put(BaseExpression* arg) {
        _args.emplace_back(arg);
    }
    Values* EvaluateAll(Context* ctx) {

    }
private:
    std::vector<BaseExpression*> _args;
};

class BaseExpression : public IExpression {
public:
    BaseExpression(ExprType type)
        : type(type)
    {}
    ExprType type;
    virtual void Assign(Context* ctx, Value* value) override {
        throw NotAssignableError();
    }
};

class UnaryExpression : public BaseExpression {
public:
    UnaryExpression()
        : BaseExpression(ExprType::Unary)
    {}
    UnaryExpression(TokenType op, BaseExpression* expr)
        : BaseExpression(ExprType::Unary)
        , _op(op)
        , _expr(expr)
    {}
    TokenType _op;
    BaseExpression* _expr;
    virtual Value* Evaluate(Context* ctx) override;
};

class IncrementExpression : public BaseExpression {
public:
    IncrementExpression()
        : BaseExpression(ExprType::Increment)
    {}
    IncrementExpression(TokenType op, bool prefix, BaseExpression* expr)
        : BaseExpression(ExprType::Increment)
        , _op(op)
        , _prefix(prefix)
        , _expr(expr)
    {}
    bool _prefix;
    TokenType _op;
    BaseExpression* _expr;
    virtual Value* Evaluate(Context* ctx) override;
};

class BinaryExpression : public BaseExpression {
public:
    BinaryExpression()
        : BaseExpression(ExprType::Binary)
    {}
    BinaryExpression(BaseExpression* left, TokenType op, BaseExpression* right)
        : BaseExpression(ExprType::Binary)
        , _left(left)
        , _op(op)
        , _right(right)
    {}
    BaseExpression* _left;
    TokenType _op;
    BaseExpression* _right;
    virtual Value* Evaluate(Context* ctx) override;
};

class TernaryExpression : public BaseExpression {
public:
    TernaryExpression()
        : BaseExpression(ExprType::Ternary)
    {}
    BaseExpression* cond;
    BaseExpression* left;
    BaseExpression* right;
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
    BaseExpression* _left;
    BaseExpression* _expr;
    virtual Value* Evaluate(Context* ctx) override;
};

class FunctionExpression : public BaseExpression, public ICallable {
public:
    FunctionExpression()
        : BaseExpression(ExprType::Function)
    {}
    std::string _name;
    virtual Value* Evaluate(Context* ctx) override;
    virtual Value* Execute(Context* ctx, Values* args) override;
};

class EvaluatedFunctionExpression : public BaseExpression, public ICallable {
public:
    EvaluatedFunctionExpression()
        : BaseExpression(ExprType::EvaluatedFunction)
    {}
    Context* _closure;
    FunctionExpression* _func;
    virtual Value* Evaluate(Context* ctx) override;
    virtual Value* Execute(Context* ctx, Values* args) override;
};

class IndexExpression : public BaseExpression {
public:
    IndexExpression()
        : BaseExpression(ExprType::Index)
    {}
    virtual Value* Evaluate(Context* ctx) override;
};

class CallExpression : public BaseExpression {
public:
    CallExpression()
        : BaseExpression(ExprType::Call)
    {}
    virtual Value* Evaluate(Context* ctx) override;
};

class ObjectExpression : public BaseExpression {
public:
    ObjectExpression()
        : BaseExpression(ExprType::Object)
    {}
    std::unordered_map<std::string, BaseExpression*> _props;
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
