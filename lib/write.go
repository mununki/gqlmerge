package lib

import (
	"strings"
)

type MergedSchema struct {
	buf    strings.Builder
	Indent string
}

func (ms *MergedSchema) StitchSchema(s *Schema) string {
	numOfDirs := len(s.DirectiveDefinitions)

	// FIXME
	ms.buf.WriteString("schema {\n")
	ms.addIndent(1)
	ms.buf.WriteString("query: Query\n")
	ms.addIndent(1)
	ms.buf.WriteString("mutation: Mutation\n")
	ms.addIndent(1)
	ms.buf.WriteString("subscription: Subscription\n")

	ms.buf.WriteString("}\n\n")

	if numOfDirs > 0 {
		for _, q := range s.DirectiveDefinitions {
			ms.buf.WriteString(`directive @`)
			ms.buf.WriteString(q.Name)
			if l := len(q.Args); l > 0 {
				ms.buf.WriteString("(")
				if l > 2 {
					ms.buf.WriteString("\n")
				}

				for i, a := range q.Args {
					ms.stitchArgument(a, l, i)
				}

				if l > 2 {
					ms.buf.WriteString("\n")
					ms.addIndent(1)
				}
				ms.buf.WriteString(")")
			}
			if q.Repeatable {
				ms.buf.WriteString(" repeatable")
			}
			ms.buf.WriteString(" on ")
			for i, a := range q.Locations {
				if i != 0 {
					ms.buf.WriteString(" | ")
				}
				ms.buf.WriteString(a)
			}
			ms.buf.WriteString("\n")
		}
		ms.buf.WriteString("\n\n")
	}

	for i, t := range s.Types {
		ms.buf.WriteString("type ")
		ms.buf.WriteString(t.Name)
		if len(t.ImplTypes) > 0 {
			ms.buf.WriteString(" implements " + strings.Join(t.ImplTypes, " & "))
		}
		ms.buf.WriteString(" {\n")
		for _, p := range t.Fields {
			ms.addIndent(1)
			ms.buf.WriteString(p.Name)

			if l := len(p.Args); l > 0 {
				ms.buf.WriteString("(")
				if l > 2 {
					ms.buf.WriteString("\n")
				}
				for i, a := range p.Args {
					ms.stitchArgument(a, l, i)
				}
				if l > 2 {
					ms.buf.WriteString("\n")
					ms.addIndent(1)
				}
				ms.buf.WriteString(")")
			}

			ms.buf.WriteString(": ")
			if p.IsList {
				ms.buf.WriteString("[")
			}
			ms.buf.WriteString(p.Type)
			if !p.Null {
				ms.buf.WriteString("!")
			}
			if p.IsList {
				ms.buf.WriteString("]")
			}
			if p.IsList && !p.IsListNull {
				ms.buf.WriteString("!")
			}

			ms.stitchDirectives(p.Directives)

			ms.buf.WriteString("\n")
		}
		ms.buf.WriteString("}\n")
		if i != len(s.Types)-1 {
			ms.buf.WriteString("\n")
		}
	}
	ms.buf.WriteString("\n")

	for i, c := range s.Scalars {
		ms.buf.WriteString("scalar " + c.Name)
		ms.stitchDirectives(c.Directives)
		ms.buf.WriteString("\n")
		if i != len(s.Scalars)-1 {
			ms.buf.WriteString("\n")
		}
	}
	ms.buf.WriteString("\n")

	for i, e := range s.Enums {
		ms.buf.WriteString("enum " + e.Name)
		ms.stitchDirectives(e.Directives)
		ms.buf.WriteString(" {\n")
		for _, n := range e.EnumValues {
			ms.addIndent(1)
			ms.buf.WriteString(n.Name)
			ms.stitchDirectives(n.Directives)
			ms.buf.WriteString("\n")
		}
		ms.buf.WriteString("}\n")
		if i != len(s.Enums)-1 {
			ms.buf.WriteString("\n")
		}
	}
	ms.buf.WriteString("\n")

	for j, i := range s.Interfaces {
		ms.buf.WriteString("interface " + i.Name)
		ms.stitchDirectives(i.Directives)
		ms.buf.WriteString(" {\n")

		for _, p := range i.Fields {
			ms.addIndent(1)
			ms.buf.WriteString(p.Name)

			if l := len(p.Args); l > 0 {
				ms.buf.WriteString("(")
				if l > 2 {
					ms.buf.WriteString("\n")
				}
				for i, a := range p.Args {
					ms.stitchArgument(a, l, i)
				}
				if l > 2 {
					ms.buf.WriteString("\n")
					ms.addIndent(1)
				}
				ms.buf.WriteString(")")
			}

			ms.buf.WriteString(": ")
			if p.IsList {
				ms.buf.WriteString("[")
			}
			ms.buf.WriteString(p.Type)
			if !p.Null {
				ms.buf.WriteString("!")
			}
			if p.IsList {
				ms.buf.WriteString("]")
			}
			if p.IsList && !p.IsListNull {
				ms.buf.WriteString("!")
			}

			ms.stitchDirectives(p.Directives)

			ms.buf.WriteString("\n")
		}
		ms.buf.WriteString("}\n")
		if j < len(s.Interfaces)-1 {
			ms.buf.WriteString("\n")
		}
	}
	ms.buf.WriteString("\n")

	for _, u := range s.Unions {
		ms.buf.WriteString("union " + u.Name)
		ms.stitchDirectives(u.Directives)
		ms.buf.WriteString(" = ")
		types := strings.Join(u.Types, " | ")
		ms.buf.WriteString(types + "\n\n")
	}

	for j, i := range s.Inputs {
		ms.buf.WriteString("input " + i.Name + " {\n")

		for _, p := range i.Fields {
			ms.addIndent(1)
			ms.buf.WriteString(p.Name + ": ")
			if p.IsList {
				ms.buf.WriteString("[")
			}
			ms.buf.WriteString(p.Type)
			if !p.Null {
				ms.buf.WriteString("!")
			}
			if p.IsList {
				ms.buf.WriteString("]")
			}
			if p.IsList && !p.IsListNull {
				ms.buf.WriteString("!")
			}

			ms.stitchDirectives(p.Directives)

			ms.buf.WriteString("\n")
		}

		ms.buf.WriteString("}\n")
		if j < len(s.Inputs)-1 {
			ms.buf.WriteString("\n")
		}
	}

	return ms.buf.String()
}

