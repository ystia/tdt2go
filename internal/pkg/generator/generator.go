package generator

import (
	"bytes"
	"fmt"
	"go/format"
	"strings"
	"text/template"

	"github.com/ystia/tdt2go/internal/pkg/model"
)

// Generator is the generator used to convert model.DataTypes into Go source file
type Generator struct {
}

// GenerateFile generates a formatted Go source file based on the given model.File representation
func (*Generator) GenerateFile(f model.File) ([]byte, error) {
	t := template.New("generator")
	t.Funcs(template.FuncMap{
		"asComment": asComment,
	})
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

func asComment(input string) string {
	return strings.ReplaceAll(input, "\n", "\n// ")
}
