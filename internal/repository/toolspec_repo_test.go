package repository

import (
	"testing"

	"github.com/postfix/cleargate/internal/models"
)

func TestToolSpecRepository(t *testing.T) {
	// Use in-memory DuckDB database for testing
	repo, err := NewToolSpecRepository("?access_mode=READ_WRITE")
	if err != nil {
		t.Fatalf("Failed to initialize repo: %v", err)
	}
	defer repo.Close()

	spec := &models.ToolSpec{
		Metadata: models.Metadata{
			Name:    "nmap",
			Version: "7.92",
		},
	}

	// Save Draft
	if err := repo.SaveDraft(spec); err != nil {
		t.Fatalf("Failed to save draft: %v", err)
	}

	// List Drafts
	drafts, err := repo.ListDrafts()
	if err != nil {
		t.Fatalf("Failed to list drafts: %v", err)
	}

	if len(drafts) != 1 {
		t.Errorf("Expected 1 draft, got %d", len(drafts))
	}

	if drafts[0].Name != "nmap" || drafts[0].Status != "draft" {
		t.Errorf("Unexpected draft content: %+v", drafts[0])
	}

	// Approve
	id := "nmap-7.92"
	if err := repo.Approve(id); err != nil {
		t.Fatalf("Failed to approve draft: %v", err)
	}

	// List Drafts again (should be empty)
	draftsAfterApprove, err := repo.ListDrafts()
	if err != nil {
		t.Fatalf("Failed to list drafts after approval: %v", err)
	}

	if len(draftsAfterApprove) != 0 {
		t.Errorf("Expected 0 drafts after approval, got %d", len(draftsAfterApprove))
	}
}
