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

package tosca

import (
	"strings"
)

// IsBuiltinType checks if a given type name corresponds to a TOSCA builtin type.
//
// Known builtin types:
// 	- string
//	- integer
//	- float
//	- boolean
//	- timestamp
//	- null
//	- list
//	- map
//	- version
//	- range
//	- scalar-unit.size
//	- scalar-unit.time
func IsBuiltinType(typeName string) bool {
	// type representation for map and list could be map:<EntrySchema> or list:<EntrySchema> (ex: list:integer)
	return strings.HasPrefix(typeName, "list") || strings.HasPrefix(typeName, "map") ||
		typeName == "string" || typeName == "integer" || typeName == "float" || typeName == "boolean" ||
		typeName == "timestamp" || typeName == "null" || typeName == "version" || typeName == "range" ||
		typeName == "scalar-unit.size" || typeName == "scalar-unit.time"
}

//go:generate go-enum -f=types.go

// TypeBase x ENUM(
// NODE,
// RELATIONSHIP,
// CAPABILITY,
// POLICY,
// ARTIFACT,
// DATA,
// )
type TypeBase int

// Type is the base type for all TOSCA types (like node types, relationship types, ...)
type Type struct {
	Base        TypeBase          `yaml:"base,omitempty" json:"base,omitempty"`
	DerivedFrom string            `yaml:"derived_from,omitempty" json:"derived_from,omitempty"`
	Version     string            `yaml:"version,omitempty" json:"version,omitempty"`
	ImportPath  string            `yaml:"import_path,omitempty" json:"import_path,omitempty"`
	Description string            `yaml:"description,omitempty" json:"description,omitempty"`
	Metadata    map[string]string `yaml:"metadata,omitempty" json:"metadata,omitempty"`
}

// An DataType is the representation of a TOSCA Data Type
//
// See http://docs.oasis-open.org/tosca/TOSCA-Simple-Profile-YAML/v1.2/TOSCA-Simple-Profile-YAML-v1.2.html#DEFN_ENTITY_DATA_TYPE
// for more details
type DataType struct {
	Type       `yaml:",inline"`
	Properties map[string]PropertyDefinition `yaml:"properties,omitempty" json:"properties,omitempty"`
	// Constraints not enforced in Yorc so we don't parse them
	// Constraints []ConstraintClause
}
