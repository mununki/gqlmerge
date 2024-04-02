package lib

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"sync"
)

// GetSchema is to parse ./schema/**/*.graphql
func (sc *Schema) ReadSchema(path string) {
	// FIX: is there any way to use a relative path?
	// currently, it works only with absolute path
	// in case of using a relative path such as '../schema', it spits out an error
	// the error says invalid memory or nil pointer deference.
	err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if p == "" {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		ext := filepath.Ext(p)
		if ext != ".graphql" && ext != ".gql" {
			return nil
		}

		file, err := os.Open(p)
		if err != nil {
			fmt.Printf("[Error] There is an error to open %s", p)
			return err
		}

		// TODO: split and get a only filename and print it to user
		// needs to handle in case of OS (windows / unix compatibles)
		sc.Files = append(sc.Files, file)

		return nil
	})
	if err != nil {
		panic(err)
	}

	rel, err := GetRelPath(path)
	if err == nil {
		// if failed print absolute path
		path = *rel
	}

	if len(sc.Files) > 0 {
		fmt.Printf("🎉 [%s] Total %d GraphQL files found!\n", path, len(sc.Files))
	}
}

func (s *Schema) Parse(p *Parser) {
	isExtended := false
	for {
		tok := p.lex.next()
		if tok.typ == tokEOF {
			break
		}
		switch tok.typ {
		case tokSchema:
			sd := SchemaDefinition{}
			sd.Filename = p.lex.filename
			sd.Line = p.lex.line
			sd.Column = p.lex.col
			sd.Descriptions = p.bufString()
			p.lex.consumeToken(tokLBrace)
			for p.lex.peek() != '}' {
				op := p.lex.next()
				p.lex.consumeToken(tokColon)
				switch op.String() {
				case "query":
					q := p.lex.next().String()
					sd.Query = &q
				case "mutation":
					m := p.lex.next().String()
					sd.Mutation = &m
				case "subscription":
					s := p.lex.next().String()
					sd.Subscription = &s
				default:
					errorf(`%s:%d:%d: unexpected "%s", one of operation types expected`, p.lex.filename, p.lex.line, p.lex.col, op.String())
				}
				p.lex.skipSpace()
			}
			p.lex.consumeToken(tokRBrace)
			s.SchemaDefinitions = append(s.SchemaDefinitions, &sd)

		case tokString:
			// description
			p.buf = append(p.buf, tok)

		case tokSingleLineComment:
			p.buf = append(p.buf, tok)

		case tokBlockString:
			p.buf = append(p.buf, tok)

		case tokDirective:
			d := DirectiveDefinition{}
			d.Filename = p.lex.filename
			d.Line = p.lex.line
			d.Column = p.lex.col
			d.Descriptions = p.bufString()
			p.lex.consumeToken(tokAt)
			name, _ := p.lex.consumeIdent()
			d.Name = name.String()
			d.Args = p.parseArgs()
			t := p.lex.next()
			if t.typ == tokRepeatable {
				d.Repeatable = true
			} else if t.typ == tokOn {
				d.Repeatable = false
				ls := []string{}
				for p.lex.peek() != EofRune {
					l, _ := p.lex.consumeIdent()
					ls = append(ls, l.String())
					if p.lex.peek() == '|' {
						p.lex.consumeToken(tokBar)
					} else {
						break
					}
				}
				d.Locations = ls
			}
			s.DirectiveDefinitions = append(s.DirectiveDefinitions, &d)

		case tokExtend:
			isExtended = true

		case tokScalar:
			c := Scalar{}
			c.Filename = p.lex.filename
			c.Line = p.lex.line
			c.Column = p.lex.col
			c.Descriptions = p.bufString()
			name, _ := p.lex.consumeIdent()
			c.Name = name.String()
			c.Directives = p.parseDirectives()
			sc := p.parseSingleLineComment()
			if c.Comments != nil && sc != nil {
				cs := append(*c.Comments, *sc)
				c.Comments = &cs
			} else if c.Comments == nil && sc != nil {
				cs := []string{*sc}
				c.Comments = &cs
			}
			s.Scalars = append(s.Scalars, &c)

		case tokEnum:
			e := Enum{}
			e.Filename = p.lex.filename
			e.Line = p.lex.line
			e.Column = p.lex.col
			e.Descriptions = p.bufString()
			name, _ := p.lex.consumeIdent()
			e.Name = name.String()
			e.Directives = p.parseDirectives()
			p.lex.consumeToken(tokLBrace)
			for p.lex.peek() != '}' {
				ev := EnumValue{}
				name, comments := p.lex.consumeIdent()
				ev.Name = name.String()
				ev.Descriptions = comments
				ev.Directives = p.parseDirectives()
				sc := p.parseSingleLineComment()
				if ev.Comments != nil && sc != nil {
					cs := append(*ev.Comments, *sc)
					ev.Comments = &cs
				} else if ev.Comments == nil && sc != nil {
					cs := []string{*sc}
					ev.Comments = &cs
				}
				e.EnumValues = append(e.EnumValues, ev)
			}
			p.lex.consumeToken(tokRBrace)
			s.Enums = append(s.Enums, &e)

		case tokInterface:
			i := Interface{}
			i.Filename = p.lex.filename
			i.Line = p.lex.line
			i.Column = p.lex.col
			name, _ := p.lex.consumeIdent()
			i.Name = name.String()
			i.Descriptions = p.bufString()
			i.Directives = p.parseDirectives()

			p.lex.consumeToken(tokLBrace)
			for p.lex.peek() != '}' {
				fd := Field{}
				fd.Filename = p.lex.filename
				fd.Line = p.lex.line
				fd.Column = p.lex.col
				name, comments := p.lex.consumeIdent(tokInput, tokType)
				fd.Name = name.String()
				fd.Descriptions = comments

				fd.Args = p.parseArgs()

				p.lex.consumeToken(tokColon)

				if p.lex.peek() == '[' {
					fd.IsList = true
					p.lex.consumeToken(tokLBracket)
					name, _ = p.lex.consumeIdent()
					fd.Type = name.String()
					if p.lex.peek() == '!' {
						fd.Null = false
						p.lex.consumeToken(tokBang)
					} else {
						fd.Null = true
					}
					p.lex.consumeToken(tokRBracket)
					if p.lex.peek() == '!' {
						fd.IsListNull = false
						p.lex.consumeToken(tokBang)
					} else {
						fd.IsListNull = true
					}
				} else {
					fd.IsList = false
					fd.IsListNull = false
					name, _ = p.lex.consumeIdent()
					fd.Type = name.String()
					if p.lex.peek() == '!' {
						fd.Null = false
						p.lex.consumeToken(tokBang)
					} else {
						fd.Null = true
					}
				}

				fd.Directives = p.parseDirectives()

				sc := p.parseSingleLineComment()
				if fd.Comments != nil && sc != nil {
					cs := append(*fd.Comments, *sc)
					fd.Comments = &cs
				} else if fd.Comments == nil && sc != nil {
					cs := []string{*sc}
					fd.Comments = &cs
				}

				i.Fields = append(i.Fields, &fd)
			}

			s.Interfaces = append(s.Interfaces, &i)
			p.lex.consumeToken(tokRBrace)

		case tokUnion:
			u := Union{}
			u.Filename = p.lex.filename
			u.Line = p.lex.line
			u.Column = p.lex.col
			u.Descriptions = p.bufString()
			name, _ := p.lex.consumeIdent()
			u.Name = name.String()
			u.Directives = p.parseDirectives()
			p.lex.consumeToken(tokEqual)
			for p.lex.peek() != '\n' || p.lex.peek() != '\r' || p.lex.peek() != EofRune {
				name, _ = p.lex.consumeIdent()
				// FIXME comments?
				u.Types = append(u.Types, name.String())
				if p.lex.peek() == '|' {
					p.lex.consumeToken(tokBar)
				} else {
					break
				}
			}
			s.Unions = append(s.Unions, &u)

		case tokInput:
			i := Input{}
			i.Filename = p.lex.filename
			i.Line = p.lex.line
			i.Column = p.lex.col
			name, _ := p.lex.consumeIdent()
			i.Name = name.String()
			i.Descriptions = p.bufString()

			p.lex.consumeToken(tokLBrace)

			for p.lex.peek() != '}' {
				fd := Field{}
				fd.Filename = p.lex.filename
				fd.Line = p.lex.line
				fd.Column = p.lex.col
				name, comments := p.lex.consumeIdent(tokInput, tokType)
				fd.Name = name.String()
				fd.Descriptions = comments
				p.lex.consumeToken(tokColon)

				if p.lex.peek() == '[' {
					fd.IsList = true
					p.lex.consumeToken(tokLBracket)
					name, _ = p.lex.consumeIdent()
					fd.Type = name.String()
					if p.lex.peek() == '!' {
						fd.Null = false
						p.lex.consumeToken(tokBang)
					} else {
						fd.Null = true
					}
					p.lex.consumeToken(tokRBracket)
					if p.lex.peek() == '!' {
						fd.IsListNull = false
						p.lex.consumeToken(tokBang)
					} else {
						fd.IsListNull = true
					}
				} else {
					fd.IsList = false
					fd.IsListNull = false
					name, _ = p.lex.consumeIdent()
					fd.Type = name.String()
					if p.lex.peek() == '!' {
						fd.Null = false
						p.lex.consumeToken(tokBang)
					} else {
						fd.Null = true
					}

					if p.lex.peek() == '=' {
						p.lex.consumeToken(tokEqual)
						tex, _ := p.lex.consumeIdentInclString(tokNumber)
						te := tex.String()
						fd.DefaultValue = &te
					}
				}

				fd.Directives = p.parseDirectives()

				sc := p.parseSingleLineComment()
				if fd.Comments != nil && sc != nil {
					cs := append(*fd.Comments, *sc)
					fd.Comments = &cs
				} else if fd.Comments == nil && sc != nil {
					cs := []string{*sc}
					fd.Comments = &cs
				}

				i.Fields = append(i.Fields, &fd)
			}

			s.Inputs = append(s.Inputs, &i)
			p.lex.consumeToken(tokRBrace)

		case tokType:
			t := Type{}
			t.Extend = isExtended
			isExtended = false
			t.Filename = p.lex.filename
			t.Line = p.lex.line
			t.Column = p.lex.col
			t.Descriptions = p.bufString()
			name, _ := p.lex.consumeIdent()
			t.Name = name.String()
			t.Directives = p.parseDirectives()

			next := p.lex.next()
			switch next.typ {
			case tokImplements:
				if len(t.Directives) > 0 {
					errorf(`%s:%d:%d: directives cann't be placed in front of implements`, p.lex.filename, p.lex.line, p.lex.col)
				}
				t.Impl = true
				name, _ := p.lex.consumeIdent()
				t.ImplTypes = append(t.ImplTypes, name.String())
				for p.lex.peek() == '&' {
					p.lex.consumeToken(tokAmpersand)
					name, _ = p.lex.consumeIdent()
					t.ImplTypes = append(t.ImplTypes, name.String())
				}
				t.Directives = p.parseDirectives()
				p.lex.consumeToken(tokLBrace)
				fallthrough
			case tokLBrace:
				for p.lex.peek() != '}' {
					fd := Field{}
					fd.Filename = p.lex.filename
					fd.Line = p.lex.line
					fd.Column = p.lex.col
					name, comments := p.lex.consumeIdent(tokInput, tokType)
					fd.Name = name.String()
					fd.Descriptions = comments

					fd.Args = p.parseArgs()

					p.lex.consumeToken(tokColon)

					if p.lex.peek() == '[' {
						fd.IsList = true
						p.lex.consumeToken(tokLBracket)
						name, _ = p.lex.consumeIdent()
						fd.Type = name.String()
						if p.lex.peek() == '!' {
							fd.Null = false
							p.lex.consumeToken(tokBang)
						} else {
							fd.Null = true
						}
						p.lex.consumeToken(tokRBracket)
						if p.lex.peek() == '!' {
							fd.IsListNull = false
							p.lex.consumeToken(tokBang)
						} else {
							fd.IsListNull = true
						}
					} else {
						fd.IsList = false
						fd.IsListNull = false
						name, _ = p.lex.consumeIdent()
						fd.Type = name.String()
						if p.lex.peek() == '!' {
							fd.Null = false
							p.lex.consumeToken(tokBang)
						} else {
							fd.Null = true
						}
					}

					fd.Directives = p.parseDirectives()

					sc := p.parseSingleLineComment()
					if fd.Comments != nil && sc != nil {
						cs := append(*fd.Comments, *sc)
						fd.Comments = &cs
					} else if fd.Comments == nil && sc != nil {
						cs := []string{*sc}
						fd.Comments = &cs
					}

					t.Fields = append(t.Fields, &fd)
				}

				s.Types = append(s.Types, &t)
				p.lex.consumeToken(tokRBrace)
			default:
				errorf(`%s:%d:%d: unexpected "%s", expected implments or {`, p.lex.filename, p.lex.line, p.lex.col, next.String())
			}

		}
	}
}

