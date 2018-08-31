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
	ttDot
	ttComma
	ttSemicolon
	ttColon
	ttLambda

	// assignment
	ttAssign

	// arithmetic
	ttAddition
	ttSubstraction
	ttMultiply
	ttDivision
	ttPercent

	// comparision
	ttGreaterThan
	ttGreaterThanOrEqual
	ttEqual
	ttNotEqual
	ttLessThan
	ttLessThanOrEqual

	// Logical
	ttNot
	ttAndAnd
	ttOrOr

	// Bit
	ttBitAnd
	ttBitOr

	// Literals
	ttNil
	ttString
	ttNumber
	ttBoolean
	ttIdentifier

	// Keywords
	ttLet
	ttFunction
	ttReturn
	ttFor
	ttBreak
	ttIf
	ttElse
)

var keywords map[string]TokenType

func init() {
	keywords = map[string]TokenType{
		"let":      ttLet,
		"function": ttFunction,
		"return":   ttReturn,
		"for":      ttFor,
		"break":    ttBreak,
		"if":       ttIf,
		"else":     ttElse,
		"nil":      ttNil,
		"true":     ttBoolean,
		"false":    ttBoolean,
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
		ttDot:          ".",
		ttComma:        ",",
		ttSemicolon:    ";",
		ttColon:        ":",
		ttLambda:       "=>",

		ttAssign: "=",

		ttAddition:     "+",
		ttSubstraction: "-",
		ttMultiply:     "*",
		ttDivision:     "/",
		ttPercent:      "%",

		ttGreaterThan:        ">",
		ttGreaterThanOrEqual: ">=",
		ttEqual:              "==",
		ttNotEqual:           "!=",
		ttLessThan:           "<",
		ttLessThanOrEqual:    "<=",

		ttNot:    "!",
		ttAndAnd: "&&",
		ttOrOr:   "||",

		ttNil: "nil",

		ttLet:      "let",
		ttFunction: "function",
		ttReturn:   "return",
		ttFor:      "for",
		ttBreak:    "break",
		ttIf:       "if",
		ttElse:     "else",
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
		ret += fmt.Sprintf(" (line:%d col:%d)", t.line, t.col)
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
		} else if ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' {
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
		case ';':
			return Token{typ: ttSemicolon}
		case '+':
			return Token{typ: ttAddition}
		case '-':
			return Token{typ: ttSubstraction}
		case '*':
			return Token{typ: ttMultiply}
		case '/':
			c := t.read()
			if c == '/' {
				for {
					c = t.read()
					if c == '\n' || c == 0 {
						if c == '\n' {
							t.line++
						}
						break
					}
				}
				continue
			} else {
				t.unread()
				return Token{typ: ttDivision}
			}
		case '%':
			return Token{typ: ttPercent}
		case '=':
			next := t.read()
			switch next {
			case '=':
				return Token{typ: ttEqual}
			case '>':
				return Token{typ: ttLambda}
			default:
				t.unread()
				return Token{typ: ttAssign}
			}
		case '>':
			return t.iif('=', ttGreaterThanOrEqual, ttGreaterThan)
		case '<':
			return t.iif('=', ttLessThanOrEqual, ttLessThan)
		case '!':
			return Token{typ: ttNot}
		case '&':
			return t.iif('&', ttAndAnd, ttBitAnd)
		case '|':
			return t.iif('|', ttOrOr, ttBitOr)
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
