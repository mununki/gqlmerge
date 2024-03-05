package lib

import (
	"bytes"
	"fmt"
	"io"
	"math/big"
	"os"
	"unicode"
)

type tokenType int

const (
	tokError             tokenType = iota
	tokEOF                         // EOF
	tokNewLine                     // \n \r
	tokLParen                      // (
	tokRParen                      // )
	tokLBrace                      // {
	tokRBrace                      // }
	tokLBracket                    // [
	tokRBracket                    // ]
	tokBar                         // |
	tokQuote                       // '
	tokSingleLineComment           // # comment
	tokBlockString                 // """ multi line comment
	tokColon                       // :
	tokComma                       // ,
	tokDot                         // .
	tokThreeDot                    // ...
	tokEqual                       // =
	tokBang                        // !
	tokAt                          // @
	tokAmpersand                   // &
	tokDollar                      // $
	tokPlus                        // +
	tokMinus                       // -
	tokMul                         // *
	tokDiv                         // /
	tokNumber                      // number
	tokIdent                       // ident
	tokString                      // "..."
	tokOn                          // on
	tokRepeatable                  // repeatable
	tokDirective                   // directive
	tokType                        // type
	tokInterface                   // interface
	tokInput                       // input
	tokImplements                  // implements
	tokEnum                        // enum
	tokScalar                      // scalar
	tokUnion                       // union
	tokExtend                      // extend
	tokSchema                      // schema
)

func (typ tokenType) String() string {
	switch typ {
	case tokError:
		return "unknown error"
	case tokEOF:
		return "EOF"
	case tokNewLine:
		return "\n"
	case tokLParen:
		return "("
	case tokRParen:
		return ")"
	case tokLBrace:
		return "{"
	case tokRBrace:
		return "}"
	case tokLBracket:
		return "["
	case tokRBracket:
		return "]"
	case tokBar:
		return "|"
	case tokQuote:
		return "'"
	case tokSingleLineComment:
		return "comment"
	case tokBlockString:
		return "multi line comment"
	case tokColon:
		return ":"
	case tokComma:
		return ","
	case tokDot:
		return "."
	case tokThreeDot:
		return "..."
	case tokEqual:
		return "="
	case tokBang:
		return "!"
	case tokAt:
		return "@"
	case tokAmpersand:
		return "&"
	case tokDollar:
		return "$"
	case tokPlus:
		return "+"
	case tokMinus:
		return "-"
	case tokMul:
		return "*"
	case tokDiv:
		return "/"
	case tokNumber:
		return "number"
	case tokIdent:
		return "ident"
	case tokString:
		return "string"
	case tokOn:
		return "on"
	case tokRepeatable:
		return "repeatable"
	case tokDirective:
		return "directive"
	case tokType:
		return "type"
	case tokInterface:
		return "interface"
	case tokInput:
		return "input"
	case tokImplements:
		return "implements"
	case tokEnum:
		return "enum"
	case tokScalar:
		return "scalar"
	case tokUnion:
		return "union"
	case tokExtend:
		return "extend"
	case tokSchema:
		return "schema"
	default:
		return "unknown token"
	}
}

const EofRune rune = -1

type token struct {
	typ  tokenType
	text *string
	num  *big.Int
}

func (t *token) String() string {
	if t.typ == tokNumber {
		return fmt.Sprint(t.num)
	}
	return *t.text
}

type lexer struct {
	filename string
	line     int
	col      int
	rd       io.RuneReader
	peeking  bool
	peekRune rune
	last     rune
	buf      bytes.Buffer
}

func newLexer(rd io.RuneReader, filename string) *lexer {
	return &lexer{
		filename: filename,
		line:     1,
		col:      0,
		rd:       rd,
	}
}

var tokens = make(map[string]*token)