func (s *Schema) mergeSchemaDefinition(wg *sync.WaitGroup) {
	defer wg.Done()
	sd := SchemaDefinition{}
	for i, v := range s.SchemaDefinitions {
		if i == 0 {
			sd = *v
			continue
		}
		if sd.Query == nil {
			sd.Query = v.Query
		} else if v.Query != nil && *sd.Query != *v.Query {
			rel1, err := GetRelPath(sd.Filename)
			if err != nil {
				panic(err)
			}
			rel2, err := GetRelPath(v.Filename)
			if err != nil {
				panic(err)
			}
			errorf("Duplicated Directive Definitions: %s(%s:%v:%v) and (%s:%v:%v)", *sd.Query, *rel1, sd.Line, sd.Column, *rel2, v.Line, v.Column)
		}
		if sd.Mutation == nil {
			sd.Mutation = v.Mutation
		} else if v.Mutation != nil && *sd.Mutation != *v.Mutation {
			rel1, err := GetRelPath(sd.Filename)
			if err != nil {
				panic(err)
			}
			rel2, err := GetRelPath(v.Filename)
			if err != nil {
				panic(err)
			}
			errorf("Duplicated Directive Definitions: %s(%s:%v:%v) and (%s:%v:%v)", *sd.Mutation, *rel1, sd.Line, sd.Column, *rel2, v.Line, v.Column)
		}
		if sd.Subscription == nil {
			sd.Subscription = v.Subscription
		} else if v.Subscription != nil && *sd.Subscription != *v.Subscription {
			rel1, err := GetRelPath(sd.Filename)
			if err != nil {
				panic(err)
			}
			rel2, err := GetRelPath(v.Filename)
			if err != nil {
				panic(err)
			}
			errorf("Duplicated Directive Definitions: %s(%s:%v:%v) and (%s:%v:%v)", *sd.Subscription, *rel1, sd.Line, sd.Column, *rel2, v.Line, v.Column)
		}

		sd.Descriptions = mergeStrings(sd.Descriptions, v.Descriptions)
	}
	sds := []*SchemaDefinition{&sd}
	s.SchemaDefinitions = sds
}

