package generator

import (
	"bytes"
	"fmt"
	"go/format"
	"text/template"

	"github.com/ystia/tdt2go/internal/pkg/model"
)

type File struct {
	Package   string
	Imports   []string
	DataTypes []model.DataType
}

type Generator struct {
}

func (*Generator) GenerateFile(f File) ([]byte, error) {
	t := template.New("generator")
	t = template.Must(t.Parse(fileTemplate))

	b := &bytes.Buffer{}
	err := t.Execute(b, f)
	if err != nil {
		return nil, fmt.Errorf("failed to generate file, templating failed: %w", err)
	}

	result, err := format.Source(b.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to format generated file: %w", err)
	}

	return result, nil
}