func mkToken(typ tokenType, text string) *token {
	if typ == tokNumber {
		var z big.Int
		num, ok := z.SetString(text, 0)
		if !ok {
			errorf("bad number syntax: %s", text)
		}
		return number(num)
	}
	tok := tokens[text]
	if tok == nil {
		tok = &token{typ: typ, text: &text}
	}
	return tok
}

func number(num *big.Int) *token {
	return &token{typ: tokNumber, num: num}
}

func (l *lexer) skipSpace() rune {
	for {
		r := l.read()
		if !isSpace(r) {
			l.back(r)
			return r
		}
	}
}

func (l *lexer) read() rune {
	if l.peeking {
		l.peeking = false
		return l.peekRune
	}
	return l.nextRune()
}

func (l *lexer) nextRune() rune {
	r, _, err := l.rd.ReadRune()
	if err != nil {
		if err != io.EOF {
			fmt.Fprintln(os.Stderr)
		}
		r = EofRune
	}
	l.last = r
	l.counter(r)
	return r
}

func (l *lexer) peek() rune {
	if l.peeking {
		return l.peekRune
	}
	r := l.read()
	l.peeking = true
	l.peekRune = r
	return r
}

func (l *lexer) back(r rune) {
	l.peeking = true
	l.peekRune = r
}

func (l *lexer) next() *token {
	for {
		r := l.read()
		switch {
		case isSpace(r):
			l.skipSpace()
		case r == EofRune:
			return mkToken(tokEOF, "EOF")
		case r == '\n' || r == '\r':
			return mkToken(tokNewLine, "\n")
		case r == '(':
			return mkToken(tokLParen, "(")
		case r == ')':
			return mkToken(tokRParen, ")")
		case r == '{':
			return mkToken(tokLBrace, "{")
		case r == '}':
			return mkToken(tokRBrace, "}")
		case r == '[':
			return mkToken(tokLBracket, "[")
		case r == ']':
			return mkToken(tokRBracket, "]")
		case r == '|':
			return mkToken(tokBar, "|")
		case r == '#':
			return l.comment(r)
		case r == '"':
			return l.stringOrBlockString(r)
		case r == '\'':
			return mkToken(tokQuote, "'")
		case r == ':':
			return mkToken(tokColon, ":")
		case r == ',':
			return mkToken(tokComma, ",")
		case r == '.':
			return mkToken(tokDot, ".")
		case r == '=':
			return mkToken(tokEqual, "=")
		case r == '!':
			return mkToken(tokBang, "!")
		case r == '@':
			return mkToken(tokAt, "@")
		case r == '&':
			return mkToken(tokAmpersand, "&")
		case r == '$':
			return mkToken(tokDollar, "$")
		case r == '*':
			return mkToken(tokMul, "*")
		case r == '/':
			return mkToken(tokDiv, "/")
		case r == '-' || r == '+':
			if !isNumber(l.peek()) {
				if r == '-' {
					return mkToken(tokMinus, "-")
				} else {
					return mkToken(tokPlus, "+")
				}
			}
			fallthrough
		case isNumber(r):
			return l.number(r)
		case r == '_' || unicode.IsLetter(r):
			s := l.alphanum(r)
			switch s {
			case "on":
				return mkToken(tokOn, "on")
			case "repeatable":
				return mkToken(tokRepeatable, "repeatable")
			case "directive":
				return mkToken(tokDirective, "directive")
			case "type":
				return mkToken(tokType, "type")
			case "interface":
				return mkToken(tokInterface, "interface")
			case "input":
				return mkToken(tokInput, "input")
			case "implements":
				return mkToken(tokImplements, "implements")
			case "enum":
				return mkToken(tokEnum, "enum")
			case "scalar":
				return mkToken(tokScalar, "scalar")
			case "union":
				return mkToken(tokUnion, "union")
			case "extend":
				return mkToken(tokExtend, "extend")
			case "schema":
				return mkToken(tokSchema, "schema")
			default:
				return mkToken(tokIdent, string(s))
			}
		default:
			return mkToken(tokIdent, string(r))
		}
	}
}

