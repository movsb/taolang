#pragma once

#include <cstdio>

#include <map>
#include <string>
#include <sstream>
#include <deque>

namespace taolang {

enum TokenType {
	ttEOF,

	// braces
	ttLeftParen,
	ttRightParen,
	ttLeftBrace,
	ttRightBrace,
	ttLeftBracket,
	ttRightBracket,

	// seperators
	ttDot,
	ttComma,
	ttSemicolon,
	ttColon,
	ttLambda,

	// assignment
	ttAssign,
	ttPlusAssign,
	ttMinusAssign,
	ttStarStarAssign,
	ttStarAssign,
	ttDivideAssign,
	ttPercentAssign,
	ttLeftShiftAssign,
	ttRightShiftAssign,
	ttAndAssign,
	ttXorAssign,
	ttOrAssign,
	ttAndNotAssign,

	// Conditional
	ttQuestion,

	// Logical
	ttLogicalNot,
	ttLogicalAnd,
	ttLogicalOr,

	// Bit
	ttBitAnd,
	ttBitOr,
	ttBitXor,
	ttBitAndNot,

	// Equality
	ttEqual,
	ttNotEqual,

	// comparision
	ttGreaterThan,
	ttGreaterThanOrEqual,
	ttLessThan,
	ttLessThanOrEqual,

	// Shift
	ttLeftShift,
	ttRightShift,

	// arithmetic
	ttAddition,
	ttSubtraction,
	ttMultiply,
	ttDivision,
	ttPercent,
	ttStarStar,

	// ++ --
	ttIncrement,
	ttDecrement,

	// Literals
	ttNil,
	ttString,
	ttNumber,
	ttBoolean,
	ttIdentifier,

	// Keywords
	ttBreak,
	ttCase,
	ttDefault,
	ttElse,
	ttFor,
	ttFunction,
	ttIf,
	ttLet,
	ttSwitch,
	ttReturn,
	ttTao,
	ttNew,
};

extern std::map<const char*,TokenType> keywords;
extern std::map<TokenType, const char*> tokenNames;

struct Token {
	TokenType   type;
	std::string str;
	int64_t     num;
	int         line;
	int         col;

	std::string string() {
		auto it = tokenNames.find(type);
		if(it != tokenNames.cend()) {
			return it->second;
		}

		std::ostringstream ss;

		ss << "`";

		switch(type) {
		case ttString:
			ss << str;
			break;
		case ttNumber:
			ss << num;
			break;
		case ttBoolean:
			ss << str;
			break;
		case ttIdentifier:
			ss << str;
			break;
		default:
			ss << "--unknown-token--";
			break;
		}

		ss << "'";
		if(line > 0 && col >0) {
			ss << " (line:" << line << " col:" << col << ")";
		}
		return std::move(ss.str());
	}
};

class Tokenizer {
public:
	Tokenizer(std::FILE* fp)
		: _fp(fp)
        , _line(1)
        , _col(1)
        , _ch(0)
    {
    }
public:
	Token Next();
	Token Peek();
	void Undo(const Token& token);

private:
	int _line;
	int _col;
	int _ch;
	FILE* _fp;
	std::deque<Token> _buf;

private:
	Token next();
	int read();
	void unread();
	void checkFollow();
	TokenType iif(uint8_t c, TokenType t1, TokenType t2);
	TokenType iiif(uint8_t c1, uint8_t c2, TokenType t1, TokenType t2, TokenType t3);
	std::string readString();
	int64_t readNumber();
	std::string readIdent();
	void readComment();
};

}
