#include "parse.h"

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


}