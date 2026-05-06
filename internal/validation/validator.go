package validation

import (
	"fmt"

	"github.com/postfix/cleargate/internal/models"
)

// ValidateJobValues checks if the provided values map matches the requirements and types defined in the ToolSpec.
func ValidateJobValues(spec *models.ToolSpec, values map[string]interface{}) error {
	// Build a map of valid flag IDs for quick lookup
	validFlags := make(map[string]models.Flag)
	for _, f := range spec.Flags {
		validFlags[f.ID] = f
	}

	// 1. Check for unknown flags in the input values
	for k := range values {
		if _, ok := validFlags[k]; !ok {
			return fmt.Errorf("unknown flag provided: %s", k)
		}
	}

	// 2. Validate each flag defined in the spec
	for _, f := range spec.Flags {
		val, provided := values[f.ID]

		if !provided {
			if f.Required {
				return fmt.Errorf("missing required flag: %s", f.ID)
			}
			continue
		}

		// Type and value validation
		switch f.Type {
		case "string":
			if _, ok := val.(string); !ok {
				return fmt.Errorf("flag %s expects a string value", f.ID)
			}
		case "boolean":
			if _, ok := val.(bool); !ok {
				return fmt.Errorf("flag %s expects a boolean value", f.ID)
			}
		case "enum":
			strVal, ok := val.(string)
			if !ok {
				return fmt.Errorf("flag %s expects a string (enum) value", f.ID)
			}
			valid := false
			for _, allowed := range f.Values {
				if strVal == allowed {
					valid = true
					break
				}
			}
			if !valid {
				return fmt.Errorf("invalid value for enum flag %s: %s", f.ID, strVal)
			}
		// Additional types (multi-enum, integer, file) can be added here
		}
	}

	return nil
}
