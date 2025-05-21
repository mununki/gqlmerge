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
	SchemaDefinitions    []*SchemaDefinition
	Types                []*Type
	Scalars              []*Scalar
	Enums                []*Enum
	Interfaces           []*Interface
	Unions               []*Union
	Inputs               []*Input
	DirectiveDefinitions []*DirectiveDefinition
}

type SchemaDefinition struct {
	BaseFileInfo
	Query        *string
	Mutation     *string
	Subscription *string
	Descriptions *[]string
}

type DirectiveDefinition struct {
	BaseFileInfo
	Name         string
	Args         []*Arg
	Repeatable   bool
	Locations    []string
	Descriptions *[]string
}

type DirectiveArg struct {
	Name         string
	Value        []string
	IsList       bool
	Descriptions *[]string
}

type Directive struct {
	Name          string
	DirectiveArgs []*DirectiveArg
	Descriptions  *[]string
}

type Type struct {
	BaseFileInfo
	Name         string
	Impl         bool
	ImplTypes    []string
	Fields       []*Field
	Directives   []*Directive
	Descriptions *[]string
	Extend       bool
}

type Arg struct {
	Name          string
	Type          string
	DefaultValues *[]string // in case of default values e.g. admin(role: Role = ADMIN): Admin!
	Null          bool
	IsList        bool
	IsListNull    bool
	Directives    []*Directive
	Descriptions  *[]string
}

type Field struct {
	BaseFileInfo
	Name          string
	Args          []*Arg
	Type          string
	Null          bool
	IsList        bool
	IsListNull    bool
	DefaultValues *[]string
	Directives    []*Directive
	Descriptions  *[]string
	Comments      *[]string
}

type Scalar struct {
	BaseFileInfo
	Name         string
	Directives   []*Directive
	Descriptions *[]string
	Comments     *[]string
}

type EnumValue struct {
	Name         string
	Directives   []*Directive
	Descriptions *[]string
	Comments     *[]string
}

type Enum struct {
	BaseFileInfo
	Name         string
	EnumValues   []EnumValue
	Directives   []*Directive
	Descriptions *[]string
}

type Interface struct {
	BaseFileInfo
	Name         string
	Fields       []*Field
	Directives   []*Directive
	Descriptions *[]string
}

type Union struct {
	BaseFileInfo
	Name         string
	Types        []string
	Directives   []*Directive
	Descriptions *[]string
}

type Input struct {
	BaseFileInfo
	Name         string
	Descriptions *[]string
	Directives   []*Directive
	Fields       []*Field
}