func (s *Schema) UniqueDirectiveDefinition(wg *sync.WaitGroup) {
	defer wg.Done()
	j := 0
	seen := make(map[string]struct{}, len(s.DirectiveDefinitions))
	for _, v := range s.DirectiveDefinitions {
		if _, ok := seen[v.Name]; ok {
			for i := 0; i < j; i++ {
				if s.DirectiveDefinitions[i].Name == v.Name {
					if IsEqualWithoutDescriptions(s.DirectiveDefinitions[i], v) {
						mergeDescriptionsAndComments(s.DirectiveDefinitions[i], v)
						break
					} else {
						rel1, err := GetRelPath(s.DirectiveDefinitions[i].Filename)
						if err != nil {
							panic(err)
						}
						rel2, err := GetRelPath(v.Filename)
						if err != nil {
							panic(err)
						}

						errorf("Duplicated Directive Definitions: %s(%s:%v:%v) and (%s:%v:%v)", s.DirectiveDefinitions[i].Name, *rel1, s.DirectiveDefinitions[i].Line, s.DirectiveDefinitions[i].Column, *rel2, v.Line, v.Column)
					}
				}
			}
			continue
		}
		seen[v.Name] = struct{}{}
		s.DirectiveDefinitions[j] = v
		j++
	}
	s.DirectiveDefinitions = s.DirectiveDefinitions[:j]
}

