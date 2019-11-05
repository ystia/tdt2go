package parser

import (
	"sort"
	"testing"

	"github.com/ystia/tdt2go/internal/pkg/model"

	"gotest.tools/v3/assert"
)

type dtSlice []model.DataType

func (p dtSlice) Len() int           { return len(p) }
func (p dtSlice) Less(i, j int) bool { return p[i].FQDTN < p[j].FQDTN }
func (p dtSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type dtFieldsSlice []model.Field

func (p dtFieldsSlice) Len() int           { return len(p) }
func (p dtFieldsSlice) Less(i, j int) bool { return p[i].OriginalName < p[j].OriginalName }
func (p dtFieldsSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func TestParser_ParseTypes(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		p       *Parser
		args    args
		want    []model.DataType
		wantErr bool
	}{
		{"NoTOSCAFile", &Parser{}, args{"testdata/donotexists.yaml"}, nil, true},
		{"InvalidTOSCAFile", &Parser{}, args{"testdata/invalid.yaml"}, nil, true},
		{"TestParseNormativeLight", &Parser{}, args{"testdata/normative-light.yaml"}, []model.DataType{
			model.DataType{
				Name:        "Credential",
				FQDTN:       "tosca.datatypes.Credential",
				DerivedFrom: "Root",
				Fields: []model.Field{
					model.Field{
						Name:         "Protocol",
						OriginalName: "protocol",
						Type:         "string",
					},
					model.Field{
						Name:         "TokenType",
						OriginalName: "token_type",
						Type:         "string",
					},
					model.Field{
						Name:         "Token",
						OriginalName: "token",
						Type:         "string",
					},
					model.Field{
						Name:         "Keys",
						OriginalName: "keys",
						Type:         "map[string]string",
					},
					model.Field{
						Name:         "User",
						OriginalName: "user",
						Type:         "string",
					},
				},
			},
			model.DataType{
				Name:   "Root",
				FQDTN:  "tosca.datatypes.Root",
				Fields: []model.Field{},
			},
			model.DataType{
				Name:        "TimeInterval",
				FQDTN:       "tosca.datatypes.TimeInterval",
				DerivedFrom: "Root",
				Fields: []model.Field{
					model.Field{
						Name:         "StartTime",
						OriginalName: "start_time",
						Type:         "time.Time",
					},
					model.Field{
						Name:         "EndTime",
						OriginalName: "end_time",
						Type:         "time.Time",
					},
				},
			},
		}, false},
		{"TestExtraToscaTypes", &Parser{}, args{"testdata/extratypes.yaml"}, []model.DataType{
			model.DataType{
				Name:        "SpecificTypes",
				FQDTN:       "tosca.datatypes.SpecificTypes",
				DerivedFrom: "Root",
				Fields: []model.Field{
					model.Field{
						Name:         "TestAList",
						OriginalName: "test_a_list",
						Type:         "[]int",
					},
					model.Field{
						Name:         "ValidBool",
						OriginalName: "valid_bool",
						Type:         "bool",
					},
					model.Field{
						Name:         "ANumber",
						OriginalName: "a_number",
						Type:         "float64",
					},
					model.Field{
						Name:         "AnotherType",
						OriginalName: "another_type",
						Type:         "Credential",
					},
				},
			},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{}
			got, err := p.ParseTypes(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.ParseTypes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				sortSlices(got, tt.want)
				assert.DeepEqual(t, got, tt.want)
			}
		})
	}
}

func sortSlices(got, want []model.DataType) {
	var cGot, cWant dtSlice
	cGot = got
	cWant = want
	sort.Sort(cGot)
	sort.Sort(cWant)

	for _, dt := range got {
		var s dtFieldsSlice = dt.Fields
		sort.Sort(s)
	}
	for _, dt := range want {
		var s dtFieldsSlice = dt.Fields
		sort.Sort(s)
	}
}
