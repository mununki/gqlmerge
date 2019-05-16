package lib

type Schema struct {
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
	Name string
	Args []*Arg
	Resp Resp
}

type Query struct {
	Name string
	Args []*Arg
	Resp Resp
}

type Subscription struct {
	Name string
	Args []*Arg
	Resp Resp
}

type TypeName struct {
	Name     string
	Impl     bool
	ImplType *string
	Props    []*Prop
}

type Arg struct {
	Param string
	Type  string
	Null  bool
}

type Resp struct {
	Name       string
	Null       bool
	IsList     bool
	IsListNull bool
}

type Prop struct {
	Name       string
	Type       string
	Null       bool
	IsList     bool
	IsListNull bool
}

type Scalar struct {
	Name string
}

type Enum struct {
	Name   string
	Fields []string
}

type Interface struct {
	Name  string
	Props []*Prop
}

type Union struct {
	Name   string
	Fields []string
}

type Input struct {
	Name  string
	Props []*Prop
}
