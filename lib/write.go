package lib

import (
	"sort"
	"strings"
)

type MergedSchema struct {
	buf    strings.Builder
	Indent string
}

func (ms *MergedSchema) WriteSchema(s *Schema) string {
	if (s.SchemaDefinitions[0].Query != nil) || (s.SchemaDefinitions[0].Mutation != nil) || (s.SchemaDefinitions[0].Subscription != nil) {
		ms.writeDescriptions(s.SchemaDefinitions[0].Descriptions, 0, true)
		ms.buf.WriteString("schema {\n")
		ms.addIndent(1)

		if s.SchemaDefinitions[0].Query != nil {
			ms.buf.WriteString("query: " + *s.SchemaDefinitions[0].Query + "\n")
		}
		ms.addIndent(1)
		if s.SchemaDefinitions[0].Mutation != nil {
			ms.buf.WriteString("mutation: " + *s.SchemaDefinitions[0].Mutation + "\n")
		}
		ms.addIndent(1)
		if s.SchemaDefinitions[0].Subscription != nil {
			ms.buf.WriteString("subscription: " + *s.SchemaDefinitions[0].Subscription + "\n")
		}

		ms.buf.WriteString("}\n\n")
	}

	numOfDirs := len(s.DirectiveDefinitions)
	if numOfDirs > 0 {
		for _, q := range s.DirectiveDefinitions {
			ms.writeDescriptions(q.Descriptions, 0, true)
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
		ms.writeDescriptions(t.Descriptions, 0, true)
		ms.buf.WriteString("type ")
		ms.buf.WriteString(t.Name)
		if len(t.ImplTypes) > 0 {
			ms.buf.WriteString(" implements " + strings.Join(t.ImplTypes, " & "))
		}
		ms.stitchDirectives(t.Directives)
		ms.buf.WriteString(" {\n")
		for _, p := range t.Fields {
			ms.writeDescriptions(p.Descriptions, 1, false)
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

			ms.writeComments(p.Comments)

			ms.buf.WriteString("\n")
		}
		ms.buf.WriteString("}\n")
		if i != len(s.Types)-1 {
			ms.buf.WriteString("\n")
		}
	}
	ms.buf.WriteString("\n")

	for i, c := range s.Scalars {
		ms.writeDescriptions(c.Descriptions, 0, true)
		ms.buf.WriteString("scalar " + c.Name)
		ms.stitchDirectives(c.Directives)
		ms.writeComments(c.Comments)
		ms.buf.WriteString("\n")
		if i != len(s.Scalars)-1 {
			ms.buf.WriteString("\n")
		}
	}
	ms.buf.WriteString("\n")

	for i, e := range s.Enums {
		ms.writeDescriptions(e.Descriptions, 0, true)
		ms.buf.WriteString("enum " + e.Name)
		ms.stitchDirectives(e.Directives)
		ms.buf.WriteString(" {\n")
		for _, n := range e.EnumValues {
			ms.addIndent(1)
			ms.buf.WriteString(n.Name)
			ms.stitchDirectives(n.Directives)
			ms.writeComments(n.Comments)
			ms.buf.WriteString("\n")
		}
		ms.buf.WriteString("}\n")
		if i != len(s.Enums)-1 {
			ms.buf.WriteString("\n")
		}
	}
	ms.buf.WriteString("\n")

	for j, i := range s.Interfaces {
		ms.writeDescriptions(i.Descriptions, 0, true)
		ms.buf.WriteString("interface " + i.Name)
		ms.stitchDirectives(i.Directives)
		ms.buf.WriteString(" {\n")

		for _, fd := range i.Fields {
			ms.writeDescriptions(fd.Descriptions, 1, true)
			ms.addIndent(1)
			ms.buf.WriteString(fd.Name)

			if l := len(fd.Args); l > 0 {
				ms.buf.WriteString("(")
				if l > 2 {
					ms.buf.WriteString("\n")
				}
				for i, a := range fd.Args {
					ms.stitchArgument(a, l, i)
				}
				if l > 2 {
					ms.buf.WriteString("\n")
					ms.addIndent(1)
				}
				ms.buf.WriteString(")")
			}

			ms.buf.WriteString(": ")
			if fd.IsList {
				ms.buf.WriteString("[")
			}
			ms.buf.WriteString(fd.Type)
			if !fd.Null {
				ms.buf.WriteString("!")
			}
			if fd.IsList {
				ms.buf.WriteString("]")
			}
			if fd.IsList && !fd.IsListNull {
				ms.buf.WriteString("!")
			}

			ms.stitchDirectives(fd.Directives)

			ms.buf.WriteString("\n")
		}
		ms.buf.WriteString("}\n")
		if j < len(s.Interfaces)-1 {
			ms.buf.WriteString("\n")
		}
	}
	ms.buf.WriteString("\n")

	for _, u := range s.Unions {
		ms.writeDescriptions(u.Descriptions, 0, true)
		ms.buf.WriteString("union " + u.Name)
		ms.stitchDirectives(u.Directives)
		ms.buf.WriteString(" = ")
		types := strings.Join(u.Types, " | ")
		ms.buf.WriteString(types + "\n\n")
	}

	for j, i := range s.Inputs {
		ms.writeDescriptions(i.Descriptions, 0, true)
		ms.buf.WriteString("input " + i.Name)
		ms.stitchDirectives(i.Directives)
		ms.buf.WriteString(" {\n")

		for _, p := range i.Fields {
			ms.writeDescriptions(p.Descriptions, 1, true)
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
			if p.DefaultValues != nil {
				if p.IsList {
					ms.buf.WriteString(" = ")
					ms.buf.WriteString("[")
					ms.stitchDefaultValues(p.DefaultValues)
					ms.buf.WriteString("]")
				} else {
					ms.buf.WriteString(" = ")
					ms.stitchDefaultValues(p.DefaultValues)
				}
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
	indent := 0
	if l > 2 {
		indent = 2
	}
	ms.addIndent(indent)
	ms.writeDescriptions(a.Descriptions, indent, false)

	ms.buf.WriteString(a.Name + ": ")

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
		if a.DefaultValues != nil {
			ms.buf.WriteString(" = ")
			ms.buf.WriteString("[")
			ms.stitchDefaultValues(a.DefaultValues)
			ms.buf.WriteString("]")
		}
		ms.stitchDirectives(a.Directives)
	} else {
		ms.buf.WriteString(a.Type)
		if !a.Null {
			ms.buf.WriteString("!")
		}
		if a.DefaultValues != nil {
			ms.buf.WriteString(" = ")
			ms.stitchDefaultValues(a.DefaultValues)
		}
		ms.stitchDirectives(a.Directives)
	}

	if l <= 2 && i != l-1 {
		ms.buf.WriteString(", ")
	}
	if l > 2 && i != l-1 {
		ms.buf.WriteString("\n")
	}
}

func (ms *MergedSchema) stitchDirectives(a []*Directive) {
	sort.SliceStable(a, func(i, j int) bool {
		return a[i].Name > a[j].Name
	})
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

func (ms *MergedSchema) stitchDefaultValues(a *[]string) {
	if l := len(*a); l > 0 {
		for i, v := range *a {
			ms.buf.WriteString(v)
			if i < l-1 {
				ms.buf.WriteString(", ")
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

func (ms *MergedSchema) writeDescriptions(descriptions *[]string, indent int, newLine bool) {
	if descriptions == nil || len(*descriptions) == 0 {
		return
	}

	if indent != 0 {
		ms.addIndent(indent)
	}

	ds := *descriptions
	ms.buf.WriteString(ds[0])
	if newLine {
		ms.buf.WriteString("\n")
	} else {
		ms.buf.WriteString(" ")
	}
}

func (ms *MergedSchema) writeComments(comments *[]string) {
	if comments != nil && len(*comments) > 0 {
		c := *comments
		ms.buf.WriteString(" " + c[0])
	}
}
