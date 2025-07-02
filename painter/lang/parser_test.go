package lang_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/roman-mazur/architecture-lab-3/painter"
	"github.com/roman-mazur/architecture-lab-3/painter/lang"
)

func TestParser_Parse(t *testing.T) {
	p := lang.Parser{}

	tests := []struct {
		name        string
		input       string
		wantOpsNum  int
		wantErr     bool
		firstOpType interface{}
	}{
		{
			name:       "empty input",
			input:      "",
			wantOpsNum: 0,
			wantErr:    false,
		},
		{
			name:        "white and green",
			input:       "white\ngreen\n",
			wantOpsNum:  2,
			wantErr:     false,
			firstOpType: painter.FillBackground{},
		},
		{
			name:        "bgrect valid",
			input:       "bgrect 0.1 0.2 0.3 0.4\n",
			wantOpsNum:  1,
			wantErr:     false,
			firstOpType: painter.BgRect{},
		},
		{
			name:        "figure valid",
			input:       "figure 100 200\n",
			wantOpsNum:  1,
			wantErr:     false,
			firstOpType: painter.DrawT180{},
		},
		{
			name:        "move valid",
			input:       "move 10 20\n",
			wantOpsNum:  1,
			wantErr:     false,
			firstOpType: painter.Move{},
		},
		{
			name:        "reset valid",
			input:       "reset\n",
			wantOpsNum:  1,
			wantErr:     false,
			firstOpType: painter.Reset{},
		},
		{
			name:        "update valid",
			input:       "update\n",
			wantOpsNum:  1,
			wantErr:     false,
			firstOpType: painter.UpdateOp,
		},
		{
			name:    "unknown command",
			input:   "foobar\n",
			wantErr: true,
		},
		{
			name:    "bgrect missing args",
			input:   "bgrect 0.1 0.2 0.3\n",
			wantErr: true,
		},
		{
			name:    "figure invalid arg",
			input:   "figure 100 abc\n",
			wantErr: true,
		},
		{
			name:        "complex valid script",
			input:       "white\nbgrect 0.2 0.2 0.6 0.6\nfigure 400 400\nupdate\n",
			wantOpsNum:  4,
			wantErr:     false,
			firstOpType: painter.FillBackground{},
		},
		{
			name:        "comments and empty lines ignored",
			input:       "# this is a comment\n\nwhite\n# another\nupdate\n",
			wantOpsNum:  2,
			wantErr:     false,
			firstOpType: painter.FillBackground{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ops, err := p.Parse(strings.NewReader(tt.input))
			if (err != nil) != tt.wantErr {
				t.Fatalf("Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}

			if len(ops) != tt.wantOpsNum {
				t.Errorf("Parse() returned %d operations, want %d", len(ops), tt.wantOpsNum)
			}

			if tt.wantOpsNum > 0 && tt.firstOpType != nil {
				gotType := ops[0]
				wantType := tt.firstOpType
				if gotTypeType := reflect.TypeOf(gotType); gotTypeType != reflect.TypeOf(wantType) {
					t.Errorf("First operation type = %T, want %T", gotType, wantType)
				}
			}

			if tt.name == "complex valid script" {
				wantTypes := []reflect.Type{
					reflect.TypeOf(painter.FillBackground{}),
					reflect.TypeOf(painter.BgRect{}),
					reflect.TypeOf(painter.DrawT180{}),
					reflect.TypeOf(painter.UpdateOp),
				}
				for i, op := range ops {
					if reflect.TypeOf(op) != wantTypes[i] {
						t.Errorf("operation %d has type %T, want %v", i, op, wantTypes[i])
					}
				}
			}
		})
	}
}
