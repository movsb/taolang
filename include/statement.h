#pragma once

#include <string>
#include <vector>

#include "expression.h"

namespace taolang {

class Context;

enum class StmtType {
    Empty,
    Let,
    Function,
    Return,
    Block,
    Expression,
    For,
    Break,
    If,
    Switch,
};

class BaseStatement {
public:
    BaseStatement(StmtType type)
        : Type(type)
    {}
    StmtType Type;
    virtual void Execute(Context* ctx) = 0;
};

class EmptyStatement : public BaseStatement {
public:
    EmptyStatement()
        : BaseStatement(StmtType::Empty)
    {}
    virtual void Execute(Context* ctx) override;
};

class LetStatement : public BaseStatement {
public:
    LetStatement()
        : BaseStatement(StmtType::Let)
    {}
    std::string _name;
    BaseExpression* _expr;
    virtual void Execute(Context* ctx) override;
};

class FunctionStatement : public BaseStatement {
public:
    FunctionStatement()
        : BaseStatement(StmtType::Function)
    {}
    BaseExpression* _expr;
    virtual void Execute(Context* ctx) override;
};

class ReturnStatement : public BaseStatement {
public:
    ReturnStatement()
        : BaseStatement(StmtType::Return)
    {}
    BaseExpression* _expr;
    virtual void Execute(Context* ctx) override;
};

class BlockStatement : public BaseStatement {
public:
    BlockStatement()
        : BaseStatement(StmtType::Block)
    {}
    std::vector<BaseStatement*> _stmts;
    virtual void Execute(Context* ctx) override;
};

class ExpressionStatement : public BaseStatement {
public:
    ExpressionStatement()
        : BaseStatement(StmtType::Expression)
    {}
    BaseExpression* _expr;
    virtual void Execute(Context* ctx) override;
};

class ForStatement : public BaseStatement {
public:
    ForStatement()
        : BaseStatement(StmtType::For)
    {}
    BaseStatement* _init;
    BaseExpression* _test;
    BaseExpression* _incr;
    BlockStatement* _block;
    virtual void Execute(Context* ctx) override;
};

class BreakStatement : public BaseStatement {
public:
    BreakStatement()
        : BaseStatement(StmtType::Break)
    {}
    virtual void Execute(Context* ctx) override;
};

class IfStatement : public BaseStatement {
public:
    IfStatement()
        : BaseStatement(StmtType::If)
    {}
    BaseExpression* _cond;
    BlockStatement* _ifBlock;
    BaseStatement* _elseBlock;
    virtual void Execute(Context* ctx) override;
};

class SwitchStatement : public BaseStatement {
public:
    SwitchStatement()
        : BaseStatement(StmtType::Switch)
    {}
    virtual void Execute(Context* ctx) override;
};

}