func isSpace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}

func isNumber(r rune) bool {
	// float?
	return '0' <= r && r <= '9'
}

func isAlphanum(r rune) bool {
	return r == '_' || isNumber(r) || unicode.IsLetter(r)
}

func (l *lexer) accum(r rune, valid func(rune) bool) {
	l.buf.Reset()
	for {
		l.buf.WriteRune(r)
		r = l.read()
		if r == EofRune {
			return
		}
		if !valid(r) {
			l.back(r)
			return
		}
	}
}

func (l *lexer) number(r rune) *token {
	l.accum(r, isNumber)
	l.endToken()
	return mkToken(tokNumber, l.buf.String())
}

func (l *lexer) alphanum(r rune) string {
	l.accum(r, isAlphanum)
	l.endToken()
	return l.buf.String()
}

func (l *lexer) comment(r rune) *token {
	l.buf.Reset()
	for {
		l.buf.WriteRune(r)
		r = l.read()
		if r == '\n' || r == EofRune {
			return mkToken(tokSingleLineComment, l.buf.String())
		}
	}
}

func (l *lexer) stringOrBlockString(r rune) *token {
	l.buf.Reset()
	l.buf.WriteRune(r)
	r = l.read()
	l.buf.WriteRune(r)
	if r == '"' {
		if l.peek() == '"' {
			// multi line comment
			r = l.read()
			l.buf.WriteRune(r)
			for {
				r = l.read()
				l.buf.WriteRune(r)
				if r == EofRune {
					errorf("%s:%d:%d: Block string literal terminated", l.filename, l.line, l.col)
				}
				if r == '"' {
					r = l.read()
					l.buf.WriteRune(r)
					if r == '"' && l.peek() == '"' {
						r = l.read()
						l.buf.WriteRune(r)
						return mkToken(tokBlockString, l.buf.String())
					}
					errorf("%s:%d:%d: Block string literal terminated", l.filename, l.line, l.col)
				}
			}
		}
		// empty string
		return mkToken(tokString, l.buf.String())
	}
	// normal string
	for {
		r = l.read()
		l.buf.WriteRune(r)
		if r == '\n' || r == EofRune {
			errorf("%s:%d:%d: String literal terminated", l.filename, l.line, l.col)
		}
		if r == '"' {
			return mkToken(tokString, l.buf.String())
		}
	}
}

func (l *lexer) endToken() {
	if r := l.peek(); isAlphanum(r) || !isSpace(r) && r != '(' && r != ')' && r != '[' && r != ']' && r != '{' && r != '}' && r != ':' && r != '!' && r != ',' && r != EofRune {
		errorf("%s:%d:%d: invalid token after %s", l.filename, l.line, l.col, &l.buf)
	}
}

func (l *lexer) counter(r rune) {
	if r == EofRune {
		return
	}
	if r == '\n' || r == '\r' {
		l.line++
		l.col = 0
		return
	}
	l.col++
}

func (l *lexer) consumeToken(expected tokenType) {
	tok := l.next()
	if tok.typ != expected {
		errorf(`%s:%d:%d: unexpected "%s", expected %s`, l.filename, l.line, l.col, tok.String(), expected.String())
	}
	l.skipSpace()
}

func (l *lexer) consumeIdent(includings ...tokenType) (*token, *[]string) {
	comments := []string{}
	for {
		tok := l.next()
		if tok.typ == tokString || tok.typ == tokSingleLineComment || tok.typ == tokBlockString {
			comments = append(comments, tok.String())
			continue
		}

		isIncluded := false
		for _, incl := range includings {
			if tok.typ == incl {
				isIncluded = true
				break
			}
		}

		if tok.typ != tokIdent && !isIncluded {
			errorf(`%s:%d:%d: unexpected "%s"`, l.filename, l.line, l.col, tok.String())
		}
		l.skipSpace()
		return tok, &comments
	}
}
