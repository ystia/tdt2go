package generator

import (
	"testing"

	"github.com/ystia/tdt2go/internal/pkg/model"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/golden"
)

func TestGenerator_GenerateFile(t *testing.T) {
	type args struct {
		f model.File
	}
	tests := []struct {
		name    string
		g       *Generator
		args    args
		wantErr bool
	}{
		{"EmptyFile", &Generator{}, args{model.File{Package: "something"}}, false},
		{"SimpleDataType", &Generator{}, args{
			model.File{
				Package: "simple",
				DataTypes: []model.DataType{
					{
						Name:  "MyDT",
						FQDTN: "org.ystia.datatypes.MyDT",
						Fields: []model.Field{
							{Name: "F1", Type: "string"},
							{Name: "F2", Type: "int"},
						},
					},
				},
			},
		}, false},
		{"FieldsTags", &Generator{}, args{
			model.File{
				Package: "simple",
				DataTypes: []model.DataType{
					{
						Name:  "MyDT",
						FQDTN: "org.ystia.datatypes.MyDT",
						Fields: []model.Field{
							{Name: "F1", OriginalName: "f1", Type: "string"},
							{Name: "F2", OriginalName: "my_f2", Type: "int"},
						},
					},
				},
			},
		}, false},
		{"WithImports", &Generator{}, args{
			model.File{
				Package: "simple",
				Imports: []string{"fmt", "time"},
				DataTypes: []model.DataType{
					{
						Name:  "MyDT",
						FQDTN: "org.ystia.datatypes.MyDT",
						Fields: []model.Field{
							{Name: "F1", OriginalName: "f1", Type: "string"},
							{Name: "F2", OriginalName: "my_f2", Type: "time.Date"},
						},
					},
				},
			},
		}, false},
		{"DerivedDataType", &Generator{}, args{
			model.File{
				Package: "simple",
				DataTypes: []model.DataType{
					{
						Name:  "MyDT",
						FQDTN: "org.ystia.datatypes.MyDT",
						Fields: []model.Field{
							{Name: "F1", Type: "string"},
							{Name: "F2", Type: "int"},
						},
					},
					{
						Name:        "MyDerivedDT",
						FQDTN:       "org.ystia.datatypes.MyDerivedDT",
						DerivedFrom: "MyDT",
						Fields: []model.Field{
							{Name: "F3", Type: "[]string"},
						},
					},
				},
			},
		}, false},
		{"DerivedFromBuildtin", &Generator{}, args{
			model.File{
				Package: "simple",
				DataTypes: []model.DataType{
					{
						Name:        "MyDerivedDT",
						FQDTN:       "org.ystia.datatypes.MyDerivedDT",
						DerivedFrom: "string",
					},
				},
			},
		}, false},
		{"WithDescriptions", &Generator{}, args{
			model.File{
				Package: "simple",
				DataTypes: []model.DataType{
					{
						Name:        "MyDT",
						FQDTN:       "org.ystia.datatypes.MyDT",
						Description: "A oneliner description",
						Fields: []model.Field{
							{Name: "F1", Type: "string"},
							{Name: "F2", Type: "int", Description: "A multiline\ndescription"},
						},
					},
					{
						Name:        "MyDerivedDT",
						FQDTN:       "org.ystia.datatypes.MyDerivedDT",
						DerivedFrom: "MyDT",
						Description: "A multiline\ndescription",
						Fields: []model.Field{
							{Name: "F3", Type: "[]string", Description: "A oneliner description"},
						},
					},
				},
			},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Generator{}
			got, err := g.GenerateFile(tt.args.f)
			if (err != nil) != tt.wantErr {
				t.Errorf("Generator.GenerateFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				assert.Assert(t, golden.String(string(got), "golden/"+tt.name))
			}
		})
	}
}
