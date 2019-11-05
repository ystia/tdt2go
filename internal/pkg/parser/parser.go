package parser

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/ystia/yorc/v4/tosca"
	"gopkg.in/yaml.v3"

	"github.com/ystia/tdt2go/internal/pkg/model"
)

type Parser struct {
}

func (p *Parser) ParseTypes(filePath string) ([]model.DataType, error) {
	topo, err := p.parseTopology(filePath)
	if err != nil {
		return nil, err
	}
	ts := make([]model.DataType, 0)
	for dtName, dt := range topo.DataTypes {

		ts = append(ts, model.DataType{
			Name:        p.convertDTName(dtName),
			FQDTN:       dtName,
			DerivedFrom: p.convertDTName(dt.DerivedFrom),
			Fields:      p.convertDTFields(dt.Properties),
		})
	}

	return ts, nil
}

func (p *Parser) parseTopology(filePath string) (*tosca.Topology, error) {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse TOSCA definition: %w", err)
	}
	topo := &tosca.Topology{}
	err = yaml.Unmarshal(b, topo)
	if err != nil {
		return nil, fmt.Errorf("failed to parse TOSCA definition: %w", err)
	}
	return topo, nil
}

func (p *Parser) convertDTFields(props map[string]tosca.PropertyDefinition) []model.Field {
	fields := make([]model.Field, 0)
	for pName, prop := range props {
		f := model.Field{
			Name:         strings.ReplaceAll(strings.Title(strings.ReplaceAll(pName, "_", " ")), " ", ""),
			OriginalName: pName,
			Type:         p.convertDTPropType(prop),
		}
		fields = append(fields, f)
	}
	return fields
}

func (p *Parser) convertDTPropType(prop tosca.PropertyDefinition) string {
	switch strings.ToLower(prop.Type) {
	case "list":
		return "[]" + p.convertTOSCAType(prop.EntrySchema.Type)
	case "map":
		return "map[string]" + p.convertTOSCAType(prop.EntrySchema.Type)
	default:
		return p.convertTOSCAType(prop.Type)
	}
}

func (p *Parser) convertTOSCAType(t string) string {
	switch t {
	case "string":
		return "string"
	case "integer":
		return "int"
	case "boolean":
		return "bool"
	case "float":
		return "float64"
	case "timestamp":
		return "time.Time"
	}
	return p.convertDTName(t)
}

func (p *Parser) convertDTName(dtName string) string {
	s := strings.Split(dtName, ".")
	name := s[len(s)-1]
	name = strings.ReplaceAll(strings.Title(strings.ReplaceAll(name, "_", " ")), " ", "")
	return name
}
