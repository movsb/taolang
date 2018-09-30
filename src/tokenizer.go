package main

import (
	"bufio"
	"bytes"
	"container/list"
	"fmt"
	"io"
	"os"
)

// TokenType is the type of a token.
// Operators are listed by precedence groups from lowest to highest.
// ttAssign -> ttQuestion -> ttIncrement
type TokenType uint

const (
	_ TokenType = iota + 127

	ttEOF

	// braces
	ttLeftParen
	ttRightParen
	ttLeftBrace
	ttRightBrace
	ttLeftBracket
	ttRightBracket

	// seperators
	ttDot
	ttComma
	ttSemicolon
	ttColon
	ttLambda

	// assignment
	ttAssign
	ttPlusAssign
	ttMinusAssign
	ttStarStarAssign
	ttStarAssign
	ttDivideAssign
	ttPercentAssign
	ttLeftShiftAssign
	ttRightShiftAssign
	ttAndAssign
	ttXorAssign
	ttOrAssign
	ttAndNotAssign

	// Conditional
	ttQuestion

	// Logical
	ttNot
	ttAndAnd
	ttOrOr

	// Bit
	ttBitAnd
	ttBitOr
	ttBitXor
	ttBitAndNot

	// Equality
	ttEqual
	ttNotEqual

	// comparision
	ttGreaterThan
	ttGreaterThanOrEqual
	ttLessThan
	ttLessThanOrEqual

	// Shift
	ttLeftShift
	ttRightShift

	// arithmetic
	ttAddition
	ttSubtraction
	ttMultiply
	ttDivision
	ttPercent
	ttStarStar

	// ++ --
	ttIncrement
	ttDecrement

	// Literals
	ttNil
	ttString
	ttNumber
	ttBoolean
	ttIdentifier

	// Keywords
	ttBreak
	ttCase
	ttDefault
	ttElse
	ttFor
	ttFunction
	ttIf
	ttLet
	ttSwitch
	ttReturn
	ttTao
)

var keywords map[string]TokenType

func init() {
	keywords = map[string]TokenType{
		"break":    ttBreak,
		"case":     ttCase,
		"default":  ttDefault,
		"else":     ttElse,
		"for":      ttFor,
		"function": ttFunction,
		"if":       ttIf,
		"let":      ttLet,
		"switch":   ttSwitch,
		"return":   ttReturn,
		"nil":      ttNil,
		"true":     ttBoolean,
		"false":    ttBoolean,
		"tao":      ttTao,
	}
}

var tokenNames map[TokenType]string

func init() {
	tokenNames = map[TokenType]string{
		ttEOF:          "EOF",
		ttLeftParen:    "(",
		ttRightParen:   ")",
		ttLeftBracket:  "[",
		ttRightBracket: "]",
		ttLeftBrace:    "{",
		ttRightBrace:   "}",

		ttDot:       ".",
		ttComma:     ",",
		ttSemicolon: ";",
		ttColon:     ":",
		ttLambda:    "=>",

		ttAssign:           "=",
		ttPlusAssign:       "+=",
		ttMinusAssign:      "-=",
		ttStarStarAssign:   "**=",
		ttStarAssign:       "*=",
		ttDivideAssign:     "/=",
		ttPercentAssign:    "%=",
		ttLeftShiftAssign:  "<<=",
		ttRightShiftAssign: ">>=",
		ttAndAssign:        "&=",
		ttXorAssign:        "^=",
		ttOrAssign:         "|=",
		ttAndNotAssign:     "&^=",

		ttQuestion: "?",

		ttNot:    "!",
		ttAndAnd: "&&",
		ttOrOr:   "||",

		ttBitAnd:    "&",
		ttBitOr:     "|",
		ttBitXor:    "^",
		ttBitAndNot: "&^",

		ttEqual:    "==",
		ttNotEqual: "!=",

		ttGreaterThan:        ">",
		ttGreaterThanOrEqual: ">=",
		ttLessThan:           "<",
		ttLessThanOrEqual:    "<=",

		ttLeftShift:  "<<",
		ttRightShift: ">>",

		ttAddition:    "+",
		ttSubtraction: "-",
		ttMultiply:    "*",
		ttDivision:    "/",
		ttPercent:     "%",
		ttStarStar:    "**",

		ttIncrement: "++",
		ttDecrement: "--",

		ttNil: "nil",

		ttBreak:    "break",
		ttCase:     "case",
		ttDefault:  "default",
		ttElse:     "else",
		ttFor:      "for",
		ttFunction: "function",
		ttIf:       "if",
		ttLet:      "let",
		ttSwitch:   "switch",
		ttReturn:   "return",
	}
}

