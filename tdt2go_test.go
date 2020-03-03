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

package tdt2go

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"gotest.tools/v3/assert"
	"gotest.tools/v3/golden"
)

type invalidWriter struct{}

func (w *invalidWriter) Write(p []byte) (n int, err error) {
	err = fmt.Errorf("this writer always fails")
	return
}

func TestGenerateFile(t *testing.T) {
	type args struct {
		toscaFile string
		opts      []Option
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"MissingFile", args{toscaFile: "testdata/donotexist"}, true},
		{"ErrorOnWrite", args{toscaFile: "testdata/normative-light.yaml", opts: []Option{Output(&invalidWriter{})}}, true},
		{"NormativeLight", args{toscaFile: "testdata/normative-light.yaml"}, false},
		{"ChangePackage", args{toscaFile: "testdata/normative-light.yaml", opts: []Option{Package("somepkg")}}, false},
		{"IncludeExcludePatterns", args{toscaFile: "testdata/normative-light.yaml", opts: []Option{
			// Exclude all but include some
			ExcludePatterns([]string{`tosca.*`}),
			IncludePatterns([]string{`tosca\.datatypes.Root`}),
		}}, false},
		{"NameMappings", args{toscaFile: "testdata/normative-light.yaml", opts: []Option{
			NameMappings(map[string]string{`tosca\.datatypes\.(.+)`: `TOSCA_${1}`}),
		}}, false},
		{"NormativeLightPlusBuiltin", args{toscaFile: "testdata/normative-light.yaml", opts: []Option{GenerateBuiltinTypes(true)}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &strings.Builder{}
			tt.args.opts = append([]Option{Output(b)}, tt.args.opts...)
			err := GenerateFile(tt.args.toscaFile, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				assert.Assert(t, golden.String(b.String(), "golden/"+tt.name))
			}
		})
	}
}

func TestOutputToFile(t *testing.T) {
	type args struct {
		outputFile string
		perm       os.FileMode
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"CheckOutputFile", args{"testdata/generated/output.test", 0666}, false},
		{"CheckOutputFile", args{"testdata/generated/subpathdoesnotexist/output.test", 0666}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := OutputToFile(tt.args.outputFile, tt.args.perm)
			if (err != nil) != tt.wantErr {
				t.Errorf("OutputToFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				o := &Options{}
				got(o)
				assert.Assert(t, o.output != nil, "output should not be nil")
				f, ok := o.output.(*os.File)
				assert.Assert(t, ok, "output is not a file")
				assert.Equal(t, tt.args.outputFile, f.Name(), "wrong file name")
				os.Remove(f.Name())
			}
		})
	}
}
