#pragma once

#include <map>

#include "tokenizer.h"
#include "expression.h"
#include "statement.h"

namespace taolang {

enum class Precedence {
    _Unspecified,
	Assignment,
	Conditional,
	LogicalOr,
	LogicalAnd,

	BitwiseAnd,
	BitwiseOr = BitwiseAnd,
	BitwiseXor = BitwiseAnd,
	BitwiseAndNot = BitwiseAnd,

	Equality,
	Comparison,
	BitwiseShift,
	Addition,
	Multiplication,
	Exponentiation,

	LogicalNot,
	BitwiseNot = LogicalNot,
	UnaryPlus = LogicalNot,
	UnaryNegation = LogicalNot,

	PrefixIncrement,
	PrefixDecrement = PrefixIncrement,

	PostfixIncrement,
	PostfixDecrement = PostfixIncrement,

	Indexing,
	New = Indexing,
	Call = Indexing,
};

extern std::map<TokenType,Precedence> precedenceTable;

class Program;

class Parser {
public:
    Parser(Tokenizer* tkz)
        : _tkz(tkz)
        {}
public:
    Program* Parse();

protected:
    Token _expect(TokenType tt);
    template<typename... Args>
    Token _match(Args... args);
    Token _next();
    void _undo(Token tk);
    bool _skip(TokenType tt);
    Token _peek();
    bool _follow(TokenType tt);
    void _enter();
    void _leave(bool putBack);

protected:
    Precedence _getPrecedence(TokenType op);

protected:
    BaseStatement*      _parseStatement(bool global);
    LetStatement*       _parseLetStatement();
    FunctionStatement*  _parseFunctionStatement();
    ReturnStatement*    _parseReturnStatememt();
    BlockStatement*     _parseBlockStatement();
    ForStatement*       _parseForStatement();
    BreakStatement*     _parseBreakStatement();
    IfStatement*        _parseIfStatement();
    //SwitchStatement*    _parseSwitchStatement();

protected:
    IExpression*            _parseExpression(Precedence prec);
    IExpression*            _parsePrimaryExpression();
    TernaryExpression*      _parseTernaryExpression(IExpression* cond);
    AssignmentExpression*   _parseAssignmentExpression(IExpression* left, TokenType op);
    NewExpression*          _parseNewExpression();
    FunctionExpression*     _tryParseLambdaExpression(bool must);
    IndexExpression*        _parseIndexExpression(IExpression* left);
    CallExpression*         _parseCallExpression(IExpression* left);
    FunctionExpression*     _parseFunctionExpression();
    ObjectExpression*       _parseObjectExpression();
    ArrayExpression*        _parseArrayExpression();

protected:
    Tokenizer* _tkz;
    int _breakCount;
};

}
