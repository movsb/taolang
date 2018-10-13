#include "parser.h"
#include "error.h"

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
        return _parseSwitchStatement();
    default:
        break;
    }

    {
        auto expr = _parseExpression(Precedence::Assignment);
        // TODO
    }

    throw SyntaxError("unknown statement at line: %d", tk.line);
}

// done
LetStatement* Parser::_parseLetStatement() {
    auto let = new LetStatement();
    _expect(ttLet);
    let->name = _expect(ttIdentifier).str;
    if(_follow(ttAssign)) {
        _next();
        let->expr = _parseExpression(Precedence::Conditional);
    }
    _expect(ttSemicolon);
    return let;
}

// done
FunctionStatement* Parser::_parseFunctionStatement() {
    auto fn = new FunctionStatement();
    fn->expr = _parseFunctionExpression();
    return fn;
}

// done
ReturnStatement* Parser::_parseReturnStatememt() {
    auto rtn = new ReturnStatement();
    _expect(ttReturn);
    if(!_follow(ttSemicolon)) {
        rtn->expr = _parseExpression(Precedence::Conditional);
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
        block->stmts.push_back(stmt);
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
        fs->init = _parseLetStatement();
    } else if(_follow(ttSemicolon)) {
        hasInit = true;
        _expect(ttSemicolon);
    } else if(!_follow(ttLeftBrace)) {
        hasInit = false;
        fs->test = _parseExpression(Precedence::Conditional);
    }

    if(hasInit) {
        // test
        if(!_follow(ttSemicolon)) {
            fs->test = _parseExpression(Precedence::Conditional);
            _expect(ttSemicolon);
        } else {
            _next();
        }
        // incr
        if(!_follow(ttLeftBrace)) {
            fs->incr = _parseExpression(Precedence::Assignment);
        }
    } else {
        if(!_follow(ttLeftBrace)) {
            throw SyntaxError("for needs body");
        }
    }

    _breakCount++;

    fs->block = _parseBlockStatement();
    
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
    stmt->cond = _parseExpression(Precedence::Conditional);
    stmt->ifBlock = _parseBlockStatement();
    if(_follow(ttElse)) {
        _next();
        switch(_peek().type) {
        case ttIf:
            stmt->elseBlock = _parseIfStatement();
            break;
        case ttLeftBrace:
            stmt->elseBlock = _parseBlockStatement();
            break;
        default:
            throw SyntaxError("else expect if or block to follow");
        }
    }
    return stmt;
}

}
