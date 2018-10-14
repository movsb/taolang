#pragma once

#include "context.h"
#include "value.h"

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

class BaseExpression {
public:
    BaseExpression(ExprType type)
        : type(type)
    {}
    ExprType type;
    virtual Value* Evaluate(Context* ctx) = 0;
    virtual void Assign(Value* value) {
        throw NotAssignableError();
    }
};

class UnaryExpression : public BaseExpression {
public:
    UnaryExpression()
        : BaseExpression(ExprType::Unary)
    {}
    TokenType op;
    BaseExpression* expr;
    virtual Value* Evaluate(Context* ctx) override;
};

class IncrementExpression : public BaseExpression {
public:
    IncrementExpression()
        : BaseExpression(ExprType::Increment)
    {}
    bool prefix;
    TokenType op;
    BaseExpression* expr;
    virtual Value* Evaluate(Context* ctx) override;
};

class BinaryExpression : public BaseExpression {
public:
    BinaryExpression()
        : BaseExpression(ExprType::Binary)
    {}
    virtual Value* Evaluate(Context* ctx) override;
};

class TernaryExpression : public BaseExpression {
public:
    TernaryExpression()
        : BaseExpression(ExprType::Ternary)
    {}
    virtual Value* Evaluate(Context* ctx) override;
};

class NewExpression : public BaseExpression {
public:
    NewExpression()
        : BaseExpression(ExprType::New)
    {}
    virtual Value* Evaluate(Context* ctx) override;
};

class AssignmentExpression : public BaseExpression {
public:
    AssignmentExpression()
        : BaseExpression(ExprType::Assignment)
    {}
    virtual Value* Evaluate(Context* ctx) override;
};

class FunctionExpression : public BaseExpression {
public:
    FunctionExpression()
        : BaseExpression(ExprType::Function)
    {}
    virtual Value* Evaluate(Context* ctx) override;
};

class EvaluatedFunctionExpression : public BaseExpression {
public:
    EvaluatedFunctionExpression()
        : BaseExpression(ExprType::EvaluatedFunction)
    {}
    virtual Value* Evaluate(Context* ctx) override;
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
    virtual Value* Evaluate(Context* ctx) override;
};

class ArrayExpression : public BaseExpression {
public:
    ArrayExpression()
        : BaseExpression(ExprType::Array)
    {}
    virtual Value* Evaluate(Context* ctx) override;
};

}
