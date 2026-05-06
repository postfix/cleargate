package validation

import (
	"testing"

	"github.com/postfix/cleargate/internal/models"
)

func TestValidateJobValues(t *testing.T) {
	spec := &models.ToolSpec{
		Flags: []models.Flag{
			{ID: "output_format", Type: "enum", Required: true, Values: []string{"html", "pdf"}},
			{ID: "standalone", Type: "boolean"},
			{ID: "title", Type: "string"},
		},
	}

	tests := []struct {
		name    string
		values  map[string]interface{}
		wantErr bool
	}{
		{
			name: "Valid input",
			values: map[string]interface{}{
				"output_format": "html",
				"standalone":    true,
				"title":         "My Document",
			},
			wantErr: false,
		},
		{
			name: "Missing required flag",
			values: map[string]interface{}{
				"standalone": true,
			},
			wantErr: true,
		},
		{
			name: "Unknown flag",
			values: map[string]interface{}{
				"output_format": "html",
				"haxor_flag":    "true",
			},
			wantErr: true,
		},
		{
			name: "Invalid enum value",
			values: map[string]interface{}{
				"output_format": "docx",
			},
			wantErr: true,
		},
		{
			name: "Invalid type for boolean",
			values: map[string]interface{}{
				"output_format": "html",
				"standalone":    "true", // Should be a bool, not a string
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateJobValues(spec, tt.values)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJobValues() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