// Token is a token.
type Token struct {
	typ  TokenType
	str  string
	num  int
	line int
	col  int
}

func (t Token) String() (ret string) {
	defer func() {
		ret = "`" + ret + "'"
		if t.line > 0 && t.col > 0 {
			ret += fmt.Sprintf(" (line:%d col:%d)", t.line, t.col)
		}
	}()

	if s, ok := tokenNames[t.typ]; ok {
		return s
	}
	switch t.typ {
	case ttString:
		return t.str
	case ttNumber:
		return fmt.Sprint(t.num)
	case ttBoolean:
		return t.str
	case ttIdentifier:
		return t.str
	}
	return "--unknown-token--"
}

// Tokenizer splits input into tokens.
type Tokenizer struct {
	input  *bufio.Reader
	buf    *list.List
	frames []*list.List
	line   int  // current line number
	col    int  // current line column
	ch     byte // current read char
}

// NewTokenizer creates a new tokenizer.
func NewTokenizer(input io.Reader) *Tokenizer {
	return &Tokenizer{
		input: bufio.NewReader(input),
		buf:   list.New(),
		line:  1,
	}
}

// Next returns next token tokenized.
// It first uses buffers or frames if there is one.
func (t *Tokenizer) Next() (token Token) {
	// use frame
	defer func() {
		if len(t.frames) > 0 {
			frame := t.frames[len(t.frames)-1]
			frame.PushBack(token)
		}
		if except := recover(); except != nil {
			fmt.Printf("%v\n", except)
			os.Exit(-1)
		}
	}()

	// use inner buffer
	if t.buf.Len() > 0 {
		tk := t.buf.Front()
		t.buf.Remove(tk)
		token = tk.Value.(Token)
		return
	}

	// use new
	token = t.next()
	return
}

// Undo undoes(put back) a token.
func (t *Tokenizer) Undo(token Token) {
	t.buf.PushFront(token)
	if len(t.frames) > 0 {
		last := t.frames[len(t.frames)-1]
		if last.Len() == 0 {
			panic("cannot undo")
		}
		last.Remove(last.Back())
	}
}

// Peek peeks the next token.
func (t *Tokenizer) Peek() Token {
	token := t.Next()
	t.Undo(token)
	return token
}

// PushFrame starts a new lookahead.
func (t *Tokenizer) PushFrame() {
	t.frames = append(t.frames, list.New())
}

// PopFrame stops a lookahead.
func (t *Tokenizer) PopFrame(putBack bool) {
	if len(t.frames) == 0 {
		panic("bad PopFrame call")
	}
	last := t.frames[len(t.frames)-1]
	t.frames = t.frames[0 : len(t.frames)-1]
	if putBack && last.Len() > 0 {
		t.buf.PushFrontList(last)
	}
}

// next produces next token.
func (t *Tokenizer) next() (token Token) {
	defer func() {
		token.line = t.line
		token.col = t.col
		// fmt.Println(token)
	}()

	for {
		ch := t.read()
		if ch == 0 {
			return Token{
				typ: ttEOF,
			}
		}

		if ch >= '0' && ch <= '9' {
			t.unread()
			n := t.readNumber()
			t.checkFollow()
			return Token{
				typ: ttNumber,
				num: n,
			}
		} else if ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch == '_' {
			t.unread()
			name := t.readIdentifier()
			t.checkFollow()
			typ := ttIdentifier
			if t, ok := keywords[name]; ok {
				typ = t
			}
			return Token{
				typ: typ,
				str: name,
			}
		} else if ch == '"' {
			t.unread()
			s := t.readString()
			t.checkFollow()
			return Token{
				typ: ttString,
				str: s,
			}
		}

		switch ch {
		case ' ', '\t', '\r', '\n':
			continue
		case '(':
			return Token{typ: ttLeftParen}
		case ')':
			return Token{typ: ttRightParen}
		case '[':
			return Token{typ: ttLeftBracket}
		case ']':
			return Token{typ: ttRightBracket}
		case '{':
			return Token{typ: ttLeftBrace}
		case '}':
			return Token{typ: ttRightBrace}
		case '.':
			return Token{typ: ttDot}
		case ',':
			return Token{typ: ttComma}
		case ':':
			return Token{typ: ttColon}
		case '?':
			return Token{typ: ttQuestion}
		case ';':
			return Token{typ: ttSemicolon}
		case '+':
			return t.iiif('+', '=', ttIncrement, ttPlusAssign, ttAddition)
		case '-':
			return t.iiif('-', '=', ttDecrement, ttMinusAssign, ttSubtraction)
		case '*':
			switch next := t.read(); next {
			case '*':
				return t.iif('=', ttStarStarAssign, ttStarStar)
			case '=':
				return Token{typ: ttStarAssign}
			default:
				t.unread()
				return Token{typ: ttMultiply}
			}
		case '/':
			switch c := t.read(); c {
			case '/':
				t.readComment()
				continue
			case '=':
				return Token{typ: ttDivideAssign}
			default:
				t.unread()
				return Token{typ: ttDivision}
			}
		case '%':
			return t.iif('=', ttPercentAssign, ttPercent)
		case '=':
			return t.iiif('=', '>', ttEqual, ttLambda, ttAssign)
		case '>':
			switch c := t.read(); c {
			case '=':
				return Token{typ: ttGreaterThanOrEqual}
			case '>':
				return t.iif('=', ttRightShiftAssign, ttRightShift)
			default:
				t.unread()
				return Token{typ: ttGreaterThan}
			}
		case '<':
			switch c := t.read(); c {
			case '=':
				return Token{typ: ttLessThanOrEqual}
			case '<':
				return t.iif('=', ttLeftShiftAssign, ttLeftShift)
			default:
				t.unread()
				return Token{typ: ttLessThan}
			}
		case '!':
			return t.iif('=', ttNotEqual, ttNot)
		case '&':
			switch c := t.read(); c {
			case '&':
				return Token{typ: ttAndAnd}
			case '=':
				return Token{typ: ttAndAssign}
			case '^':
				return t.iif('=', ttAndNotAssign, ttBitAndNot)
			default:
				t.unread()
				return Token{typ: ttBitAnd}
			}
		case '|':
			return t.iiif('|', '=', ttOrOr, ttOrAssign, ttBitOr)
		case '^':
			return t.iif('=', ttXorAssign, ttBitXor)
		}

		panic(fmt.Sprintf("unhandled character `%c' at: line:%d,col:%d", ch, t.line, t.col))
	}
}