func (s *Schema) MergeTypeName(wg *sync.WaitGroup) {
	defer wg.Done()
	j := 0
	seen := make(map[string]struct{}, len(s.Types))
	sort.SliceStable(s.Types, func(i, j int) bool {
		return !s.Types[i].Extend && s.Types[j].Extend
	})
	for _, v := range s.Types {
		if _, ok := seen[v.Name]; ok {
			for i := 0; i < j; i++ {
				if s.Types[i].Name == v.Name {
					if v.Extend {
						s.Types[i].Fields = mergeFields(s.Types[i].Fields, v.Fields)
						s.Types[i].Directives = mergeDirectives(s.Types[i].Directives, v.Directives)
						break
					} else {
						if reflect.DeepEqual(s.Types[i].ImplTypes, v.ImplTypes) && IsEqualWithoutDescriptions(s.Types[i].Directives, v.Directives) {
							s.Types[i].Fields = mergeFields(s.Types[i].Fields, v.Fields)
							mergeDescriptionsAndComments(s.Types[i].Directives, v.Directives)
							break
						} else {

							rel1, err := GetRelPath(s.Types[i].Filename)
							if err != nil {
								panic(err)
							}
							rel2, err := GetRelPath(v.Filename)
							if err != nil {
								panic(err)
							}

							errorf("Duplicated Types: %s(%s:%v:%v) and (%s:%v:%v)", s.Types[i].Name, *rel1, s.Types[i].Line, s.Types[i].Column, *rel2, v.Line, v.Column)
						}
					}
				}
			}
			continue
		}
		seen[v.Name] = struct{}{}
		s.Types[j] = v
		j++
	}
	s.Types = s.Types[:j]
}

