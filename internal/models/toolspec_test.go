package models

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestParseToolSpec(t *testing.T) {
	yamlSpec := `
apiVersion: cleargate.dev/v1
kind: ToolSpec
metadata:
  name: pandoc
  displayName: Pandoc
  version: "3.1.12"
runtime:
  executable: /usr/bin/pandoc
  argv0: pandoc
inputs:
  - id: input_file
    type: file
    required: true
    destination: input/
flags:
  - id: output_format
    type: enum
    required: true
    default: html
    values: [html, pdf]
    cli:
      render:
        sequence: ["-t", "{{value}}"]
`

	decoder := yaml.NewDecoder(strings.NewReader(yamlSpec))
	decoder.KnownFields(true)

	var spec ToolSpec
	err := decoder.Decode(&spec)
	if err != nil {
		t.Fatalf("Failed to parse YAML: %v", err)
	}

	if spec.Metadata.Name != "pandoc" {
		t.Errorf("Expected metadata.name to be 'pandoc', got %s", spec.Metadata.Name)
	}

	if spec.Runtime.Executable != "/usr/bin/pandoc" {
		t.Errorf("Expected runtime.executable to be '/usr/bin/pandoc', got %s", spec.Runtime.Executable)
	}

	if len(spec.Inputs) != 1 || spec.Inputs[0].ID != "input_file" {
		t.Errorf("Expected 1 input with ID 'input_file'")
	}

	if len(spec.Flags) != 1 || spec.Flags[0].ID != "output_format" {
		t.Errorf("Expected 1 flag with ID 'output_format'")
	}

	if spec.Flags[0].Cli.Render.Sequence[0] != "-t" {
		t.Errorf("Expected cli.render.sequence to contain '-t'")
	}
}
