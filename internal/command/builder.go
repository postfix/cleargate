package command

import (
	"fmt"
	"sort"
	"strings"

	"github.com/postfix/cleargate/internal/models"
)

// BuildCommand constructs a safe argv array based on the ToolSpec and provided job values.
func BuildCommand(spec *models.ToolSpec, values map[string]interface{}) ([]string, error) {
	argv := []string{spec.Runtime.Argv0}

	// Process flags
	for _, f := range spec.Flags {
		val, provided := values[f.ID]
		if !provided {
			// If not provided, should we use default? Let's assume validation handled defaults,
			// or we only render if provided. For simplicity, we only render provided values.
			continue
		}

		renderRules := f.Cli.Render

		// Boolean flags with whenTrue
		if f.Type == "boolean" {
			bVal, ok := val.(bool)
			if ok && bVal && len(renderRules.WhenTrue) > 0 {
				argv = append(argv, renderRules.WhenTrue...)
			}
			continue
		}

		// Sequence rendering (e.g. ["-t", "{{value}}"])
		if len(renderRules.Sequence) > 0 {
			strVal := fmt.Sprintf("%v", val)
			for _, seqPart := range renderRules.Sequence {
				renderedPart := strings.ReplaceAll(seqPart, "{{value}}", strVal)
				argv = append(argv, renderedPart)
			}
			continue
		}

		// KeyValue rendering (e.g. "--define={{key}}={{value}}")
		// Not fully implementing map values in this MVP, but providing the structure
		if renderRules.KeyValue != "" {
			// e.g., val could be a map[string]string
			if mapVal, ok := val.(map[string]string); ok {
				for k, v := range mapVal {
					rendered := strings.ReplaceAll(renderRules.KeyValue, "{{key}}", k)
					rendered = strings.ReplaceAll(rendered, "{{value}}", v)
					argv = append(argv, rendered)
				}
			}
		}
	}

	// Process positionals (sort by order first)
	pos := append([]models.Positional(nil), spec.Positionals...)
	sort.Slice(pos, func(i, j int) bool {
		return pos[i].Order < pos[j].Order
	})

	for _, p := range pos {
		// e.g. source: "input_file.path"
		// Simplified: we'll look for input_file in values and assume it's a path string for this MVP.
		// In a real system, the file handler would resolve the path.
		parts := strings.Split(p.Source, ".")
		if len(parts) > 0 {
			inputId := parts[0]
			if val, provided := values[inputId]; provided {
				argv = append(argv, fmt.Sprintf("%v", val))
			}
		}
	}

	return argv, nil
}
