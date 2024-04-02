package lib

import (
	"fmt"
	"io"
)

type (
	Error string
	EOF   string
)

func errorf(format string, args ...interface{}) {
	panic(Error(fmt.Sprintf(format, args...)))
}

type Parser struct {
	lex *lexer
	buf []*token
}

func NewParser(r io.RuneReader, filename string) *Parser {
	return &Parser{
		lex: newLexer(r, filename),
	}
}

func (p *Parser) bufString() *[]string {
	ss := []string{}
	for _, t := range p.buf {
		ss = append(ss, t.String())
	}
	p.buf = []*token{}
	return &ss
}

func (p *Parser) parseArgs() []*Arg {
	args := []*Arg{}
	for p.lex.peek() == '(' {
		p.lex.consumeToken(tokLParen)
		for p.lex.peek() != ')' {
			arg := Arg{}
			name, comments := p.lex.consumeIdent(tokInput, tokType)
			arg.Name = name.String()
			arg.Descriptions = comments
			p.lex.consumeToken(tokColon)

			if p.lex.peek() == '[' {
				arg.IsList = true
				p.lex.consumeToken(tokLBracket)
				typ, _ := p.lex.consumeIdent()
				arg.Type = typ.String()
				if p.lex.peek() == '!' {
					arg.Null = false
					p.lex.consumeToken(tokBang)
				} else {
					arg.Null = true
				}
				p.lex.consumeToken(tokRBracket)

				if p.lex.peek() == '!' {
					arg.IsListNull = false
					p.lex.consumeToken(tokBang)
				} else {
					arg.IsListNull = true
				}

				if p.lex.peek() == '=' {
					p.lex.consumeToken(tokEqual)
					defaultValues := []string{}
					for p.lex.peek() == '[' {
						p.lex.consumeIdent(tokLBracket)
						for p.lex.peek() != ']' {
							tex, _ := p.lex.consumeIdentInclString(tokNumber)
							te := tex.String()
							defaultValues = append(defaultValues, te)
							if p.lex.peek() == ',' {
								p.lex.consumeToken(tokComma)
							}
						}
						p.lex.consumeIdent(tokRBracket)
						arg.DefaultValues = &defaultValues
					}
				}
			} else {
				typ, _ := p.lex.consumeIdent()
				arg.Type = typ.String()

				if p.lex.peek() == '!' {
					arg.Null = false
					p.lex.consumeToken(tokBang)
				} else {
					arg.Null = true
				}

				if p.lex.peek() == '=' {
					p.lex.consumeToken(tokEqual)
					tex, _ := p.lex.consumeIdentInclString(tokNumber)
					te := tex.String()
					defaultValues := []string{}
					defaultValues = append(defaultValues, te)
					arg.DefaultValues = &defaultValues
				}
			}
			arg.Directives = p.parseDirectives()

			args = append(args, &arg)

			if p.lex.peek() == ',' {
				p.lex.consumeToken(tokComma)
			}
		}
		p.lex.consumeToken(tokRParen)
	}
	return args
}

func (p *Parser) parseDirectives() []*Directive {
	ds := []*Directive{}

	for p.lex.peek() == '@' {
		p.lex.consumeToken(tokAt)
		d := Directive{}
		name, comments := p.lex.consumeIdent()
		d.Name = name.String()
		d.Descriptions = comments

		for p.lex.peek() == '(' {
			p.lex.consumeToken(tokLParen)
			for p.lex.peek() != ')' {
				da := DirectiveArg{}
				name, comments = p.lex.consumeIdent()
				da.Name = name.String()
				da.Descriptions = comments
				p.lex.consumeToken(tokColon)

				if p.lex.peek() == '[' {
					da.IsList = true
					p.lex.consumeToken(tokLBracket)
					da.Value = p.parseList()
					p.lex.consumeToken(tokRBracket)
				} else {
					da.IsList = false
					tok := p.lex.next()
					switch tok.typ {
					case tokString, tokIdent, tokNumber:
						da.Value = append(da.Value, tok.String())
					}
				}

				d.DirectiveArgs = append(d.DirectiveArgs, &da)

				if p.lex.peek() == ',' {
					p.lex.consumeToken(tokComma)
				}
			}
			p.lex.consumeToken(tokRParen)
		}
		ds = append(ds, &d)
	}
	return ds
}

func (p *Parser) parseList() []string {
	ss := []string{}
	for p.lex.peek() != ']' {
		tok := p.lex.next()
		switch tok.typ {
		case tokString, tokIdent, tokNumber:
			ss = append(ss, tok.String())
		}
		if p.lex.peek() == ',' {
			p.lex.consumeToken(tokComma)
		}
	}
	return ss
}

func (p *Parser) parseSingleLineComment() *string {
	if p.lex.peek() == '#' {
		tok := p.lex.next().String()
		p.lex.skipSpace()
		return &tok
	}
	p.lex.skipSpace()
	return nil
}
