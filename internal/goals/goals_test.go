package goals

import (
	"context"
	"testing"

	"github.com/thornzero/mcp-server-go/internal/server"
	"github.com/thornzero/mcp-server-go/internal/types"
)

func TestGoalsHandler_GoalsList(t *testing.T) {
	// Setup
	tempDir := t.TempDir()
	srv, err := server.NewServer(tempDir)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}
	defer srv.Close()

	handler := NewGoalsHandler(srv)

	tests := []struct {
		name      string
		input     types.GoalsListInput
		wantError bool
	}{
		{
			name: "List goals with no limit",
			input: types.GoalsListInput{
				Limit: 0,
			},
			wantError: false,
		},
		{
			name: "List goals with limit",
			input: types.GoalsListInput{
				Limit: 10,
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := handler.GoalsList(context.Background(), nil, tt.input)

			if tt.wantError {
				if err == nil {
					t.Errorf("GoalsList() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("GoalsList() unexpected error: %v", err)
				return
			}

			// Goals slice can be nil or empty when no goals exist - both are valid
			// We just verify the function didn't error
		})
	}
}

func TestGoalsHandler_GoalsAdd(t *testing.T) {
	// Setup
	tempDir := t.TempDir()
	srv, err := server.NewServer(tempDir)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}
	defer srv.Close()

	handler := NewGoalsHandler(srv)

	tests := []struct {
		name      string
		input     types.GoalsAddInput
		wantError bool
	}{
		{
			name: "Add valid goal",
			input: types.GoalsAddInput{
				Title: "Test Goal",
			},
			wantError: false,
		},
		{
			name: "Add goal with priority",
			input: types.GoalsAddInput{
				Title:    "High Priority Goal",
				Priority: func() *int { p := 1; return &p }(),
			},
			wantError: false,
		},
		{
			name: "Add goal with notes",
			input: types.GoalsAddInput{
				Title: "Goal with Notes",
				Notes: func() *string { n := "Some notes"; return &n }(),
			},
			wantError: false,
		},
		{
			name: "Add goal with empty title",
			input: types.GoalsAddInput{
				Title: "",
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, output, err := handler.GoalsAdd(context.Background(), nil, tt.input)

			if tt.wantError {
				if err == nil {
					t.Errorf("GoalsAdd() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("GoalsAdd() unexpected error: %v", err)
				return
			}

			if output.ID <= 0 {
				t.Errorf("GoalsAdd() returned invalid ID: %v", output.ID)
			}
		})
	}
}

func TestGoalsHandler_GoalsUpdate(t *testing.T) {
	// Setup
	tempDir := t.TempDir()
	srv, err := server.NewServer(tempDir)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}
	defer srv.Close()

	handler := NewGoalsHandler(srv)

	// First, add a goal to update
	addInput := types.GoalsAddInput{
		Title: "Goal to Update",
	}
	_, addOutput, err := handler.GoalsAdd(context.Background(), nil, addInput)
	if err != nil {
		t.Fatalf("Failed to add goal for update test: %v", err)
	}

	tests := []struct {
		name      string
		input     types.GoalsUpdateInput
		wantError bool
	}{
		{
			name: "Update goal status",
			input: types.GoalsUpdateInput{
				ID:     addOutput.ID,
				Status: func() *string { s := "done"; return &s }(),
			},
			wantError: false,
		},
		{
			name: "Update goal notes",
			input: types.GoalsUpdateInput{
				ID:    addOutput.ID,
				Notes: func() *string { n := "Updated notes"; return &n }(),
			},
			wantError: false,
		},
		{
			name: "Update non-existent goal",
			input: types.GoalsUpdateInput{
				ID:     99999,
				Status: func() *string { s := "done"; return &s }(),
			},
			wantError: false, // Should not error, just update 0 rows
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, output, err := handler.GoalsUpdate(context.Background(), nil, tt.input)

			if tt.wantError {
				if err == nil {
					t.Errorf("GoalsUpdate() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("GoalsUpdate() unexpected error: %v", err)
				return
			}

			if output.Updated < 0 {
				t.Errorf("GoalsUpdate() returned negative updated count: %v", output.Updated)
			}
		})
	}
}
