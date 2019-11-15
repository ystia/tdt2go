# tdt2go

[![Go Report Card](https://goreportcard.com/badge/github.com/ystia/tdt2go)](https://goreportcard.com/report/github.com/ystia/tdt2go) [![GoDoc](https://godoc.org/github.com/ystia/tdt2go?status.svg)](https://godoc.org/github.com/ystia/tdt2go) [![Build Status](https://github.com/ystia/tdt2go/workflows/Build/badge.svg)](https://github.com/ystia/tdt2go/actions?query=workflow%3ABuild) [![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fystia%2Ftdt2go.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fystia%2Ftdt2go?ref=badge_shield)  [![license](https://img.shields.io/github/license/ystia/tdt2go.svg)](https://github.com/ystia/tdt2go/blob/master/LICENSE) [![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square)](http://makeapullrequest.com) [![Join the chat at https://gitter.im/ystia/tdt2go](https://badges.gitter.im/ystia/tdt2go.svg)](https://gitter.im/ystia/tdt2go?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

[![Quality gate](https://sonarcloud.io/api/project_badges/quality_gate?project=ystia_tdt2go)](https://sonarcloud.io/dashboard?id=ystia_tdt2go)

TOSCA DataTypes to Go (tdt2go) structures generator

## How it works

The goal of tdt2go is to create a generator that takes Data Types defined in a [TOSCA](https://docs.oasis-open.org/tosca/TOSCA-Simple-Profile-YAML/v1.3/TOSCA-Simple-Profile-YAML-v1.3.html) definition file and generate a Go source file with corresponding Go data structures.

## Installation

### Install using go get

tdt2go can be installed using `go get`

```bash
go get -u https://github.com/ystia/tdt2go/cmd/tdt2go
```

## Binaries distributions

Binaries distributions could be downloaded from [github](https://github.com/ystia/tdt2go/releases)

## Command options

```bash
$ tdt2go --help
tdt2go allows to generate Go source files containing data structures generated from files containing TOSCA data types

Usage:
  tdt2go <tosca_file> [flags]

Flags:
  -e, --exclude strings                regexp patterns of data types fully qualified names to exclude. Only non-matching datatypes will be transformed. Include patterns have the precedence over exclude patterns.
  -f, --file string                    file to be generated, if not defined resulting generated file will be printed on default output.
  -b, --generate-builtin               Generate tosca builtin types as 'range' or 'scalar-unit' for instance along with datatypes in this file. (default: false)
  -h, --help                           help for tdt2go
  -i, --include strings                regexp patterns of data types fully qualified names to include. Only matching datatypes will be transformed. Include patterns have the precedence over exclude patterns.
  -m, --name-mappings stringToString   map of regular expressions and their corresponding remplacements that will be applied to TOSCA datatypes fully qualified names to transform them into Go struct names. This is generally used to keep information from the fully qualified name into the generated name. (default [])
  -p, --package string                 package name as it should appear in source file, defaults to the package name of the current directory.
```

## Features & Roadmap

- [x] Support of types inheritance via Go composition
- [x] Support basic types conversions:
  - `integer` :arrow_right: `int`
  - `string` :arrow_right: `string`
  - `boolean` :arrow_right: `bool`
  - `float` :arrow_right: `float64`
  - `timestamp` :arrow_right: `time.Time`
  - `list` :arrow_right: slice with entry_schema support
  - `map` :arrow_right: map with entry_schema support
- [x] Generation of TOSCA builtin types such as `version`, `range`, `scalar-unit`s ...
- [x] include/exclude filters
- [x] Type name mapping like `tosca\.datatypes\.(.+)` :arrow_right: `Normative${1}` so `tosca.datatypes.Credential` become `NormativeCredential`
- [x] Use type or property description on generated comments
- [ ] Make use of TOSCA `constraints` and `default`

## Example

This example is an extract of the TOSCA normative types.

```yaml
tosca_definitions_version: tosca_simple_yaml_1_0

metadata:
  template_name: tosca-normative-types
  template_author: TOSCA TC
  template_version: 1.3.0

data_types:
  tosca.datatypes.Root:
    description: The TOSCA root Data Type all other TOSCA base Data Types derive from

  tosca.datatypes.json:
    derived_from: string
  tosca.datatypes.xml:
    derived_from: string

  tosca.datatypes.Credential:
    derived_from: tosca.datatypes.Root
    description: >
      The Credential type is a complex TOSCA data Type used when describing authorization credentials used to access network accessible resources.
    properties:
      protocol:
        type: string
        description: The optional protocol name.
        required: false
      token_type:
        type: string
        description: The required token type.
        default: password
      token:
        type: string
        description: The required token used as a credential for authorization or access to a networked resource.
      keys:
        type: map
        description: The optional list of protocol-specific keys or assertions.
        required: false
        entry_schema:
          type: string
      user:
        type: string
        description: The optional user (name or ID) used for non-token based credentials.
        required: false

  tosca.datatypes.TimeInterval:
    derived_from: tosca.datatypes.Root
    properties:
      start_time:
        type: timestamp
        required: true
      end_time:
        type: timestamp
        required: true

  tosca.datatypes.network.NetworkInfo:
    derived_from: tosca.datatypes.Root
    description: The Network type is a complex TOSCA data type used to describe logical network information.
    properties:
      network_name:
        type: string
        description: >
          The name of the logical network.
          e.g., “public”, “private”, “admin”. etc.
      network_id:
        type: string
        description: The unique ID of for the network generated by the network provider.
      addresses:
        type: list
        description: The list of IP addresses assigned from the underlying network.
        entry_schema:
          type: string

  tosca.datatypes.network.PortInfo:
    derived_from: tosca.datatypes.Root
    description: The PortInfo type is a complex TOSCA data type used to describe network port information.
    properties:
      port_name:
        type: string
        description: The logical network port name.
      port_id:
        type: string
        description: The unique ID for the network port generated by the network provider.
      network_id:
        type: string
        description: The unique ID for the network.
      mac_address:
        type: string
        description: The unique media access control address (MAC address) assigned to the port.
      addresses:
        type: list
        description: The list of IP address(es) assigned to the port.
        entry_schema:
          type: string

  tosca.datatypes.network.PortDef:
    derived_from: integer
    description: The PortDef type is a TOSCA data Type used to define a network port.
    constraints:
      - in_range: [ 1, 65535 ]

  tosca.datatypes.network.PortSpec:
    derived_from: tosca.datatypes.Root
    description: The PortSpec type is a complex TOSCA data Type used when describing port specifications for a network connection.
    properties:
      protocol:
        type: string
        description: The required protocol used on the port.
        required: true
        default: tcp
        constraints:
          - valid_values: [ udp, tcp, igmp ]
      target:
        type: tosca.datatypes.network.PortDef
        description: The optional target port.
      target_range:
        type: range
        description: The optional range for target port.
        constraints:
          - in_range: [ 1, 65535 ]
      source:
        type: tosca.datatypes.network.PortDef
        description: The optional target port.
      source_range:
        type: range
        description: The optional range for source port.
        constraints:
          - in_range: [ 1, 65535 ]

```

Here are the resulting generated structures (with the builtin types option):

```go
// Code generated by tdt2go
// DO NOT EDIT! ANY CHANGES MAY BE OVERWRITTEN.

package datatypes

import (
   "time"
)

// Credential is the generated representation of tosca.datatypes.Credential data type
type Credential struct {
    Root
    Keys      map[string]string `mapstructure:"keys" json:"keys,omitempty"`
    Protocol  string            `mapstructure:"protocol" json:"protocol,omitempty"`
    Token     string            `mapstructure:"token" json:"token,omitempty"`
    TokenType string            `mapstructure:"token_type" json:"token_type,omitempty"`
    User      string            `mapstructure:"user" json:"user,omitempty"`
}

// Root is the generated representation of tosca.datatypes.Root data type
type Root struct {
}

// TimeInterval is the generated representation of tosca.datatypes.TimeInterval data type
type TimeInterval struct {
    Root
    EndTime   time.Time `mapstructure:"end_time" json:"end_time,omitempty"`
    StartTime time.Time `mapstructure:"start_time" json:"start_time,omitempty"`
}

// Json is the generated representation of tosca.datatypes.json data type
type Json string

// NetworkInfo is the generated representation of tosca.datatypes.network.NetworkInfo data type
type NetworkInfo struct {
    Root
    Addresses   []string `mapstructure:"addresses" json:"addresses,omitempty"`
    NetworkId   string   `mapstructure:"network_id" json:"network_id,omitempty"`
    NetworkName string   `mapstructure:"network_name" json:"network_name,omitempty"`
}

// PortDef is the generated representation of tosca.datatypes.network.PortDef data type
type PortDef int

// PortInfo is the generated representation of tosca.datatypes.network.PortInfo data type
type PortInfo struct {
    Root
    Addresses  []string `mapstructure:"addresses" json:"addresses,omitempty"`
    MacAddress string   `mapstructure:"mac_address" json:"mac_address,omitempty"`
    NetworkId  string   `mapstructure:"network_id" json:"network_id,omitempty"`
    PortId     string   `mapstructure:"port_id" json:"port_id,omitempty"`
    PortName   string   `mapstructure:"port_name" json:"port_name,omitempty"`
}

// PortSpec is the generated representation of tosca.datatypes.network.PortSpec data type
type PortSpec struct {
    Root
    Protocol    string  `mapstructure:"protocol" json:"protocol,omitempty"`
    Source      PortDef `mapstructure:"source" json:"source,omitempty"`
    SourceRange Range   `mapstructure:"source_range" json:"source_range,omitempty"`
    Target      PortDef `mapstructure:"target" json:"target,omitempty"`
    TargetRange Range   `mapstructure:"target_range" json:"target_range,omitempty"`
}

// Xml is the generated representation of tosca.datatypes.xml data type
type Xml string

// Range is the generated representation of tosca:range data type
type Range []uint64

// ScalarUnit is the generated representation of tosca:scalar-unit data type
type ScalarUnit string

// ScalarUnitBitRate is the generated representation of tosca:scalar-unit.bitrate data type
type ScalarUnitBitRate ScalarUnit

// ScalarUnitFrequency is the generated representation of tosca:scalar-unit.frequency data type
type ScalarUnitFrequency ScalarUnit

// ScalarUnitSize is the generated representation of tosca:scalar-unit.size data type
type ScalarUnitSize ScalarUnit

// ScalarUnitTim is the generated representation of tosca:scalar-unit.time data type
type ScalarUnitTim ScalarUnit

// Version is the generated representation of tosca:version data type
type Version string

```

## Gotchas on names mappings

Name mappings allows to rename a generated Go struct name based on its TOSCA fully qualified name using regular expressions.
It is also very common to use the go generate command to generate Go code. Here are some gotchas to take into account.

- the `$` is used in regexp replacements to identify capturing groups. But `$` is interpreted by go generate and should be replaced by `$DOLLAR`
- mappings should not be single quoted as it should be using a regular shell

Here is an example:

```go
package mytoscatypes

//go:generate tdt2go -m tosca\.datatypes\.(.*)=TOSCA_$DOLLAR{1} -f struct_normative.go normative-types.yml

```

## License

tdt2go is distributed under Apache 2.0 License.

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fystia%2Ftdt2go.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fystia%2Ftdt2go?ref=badge_large)

### Open-source dependencies & attributions

tdt2go dependencies and attributions could be found here: <https://app.fossa.io/reports/9d5e1c45-ef44-4ed7-a433-fe47474d9d56>
