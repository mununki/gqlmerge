package lib

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sync"
)

// GetSchema is to parse ./schema/**/*.graphql
func (sc *Schema) GetSchema(path string) {
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
	for {
		tok := p.lex.next()
		if tok.typ == tokEOF {
			break
		}
		switch tok.typ {
		case tokSchema:
			// skip the schema {...}
			// it will be generated after parsing all
			for {
				t := p.lex.next()
				if t.typ == tokRBrace {
					break
				}
			}

		case tokString:
			// description
			p.buf = append(p.buf, tok)

		case tokSingleLineComment:
			p.buf = append(p.buf, tok)

		case tokMultiLineComment:
			p.buf = append(p.buf, tok)

		case tokDirective:
			d := DirectiveDefinition{}
			d.Filename = p.lex.filename
			d.Line = p.lex.line
			d.Column = p.lex.col
			d.Comments = p.bufString()
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

		case tokScalar:
			c := Scalar{}
			c.Filename = p.lex.filename
			c.Line = p.lex.line
			c.Column = p.lex.col
			c.Comments = p.bufString()
			name, _ := p.lex.consumeIdent()
			c.Name = name.String()
			c.Directives = p.parseDirectives()
			s.Scalars = append(s.Scalars, &c)

		case tokEnum:
			e := Enum{}
			e.Filename = p.lex.filename
			e.Line = p.lex.line
			e.Column = p.lex.col
			e.Comments = p.bufString()
			name, _ := p.lex.consumeIdent()
			e.Name = name.String()
			e.Directives = p.parseDirectives()
			p.lex.consumeToken(tokLBrace)
			for p.lex.peek() != '}' {
				ev := EnumValue{}
				name, comments := p.lex.consumeIdent()
				ev.Name = name.String()
				ev.Comments = comments
				ev.Directives = p.parseDirectives()
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
			i.Comments = p.bufString()
			i.Directives = p.parseDirectives()

			p.lex.consumeToken(tokLBrace)
			for p.lex.peek() != '}' {
				pr := Prop{}
				name, comments := p.lex.consumeIdent()
				pr.Name = name.String()
				pr.Comments = comments

				pr.Args = p.parseArgs()

				p.lex.consumeToken(tokColon)

				if p.lex.peek() == '[' {
					pr.IsList = true
					p.lex.consumeToken('[')
					name, _ = p.lex.consumeIdent()
					pr.Type = name.String()
					if p.lex.peek() == '!' {
						pr.Null = false
						p.lex.consumeToken(tokBang)
					} else {
						pr.Null = true
					}
					p.lex.consumeToken(tokRBracket)
					if p.lex.peek() == '!' {
						pr.IsListNull = false
						p.lex.consumeToken(tokBang)
					} else {
						pr.IsListNull = true
					}
				} else {
					pr.IsList = false
					pr.IsListNull = false
					name, _ = p.lex.consumeIdent()
					pr.Type = name.String()
					if p.lex.peek() == '!' {
						pr.Null = false
						p.lex.consumeToken(tokBang)
					} else {
						pr.Null = true
					}
				}

				sc := p.parseSingleLineComment()
				if sc != nil {
					cs := append(*pr.Comments, *sc)
					pr.Comments = &cs
				}

				pr.Directives = p.parseDirectives()

				i.Props = append(i.Props, &pr)
			}

			s.Interfaces = append(s.Interfaces, &i)
			p.lex.consumeToken(tokRBrace)

		case tokUnion:
			u := Union{}
			u.Filename = p.lex.filename
			u.Line = p.lex.line
			u.Column = p.lex.col
			name, _ := p.lex.consumeIdent()
			u.Name = name.String()
			u.Directives = p.parseDirectives()
			p.lex.consumeToken(tokEqual)
			for p.lex.peek() != '\n' || p.lex.peek() != '\r' || p.lex.peek() != EofRune {
				name, _ = p.lex.consumeIdent()
				// FIXME comments?
				u.Fields = append(u.Fields, name.String())
				if p.lex.peek() == '|' {
					p.lex.consumeToken(tokBang)
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
			i.Comments = p.bufString()

			p.lex.consumeToken(tokLBrace)

			for p.lex.peek() != '}' {
				pr := Prop{}
				name, comments := p.lex.consumeIdent()
				pr.Name = name.String()
				pr.Comments = comments
				p.lex.consumeToken(tokColon)

				if p.lex.peek() == '[' {
					pr.IsList = true
					p.lex.consumeToken(tokLBracket)
					name, _ = p.lex.consumeIdent()
					pr.Type = name.String()
					if p.lex.peek() == '!' {
						pr.Null = false
						p.lex.consumeToken(tokBang)
					} else {
						pr.Null = true
					}
					p.lex.consumeToken(tokRBracket)
					if p.lex.peek() == '!' {
						pr.IsListNull = false
						p.lex.consumeToken(tokBang)
					} else {
						pr.IsListNull = true
					}
				} else {
					pr.IsList = false
					pr.IsListNull = false
					name, _ = p.lex.consumeIdent()
					pr.Type = name.String()
					if p.lex.peek() == '!' {
						pr.Null = false
						p.lex.consumeToken(tokBang)
					} else {
						pr.Null = true
					}
				}

				pr.Directives = p.parseDirectives()

				i.Props = append(i.Props, &pr)
			}

			s.Inputs = append(s.Inputs, &i)
			p.lex.consumeToken(tokRBrace)

		case tokType:
			tok := p.lex.next()
			switch tok.typ {
			case tokQuery:
				p.lex.consumeToken(tokLBrace)

				for p.lex.peek() != '}' {
					q := Query{}
					q.Filename = p.lex.filename
					q.Line = p.lex.line
					q.Column = p.lex.col
					name, comments := p.lex.consumeIdent()
					q.Name = name.String()
					q.Comments = comments

					q.Args = p.parseArgs()

					p.lex.consumeToken(tokColon)
					r := Resp{}
					if p.lex.peek() == '[' {
						r.IsList = true
						p.lex.consumeToken(tokLBracket)
						name, _ = p.lex.consumeIdent()
						r.Name = name.String()
						if p.lex.peek() == '!' {
							r.Null = false
							p.lex.consumeToken(tokBang)
						} else {
							r.Null = true
						}
						p.lex.consumeToken(tokRBracket)
						if p.lex.peek() == '!' {
							r.IsListNull = false
							p.lex.consumeToken(tokBang)
						} else {
							r.IsListNull = true
						}
					} else {
						r.IsList = false
						r.IsListNull = false
						name, _ = p.lex.consumeIdent()
						r.Name = name.String()
						if p.lex.peek() == '!' {
							r.Null = false
							p.lex.consumeToken(tokBang)
						} else {
							r.Null = true
						}
					}
					q.Resp = r
					q.Directives = p.parseDirectives()
					sc := p.parseSingleLineComment()
					if sc != nil {
						cs := append(*q.Comments, *sc)
						q.Comments = &cs
					}
					s.Queries = append(s.Queries, &q)
				}
				p.lex.consumeToken(tokRBrace)

			case tokMutation:
				p.lex.consumeToken(tokLBrace)

				for p.lex.peek() != '}' {
					m := Mutation{}
					m.Filename = p.lex.filename
					m.Line = p.lex.line
					m.Column = p.lex.col
					name, comments := p.lex.consumeIdent()
					m.Name = name.String()
					m.Comments = comments

					m.Args = p.parseArgs()

					p.lex.consumeToken(tokColon)
					r := Resp{}
					if p.lex.peek() == '[' {
						r.IsList = true
						p.lex.consumeToken(tokLBracket)
						name, _ = p.lex.consumeIdent()
						r.Name = name.String()
						if p.lex.peek() == '!' {
							r.Null = false
							p.lex.consumeToken(tokBang)
						} else {
							r.Null = true
						}
						p.lex.consumeToken(tokRBracket)
						if p.lex.peek() == '!' {
							r.IsListNull = false
							p.lex.consumeToken(tokBang)
						} else {
							r.IsListNull = true
						}
					} else {
						r.IsList = false
						r.IsListNull = false
						name, _ = p.lex.consumeIdent()
						r.Name = name.String()
						if p.lex.peek() == '!' {
							r.Null = false
							p.lex.consumeToken(tokBang)
						} else {
							r.Null = true
						}
					}

					m.Resp = r
					m.Directives = p.parseDirectives()
					s.Mutations = append(s.Mutations, &m)
				}
				p.lex.consumeToken(tokRBrace)

			case tokSubscription:
				p.lex.consumeToken(tokLBrace)

				for p.lex.peek() != '}' {
					c := Subscription{}
					c.Filename = p.lex.filename
					c.Line = p.lex.line
					c.Column = p.lex.col
					name, comments := p.lex.consumeIdent()
					c.Name = name.String()
					c.Comments = comments

					c.Args = p.parseArgs()

					p.lex.consumeToken(tokColon)
					r := Resp{}
					if p.lex.peek() == '[' {
						r.IsList = true
						p.lex.consumeToken(tokLBracket)
						name, _ = p.lex.consumeIdent()
						r.Name = name.String()
						if p.lex.peek() == '!' {
							r.Null = false
							p.lex.consumeToken(tokBang)
						} else {
							r.Null = true
						}
						p.lex.consumeToken(tokRBracket)
						if p.lex.peek() == '!' {
							r.IsListNull = false
							p.lex.consumeToken(tokBang)
						} else {
							r.IsListNull = true
						}
					} else {
						r.IsList = false
						r.IsListNull = false
						name, _ = p.lex.consumeIdent()
						r.Name = name.String()
						if p.lex.peek() == '!' {
							r.Null = false
							p.lex.consumeToken(tokBang)
						} else {
							r.Null = true
						}
					}

					c.Resp = r
					c.Directives = p.parseDirectives()
					s.Subscriptions = append(s.Subscriptions, &c)
					p.lex.skipSpace()
				}
				p.lex.consumeToken(tokRBrace)

			default:
				t := TypeName{}
				t.Filename = p.lex.filename
				t.Line = p.lex.line
				t.Column = p.lex.col
				t.Name = tok.String()
				t.Comments = p.bufString()

				next := p.lex.next()
				switch next.typ {
				case tokImplements:
					t.Impl = true
					name, _ := p.lex.consumeIdent()
					t.ImplTypes = append(t.ImplTypes, name.String())
					for p.lex.peek() == '&' {
						p.lex.consumeToken(tokAmpersand)
						name, _ = p.lex.consumeIdent()
						t.ImplTypes = append(t.ImplTypes, name.String())
					}
				case tokLBrace:
					for p.lex.peek() != '}' {
						pr := Prop{}
						name, comments := p.lex.consumeIdent()
						pr.Name = name.String()
						pr.Comments = comments

						pr.Args = p.parseArgs()

						p.lex.consumeToken(tokColon)

						if p.lex.peek() == '[' {
							pr.IsList = true
							p.lex.consumeToken(tokLBracket)
							name, _ = p.lex.consumeIdent()
							pr.Type = name.String()
							if p.lex.peek() == '!' {
								pr.Null = false
								p.lex.consumeToken(tokBang)
							} else {
								pr.Null = true
							}
							p.lex.consumeToken(tokRBracket)
							if p.lex.peek() == '!' {
								pr.IsListNull = false
								p.lex.consumeToken(tokBang)
							} else {
								pr.IsListNull = true
							}
						} else {
							pr.IsList = false
							pr.IsListNull = false
							name, _ = p.lex.consumeIdent()
							pr.Type = name.String()
							if p.lex.peek() == '!' {
								pr.Null = false
								p.lex.consumeToken(tokBang)
							} else {
								pr.Null = true
							}
						}

						pr.Directives = p.parseDirectives()

						t.Props = append(t.Props, &pr)
					}

					s.TypeNames = append(s.TypeNames, &t)
					p.lex.consumeToken(tokRBrace)
				}
			}
		}
	}
}

func (s *Schema) UniqueDirectiveDefinition(wg *sync.WaitGroup) {
	defer wg.Done()
	j := 0
	seen := make(map[string]struct{}, len(s.DirectiveDefinitions))
	for _, v := range s.DirectiveDefinitions {
		if _, ok := seen[v.Name]; ok {
			for i := 0; i < j; i++ {
				if s.DirectiveDefinitions[i].Name == v.Name {
					if reflect.DeepEqual(s.DirectiveDefinitions[i].Args, v.Args) && reflect.DeepEqual(s.DirectiveDefinitions[i].Repeatable, v.Repeatable) && reflect.DeepEqual(s.DirectiveDefinitions[i].Locations, v.Locations) {
						break
					} else {
						rel1, err := GetRelPath(s.Mutations[i].Filename)
						if err != nil {
							panic(err)
						}
						rel2, err := GetRelPath(v.Filename)
						if err != nil {
							panic(err)
						}

						panic(fmt.Sprintf("Duplicated Directive Definitions: %s(%s, %v:%v) and (%s, %v:%v)", s.DirectiveDefinitions[i].Name, *rel1, s.DirectiveDefinitions[i].Line, s.DirectiveDefinitions[i].Column, *rel2, v.Line, v.Column))
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

func (s *Schema) UniqueMutation(wg *sync.WaitGroup) {
	defer wg.Done()
	j := 0
	seen := make(map[string]struct{}, len(s.Mutations))
	for _, v := range s.Mutations {
		if _, ok := seen[v.Name]; ok {
			for i := 0; i < j; i++ {
				if s.Mutations[i].Name == v.Name {
					if reflect.DeepEqual(s.Mutations[i].Args, v.Args) && reflect.DeepEqual(s.Mutations[i].Resp, v.Resp) && reflect.DeepEqual(s.Mutations[i].Directives, v.Directives) {
						break
					} else {

						rel1, err := GetRelPath(s.Mutations[i].Filename)
						if err != nil {
							panic(err)
						}
						rel2, err := GetRelPath(v.Filename)
						if err != nil {
							panic(err)
						}

						panic(fmt.Sprintf("Duplicated Mutations: %s(%s, %v:%v) and (%s, %v:%v)", s.Mutations[i].Name, *rel1, s.Mutations[i].Line, s.Mutations[i].Column, *rel2, v.Line, v.Column))
					}
				}
			}
			continue
		}
		seen[v.Name] = struct{}{}
		s.Mutations[j] = v
		j++
	}
	s.Mutations = s.Mutations[:j]
}

func (s *Schema) UniqueQuery(wg *sync.WaitGroup) {
	defer wg.Done()
	j := 0
	seen := make(map[string]struct{}, len(s.Queries))
	for _, v := range s.Queries {
		if _, ok := seen[v.Name]; ok {
			for i := 0; i < j; i++ {
				if s.Queries[i].Name == v.Name {
					if reflect.DeepEqual(s.Queries[i].Args, v.Args) && reflect.DeepEqual(s.Queries[i].Resp, v.Resp) && reflect.DeepEqual(s.Queries[i].Directives, v.Directives) {
						break
					} else {

						rel1, err := GetRelPath(s.Queries[i].Filename)
						if err != nil {
							panic(err)
						}
						rel2, err := GetRelPath(v.Filename)
						if err != nil {
							panic(err)
						}

						panic(fmt.Sprintf("Duplicated Queries: %s(%s, %v:%v) and (%s, %v:%v)", s.Queries[i].Name, *rel1, s.Queries[i].Line, s.Queries[i].Column, *rel2, v.Line, v.Column))
					}
				}
			}
			continue
		}
		seen[v.Name] = struct{}{}
		s.Queries[j] = v
		j++
	}
	s.Queries = s.Queries[:j]
}

func (s *Schema) UniqueSubscription(wg *sync.WaitGroup) {
	defer wg.Done()
	j := 0
	seen := make(map[string]struct{}, len(s.Subscriptions))
	for _, v := range s.Subscriptions {
		if _, ok := seen[v.Name]; ok {
			for i := 0; i < j; i++ {
				if s.Subscriptions[i].Name == v.Name {
					if reflect.DeepEqual(s.Subscriptions[i].Args, v.Args) && reflect.DeepEqual(s.Subscriptions[i].Resp, v.Resp) && reflect.DeepEqual(s.Subscriptions[i].Directives, v.Directives) {
						break
					} else {

						rel1, err := GetRelPath(s.Subscriptions[i].Filename)
						if err != nil {
							panic(err)
						}
						rel2, err := GetRelPath(v.Filename)
						if err != nil {
							panic(err)
						}

						panic(fmt.Sprintf("Duplicated Subscriptions: %s(%s, %v:%v) and (%s, %v:%v)", s.Subscriptions[i].Name, *rel1, s.Subscriptions[i].Line, s.Subscriptions[i].Column, *rel2, v.Line, v.Column))
					}
				}
			}
			continue
		}
		seen[v.Name] = struct{}{}
		s.Subscriptions[j] = v
		j++
	}
	s.Subscriptions = s.Subscriptions[:j]
}

func (s *Schema) UniqueTypeName(wg *sync.WaitGroup) {
	defer wg.Done()
	j := 0
	seen := make(map[string]struct{}, len(s.TypeNames))
	for _, v := range s.TypeNames {
		if _, ok := seen[v.Name]; ok {
			for i := 0; i < j; i++ {
				if s.TypeNames[i].Name == v.Name {
					if reflect.DeepEqual(s.TypeNames[i].ImplTypes, v.ImplTypes) && reflect.DeepEqual(s.TypeNames[i].Props, v.Props) && reflect.DeepEqual(s.TypeNames[i].Directives, v.Directives) {
						break
					} else {

						rel1, err := GetRelPath(s.TypeNames[i].Filename)
						if err != nil {
							panic(err)
						}
						rel2, err := GetRelPath(v.Filename)
						if err != nil {
							panic(err)
						}

						panic(fmt.Sprintf("Duplicated Types: %s(%s, %v:%v) and (%s, %v:%v)", s.TypeNames[i].Name, *rel1, s.TypeNames[i].Line, s.TypeNames[i].Column, *rel2, v.Line, v.Column))
					}
				}
			}
			continue
		}
		seen[v.Name] = struct{}{}
		s.TypeNames[j] = v
		j++
	}
	s.TypeNames = s.TypeNames[:j]
}

func (s *Schema) UniqueScalar(wg *sync.WaitGroup) {
	defer wg.Done()
	j := 0
	seen := make(map[string]struct{}, len(s.Scalars))
	for _, v := range s.Scalars {
		if _, ok := seen[v.Name]; ok {
			for i := 0; i < j; i++ {
				if s.Scalars[i].Name == v.Name {
					if reflect.DeepEqual(s.Scalars[i].Directives, v.Directives) {
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

						panic(fmt.Sprintf("Duplicated Scalars: %s(%s, %v:%v) and (%s, %v:%v)", s.Scalars[i].Name, *rel1, s.Scalars[i].Line, s.Scalars[i].Column, *rel2, v.Line, v.Column))
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
					if reflect.DeepEqual(s.Enums[i].Directives, v.Directives) && reflect.DeepEqual(s.Enums[i].EnumValues, v.EnumValues) {
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

						panic(fmt.Sprintf("Duplicated Enums: %s(%s, %v:%v) and (%s, %v:%v)", s.Enums[i].Name, *rel1, s.Enums[i].Line, s.Enums[i].Column, *rel2, v.Line, v.Column))
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
					if reflect.DeepEqual(s.Interfaces[i].Directives, v.Directives) && reflect.DeepEqual(s.Interfaces[i].Props, v.Props) {
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

						panic(fmt.Sprintf("Duplicated Interfaces: %s(%s, %v:%v) and (%s, %v:%v)", s.Interfaces[i].Name, *rel1, s.Interfaces[i].Line, s.Interfaces[i].Column, *rel2, v.Line, v.Column))
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
					if reflect.DeepEqual(s.Unions[i].Directives, v.Directives) && reflect.DeepEqual(s.Unions[i].Fields, v.Fields) {
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

						panic(fmt.Sprintf("Duplicated Unions: %s(%s, %v:%v) and (%s, %v:%v)", s.Unions[i].Name, *rel1, s.Unions[i].Line, s.Unions[i].Column, *rel2, v.Line, v.Column))
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
					if reflect.DeepEqual(s.Inputs[i].Props, v.Props) {
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

						panic(fmt.Sprintf("Duplicated Inputs: %s(%s, %v:%v) and (%s, %v:%v)", s.Inputs[i].Name, *rel1, s.Inputs[i].Line, s.Inputs[i].Column, *rel2, v.Line, v.Column))
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
