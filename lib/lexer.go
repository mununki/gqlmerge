package lib

import (
	"fmt"
	"strings"
	"text/scanner"
)

type Lexer struct {
	sc   *scanner.Scanner
	next rune
}

func NewLexer(s string) *Lexer {
	sc := scanner.Scanner{
		Mode: scanner.ScanIdents | scanner.ScanInts | scanner.ScanFloats | scanner.ScanStrings,
	}

	sc.Init(strings.NewReader(s))

	return &Lexer{sc: &sc}
}

func (l *Lexer) ConsumeWhitespace() {
	for {
		l.next = l.sc.Scan()

		if l.next == ',' {
			continue
		}

		if l.next == '#' {
			l.ConsumeComment()
			continue
		}

		break
	}
}

func (l *Lexer) ConsumeComment() {
	for {
		next := l.sc.Next()
		if next == '\r' || next == '\n' || next == scanner.EOF {
			break
		}
	}
}

func (l *Lexer) ConsumeIdent() string {
	name := l.sc.TokenText()
	l.ConsumeToken(scanner.Ident)
	return name
}

func (l *Lexer) ConsumeToken(expected rune) {
	if l.next != expected {
		fmt.Printf("syntax error: unexpected %s, expected %s", l.sc.TokenText(), scanner.TokenString(expected))
		panic("syntax error found")
	}
	l.ConsumeWhitespace()
}

func (l *Lexer) Peek() rune {
	return l.next
}