func (s *Schema) UniqueScalar(wg *sync.WaitGroup) {
	defer wg.Done()
	j := 0
	seen := make(map[string]struct{}, len(s.Scalars))
	for _, v := range s.Scalars {
		if _, ok := seen[v.Name]; ok {
			for i := 0; i < j; i++ {
				if s.Scalars[i].Name == v.Name {
					if IsEqualWithoutDescriptions(s.Scalars[i].Directives, v.Directives) {
						mergeDescriptionsAndComments(s.Scalars[i], v)
						break
					} else {
						rel1, err := GetRelPath(s.Scalars[i].Filename)
						if err != nil {
							panic(err)
						}
						rel2, err := GetRelPath(v.Filename)
						if err != nil {
							panic(err)
						}

						errorf("Duplicated Scalars: %s(%s:%v:%v) and (%s:%v:%v)", s.Scalars[i].Name, *rel1, s.Scalars[i].Line, s.Scalars[i].Column, *rel2, v.Line, v.Column)
					}
				}
			}
			continue
		}
		seen[v.Name] = struct{}{}
		s.Scalars[j] = v
		j++
	}
	s.Scalars = s.Scalars[:j]
}

func (s *Schema) UniqueEnum(wg *sync.WaitGroup) {
	defer wg.Done()
	j := 0
	seen := make(map[string]struct{}, len(s.Enums))
	for _, v := range s.Enums {
		if _, ok := seen[v.Name]; ok {
			for i := 0; i < j; i++ {
				if s.Enums[i].Name == v.Name {
					if IsEqualWithoutDescriptions(s.Enums[i].Directives, v.Directives) && IsEqualWithoutDescriptions(s.Enums[i].EnumValues, v.EnumValues) {
						mergeDescriptionsAndComments(s.Enums[i], v)
						break
					} else {

						rel1, err := GetRelPath(s.Enums[i].Filename)
						if err != nil {
							panic(err)
						}
						rel2, err := GetRelPath(v.Filename)
						if err != nil {
							panic(err)
						}

						errorf("Duplicated Enums: %s(%s:%v:%v) and (%s:%v:%v)", s.Enums[i].Name, *rel1, s.Enums[i].Line, s.Enums[i].Column, *rel2, v.Line, v.Column)
					}
				}
			}
			continue
		}
		seen[v.Name] = struct{}{}
		s.Enums[j] = v
		j++
	}
	s.Enums = s.Enums[:j]
}

func (s *Schema) UniqueInterface(wg *sync.WaitGroup) {
	defer wg.Done()
	j := 0
	seen := make(map[string]struct{}, len(s.Interfaces))
	for _, v := range s.Interfaces {
		if _, ok := seen[v.Name]; ok {
			for i := 0; i < j; i++ {
				if s.Interfaces[i].Name == v.Name {
					if IsEqualWithoutDescriptions(s.Interfaces[i].Directives, v.Directives) && IsEqualWithoutDescriptions(s.Interfaces[i].Fields, v.Fields) {
						mergeDescriptionsAndComments(s.Interfaces[i], v)
						break
					} else {

						rel1, err := GetRelPath(s.Interfaces[i].Filename)
						if err != nil {
							panic(err)
						}
						rel2, err := GetRelPath(v.Filename)
						if err != nil {
							panic(err)
						}

						errorf("Duplicated Interfaces: %s(%s:%v:%v) and (%s:%v:%v)", s.Interfaces[i].Name, *rel1, s.Interfaces[i].Line, s.Interfaces[i].Column, *rel2, v.Line, v.Column)
					}
				}
			}
			continue
		}
		seen[v.Name] = struct{}{}
		s.Interfaces[j] = v
		j++
	}
	s.Interfaces = s.Interfaces[:j]
}

