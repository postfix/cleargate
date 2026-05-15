package models

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestToolSpecValidation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name    string
		spec    ToolSpec
		wantErr bool
	}{
		{
			name: "Valid ToolSpec",
			spec: ToolSpec{
				APIVersion: "v1",
				Kind:       "Tool",
				Metadata: Metadata{
					Name:    "testtool",
					Version: "1.0",
				},
				Runtime: Runtime{
					Executable:     "echo",
					TimeoutSeconds: 5,
				},
			},
			wantErr: false,
		},
		{
			name: "Missing required Metadata Name",
			spec: ToolSpec{
				APIVersion: "v1",
				Kind:       "Tool",
				Metadata: Metadata{
					Version: "1.0",
				},
				Runtime: Runtime{
					Executable:     "echo",
					TimeoutSeconds: 5,
				},
			},
			wantErr: true,
		},
		{
			name: "TimeoutSeconds too low",
			spec: ToolSpec{
				APIVersion: "v1",
				Kind:       "Tool",
				Metadata: Metadata{
					Name:    "testtool",
					Version: "1.0",
				},
				Runtime: Runtime{
					Executable:     "echo",
					TimeoutSeconds: 0,
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid Input missing ID",
			spec: ToolSpec{
				APIVersion: "v1",
				Kind:       "Tool",
				Metadata: Metadata{
					Name:    "testtool",
					Version: "1.0",
				},
				Runtime: Runtime{
					Executable:     "echo",
					TimeoutSeconds: 5,
				},
				Inputs: []Input{
					{Type: "file"},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.spec)
			if (err != nil) != tt.wantErr {
				t.Errorf("validator.Struct() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
