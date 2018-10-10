#include "tokenizer.h"
#include "error.h"

namespace taolang {

std::map<const char*,TokenType> keywords = {
	{"break",    ttBreak},
	{"case",     ttCase},
	{"default",  ttDefault},
	{"else",     ttElse},
	{"for",      ttFor},
	{"function", ttFunction},
	{"if",       ttIf},
	{"let",      ttLet},
	{"switch",   ttSwitch},
	{"return",   ttReturn},
	{"nil",      ttNil},
	{"true",     ttBoolean},
	{"false",    ttBoolean},
	{"tao",      ttTao},
	{"new",      ttNew},
};

std::map<TokenType, const char*> tokenNames = {
	{ttEOF,          "EOF"},
	{ttLeftParen,    "("},
	{ttRightParen,   ")"},
	{ttLeftBracket,  "["},
	{ttRightBracket, "]"},
	{ttLeftBrace,    "{"},
	{ttRightBrace,   "}"},

	{ttDot,       "."},
	{ttComma,     ","},
	{ttSemicolon, ";"},
	{ttColon,     ":"},
	{ttLambda,    "=>"},

	{ttAssign,           "="},
	{ttPlusAssign,       "+="},
	{ttMinusAssign,      "-="},
	{ttStarStarAssign,   "**="},
	{ttStarAssign,       "*="},
	{ttDivideAssign,     "/="},
	{ttPercentAssign,    "%="},
	{ttLeftShiftAssign,  "<<="},
	{ttRightShiftAssign, ">>="},
	{ttAndAssign,        "&="},
	{ttXorAssign,        "^="},
	{ttOrAssign,         "|="},
	{ttAndNotAssign,     "&^="},

	{ttQuestion, "?"},

	{ttLogicalNot, "!"},
	{ttLogicalAnd, "&&"},
	{ttLogicalOr,  "||"},

	{ttBitAnd,    "&"},
	{ttBitOr,     "|"},
	{ttBitXor,    "^"},
	{ttBitAndNot, "&^"},

	{ttEqual,    "=="},
	{ttNotEqual, "!="},

	{ttGreaterThan,        ">"},
	{ttGreaterThanOrEqual, ">="},
	{ttLessThan,           "<"},
	{ttLessThanOrEqual,    "<="},

	{ttLeftShift,  "<<"},
	{ttRightShift, ">>"},

	{ttAddition,    "+"},
	{ttSubtraction, "-"},
	{ttMultiply,    "*"},
	{ttDivision,    "/"},
	{ttPercent,     "%"},
	{ttStarStar,    "**"},

	{ttIncrement, "++"},
	{ttDecrement, "--"},

	{ttNil, "nil"},

	{ttBreak,    "break"},
	{ttCase,     "case"},
	{ttDefault,  "default"},
	{ttElse,     "else"},
	{ttFor,      "for"},
	{ttFunction, "function"},
	{ttIf,       "if"},
	{ttLet,      "let"},
	{ttSwitch,   "switch"},
	{ttReturn,   "return"},
	{ttNew,      "new"},
};

Token Tokenizer::Next() {
	if(!_buf.empty()) {
		auto tk = _buf.front();
		_buf.pop_front();
		return tk;
	}
	return next();
}

Token Tokenizer::Peek() {
	auto tk = Next();
	Undo(tk);
	return tk;
}

void Tokenizer::Undo(const Token& token) {
	_buf.push_front(token);
}

Token Tokenizer::next() {
	Token tk;

	while((true)) {
		auto ch = read();
		if(ch == '\0' || ch == -1) {
			tk.type = ttEOF;
			goto exit;
		}

		if(ch >= '0' && ch <= '9'){
			unread();
			auto n = readNumber();
			checkFollow();
			tk.type = ttNumber;
			tk.num = n;
			goto exit;
		} else if((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'){
			unread();
			auto name = readIdent();
			checkFollow();
			auto type = ttIdentifier;
			auto it = keywords.find(name.c_str());
			if(it != keywords.cend()) {
				tk.type = it->second;
			} else {
				tk.type = ttIdentifier;
				tk.str = name;
			}
			goto exit;
		} else if(ch == '"') {
			unread();
			auto s = readString();
			checkFollow();
			tk.type = ttString;
			tk.str = s;
			goto exit;
		}

		switch(ch) {
		case ' ':
		case '\t':
		case '\r':
		case '\n':
			continue;
		case '(':
			tk.type = ttLeftParen;
			goto exit;
		case ')':
			tk.type = ttRightParen;
			goto exit;
		case '[':
			tk.type = ttLeftBracket;
			goto exit;
		case ']':
			tk.type = ttRightBracket;
			goto exit;
		case '{':
			tk.type = ttLeftBrace;
			goto exit;
		case '}':
			tk.type = ttRightBrace;
			goto exit;
		case '.':
			tk.type = ttDot;
			goto exit;
		case ',':
			tk.type = ttComma;
			goto exit;
		case ':':
			tk.type = ttColon;
			goto exit;
		case '?':
			tk.type = ttQuestion;
			goto exit;
		case ';':
			tk.type = ttSemicolon;
			goto exit;
		case '+':
			tk.type = iiif('+', '=', ttIncrement, ttPlusAssign, ttAddition);
			goto exit;
		case '-':
			tk.type = iiif('-', '=', ttDecrement, ttMinusAssign, ttSubtraction);
			goto exit;
		case '*':
			switch(read()) {
			case '*':
				tk.type = iif('=', ttStarStarAssign, ttStarStar);
				break;
			case '=':
				tk.type = ttStarAssign;
				break;
			default:
				unread();
				tk.type = ttMultiply;
				break;
			}
			goto exit;
		case '/':
			switch(read()) {
			case '/':
				readComment();
				continue;
			case '=':
				tk.type = ttDivideAssign;
				break;
			default:
				unread();
				tk.type = ttDivision;
				break;
			}
			goto exit;
		case '%':
			tk.type = iif('=', ttPercentAssign, ttPercent);
			goto exit;
		case '=':
			tk.type = iiif('=', '>', ttEqual, ttLambda, ttAssign);
			goto exit;
		case '>':
			switch(read()) {
			case '=':
				tk.type = ttGreaterThanOrEqual;
				break;
			case '>':
				tk.type = iif('=', ttRightShiftAssign, ttRightShift);
				break;
			default:
				unread();
				tk.type = ttGreaterThan;
				break;
			}
			goto exit;
		case '<':
			switch(read()){
			case '=':
				tk.type = ttLessThanOrEqual;
				break;
			case '<':
				tk.type = iif('=', ttLeftShiftAssign, ttLeftShift);
				break;
			default:
				unread();
				tk.type = ttLessThan;
				break;
			}
			goto exit;
		case '!':
			tk.type = iif('=', ttNotEqual, ttLogicalNot);
			goto exit;
		case '&':
			switch(read()) {
			case '&':
				tk.type = ttLogicalAnd;
				break;
			case '=':
				tk.type = ttAndAssign;
				break;
			case '^':
				tk.type = iif('=', ttAndNotAssign, ttBitAndNot);
				break;
			default:
				unread();
				tk.type = ttBitAnd;
				break;
			}
			goto exit;
		case '|':
			tk.type = iiif('|', '=', ttLogicalOr, ttOrAssign, ttBitOr);
			goto exit;
		case '^':
			tk.type = iif('=', ttXorAssign, ttBitXor);
			goto exit;
		}
		throw SyntaxError("unhandled character `%c' at: line:%d,col:%d", ch, _line, _col);
	}
exit:
	tk.line = _line;
	tk.col = _col;
	return std::move(tk);
}

int Tokenizer::read() {
	_ch = std::fgetc(_fp);
	_col++;
	if(_ch == '\n') {
		_line++;
		_col = 1;
	}
	return _ch;
}

void Tokenizer::unread() {
	if(_ch != '\0') {
		std::ungetc(_ch, _fp);
		if(_ch == '\n') {
			_line--;
			_col = 1;
		} else {
			_col--;
		}
	}
}

TokenType Tokenizer::iif(uint8_t c, TokenType t1, TokenType t2) {
	auto ch = read();
	if(ch == c) {
		return t1;
	} else {
		unread();
		return t2;
	}
}

TokenType Tokenizer::iiif(uint8_t c1, uint8_t c2, TokenType t1, TokenType t2, TokenType t3) {
	auto ch = read();
	if(ch == c1) {
		return t1;
	} else if(ch ==c2) {
		return t2;
	} else {
		unread();
		return t3;
	}
}

void Tokenizer::checkFollow() {
	auto ch = read();
	unread();

	if((ch >= '0' && ch <= '9') ||
		(ch >= 'a' && ch <= 'z') ||
		(ch >= 'A' && ch <= 'Z') ||
		ch == '"')
	{
		throw SyntaxError("unexpected follow character %c at line:%d,col:%d", ch, _line, _col);
	}
}

std::string Tokenizer::readString() {
	std::ostringstream ss;
	read();
	while((true)) {
		auto ch = read();
		switch(ch) {
		case '"':
			goto exit;
		case '\0':
			throw SyntaxError("unterminated string literal at: line %d, col %d", _line, _col);
			break;
		default:
			ss << char(ch);
			break;
		}
	}
exit:
	return std::move(ss.str());
}

int64_t Tokenizer::readNumber() {
	int64_t i = 0;
	while((true)) {
		auto ch = read();
		if(ch>='0' && ch<='9') {
			i = i*10 + ch-'0';
		} else {
			unread();
			break;
		}
	}
	return i;
}

std::string Tokenizer::readIdent() {
	std::ostringstream ss;
	while((true)) {
		auto ch = read();
		if((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '_')
		{
			ss << char(ch);
		} else {
			unread();
			break;
		}
	}
	return std::move(ss.str());
}

void Tokenizer::readComment() {
	while((true)) {
		auto c = read();
		if(c=='\n' || c == '\0') {
			if(c=='\n') {
				_line++;
			}
			break;
		}
	}
}

}
