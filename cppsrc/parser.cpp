#include "parser.h"
#include "error.h"
#include "expression.h"
#include "program.h"

namespace taolang {

std::map<TokenType,Precedence> precedenceTable = {
	{ttQuestion,           Precedence::Conditional},
	{ttLogicalNot,         Precedence::LogicalNot},
	{ttLogicalOr,          Precedence::LogicalOr},
	{ttLogicalAnd,         Precedence::LogicalAnd},
	{ttBitAnd,             Precedence::BitwiseAnd},
	{ttBitOr,              Precedence::BitwiseOr},
	{ttBitXor,             Precedence::BitwiseXor},
	{ttBitAndNot,          Precedence::BitwiseAndNot},
	{ttEqual,              Precedence::Equality},
	{ttNotEqual,           Precedence::Equality},
	{ttGreaterThan,        Precedence::Comparison},
	{ttGreaterThanOrEqual, Precedence::Comparison},
	{ttLessThan,           Precedence::Comparison},
	{ttLessThanOrEqual,    Precedence::Comparison},
	{ttLeftShift,          Precedence::BitwiseShift},
	{ttRightShift,         Precedence::BitwiseShift},
	{ttAddition,           Precedence::Addition},
	{ttSubtraction,        Precedence::Addition},
	{ttMultiply,           Precedence::Multiplication},
	{ttDivision,           Precedence::Multiplication},
	{ttPercent,            Precedence::Multiplication},
	{ttStarStar,           Precedence::Exponentiation},
	{ttIncrement,          Precedence::PrefixIncrement},
	{ttDecrement,          Precedence::PrefixDecrement},
	{ttLeftBracket,        Precedence::Indexing},
	{ttDot,                Precedence::Indexing},
	{ttNew,                Precedence::New},
	{ttLeftParen,          Precedence::Call},
};

Program* Parser::Parse() {
    auto program = new Program();
    for(;!_follow(ttEOF);) {
        auto stmt = _parseStatement(true);
        program->_stmts.push_back(stmt);
    }
    if(_next().type != ttEOF) {
        throw SyntaxError("unexpected token");
    }
    return program;
}

Token Parser::_expect(TokenType tt) {
    auto next = _next();
    if(next.type != tt) {
        Token tk;
        tk.type = tt;
        auto exp = tk.string();
        if(tt == ttIdentifier) {
            exp ="`identifier'";
        }
        throw SyntaxError(
            "unexpected token: %s (expect: %s)",
            next.string().c_str(), tk.string().c_str()
        );
    }
    return next;
}

Token Parser::_next() {
    return _tkz->Next();
}

void Parser::_undo(Token tk) {
    _tkz->Undo(tk);
}

bool Parser::_skip(TokenType tt) {
    if(_follow(tt)) {
        _next();
        return true;
    }
    return false;
}

Token Parser::_peek() {
    return _tkz->Peek();
}

bool Parser::_follow(TokenType tt) {
    return _peek().type == tt;
}

void Parser::_enter() {

}

void Parser::_leave(bool putBack) {

}

Precedence Parser::_getPrecedence(TokenType op) {
    if(op >= ttAssign && op <= ttAndNotAssign) {
        return Precedence::Assignment;
    }
    auto iter = precedenceTable.find(op);
    if(iter != precedenceTable.cend()) {
        return iter->second;
    }
    return Precedence::_Unspecified;
}

BaseStatement* Parser::_parseStatement(bool global) {
    auto tk = _peek();

    switch(tk.type) {
    case ttLet:
        return _parseLetStatement();
    case ttFunction:
        return _parseFunctionStatement();
    case ttSemicolon:
        _next();
        return new EmptyStatement();
    default:
        break;
    }

    if (global) {
        throw SyntaxError("non-global statement");
    }

    switch(tk.type) {
    case ttReturn:
        return _parseReturnStatememt();
    case ttLeftBrace:
        return _parseBlockStatement();
    case ttFor:
        return _parseForStatement();
    case ttBreak:
        return _parseBreakStatement();
    case ttIf:
        return _parseIfStatement();
    case ttSwitch:
        //return _parseSwitchStatement();
    default:
        break;
    }

    {
        auto expr = _parseExpression(Precedence::Assignment);
        Token tk;
        bool match;
        std::tie(tk, match) = _match(ttSemicolon);
        if(match) {
            auto es = new ExpressionStatement();
            es->_expr = expr;
            return es;
        }
        _expect(ttSemicolon);
    }

    throw SyntaxError("unknown statement at line: %d", tk.line);
}

// done
LetStatement* Parser::_parseLetStatement() {
    auto let = new LetStatement();
    _expect(ttLet);
    let->_name = _expect(ttIdentifier).str;
    if(_follow(ttAssign)) {
        _next();
        let->_expr = _parseExpression(Precedence::Conditional);
    }
    _expect(ttSemicolon);
    return let;
}

// done
FunctionStatement* Parser::_parseFunctionStatement() {
    auto fn = new FunctionStatement();
    fn->_expr = _parseFunctionExpression();
    return fn;
}

// done
ReturnStatement* Parser::_parseReturnStatememt() {
    auto rtn = new ReturnStatement();
    _expect(ttReturn);
    if(!_follow(ttSemicolon)) {
        rtn->_expr = _parseExpression(Precedence::Conditional);
    }
    _expect(ttSemicolon);
    return rtn;
}

// done
BlockStatement* Parser::_parseBlockStatement() {
    auto block = new BlockStatement();
    _expect(ttLeftBrace);
    for(;;) {
        if(_follow(ttRightBrace)) {
            break;
        }
        auto stmt = _parseStatement(false);
        block->_stmts.push_back(stmt);
    }
    _expect(ttRightBrace);
    return block;
}

// done
ForStatement* Parser::_parseForStatement() {
    auto fs = new ForStatement();
    auto hasInit = false;

    _expect(ttFor);

    if(_follow(ttLet)) {
        hasInit = true;
        fs->_init = _parseLetStatement();
    } else if(_follow(ttSemicolon)) {
        hasInit = true;
        _expect(ttSemicolon);
    } else if(!_follow(ttLeftBrace)) {
        hasInit = false;
        fs->_test = _parseExpression(Precedence::Conditional);
    }

    if(hasInit) {
        // test
        if(!_follow(ttSemicolon)) {
            fs->_test = _parseExpression(Precedence::Conditional);
            _expect(ttSemicolon);
        } else {
            _next();
        }
        // incr
        if(!_follow(ttLeftBrace)) {
            fs->_incr = _parseExpression(Precedence::Assignment);
        }
    } else {
        if(!_follow(ttLeftBrace)) {
            throw SyntaxError("for needs body");
        }
    }

    _breakCount++;

    fs->_block = _parseBlockStatement();
    
    _breakCount--;

    return fs;
}

// done
BreakStatement* Parser::_parseBreakStatement() {
    _expect(ttBreak);
    _expect(ttSemicolon);
    return new BreakStatement();
}

// done
IfStatement* Parser::_parseIfStatement() {
    auto stmt = new IfStatement();
    _expect(ttIf);
    stmt->_cond = _parseExpression(Precedence::Conditional);
    stmt->_ifBlock = _parseBlockStatement();
    if(_follow(ttElse)) {
        _next();
        switch(_peek().type) {
        case ttIf:
            stmt->_elseBlock = _parseIfStatement();
            break;
        case ttLeftBrace:
            stmt->_elseBlock = _parseBlockStatement();
            break;
        default:
            throw SyntaxError("else expect if or block to follow");
        }
    }
    return stmt;
}

IExpression* Parser::_parseExpression(Precedence prec) {
    IExpression* left = nullptr;
    auto peek = _peek();
    if(prec <= Precedence::UnaryPlus) {
        switch(peek.type) {
        case ttLogicalNot:
        case ttBitXor:
        case ttAddition:
        case ttSubtraction: {
            _next();
            auto right = _parseExpression(Precedence::UnaryPlus);
            left = new UnaryExpression(peek.type, right);
            break;
        }
        default:
            break;
        }
    }
    if(prec <= Precedence::PrefixIncrement) {
        switch(peek.type) {
        case ttIncrement:
        case ttDecrement: {
            _next();
            auto right = _parseExpression(Precedence::PrefixIncrement);
            left = new IncrementExpression(peek.type, true, right);
            break;
        }
        default:
            break;
        }
    }
    if(prec <= Precedence::New) {
        if(peek.type == ttNew) {
            left = _parseNewExpression();
        }
    }
    if(left == nullptr) {
        left = _parsePrimaryExpression();
    }

    while((1)) {
        auto op = _next();
        auto nextPrec = _getPrecedence(op.type);
        if(nextPrec == Precedence(0) || nextPrec < prec) {
            _undo(op);
            break;
        }

        if(op.type >= ttAssign && op.type < ttQuestion) {
            return _parseAssignmentExpression(left, op.type);
        }

        switch(op.type) {
        case ttQuestion:
            left = _parseTernaryExpression(left);
            continue;
        case ttIncrement:
        case ttDecrement:
            left = new IncrementExpression(op.type, false, left);
            continue;
        case ttLeftParen:
            _undo(op);
            left = _parseCallExpression(left);
            continue;
        case ttLeftBracket:
        case ttDot:
            _undo(op);
            left = _parseIndexExpression(left);
            continue;
        default:
            break;
        }

        IExpression* right = nullptr;

        switch(op.type) {
        case ttLogicalOr:
            right = _parseExpression(Precedence::LogicalAnd);
            break;
        case ttLogicalAnd:
            right = _parseExpression(Precedence::BitwiseAnd);
            break;
        case ttBitAnd:
        case ttBitOr:
        case ttBitXor:
        case ttBitAndNot:
            right = _parseExpression(Precedence::Equality);
            break;
        case ttEqual:
        case ttNotEqual:
            right = _parseExpression(Precedence::Comparison);
            break;
        case ttGreaterThan:
        case ttGreaterThanOrEqual:
        case ttLessThan:
        case ttLessThanOrEqual:
            right = _parseExpression(Precedence::BitwiseShift);
            break;
        case ttLeftShift:
        case ttRightShift:
            right = _parseExpression(Precedence::Addition);
            break;
        case ttAddition:
        case ttSubtraction:
            right = _parseExpression(Precedence::Multiplication);
            break;
        case ttMultiply:
        case ttDivision:
        case ttPercent:
            right = _parseExpression(Precedence::Exponentiation);
            break;
        case ttStarStar:
            right = _parseExpression(Precedence::Indexing);
            break;
        default:
            break;
        }

        if(right != nullptr) {
            left = new BinaryExpression(left, op.type, right);
            continue;
        }

        throw SyntaxError("unknown operator");
    }

    return left;
}

IExpression* Parser::_parsePrimaryExpression() {
    IExpression* expr;
    auto next = _next();
    switch(next.type) {
    case ttNil:
        expr = Value::fromNil();
        break;
    case ttBoolean:
        expr = Value::fromBoolean(next.str == "true");
        break;
    case ttNumber:
        expr = Value::fromNumber(next.num);
        break;
    case ttString:
        expr = Value::fromString(next.str);
        break;
    case ttLeftParen:
        _undo(next);
        if(auto lambda = _tryParseLambdaExpression(false)) {
            return lambda;
        }
        _next();
        expr = _parseExpression(Precedence::Conditional);
        _expect(ttRightParen);
        break;
    case ttIdentifier:
        if(_follow(ttLambda)) {
            _undo(next);
            return _tryParseLambdaExpression(true);
        }
        expr = Value::fromVariable(next.str);
        break;
    case ttFunction:
        _undo(next);
        expr = _parseFunctionExpression();
        break;
    case ttLeftBrace:
        _undo(next);
        expr = _parseObjectExpression();
        break;
    case ttLeftBracket:
        _undo(next);
        expr = _parseArrayExpression();
        break;
    default:
        throw SyntaxError("unexpected token");
    }

    return expr;
}

TernaryExpression* Parser::_parseTernaryExpression(IExpression* cond) {
    IExpression *left;
    IExpression *right;

    left = _parseExpression(Precedence::Conditional);
    _expect(ttColon);
    right = _parseExpression(Precedence::Conditional);

/*
    static const char* err = "nested `?:' is not allowed";
    if(left->type == ExprType::Ternary) {
        throw SyntaxError(err);
    }
    if(right->type == ExprType::Ternary) {
        throw SyntaxError(err);
    }
    */

    auto expr = new TernaryExpression();
    expr->left = left;
    expr->cond = cond;
    expr->right = right;

    return expr;
}

AssignmentExpression* Parser::_parseAssignmentExpression(IExpression* left, TokenType op) {
    auto expr = new AssignmentExpression();
    expr->_left = left;

	// ttQuestion: disable continuous assignment style
    expr->_expr = _parseExpression(Precedence::Conditional);

	if(op == ttAssign) {
		return expr;
	}

    TokenType binOp;

	switch(op) {
	case ttStarStarAssign:
		binOp = ttStarStar;
        break;
	case ttStarAssign:
		binOp = ttMultiply;
        break;
	case ttDivideAssign:
		binOp = ttDivision;
        break;
	case ttPercentAssign:
		binOp = ttPercent;
        break;
	case ttPlusAssign:
		binOp = ttAddition;
        break;
	case ttMinusAssign:
		binOp = ttSubtraction;
        break;
	case ttLeftShiftAssign:
		binOp = ttLeftShift;
        break;
	case ttRightShiftAssign:
		binOp = ttRightShift;
        break;
	case ttAndAssign:
		binOp = ttBitAnd;
        break;
	case ttOrAssign:
		binOp = ttBitOr;
        break;
	case ttXorAssign:
		binOp = ttBitXor;
        break;
	case ttAndNotAssign:
		binOp = ttBitAndNot;
        break;
	default:
        throw Error("won't go here");
	}

    auto bin = new BinaryExpression();
    bin->_left = expr->_left;
    bin->_op = binOp;
    bin->_right = expr->_expr;

    expr->_expr = bin;
    return expr;
}

NewExpression* Parser::_parseNewExpression() {
    auto expr = new NewExpression();
    _expect(ttNew);
    expr->_name = _expect(ttIdentifier).str;
    _expect(ttLeftParen);
    if(!_follow(ttRightParen)) {
        for (;;) {
            auto arg = _parseExpression(Precedence::Conditional);
            expr->_args.Put(arg);
            auto sep = _next();
            if(sep.type == ttComma) {
                continue;
            } else if(sep.type == ttRightParen) {
                _undo(sep);
                break;
            } else {
                throw SyntaxError("unexpected token");
            }
        }
    }
    _expect(ttRightParen);
    return expr;
}

FunctionExpression* Parser::_tryParseLambdaExpression(bool must) {
    _enter();
}

IndexExpression* Parser::_parseIndexExpression(IExpression* left) {
    auto ie = new IndexExpression();
    ie->_indexable = left;
    auto next = _next();
    if(next.type == ttDot) {
        auto key = _next();
        if(key.type == ttIdentifier) {
            ie->_key = Value::fromString(key.str);
            return ie;
        }
        throw SyntaxError("unexpected token: %s", key.string().c_str());
    } else if(next.type == ttLeftBracket) {
        auto key = _parseExpression(Precedence::Conditional);
        _expect(ttRightBracket);
        ie->_key = key;
        return ie;
    } else {
        throw Error("won't go here");
    }
}

CallExpression* Parser::_parseCallExpression(IExpression* left) {
    auto ce = new CallExpression();
    ce->_callable = left;
    _expect(ttLeftParen);
    if(!_follow(ttRightParen)) {
        for(;;) {
            auto arg = _parseExpression(Precedence::Conditional);
            ce->_args.Put(arg);
            auto sep = _next();
            if(sep.type == ttComma) {
                continue;
            } else if(sep.type == ttRightParen) {
                _undo(sep);
                break;
            } else {
                SyntaxError(
                    "unexpected token: %s",
                    Token(sep).string().c_str()
                );
            }
        }
    }
    _expect(ttRightParen);
    return ce;
}

FunctionExpression* Parser::_parseFunctionExpression() {
    auto fe = new FunctionExpression();

    _expect(ttFunction);

    if(_follow(ttIdentifier)) {
        fe->_name = _next().str;
    }

    _expect(ttLeftParen);
    if(!_follow(ttRightParen)) {
        for(;;) {
            auto name = _expect(ttIdentifier).str;
            fe->_params.Put(name);
            auto sep = _next();
            if(sep.type == ttComma) {
                continue;
            } else if(sep.type == ttRightParen) {
                _undo(sep);
                break;
            } else {
                throw SyntaxError(
                    "unexpected token: %s",
                    Token(sep).string().c_str()
                );
            }
        }
    }
    _expect(ttRightParen);

    if(!_follow(ttLeftBrace)) {
        throw SyntaxError("function needs a body");
    }

    auto savedBreakCount = _breakCount;
    _breakCount = 0;
    fe->_body = _parseBlockStatement();
    _breakCount = savedBreakCount;
    return fe;
}

ObjectExpression* Parser::_parseObjectExpression() {
    auto obj = new ObjectExpression();
    _expect(ttLeftBrace);
    for(;;) {
        if(!_follow(ttRightBrace)) {
            break;
        }

        std::string key;
        IExpression* val;

        auto next = _next();
        switch(next.type) {
        case ttString:
        case ttIdentifier:
            key = next.str;
            break;
        default:
            throw TypeError("unsupported key type");
        }

        _expect(ttColon);

        val = _parseExpression(Precedence::Conditional);
        auto iter = obj->_props.find(key);
        if(iter == obj->_props.cend()) {
            obj->_props[key] = val;
        } else {
            throw SyntaxError("duplicate key");
        }

        _skip(ttComma);
        if(_follow(ttRightBrace)) {
            break;
        }
    }

    _expect(ttRightBrace);
    return obj;
}

ArrayExpression* Parser::_parseArrayExpression() {
    auto arr = new ArrayExpression();
    _expect(ttLeftBracket);
    for(;;) {
        if(_follow(ttRightBracket)) {
            break;
        }
        auto elem = _parseExpression(Precedence::Conditional);
        arr->_elems.Put(elem);
        _skip(ttComma);
        if(_follow(ttRightBracket)) {
            break;
        }
    }
    _expect(ttRightBracket);
    return arr;
}

}
