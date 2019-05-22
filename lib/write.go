package lib

import (
	"strings"
)

type MergedSchema struct {
	strings.Builder
}

func (ms *MergedSchema) addIndent(n int) {
	for i := 0; i < n; i++ {
		ms.WriteString(" ")
	}
}

func (ms *MergedSchema) stitchArgument(a *Arg, l int, i int) {
	if l > 2 {
		ms.addIndent(4)
	}
	ms.WriteString(a.Param + ": ")

	if a.IsList {
		ms.WriteString("[")
		ms.WriteString(a.Type)

		if !a.Null {
			ms.WriteString("!")
		}
		ms.WriteString("]")
		if !a.IsListNull {
			ms.WriteString("!")
		}
	} else {
		ms.WriteString(a.Type)
		if a.TypeExt != nil {
			ms.WriteString(" = " + *a.TypeExt)
		}
		if !a.Null {
			ms.WriteString("!")
		}
	}

	if l <= 2 && i != l-1 {
		ms.WriteString(", ")
	}
	if l > 2 && i != l-1 {
		ms.WriteString("\n")
	}
}

func (ms *MergedSchema) StitchSchema(s *Schema) string {
	numOfQurs := len(s.Queries)
	numOfMuts := len(s.Mutations)
	numOfSubs := len(s.Subscriptions)

	ms.WriteString("schema {\n")
	if numOfQurs > 0 {
		ms.addIndent(2)
		ms.WriteString("query: Query\n")
	}
	if numOfMuts > 0 {
		ms.addIndent(2)
		ms.WriteString("mutation: Mutation\n")
	}
	if numOfSubs > 0 {
		ms.addIndent(2)
		ms.WriteString("subscription: Subscription\n")
	}
	ms.WriteString("}\n")

	if numOfQurs > 0 {
		ms.WriteString(`type Query {
`)
		for _, q := range s.Queries {
			ms.addIndent(2)
			ms.WriteString(q.Name)
			if l := len(q.Args); l > 0 {
				ms.WriteString("(")
				if l > 2 {
					ms.WriteString("\n")
				}

				for i, a := range q.Args {
					ms.stitchArgument(a, l, i)
				}

				if l > 2 {
					ms.WriteString("\n")
					ms.addIndent(2)
				}
				ms.WriteString(")")
			}
			ms.WriteString(": ")
			if q.Resp.IsList {
				ms.WriteString("[")
			}
			ms.WriteString(q.Resp.Name)
			if !q.Resp.Null {
				ms.WriteString("!")
			}
			if q.Resp.IsList {
				ms.WriteString("]")
			}
			if q.Resp.IsList && !q.Resp.IsListNull {
				ms.WriteString("!")
			}

			if q.Directive != nil {
				ms.WriteString(" @" + q.Directive.string)
			}

			ms.WriteString("\n")
		}
		ms.WriteString("}\n")
	}

	if numOfMuts > 0 {
		ms.WriteString(`type Mutation {
`)
		for _, m := range s.Mutations {
			ms.addIndent(2)
			ms.WriteString(m.Name)
			if l := len(m.Args); l > 0 {
				ms.WriteString("(")
				if l > 2 {
					ms.WriteString("\n")
				}

				for i, a := range m.Args {
					ms.stitchArgument(a, l, i)
				}

				if l > 2 {
					ms.WriteString("\n")
					ms.addIndent(2)
				}
				ms.WriteString(")")
			}
			ms.WriteString(": ")
			if m.Resp.IsList {
				ms.WriteString("[")
			}
			ms.WriteString(m.Resp.Name)
			if !m.Resp.Null {
				ms.WriteString("!")
			}
			if m.Resp.IsList {
				ms.WriteString("]")
			}
			if m.Resp.IsList && !m.Resp.IsListNull {
				ms.WriteString("!")
			}

			if m.Directive != nil {
				ms.WriteString(" @" + m.Directive.string)
			}

			ms.WriteString("\n")
		}
		ms.WriteString("}\n")
	}

	if numOfSubs > 0 {
		ms.WriteString(`type Subscription {
`)
		for _, c := range s.Subscriptions {
			ms.addIndent(2)
			ms.WriteString(c.Name)
			if l := len(c.Args); l > 0 {
				ms.WriteString("(")
				if l > 2 {
					ms.WriteString("\n")
				}

				for i, a := range c.Args {
					ms.stitchArgument(a, l, i)
				}

				if l > 2 {
					ms.WriteString("\n")
					ms.addIndent(2)
				}
				ms.WriteString(")")
			}
			ms.WriteString(": ")
			if c.Resp.IsList {
				ms.WriteString("[")
			}
			ms.WriteString(c.Resp.Name)
			if !c.Resp.Null {
				ms.WriteString("!")
			}
			if c.Resp.IsList {
				ms.WriteString("]")
			}
			if c.Resp.IsList && !c.Resp.IsListNull {
				ms.WriteString("!")
			}

			if c.Directive != nil {
				ms.WriteString(" @" + c.Directive.string)
			}

			ms.WriteString("\n")
		}
		ms.WriteString("}\n")
	}

	for i, t := range s.TypeNames {
		ms.WriteString("type ")
		ms.WriteString(t.Name)
		if t.Impl {
			ms.WriteString(" implements " + *t.ImplType)
		}
		ms.WriteString(" {\n")
		for _, p := range t.Props {
			ms.addIndent(2)
			ms.WriteString(p.Name)

			if l := len(p.Args); l > 0 {
				ms.WriteString("(")
				if l > 2 {
					ms.WriteString("\n")
				}
				for i, a := range p.Args {
					ms.stitchArgument(a, l, i)
				}
				if l > 2 {
					ms.WriteString("\n")
					ms.addIndent(2)
				}
				ms.WriteString(")")
			}

			ms.WriteString(": ")
			if p.IsList {
				ms.WriteString("[")
			}
			ms.WriteString(p.Type)
			if !p.Null {
				ms.WriteString("!")
			}
			if p.IsList {
				ms.WriteString("]")
			}
			if p.IsList && !p.IsListNull {
				ms.WriteString("!")
			}

			if p.Directive != nil {
				ms.WriteString(" @" + p.Directive.string)
			}

			ms.WriteString("\n")
		}
		ms.WriteString("}")
		if i != len(s.TypeNames)-1 {
			ms.WriteString("\n")
		}
	}
	ms.WriteString("\n")

	for i, c := range s.Scalars {
		ms.WriteString("scalar " + c.Name)
		if i != len(s.Scalars)-1 {
			ms.WriteString("\n")
		}
	}
	ms.WriteString("\n")

	for i, e := range s.Enums {
		ms.WriteString("enum " + e.Name + " {\n")
		for _, n := range e.Fields {
			ms.addIndent(2)
			ms.WriteString(n + "\n")
		}
		ms.WriteString("}")
		if i != len(s.Enums)-1 {
			ms.WriteString("\n")
		}
	}
	ms.WriteString("\n")

	for j, i := range s.Interfaces {
		ms.WriteString("interface " + i.Name + " {\n")

		for _, p := range i.Props {
			ms.addIndent(2)
			ms.WriteString(p.Name)

			if l := len(p.Args); l > 0 {
				ms.WriteString("(")
				if l > 2 {
					ms.WriteString("\n")
				}
				for i, a := range p.Args {
					ms.stitchArgument(a, l, i)
				}
				if l > 2 {
					ms.WriteString("\n")
					ms.addIndent(2)
				}
				ms.WriteString(")")
			}

			ms.WriteString(": ")
			if p.IsList {
				ms.WriteString("[")
			}
			ms.WriteString(p.Type)
			if !p.Null {
				ms.WriteString("!")
			}
			if p.IsList {
				ms.WriteString("]")
			}
			if p.IsList && !p.IsListNull {
				ms.WriteString("!")
			}

			if p.Directive != nil {
				ms.WriteString(" @" + p.Directive.string)
			}

			ms.WriteString("\n")
		}
		ms.WriteString("}")
		if j < len(s.Interfaces)-1 {
			ms.WriteString("\n")
		}
	}
	ms.WriteString("\n")

	for _, u := range s.Unions {
		ms.WriteString("union " + u.Name + " = ")
		for j, f := range u.Fields {
			ms.WriteString(f)
			if j < len(u.Fields)-1 {
				ms.WriteString(" | ")
			}
		}
	}
	ms.WriteString("\n")

	for j, i := range s.Inputs {
		ms.WriteString("input " + i.Name + " {\n")

		for _, p := range i.Props {
			ms.addIndent(2)
			ms.WriteString(p.Name + ": ")
			if p.IsList {
				ms.WriteString("[")
			}
			ms.WriteString(p.Type)
			if !p.Null {
				ms.WriteString("!")
			}
			if p.IsList {
				ms.WriteString("]")
			}
			if p.IsList && !p.IsListNull {
				ms.WriteString("!")
			}

			if p.Directive != nil {
				ms.WriteString(" @" + p.Directive.string)
			}

			ms.WriteString("\n")
		}

		ms.WriteString("}\n")
		if j < len(s.Interfaces)-1 {
			ms.WriteString("\n")
		}
	}

	return ms.String()
}
