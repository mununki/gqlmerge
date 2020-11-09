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
	Files         []*os.File
	Mutations     []*Mutation
	Queries       []*Query
	Subscriptions []*Subscription
	TypeNames     []*TypeName
	Scalars       []*Scalar
	Enums         []*Enum
	Interfaces    []*Interface
	Unions        []*Union
	Inputs        []*Input
}

type Mutation struct {
	BaseFileInfo
	Name      string
	Args      []*Arg
	Resp      Resp
	Directive *Directive
}

type Query struct {
	BaseFileInfo
	Name      string
	Args      []*Arg
	Resp      Resp
	Directive *Directive
}

type Subscription struct {
	BaseFileInfo
	Name      string
	Args      []*Arg
	Resp      Resp
	Directive *Directive
}

type TypeName struct {
	BaseFileInfo
	Name      string
	Impl      bool
	ImplType  *string // deprecated, use ImplTypes
	ImplTypes []string
	Props     []*Prop
}

type Arg struct {
	Param      string
	Type       string
	TypeExt    *string // in case of enum e.g. admin(role: Role = ADMIN): Admin!
	Null       bool
	IsList     bool
	IsListNull bool
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
	Directive  *Directive
}

type Scalar struct {
	BaseFileInfo
	Name string
}

type Enum struct {
	BaseFileInfo
	Name   string
	Fields []string
}

type Interface struct {
	BaseFileInfo
	Name  string
	Props []*Prop
}

type Union struct {
	BaseFileInfo
	Name   string
	Fields []string
}

type Input struct {
	BaseFileInfo
	Name  string
	Props []*Prop
}

type Directive struct {
	string
}
