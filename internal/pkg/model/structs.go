package model

// File is the representation of a Go source file to be generated
type File struct {
	// Package is the short package name as it should appear in source file
	Package string
	// Imports is a list of additional imports to add to the source file
	Imports []string
	// DataTypes are TOSCA DataTypes to be included in this source file
	DataTypes []DataType
}

// DataType is the representation of a TOSCA datatype
type DataType struct {
	// Name is the Go struct identifier name
	Name string
	// FQDTN is the Fully Qualified DataType Name in TOSCA
	FQDTN string
	// DerivedFrom is the parent Go struct identifier name
	DerivedFrom string
	// Description is the data type description field
	Description string
	// Fields are DataType fields (aka properties in TOSCA)
	Fields []Field
}

// Field is the representation of a TOSCA datatype property
type Field struct {
	// Name is the Go struct field name
	Name string
	// OriginalName is the field name as it appear in the TOSCA definition
	OriginalName string
	// Type is the Go struct field type
	Type string
	// Description is the property description field
	Description string
}
