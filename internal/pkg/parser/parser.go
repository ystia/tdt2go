package parser

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"sort"
	"strings"
	"unicode"

	"github.com/serenize/snaker"
	"github.com/ystia/yorc/v4/tosca"
	"gopkg.in/yaml.v3"

	"github.com/ystia/tdt2go/internal/pkg/model"
)

type dtSlice []model.DataType

func (p dtSlice) Len() int           { return len(p) }
func (p dtSlice) Less(i, j int) bool { return p[i].FQDTN < p[j].FQDTN }
func (p dtSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type dtFieldsSlice []model.Field

func (p dtFieldsSlice) Len() int           { return len(p) }
func (p dtFieldsSlice) Less(i, j int) bool { return p[i].OriginalName < p[j].OriginalName }
func (p dtFieldsSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Parser is the TOSCA parser used to extract model.DataType from a TOSCA definition file
type Parser struct {
	// IncludePatterns is a list of regular expression patterns that fully qualified names of TOSCA datatypes should
	// validates to be included.
	// If no patterns are provided then all datatypes are considered.
	// IncludePatterns have the precedence over ExcludePatterns.
	IncludePatterns []string
	// ExcludePatterns is a list of regular expression patterns that fully qualified names of TOSCA datatypes should
	// validates to be excluded.
	// If no patterns are provided then all datatypes are considered.
	// IncludePatterns have the precedence over ExcludePatterns.
	ExcludePatterns []string
	// NameMappings are regular expressions applied to TOSCA datatype fully qualified names to transform them into Go struct names
	NameMappings map[string]string
}

func (p *Parser) nameValidatesPatterns(dtName string) (bool, error) {
	if len(p.IncludePatterns) > 0 {
		for _, i := range p.IncludePatterns {
			matched, err := regexp.MatchString(i, dtName)
			if err != nil {
				return false, fmt.Errorf("invalid include pattern %q: %w", i, err)
			}
			if matched {
				return true, nil
			}
		}
		return false, nil
	}
	if len(p.ExcludePatterns) > 0 {
		for _, e := range p.ExcludePatterns {
			matched, err := regexp.MatchString(e, dtName)
			if err != nil {
				return false, fmt.Errorf("invalid exclude pattern %q: %w", e, err)
			}
			if matched {
				return false, nil
			}
		}
	}
	return true, nil
}

// ParseTypes parses a TOSCA definition file and extracts a list of model.DataType.
//
// Only the given TOSCA file is analyzed, TOSCA imports are not taken into account.
func (p *Parser) ParseTypes(filePath string) ([]model.DataType, error) {
	topo, err := p.parseTopology(filePath)
	if err != nil {
		return nil, err
	}
	ts := make(dtSlice, 0)
	for dtName, dt := range topo.DataTypes {
		valid, err := p.nameValidatesPatterns(dtName)
		if err != nil {
			return nil, err
		}
		if !valid {
			continue
		}
		ts = append(ts, model.DataType{
			Name:        p.convertDTName(dtName),
			FQDTN:       dtName,
			DerivedFrom: p.convertTOSCAType(dt.DerivedFrom),
			Fields:      p.convertDTFields(dt.Properties),
		})
	}
	sort.Sort(ts)
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
	fields := make(dtFieldsSlice, 0)
	for pName, prop := range props {
		f := model.Field{
			Name:         convertToGoIdentifier(pName),
			OriginalName: pName,
			Type:         p.convertDTPropType(prop),
		}
		fields = append(fields, f)
	}
	sort.Sort(fields)
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
	case "version":
		return "Version"
	case "range":
		return "Range"
	case "scalar-unit":
		return "ScalarUnit"
	case "scalar-unit.size":
		return "ScalarUnitSize"
	case "scalar-unit.time":
		return "ScalarUnitTime"
	case "scalar-unit.frequency":
		return "ScalarUnitFrequency"
	case "scalar-unit.bitrate":
		return "ScalarUnitBitRate"
	}
	return p.convertDTName(t)
}

func convertToGoIdentifier(name string) string {
	// Replace all non letter non digit caracters to _
	b := strings.Builder{}
	for i, c := range name {
		if i == 0 && !unicode.IsLetter(c) {
			// Use canonical exporting prefix for such rare cases
			b.WriteRune('X')
			if unicode.IsDigit(c) {
				b.WriteRune(c)
			}
		} else if !unicode.IsLetter(c) && !unicode.IsDigit(c) {
			b.WriteRune('_')
		} else {
			b.WriteRune(c)
		}
	}
	// Convert using snaker
	return snaker.SnakeToCamel(b.String())
}

func (p *Parser) convertDTName(dtName string) string {
	name := p.applyNameMappings(dtName)
	s := strings.Split(name, ".")
	name = s[len(s)-1]
	name = convertToGoIdentifier(name)
	return name
}

func (p *Parser) applyNameMappings(dtName string) string {
	for pattern, subst := range p.NameMappings {
		re := regexp.MustCompile(pattern)
		dtName = re.ReplaceAllString(dtName, subst)
	}
	return dtName
}