func (ms *MergedSchema) addIndent(n int) {
	i := strings.Repeat(ms.Indent, n)
	ms.buf.WriteString(i)
}

func (ms *MergedSchema) stitchArgument(a *Arg, l int, i int) {
	if l > 2 {
		ms.addIndent(2)
	}
	ms.buf.WriteString(a.Param + ": ")

	if a.IsList {
		ms.buf.WriteString("[")
		ms.buf.WriteString(a.Type)

		if !a.Null {
			ms.buf.WriteString("!")
		}
		ms.buf.WriteString("]")
		if !a.IsListNull {
			ms.buf.WriteString("!")
		}
	} else {
		ms.buf.WriteString(a.Type)
		if a.TypeExt != nil {
			ms.buf.WriteString(" = " + *a.TypeExt)
		}
		if !a.Null {
			ms.buf.WriteString("!")
		}
	}

	if l <= 2 && i != l-1 {
		ms.buf.WriteString(", ")
	}
	if l > 2 && i != l-1 {
		ms.buf.WriteString("\n")
	}
}

func (ms *MergedSchema) stitchDirectives(a []*Directive) {
	if l := len(a); l > 0 {
		for _, a := range a {
			ms.buf.WriteString(" @" + a.Name)
			if m := len(a.DirectiveArgs); m > 0 {
				ms.buf.WriteString("(")
				for i, b := range a.DirectiveArgs {
					ms.stitchDirectiveArgument(b, m, i)
				}
				ms.buf.WriteString(")")
			}
		}
	}
}

func (ms *MergedSchema) stitchDirectiveArgument(a *DirectiveArg, l int, i int) {
	if l > 2 {
		ms.addIndent(2)
	}
	ms.buf.WriteString(a.Name + ": ")

	if a.IsList {
		ms.buf.WriteString("[")
		for i, v := range a.Value {
			if i != 0 {
				ms.buf.WriteString(",")
			}
			ms.buf.WriteString(v)
		}
		ms.buf.WriteString("]")
	} else {
		for i, v := range a.Value {
			if i != 0 {
				ms.buf.WriteString(",")
			}
			ms.buf.WriteString(v)
		}
	}

	if l <= 2 && i != l-1 {
		ms.buf.WriteString(", ")
	}
	if l > 2 && i != l-1 {
		ms.buf.WriteString("\n")
	}
}
