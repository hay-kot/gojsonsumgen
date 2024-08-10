package gojsonsum

type SumTypeDef struct {
	TypeName      string
	TypeMap       map[string]string
	Opts          SumTypeOptions
	Discriminator string
}

type SumTypeOptions struct {
	Tag  string
	Name string
}
