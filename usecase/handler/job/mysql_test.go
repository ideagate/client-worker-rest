package job

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_mysql_Start(t *testing.T) {
	type fields struct {
		Input StartInput
	}
	tests := []struct {
		name       string
		fields     fields
		wantOutput StartOutput
		wantErr    assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &mysql{
				Input: tt.fields.Input,
			}
			gotOutput, err := j.Start()
			if !tt.wantErr(t, err, fmt.Sprintf("Start()")) {
				return
			}
			assert.Equalf(t, tt.wantOutput, gotOutput, "Start()")
		})
	}
}
