package llm

import (
	"context"
	"strings"

	"github.com/postfix/cleargate/internal/models"
	"gopkg.in/yaml.v3"
)

// TemplateAssistant defines the interface for generating ToolSpecs via LLM.
type TemplateAssistant interface {
	GenerateDraft(ctx context.Context, helpText string) (*models.ToolSpec, error)
}

// MockAssistant is a dummy implementation for testing without an actual LLM.
type MockAssistant struct {
	StaticDraftYAML string
}

func NewMockAssistant(staticYAML string) *MockAssistant {
	return &MockAssistant{StaticDraftYAML: staticYAML}
}

func (m *MockAssistant) GenerateDraft(ctx context.Context, helpText string) (*models.ToolSpec, error) {
	decoder := yaml.NewDecoder(strings.NewReader(m.StaticDraftYAML))
	decoder.KnownFields(true)

	var spec models.ToolSpec
	if err := decoder.Decode(&spec); err != nil {
		return nil, err
	}

	return &spec, nil
}
