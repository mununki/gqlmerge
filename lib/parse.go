package lib

import (
	"fmt"
	"reflect"
	"sync"
	"text/scanner"
)

func (s *Schema) ParseSchema(l *Lexer) {
	l.ConsumeWhitespace()

	for l.Peek() != scanner.EOF {

		switch x := l.ConsumeIdent(); x {

		case "schema":
			// skip the schema { ... }
			// it will be generated after parsing all
			for {
				l.next = l.sc.Scan()

				if l.next == '}' {
					break
				}
			}
			l.ConsumeToken('}')

		case "directive":
			d := DirectiveDefinition{}
			d.Filename = l.sc.Filename
			d.Line = l.sc.Line
			d.Column = l.sc.Column
			l.ConsumeToken('@')
			d.Name = l.ConsumeIdent()
			d.Args = ParseArgument(l)
			if l.ConsumeIdent() == "repeatable" {
				d.Repeatable = true
			} else {
				// "on" should be consumed already here
				d.Repeatable = false
				ls := []string{}
				for l.Peek() != '\r' || l.Peek() != '\n' || l.Peek() != scanner.EOF {
					ls = append(ls, l.ConsumeIdent())
					if l.Peek() == '|' {
						l.ConsumeToken('|')
					} else {
						break
					}
				}
				d.Locations = ls
			}
			s.DirectiveDefinitions = append(s.DirectiveDefinitions, &d)

		case "scalar":
			c := Scalar{}
			c.Filename = l.sc.Filename
			c.Line = l.sc.Line
			c.Column = l.sc.Column
			c.Name = l.ConsumeIdent()
			c.Directives = ParseDirectives(l)
			s.Scalars = append(s.Scalars, &c)

		case "enum":
			e := Enum{}
			e.Filename = l.sc.Filename
			e.Line = l.sc.Line
			e.Column = l.sc.Column
			e.Name = l.ConsumeIdent()
			e.Directives = ParseDirectives(l)
			l.ConsumeToken('{')
			for l.Peek() != '}' {
				ev := EnumValue{}
				ev.Name = l.ConsumeIdent()
				ev.Directives = ParseDirectives(l)
				e.EnumValues = append(e.EnumValues, ev)
			}
			l.ConsumeToken('}')
			s.Enums = append(s.Enums, &e)

		case "interface":
			i := Interface{}
			i.Filename = l.sc.Filename
			i.Line = l.sc.Line
			i.Column = l.sc.Column
			i.Name = l.ConsumeIdent()
			i.Directives = ParseDirectives(l)

			l.ConsumeToken('{')

			for l.Peek() != '}' {
				p := Prop{}
				p.Name = l.ConsumeIdent()

				p.Args = ParseArgument(l)

				l.ConsumeToken(':')

				if l.Peek() == '[' {
					p.IsList = true
					l.ConsumeToken('[')
					p.Type = l.ConsumeIdent()
					if x := l.sc.TokenText(); x == "!" {
						p.Null = false
						l.ConsumeToken('!')
					} else {
						p.Null = true
					}
					l.ConsumeToken(']')
					if x := l.sc.TokenText(); x == "!" {
						p.IsListNull = false
						l.ConsumeToken('!')
					} else {
						p.IsListNull = true
					}
				} else {
					p.IsList = false
					p.IsListNull = false
					p.Type = l.ConsumeIdent()
					if x := l.sc.TokenText(); x == "!" {
						p.Null = false
						l.ConsumeToken('!')
					} else {
						p.Null = true
					}
				}

				if l.Peek() == '@' {
					p.Directives = ParseDirectives(l)
				}

				i.Props = append(i.Props, &p)
			}

			s.Interfaces = append(s.Interfaces, &i)
			l.ConsumeToken('}')

		case "union":
			u := Union{}
			u.Filename = l.sc.Filename
			u.Line = l.sc.Line
			u.Column = l.sc.Column
			u.Name = l.ConsumeIdent()
			u.Directives = ParseDirectives(l)
			l.ConsumeToken('=')
			for l.Peek() != '\r' || l.Peek() != '\n' || l.Peek() != scanner.EOF {
				u.Fields = append(u.Fields, l.ConsumeIdent())
				if l.Peek() == '|' {
					l.ConsumeToken('|')
				} else {
					break
				}
			}
			s.Unions = append(s.Unions, &u)

		case "input":
			i := Input{}
			i.Filename = l.sc.Filename
			i.Line = l.sc.Line
			i.Column = l.sc.Column
			i.Name = l.ConsumeIdent()

			l.ConsumeToken('{')

			for l.Peek() != '}' {
				p := Prop{}
				p.Name = l.ConsumeIdent()
				l.ConsumeToken(':')

				if l.Peek() == '[' {
					p.IsList = true
					l.ConsumeToken('[')
					p.Type = l.ConsumeIdent()
					if x := l.sc.TokenText(); x == "!" {
						p.Null = false
						l.ConsumeToken('!')
					} else {
						p.Null = true
					}
					l.ConsumeToken(']')
					if x := l.sc.TokenText(); x == "!" {
						p.IsListNull = false
						l.ConsumeToken('!')
					} else {
						p.IsListNull = true
					}
				} else {
					p.IsList = false
					p.IsListNull = false
					p.Type = l.ConsumeIdent()
					if x := l.sc.TokenText(); x == "!" {
						p.Null = false
						l.ConsumeToken('!')
					} else {
						p.Null = true
					}
				}

				p.Directives = ParseDirectives(l)

				i.Props = append(i.Props, &p)
			}

			s.Inputs = append(s.Inputs, &i)
			l.ConsumeToken('}')

		case "type":

			switch x := l.ConsumeIdent(); x {

			case "Query":
				l.ConsumeToken('{')

				for l.Peek() != '}' {
					q := Query{}
					q.Filename = l.sc.Filename
					q.Line = l.sc.Line
					q.Column = l.sc.Column
					q.Name = l.ConsumeIdent()

					q.Args = ParseArgument(l)

					l.ConsumeToken(':')
					r := Resp{}
					if l.Peek() == '[' {
						r.IsList = true
						l.ConsumeToken('[')
						r.Name = l.ConsumeIdent()
						if x := l.sc.TokenText(); x == "!" {
							r.Null = false
							l.ConsumeToken('!')
						} else {
							r.Null = true
						}
						l.ConsumeToken(']')
						if x := l.sc.TokenText(); x == "!" {
							r.IsListNull = false
							l.ConsumeToken('!')
						} else {
							r.IsListNull = true
						}
					} else {
						r.IsList = false
						r.IsListNull = false
						r.Name = l.ConsumeIdent()
						if x := l.sc.TokenText(); x == "!" {
							r.Null = false
							l.ConsumeToken('!')
						} else {
							r.Null = true
						}
					}
					q.Resp = r

					q.Directives = ParseDirectives(l)

					s.Queries = append(s.Queries, &q)
				}
				l.ConsumeToken('}')

			case "Mutation":
				l.ConsumeToken('{')

				for l.Peek() != '}' {
					m := Mutation{}
					m.Filename = l.sc.Filename
					m.Line = l.sc.Line
					m.Column = l.sc.Column
					m.Name = l.ConsumeIdent()

					m.Args = ParseArgument(l)

					l.ConsumeToken(':')
					r := Resp{}
					if l.Peek() == '[' {
						r.IsList = true
						l.ConsumeToken('[')
						r.Name = l.ConsumeIdent()
						if x := l.sc.TokenText(); x == "!" {
							r.Null = false
							l.ConsumeToken('!')
						} else {
							r.Null = true
						}
						l.ConsumeToken(']')
						if x := l.sc.TokenText(); x == "!" {
							r.IsListNull = false
							l.ConsumeToken('!')
						} else {
							r.IsListNull = true
						}
					} else {
						r.IsList = false
						r.IsListNull = false
						r.Name = l.ConsumeIdent()
						if x := l.sc.TokenText(); x == "!" {
							r.Null = false
							l.ConsumeToken('!')
						} else {
							r.Null = true
						}
					}

					m.Resp = r

					m.Directives = ParseDirectives(l)

					s.Mutations = append(s.Mutations, &m)
				}
				l.ConsumeToken('}')

			case "Subscription":
				l.ConsumeToken('{')

				for l.Peek() != '}' {
					c := Subscription{}
					c.Filename = l.sc.Filename
					c.Line = l.sc.Line
					c.Column = l.sc.Column
					c.Name = l.ConsumeIdent()

					c.Args = ParseArgument(l)

					l.ConsumeToken(':')
					r := Resp{}
					if l.Peek() == '[' {
						r.IsList = true
						l.ConsumeToken('[')
						r.Name = l.ConsumeIdent()
						if x := l.sc.TokenText(); x == "!" {
							r.Null = false
							l.ConsumeToken('!')
						} else {
							r.Null = true
						}
						l.ConsumeToken(']')
						if x := l.sc.TokenText(); x == "!" {
							r.IsListNull = false
							l.ConsumeToken('!')
						} else {
							r.IsListNull = true
						}
					} else {
						r.IsList = false
						r.IsListNull = false
						r.Name = l.ConsumeIdent()
						if x := l.sc.TokenText(); x == "!" {
							r.Null = false
							l.ConsumeToken('!')
						} else {
							r.Null = true
						}
					}
					c.Resp = r

					c.Directives = ParseDirectives(l)

					s.Subscriptions = append(s.Subscriptions, &c)
				}
				l.ConsumeToken('}')

			default:
				t := TypeName{}
				t.Filename = l.sc.Filename
				t.Line = l.sc.Line
				t.Column = l.sc.Column
				t.Name = x

				// handling in case of type has implements
				if l.Peek() == scanner.Ident {
					l.ConsumeIdent()
					t.Impl = true
					x := l.ConsumeIdent()
					t.ImplType = &x
					t.ImplTypes = append(t.ImplTypes, x)
					for l.Peek() == '&' {
						l.ConsumeToken('&')
						x := l.ConsumeIdent()
						t.ImplTypes = append(t.ImplTypes, x)
					}
				} else {
					t.Impl = false
				}

				l.ConsumeToken('{')

				for l.Peek() != '}' {
					p := Prop{}
					p.Name = l.ConsumeIdent()

					p.Args = ParseArgument(l)

					l.ConsumeToken(':')

					if l.Peek() == '[' {
						p.IsList = true
						l.ConsumeToken('[')
						p.Type = l.ConsumeIdent()
						if x := l.sc.TokenText(); x == "!" {
							p.Null = false
							l.ConsumeToken('!')
						} else {
							p.Null = true
						}
						l.ConsumeToken(']')
						if x := l.sc.TokenText(); x == "!" {
							p.IsListNull = false
							l.ConsumeToken('!')
						} else {
							p.IsListNull = true
						}
					} else {
						p.IsList = false
						p.IsListNull = false
						p.Type = l.ConsumeIdent()
						if x := l.sc.TokenText(); x == "!" {
							p.Null = false
							l.ConsumeToken('!')
						} else {
							p.Null = true
						}
					}

					p.Directives = ParseDirectives(l)

					t.Props = append(t.Props, &p)
				}

				s.TypeNames = append(s.TypeNames, &t)
				l.ConsumeToken('}')
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

func ParseArgument(l *Lexer) []*Arg {
	args := []*Arg{}

	for l.Peek() == '(' {
		l.ConsumeToken('(')
		for l.Peek() != ')' {
			arg := Arg{}
			arg.Param = l.ConsumeIdent()
			l.ConsumeToken(':')

			if l.Peek() == '[' {
				arg.IsList = true
				l.ConsumeToken('[')
				arg.Type = l.ConsumeIdent()
				if l.Peek() == '!' {
					arg.Null = false
					l.ConsumeToken('!')
				} else {
					arg.Null = true
				}
				l.ConsumeToken(']')

				if x := l.sc.TokenText(); x == "!" {
					arg.IsListNull = false
					l.ConsumeToken('!')
				} else {
					arg.IsListNull = true
				}
			} else {
				arg.Type = l.ConsumeIdent()

				if l.Peek() == '=' {
					l.ConsumeToken('=')
					ext := l.ConsumeIdent()
					arg.TypeExt = &ext
				}

				if x := l.sc.TokenText(); x == "!" {
					arg.Null = false
					l.ConsumeToken('!')
				} else {
					arg.Null = true
				}
			}
			if l.Peek() == '@' {
				arg.Directives = ParseDirectives(l)
			}

			args = append(args, &arg)
		}
		l.ConsumeToken(')')
	}
	return args
}

func ParseDirectives(l *Lexer) []*Directive {
	ds := []*Directive{}

	for l.Peek() == '@' {
		l.ConsumeToken('@')
		d := Directive{}
		d.Name = l.ConsumeIdent()

		for l.Peek() == '(' {
			l.ConsumeToken('(')
			for l.Peek() != ')' {
				da := DirectiveArg{}
				da.Name = l.ConsumeIdent()
				l.ConsumeToken(':')

				if l.Peek() == '[' {
					da.IsList = true
					l.ConsumeToken('[')
					da.Value = ParseList(l)
					l.ConsumeToken(']')
				} else {
					da.IsList = false
					if l.Peek() == scanner.String {
						da.Value = append(da.Value, l.ConsumeString())
					} else {
						da.Value = append(da.Value, l.ConsumeIdent())
					}
				}

				d.DirectiveArgs = append(d.DirectiveArgs, &da)
			}
			l.ConsumeToken(')')
		}
		ds = append(ds, &d)
	}
	return ds
}

func ParseList(l *Lexer) []string {
	ss := []string{}
	for l.Peek() != ']' {
		if l.Peek() == scanner.String {
			ss = append(ss, l.ConsumeString())
		} else {
			ss = append(ss, l.ConsumeIdent())
		}
	}
	return ss
}
