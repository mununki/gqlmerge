package lib

import (
	"bytes"
	"fmt"
	"os"
	"text/scanner"
)

type Lexer struct {
	sc     *scanner.Scanner
	next   rune
	buffer bytes.Buffer
}

func NewLexer(file *os.File) *Lexer {
	sc := scanner.Scanner{
		Mode: scanner.ScanIdents | scanner.ScanInts | scanner.ScanFloats | scanner.ScanStrings,
	}

	sc.Init(file)
	sc.Filename = file.Name()

	return &Lexer{sc: &sc}
}

func (l *Lexer) ConsumeWhitespace() {
	l.buffer.Reset()
	for {
		l.next = l.sc.Scan()

		if l.next == '#' {
			l.ConsumeComment()
			continue
		}

		if l.next == scanner.String {
			if l.sc.Peek() != '"' {
				break
			}
			l.ConsumeMultiLineComment()
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
		l.buffer.WriteRune(next)
	}
}

func (l *Lexer) ConsumeMultiLineComment() {
	for {
		next := l.sc.Scan()
		if next == scanner.EOF {
			break
		}
		if next == scanner.String {
			next := l.sc.Next()
			if next == '"' {
				break
			}
		}
		l.buffer.WriteRune(next)
	}
}

func (l *Lexer) GetBuffer() string {
	return l.buffer.String()
}

func (l *Lexer) ConsumeIdent() string {
	name := l.sc.TokenText()
	l.ConsumeToken(scanner.Ident)
	return name
}

func (l *Lexer) ConsumeString() string {
	str := l.sc.TokenText()
	l.ConsumeWhitespace()
	return str
}

func (l *Lexer) ConsumeToken(expected rune) {
	if l.next != expected {
		msg := fmt.Sprintf(
			// doesn't quote expected because scanner.TokenString
			// do it itself
			`%s:%d:%d: unexpected "%s", expected %s`,
			l.sc.Filename,
			l.sc.Line,
			l.sc.Column,
			l.sc.TokenText(),
			scanner.TokenString(expected),
		)
		panic(msg)
	}
	l.ConsumeWhitespace()
}

func (l *Lexer) Peek() rune {
	return l.next
}
