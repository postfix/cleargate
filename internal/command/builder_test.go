package command

import (
	"reflect"
	"testing"

	"github.com/postfix/cleargate/internal/models"
)

func TestBuildCommand(t *testing.T) {
	spec := &models.ToolSpec{
		Runtime: models.Runtime{
			Argv0: "pandoc",
		},
		Flags: []models.Flag{
			{
				ID:   "output_format",
				Type: "enum",
				Cli: models.CliRender{
					Render: models.RenderRules{
						Sequence: []string{"-t", "{{value}}"},
					},
				},
			},
			{
				ID:   "standalone",
				Type: "boolean",
				Cli: models.CliRender{
					Render: models.RenderRules{
						WhenTrue: []string{"--standalone"},
					},
				},
			},
		},
		Positionals: []models.Positional{
			{ID: "input_path", Source: "input_file.path", Order: 100},
		},
	}

	values := map[string]interface{}{
		"output_format": "html",
		"standalone":    true,
		"input_file":    "/workspace/input/doc.md",
	}

	expected := []string{
		"pandoc",
		"-t",
		"html",
		"--standalone",
		"/workspace/input/doc.md",
	}

	argv, err := BuildCommand(spec, values)
	if err != nil {
		t.Fatalf("BuildCommand failed: %v", err)
	}

	if !reflect.DeepEqual(argv, expected) {
		t.Errorf("Expected argv %v, got %v", expected, argv)
	}
}