func (s *Schema) UniqueUnion(wg *sync.WaitGroup) {
	defer wg.Done()
	j := 0
	seen := make(map[string]struct{}, len(s.Unions))
	for _, v := range s.Unions {
		if _, ok := seen[v.Name]; ok {
			for i := 0; i < j; i++ {
				if s.Unions[i].Name == v.Name {
					if IsEqualWithoutDescriptions(s.Unions[i].Directives, v.Directives) && IsEqualWithoutDescriptions(s.Unions[i].Types, v.Types) {
						mergeDescriptionsAndComments(s.Unions[i], v)
						break
					} else {

						rel1, err := GetRelPath(s.Unions[i].Filename)
						if err != nil {
							panic(err)
						}
						rel2, err := GetRelPath(v.Filename)
						if err != nil {
							panic(err)
						}

						errorf("Duplicated Unions: %s(%s:%v:%v) and (%s:%v:%v)", s.Unions[i].Name, *rel1, s.Unions[i].Line, s.Unions[i].Column, *rel2, v.Line, v.Column)
					}
				}
			}
			continue
		}
		seen[v.Name] = struct{}{}
		s.Unions[j] = v
		j++
	}
	s.Unions = s.Unions[:j]
}

func (s *Schema) UniqueInput(wg *sync.WaitGroup) {
	defer wg.Done()
	j := 0
	seen := make(map[string]struct{}, len(s.Inputs))
	for _, v := range s.Inputs {
		if _, ok := seen[v.Name]; ok {
			for i := 0; i < j; i++ {
				if s.Inputs[i].Name == v.Name {
					if IsEqualWithoutDescriptions(s.Inputs[i].Fields, v.Fields) {
						mergeDescriptionsAndComments(s.Inputs[i], v)
						break
					} else {

						rel1, err := GetRelPath(s.Inputs[i].Filename)
						if err != nil {
							panic(err)
						}
						rel2, err := GetRelPath(v.Filename)
						if err != nil {
							panic(err)
						}

						errorf("Duplicated Inputs: %s(%s:%v:%v) and (%s:%v:%v)", s.Inputs[i].Name, *rel1, s.Inputs[i].Line, s.Inputs[i].Column, *rel2, v.Line, v.Column)
					}
				}
			}
			continue
		}
		seen[v.Name] = struct{}{}
		s.Inputs[j] = v
		j++
	}
	s.Inputs = s.Inputs[:j]
}

func mergeFields(a []*Field, b []*Field) []*Field {
	ps := make([]*Field, len(a)+len(b))
	j := 0
	seen := make(map[string]struct{}, len(a)+len(b))
	combined := append(a, b...)
	for _, v := range combined {
		if _, ok := seen[v.Name]; ok {
			for i := 0; i < j; i++ {
				if combined[i].Name == v.Name {
					if IsEqualWithoutDescriptions(combined[i].Args, v.Args) && combined[i].Type == v.Type && combined[i].Null == v.Null && combined[i].IsList == v.IsList && combined[i].IsListNull == v.IsListNull && IsEqualWithoutDescriptions(combined[i].Directives, v.Directives) {
						mergeDescriptionsAndComments(combined[i], v)
						break
					} else {
						rel1, err := GetRelPath(combined[i].Filename)
						if err != nil {
							panic(err)
						}
						rel2, err := GetRelPath(v.Filename)
						if err != nil {
							panic(err)
						}

						errorf("Duplicated Types: %s(%s:%v:%v) and (%s:%v:%v)", combined[i].Name, *rel1, combined[i].Line, combined[i].Column, *rel2, v.Line, v.Column)
					}
				}
			}
			continue
		}
		seen[v.Name] = struct{}{}
		ps[j] = v
		j++
	}
	return ps[:j]
}
