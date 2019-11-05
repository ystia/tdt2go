package model

type DataType struct {
	Name        string
	FQDTN       string
	DerivedFrom string
	Fields      []Field
}

type Field struct {
	Name         string
	OriginalName string
	Type         string
}