func (t *Tokenizer) read() byte {
	ch, _ := t.input.ReadByte()
	t.ch = ch
	//if t.col > 0 {
	t.col++
	//}
	if t.ch == '\n' {
		t.line++
		t.col = 1
	}
	return ch
}

func (t *Tokenizer) unread() {
	if t.ch != 0 {
		t.input.UnreadByte()
		if t.ch == '\n' {
			t.line--
			t.col = 1 // invalid
		} else {
			t.col--
		}
	}
}

func (t *Tokenizer) checkFollow() {
	ch := t.read()
	t.unread()

	if ch >= '0' && ch <= '9' ||
		ch >= 'a' && ch <= 'z' ||
		ch >= 'A' && ch <= 'Z' ||
		ch == '"' {
		panic(fmt.Sprintf("unexpected follow character %c at line:%d,col:%d", ch, t.line, t.col))
	}
}

// iif returns tt1 if next char is ch, else returns tt2.
func (t *Tokenizer) iif(ch byte, tt1 TokenType, tt2 TokenType) Token {
	c := t.read()
	if c == ch {
		return Token{typ: tt1}
	}
	t.unread()
	return Token{typ: tt2}
}

// if ch1 return tt1
// else if ch2 return tt2
// else return tt3
func (t *Tokenizer) iiif(ch1, ch2 byte, tt1, tt2, tt3 TokenType) Token {
	c := t.read()
	if c == ch1 {
		return Token{typ: tt1}
	} else if c == ch2 {
		return Token{typ: tt2}
	}
	t.unread()
	return Token{typ: tt3}
}

// readString reads a quoted string.
func (t *Tokenizer) readString() string {
	buf := bytes.NewBuffer(nil)
	t.read() // eat '"'
	for {
		ch := t.read()
		if ch == '"' {
			break
		}
		buf.WriteByte(ch)
	}
	return buf.String()
}

// readNumber reads a int32 number.
func (t *Tokenizer) readNumber() int {
	num := 0
	for {
		ch := t.read()
		if ch >= '0' && ch <= '9' {
			num = num*10 + (int(ch) - '0')
		} else {
			t.unread()
			break
		}
	}
	return num
}

// readIdentifier reads a identifier.
func (t *Tokenizer) readIdentifier() string {
	buf := bytes.NewBuffer(nil)
	for {
		ch := t.read()
		if ch >= 'a' && ch <= 'z' ||
			ch >= 'A' && ch <= 'Z' ||
			ch >= '0' && ch <= '9' ||
			ch == '_' {
			buf.WriteByte(ch)
		} else {
			t.unread()
			break
		}
	}
	return buf.String()
}

// readComment eats out comment
func (t *Tokenizer) readComment() {
	for {
		c := t.read()
		if c == '\n' || c == 0 {
			if c == '\n' {
				t.line++
			}
			break
		}
	}
}
