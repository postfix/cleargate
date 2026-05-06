package llm

import (
	"context"
	"testing"
)

func TestMockAssistant(t *testing.T) {
	validYAML := `
apiVersion: cleargate.dev/v1
kind: ToolSpec
metadata:
  name: mocktool
  displayName: Mock Tool
  version: "1.0.0"
runtime:
  executable: /bin/mock
  argv0: mock
`
	assistant := NewMockAssistant(validYAML)

	spec, err := assistant.GenerateDraft(context.Background(), "mock --help")
	if err != nil {
		t.Fatalf("GenerateDraft failed: %v", err)
	}

	if spec.Metadata.Name != "mocktool" {
		t.Errorf("Expected name 'mocktool', got %s", spec.Metadata.Name)
	}

	if spec.Runtime.Executable != "/bin/mock" {
		t.Errorf("Expected executable '/bin/mock', got %s", spec.Runtime.Executable)
	}
}
