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

// Type is the base type for all TOSCA types (like node types, relationship types, ...)
type Type struct {
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
