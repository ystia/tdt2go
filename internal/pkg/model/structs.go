// Copyright 2018 Bull S.A.S. Atos Technologies - Bull, Rue Jean Jaures, B.P.68, 78340, Les Clayes-sous-Bois, France.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
