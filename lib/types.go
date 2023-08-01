package lib

import (
	"os"
)

type BaseFileInfo struct {
	Filename string
	Line     int
	Column   int
}

type Schema struct {
	Files                []*os.File
	Mutations            []*Mutation
	Queries              []*Query
	Subscriptions        []*Subscription
	TypeNames            []*TypeName
	Scalars              []*Scalar
	Enums                []*Enum
	Interfaces           []*Interface
	Unions               []*Union
	Inputs               []*Input
	DirectiveDefinitions []*DirectiveDefinition
}

type Mutation struct {
	BaseFileInfo
	Name       string
	Args       []*Arg
	Resp       Resp
	Directives []*Directive
}

type Query struct {
	BaseFileInfo
	Name       string
	Args       []*Arg
	Resp       Resp
	Directives []*Directive
}

type Subscription struct {
	BaseFileInfo
	Name       string
	Args       []*Arg
	Resp       Resp
	Directives []*Directive
}

type TypeName struct {
	BaseFileInfo
	Name       string
	Impl       bool
	ImplType   *string // deprecated, use ImplTypes
	ImplTypes  []string
	Props      []*Prop
	Directives []*Directive
}

type Arg struct {
	Param      string
	Type       string
	TypeExt    *string // in case of enum e.g. admin(role: Role = ADMIN): Admin!
	Null       bool
	IsList     bool
	IsListNull bool
	Directives []*Directive
}

type Resp struct {
	Name       string
	Null       bool
	IsList     bool
	IsListNull bool
}

type Prop struct {
	Name       string
	Args       []*Arg // in case of having args e.g. city(page: Pagination): String
	Type       string
	Null       bool
	IsList     bool
	IsListNull bool
	Directives []*Directive
}

type Scalar struct {
	BaseFileInfo
	Name       string
	Directives []*Directive
}

type Enum struct {
	BaseFileInfo
	Name       string
	EnumValues []EnumValue
	Directives []*Directive
}

type EnumValue struct {
	Name       string
	Directives []*Directive
}

type Interface struct {
	BaseFileInfo
	Name       string
	Props      []*Prop
	Directives []*Directive
}

type Union struct {
	BaseFileInfo
	Name       string
	Fields     []string
	Directives []*Directive
}

type Input struct {
	BaseFileInfo
	Name  string
	Props []*Prop
}

type DirectiveDefinition struct {
	BaseFileInfo
	Name       string
	Args       []*Arg
	Repeatable bool
	Locations  []string
}

type DirectiveArg struct {
	Name   string
	Value  []string
	IsList bool
}

type Directive struct {
	Name          string
	DirectiveArgs []*DirectiveArg
}
