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

// An Topology is the representation of a TOSCA Service Template definition
//
// See http://docs.oasis-open.org/tosca/TOSCA-Simple-Profile-YAML/v1.2/TOSCA-Simple-Profile-YAML-v1.2.html#DEFN_ELEMENT_SERVICE_TEMPLATE for more details
type Topology struct {
	TOSCAVersion string            `yaml:"tosca_definitions_version" json:"tosca_definitions_version"`
	Description  string            `yaml:"description,omitempty" json:"description,omitempty"`
	Metadata     map[string]string `yaml:"metadata,omitempty" json:"metadata,omitempty"`

	DataTypes map[string]DataType `yaml:"data_types,omitempty" json:"data_types,omitempty"`
}
