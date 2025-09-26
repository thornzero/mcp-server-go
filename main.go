// /home/thornzero/Repositories/mewling-goat-tavern/mcp/main.go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	_ "modernc.org/sqlite"
)

type Server struct {
	db       *sql.DB
	repoRoot string
}

// Goal management inputs and outputs
type GoalsListInput struct {
	Limit int `json:"limit,omitempty" jsonschema:"Maximum number of goals to return"`
}

type GoalsListOutput struct {
	Goals []Goal `json:"goals" jsonschema:"List of active goals"`
}

type Goal struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Priority  int    `json:"priority"`
	Status    string `json:"status"`
	Notes     string `json:"notes"`
	UpdatedAt string `json:"updated_at"`
}

type GoalsAddInput struct {
	Title    string  `json:"title" jsonschema:"Goal title"`
	Priority *int    `json:"priority,omitempty" jsonschema:"Goal priority (lower = higher priority)"`
	Notes    *string `json:"notes,omitempty" jsonschema:"Additional notes for the goal"`
}

type GoalsAddOutput struct {
	ID int `json:"id" jsonschema:"ID of the created goal"`
}

type GoalsUpdateInput struct {
	ID       int     `json:"id" jsonschema:"Goal ID to update"`
	Status   *string `json:"status,omitempty" jsonschema:"New status (active paused done)"`
	Notes    *string `json:"notes,omitempty" jsonschema:"Updated notes"`
	Priority *int    `json:"priority,omitempty" jsonschema:"Updated priority"`
}

type GoalsUpdateOutput struct {
	Updated int `json:"updated" jsonschema:"Number of rows updated"`
}

// ADR management inputs and outputs
type ADRsListInput struct {
	Query *string `json:"query,omitempty" jsonschema:"Search query to filter ADRs"`
}

type ADRsListOutput struct {
	ADRs []ADR `json:"adrs" jsonschema:"List of ADRs"`
}

type ADR struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Path      string `json:"path"`
	UpdatedAt string `json:"updated_at"`
}

type ADRsGetInput struct {
	ID string `json:"id" jsonschema:"ADR ID to retrieve"`
}

type ADRsGetOutput struct {
	ID      string `json:"id"`
	Path    string `json:"path"`
	Content string `json:"content"`
}

// CI management inputs and outputs
type CIRunTestsInput struct {
	Scope *string `json:"scope,omitempty" jsonschema:"Test scope (e.g. ./cmd/jukebox)"`
}

type CIRunTestsOutput struct {
	Status string `json:"status"`
	Output string `json:"output"`
}

type CILastFailureInput struct {
	RandomString *string `json:"random_string,omitempty" jsonschema:"Dummy parameter for no-parameter tools"`
}

type CILastFailureOutput struct {
	Status    string  `json:"status"`
	Scope     *string `json:"scope,omitempty"`
	StartedAt *string `json:"started_at,omitempty"`
}

// Repository search inputs and outputs
type RepoSearchInput struct {
	Q    string  `json:"q" jsonschema:"Search query"`
	Path *string `json:"path,omitempty" jsonschema:"Path to search within"`
	Max  *int    `json:"max,omitempty" jsonschema:"Maximum number of results"`
}

type RepoSearchOutput struct {
	Results []SearchResult `json:"results"`
}

type SearchResult struct {
	File  string `json:"file"`
	Line  string `json:"line"`
	Match string `json:"match"`
}

// Change logging inputs and outputs
type StateLogChangeInput struct {
	Summary string   `json:"summary" jsonschema:"Summary of the change"`
	Files   []string `json:"files,omitempty" jsonschema:"Files that were changed"`
}

type StateLogChangeOutput struct {
	OK bool `json:"ok" jsonschema:"Whether the change was logged successfully"`
}

// Markdown linting inputs and outputs
type MarkdownLintInput struct {
	Path   *string `json:"path,omitempty" jsonschema:"Path to lint (file or directory)"`
	Fix    *bool   `json:"fix,omitempty" jsonschema:"Whether to automatically fix issues"`
	Config *string `json:"config,omitempty" jsonschema:"Path to markdownlint config file"`
}

type MarkdownLintOutput struct {
	Issues []LintIssue `json:"issues"`
	Fixed  bool        `json:"fixed"`
	Path   string      `json:"path"`
}

type LintIssue struct {
	File    string `json:"file"`
	Line    int    `json:"line"`
	Column  int    `json:"column"`
	Rule    string `json:"rule"`
	Message string `json:"message"`
}

// Template system inputs and outputs
type TemplateListInput struct {
	Category *string `json:"category,omitempty"`
}

type TemplateListOutput struct {
	Templates []Template `json:"templates"`
}

type Template struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Category    string             `json:"category"`
	Content     string             `json:"content"`
	Variables   []TemplateVariable `json:"variables"`
	CreatedAt   string             `json:"created_at"`
	UpdatedAt   string             `json:"updated_at"`
}

type TemplateVariable struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	Required     bool   `json:"required"`
	DefaultValue string `json:"default_value"`
	Description  string `json:"description"`
}

type TemplateRegisterInput struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Description *string            `json:"description,omitempty"`
	Category    string             `json:"category"`
	Content     string             `json:"content"`
	Variables   []TemplateVariable `json:"variables,omitempty"`
}

type TemplateRegisterOutput struct {
	ID      string `json:"id"`
	Success bool   `json:"success"`
}

type TemplateGetInput struct {
	ID string `json:"id"`
}

type TemplateGetOutput struct {
	Template Template `json:"template"`
}

type TemplateUpdateInput struct {
	ID          string             `json:"id"`
	Name        *string            `json:"name,omitempty"`
	Description *string            `json:"description,omitempty"`
	Category    *string            `json:"category,omitempty"`
	Content     *string            `json:"content,omitempty"`
	Variables   []TemplateVariable `json:"variables,omitempty"`
}

type TemplateUpdateOutput struct {
	Updated bool `json:"updated"`
}

type TemplateDeleteInput struct {
	ID string `json:"id"`
}

type TemplateDeleteOutput struct {
	Deleted bool `json:"deleted"`
}

type TemplateApplyInput struct {
	TemplateID string                 `json:"template_id"`
	Variables  map[string]interface{} `json:"variables"`
	OutputPath *string                `json:"output_path,omitempty"`
}

type TemplateApplyOutput struct {
	Content string `json:"content"`
	Path    string `json:"path"`
}

func main() {
	// Get the directory where the executable is located
	execPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	repoRoot := filepath.Dir(execPath)

	// Initialize database
	dbPath := filepath.Join(repoRoot, ".agent", "state.db")
	_ = os.MkdirAll(filepath.Dir(dbPath), 0o755)
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Run schema
	schemaPath := filepath.Join(repoRoot, "schema.sql")
	if schema, err := os.ReadFile(schemaPath); err == nil {
		if _, err := db.Exec(string(schema)); err != nil {
			log.Printf("Warning: failed to run schema: %v", err)
		}
	}

	server := &Server{db: db, repoRoot: repoRoot}
	server.scanADRs()

	// Create MCP server
	mcpServer := mcp.NewServer(&mcp.Implementation{
		Name:    "mewling-goat-tavern-mcp",
		Version: "1.0.0",
	}, nil)

	// Add tools
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "goals_list",
		Description: "List active goals from the project",
	}, server.GoalsList)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "goals_add",
		Description: "Add a new goal to the project",
	}, server.GoalsAdd)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "goals_update",
		Description: "Update an existing goal",
	}, server.GoalsUpdate)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "adrs_list",
		Description: "List Architecture Decision Records (ADRs)",
	}, server.ADRsList)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "adrs_get",
		Description: "Get the content of a specific ADR",
	}, server.ADRsGet)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "ci_run_tests",
		Description: "Run tests for the project",
	}, server.CIRunTests)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "ci_last_failure",
		Description: "Get information about the last test failure",
	}, server.CILastFailure)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "repo_search",
		Description: "Search the repository for text patterns",
	}, server.RepoSearch)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "state_log_change",
		Description: "Log a change to the project changelog",
	}, server.StateLogChange)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "markdown_lint",
		Description: "Lint markdown files for formatting issues",
	}, server.MarkdownLint)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "template_list",
		Description: "List available markdown templates",
	}, server.TemplateList)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "template_register",
		Description: "Register a new markdown template",
	}, server.TemplateRegister)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "template_get",
		Description: "Get template details by ID",
	}, server.TemplateGet)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "template_update",
		Description: "Update an existing markdown template",
	}, server.TemplateUpdate)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "template_delete",
		Description: "Delete a markdown template",
	}, server.TemplateDelete)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "template_apply",
		Description: "Apply a template to generate markdown content",
	}, server.TemplateApply)

	// Run the server over stdin/stdout
	if err := mcpServer.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}

func (s *Server) scanADRs() {
	adrDir := filepath.Join(s.repoRoot, "ADR")
	entries, err := os.ReadDir(adrDir)
	if err != nil {
		return // ADR directory doesn't exist, skip
	}

	tx, _ := s.db.Begin()
	defer tx.Commit()
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		p := filepath.Join("ADR", e.Name())
		id := strings.ToUpper(strings.TrimSuffix(strings.ReplaceAll(e.Name(), "-", ""), ".md"))
		title := strings.TrimSuffix(strings.TrimPrefix(e.Name(), "0000-"), ".md")
		_, _ = tx.Exec(`INSERT OR REPLACE INTO adrs(id,title,path,updated_at) VALUES(?,?,?,CURRENT_TIMESTAMP)`, id, title, p)
	}
}

// Tool implementations

func (s *Server) GoalsList(ctx context.Context, req *mcp.CallToolRequest, input GoalsListInput) (*mcp.CallToolResult, GoalsListOutput, error) {
	limit := input.Limit
	if limit == 0 {
		limit = 10
	}

	query := `SELECT id,title,priority,status,notes,updated_at FROM goals WHERE status!='done' ORDER BY priority ASC, updated_at DESC LIMIT ?`
	rows, err := s.db.Query(query, limit)
	if err != nil {
		return nil, GoalsListOutput{}, err
	}
	defer rows.Close()

	var goals []Goal
	for rows.Next() {
		var g Goal
		if err := rows.Scan(&g.ID, &g.Title, &g.Priority, &g.Status, &g.Notes, &g.UpdatedAt); err != nil {
			return nil, GoalsListOutput{}, err
		}
		goals = append(goals, g)
	}

	return nil, GoalsListOutput{Goals: goals}, nil
}

func (s *Server) GoalsAdd(ctx context.Context, req *mcp.CallToolRequest, input GoalsAddInput) (*mcp.CallToolResult, GoalsAddOutput, error) {
	if strings.TrimSpace(input.Title) == "" {
		return nil, GoalsAddOutput{}, fmt.Errorf("title required")
	}

	prio := 100
	if input.Priority != nil {
		prio = *input.Priority
	}
	notes := ""
	if input.Notes != nil {
		notes = *input.Notes
	}

	res, err := s.db.Exec(`INSERT INTO goals(title,priority,notes) VALUES(?,?,?)`, input.Title, prio, notes)
	if err != nil {
		return nil, GoalsAddOutput{}, err
	}

	id, _ := res.LastInsertId()
	return nil, GoalsAddOutput{ID: int(id)}, nil
}

func (s *Server) GoalsUpdate(ctx context.Context, req *mcp.CallToolRequest, input GoalsUpdateInput) (*mcp.CallToolResult, GoalsUpdateOutput, error) {
	if input.ID == 0 {
		return nil, GoalsUpdateOutput{}, fmt.Errorf("id required")
	}

	set := []string{}
	args := []interface{}{}

	if input.Status != nil {
		set = append(set, "status=?")
		args = append(args, *input.Status)
	}
	if input.Notes != nil {
		set = append(set, "notes=?")
		args = append(args, *input.Notes)
	}
	if input.Priority != nil {
		set = append(set, "priority=?")
		args = append(args, *input.Priority)
	}

	if len(set) == 0 {
		return nil, GoalsUpdateOutput{Updated: 0}, nil
	}

	args = append(args, input.ID)
	_, err := s.db.Exec(`UPDATE goals SET `+strings.Join(set, ",")+`, updated_at=CURRENT_TIMESTAMP WHERE id=?`, args...)
	if err != nil {
		return nil, GoalsUpdateOutput{}, err
	}

	return nil, GoalsUpdateOutput{Updated: 1}, nil
}

func (s *Server) ADRsList(ctx context.Context, req *mcp.CallToolRequest, input ADRsListInput) (*mcp.CallToolResult, ADRsListOutput, error) {
	query := `SELECT id,title,path,updated_at FROM adrs`
	var rows *sql.Rows
	var err error

	if input.Query != nil && strings.TrimSpace(*input.Query) != "" {
		query += ` WHERE title LIKE '%'||?||'%' OR id LIKE '%'||?||'%'`
		rows, err = s.db.Query(query, *input.Query, *input.Query)
	} else {
		rows, err = s.db.Query(query)
	}

	if err != nil {
		return nil, ADRsListOutput{}, err
	}
	defer rows.Close()

	var adrs []ADR
	for rows.Next() {
		var adr ADR
		if err := rows.Scan(&adr.ID, &adr.Title, &adr.Path, &adr.UpdatedAt); err != nil {
			return nil, ADRsListOutput{}, err
		}
		adrs = append(adrs, adr)
	}

	return nil, ADRsListOutput{ADRs: adrs}, nil
}

func (s *Server) ADRsGet(ctx context.Context, req *mcp.CallToolRequest, input ADRsGetInput) (*mcp.CallToolResult, ADRsGetOutput, error) {
	var path string
	if err := s.db.QueryRow(`SELECT path FROM adrs WHERE id=?`, input.ID).Scan(&path); err != nil {
		return nil, ADRsGetOutput{}, err
	}

	content, err := os.ReadFile(filepath.Join(s.repoRoot, path))
	if err != nil {
		return nil, ADRsGetOutput{}, err
	}

	return nil, ADRsGetOutput{
		ID:      input.ID,
		Path:    path,
		Content: string(content),
	}, nil
}

func (s *Server) CIRunTests(ctx context.Context, req *mcp.CallToolRequest, input CIRunTestsInput) (*mcp.CallToolResult, CIRunTestsOutput, error) {
	scope := "./..."
	if input.Scope != nil && *input.Scope != "" {
		scope = *input.Scope
	}

	start := time.Now()
	cmd := exec.Command("go", "test", scope, "-count=1")
	cmd.Dir = s.repoRoot
	output, err := cmd.CombinedOutput()

	status := "pass"
	if err != nil {
		status = "fail"
	}

	// Log to database
	_, _ = s.db.Exec(`INSERT INTO ci_runs(scope,status,started_at,finished_at) VALUES(?,?,?,?)`,
		scope, status, start, time.Now())

	return nil, CIRunTestsOutput{
		Status: status,
		Output: string(output),
	}, nil
}

func (s *Server) CILastFailure(ctx context.Context, req *mcp.CallToolRequest, input CILastFailureInput) (*mcp.CallToolResult, CILastFailureOutput, error) {
	row := s.db.QueryRow(`SELECT scope,status,started_at FROM ci_runs WHERE status='fail' ORDER BY started_at DESC LIMIT 1`)

	var scope, status, started string
	if err := row.Scan(&scope, &status, &started); err != nil {
		return nil, CILastFailureOutput{Status: "none"}, nil
	}

	return nil, CILastFailureOutput{
		Status:    status,
		Scope:     &scope,
		StartedAt: &started,
	}, nil
}

func (s *Server) RepoSearch(ctx context.Context, req *mcp.CallToolRequest, input RepoSearchInput) (*mcp.CallToolResult, RepoSearchOutput, error) {
	if strings.TrimSpace(input.Q) == "" {
		return nil, RepoSearchOutput{}, fmt.Errorf("query required")
	}

	max := 50
	if input.Max != nil {
		max = *input.Max
	}

	target := s.repoRoot
	if input.Path != nil && *input.Path != "" {
		target = filepath.Join(s.repoRoot, *input.Path)
	}

	// Use ripgrep if available, else fallback to grep
	cmd := exec.Command("rg", "--line-number", "--no-heading", "--max-count", fmt.Sprint(max), input.Q, target)
	output, err := cmd.CombinedOutput()
	if err != nil && len(output) == 0 {
		cmd = exec.Command("grep", "-Rn", input.Q, target)
		output, _ = cmd.CombinedOutput()
	}

	lines := strings.Split(string(output), "\n")
	results := []SearchResult{} // Initialize as empty slice, not nil
	re := regexp.MustCompile(`^(.+?):(\d+):(.*)$`)

	for _, line := range lines {
		if matches := re.FindStringSubmatch(line); matches != nil {
			results = append(results, SearchResult{
				File:  strings.TrimPrefix(matches[1], s.repoRoot+"/"),
				Line:  matches[2],
				Match: strings.TrimSpace(matches[3]),
			})
			if len(results) >= max {
				break
			}
		}
	}

	return nil, RepoSearchOutput{Results: results}, nil
}

func (s *Server) StateLogChange(ctx context.Context, req *mcp.CallToolRequest, input StateLogChangeInput) (*mcp.CallToolResult, StateLogChangeOutput, error) {
	if strings.TrimSpace(input.Summary) == "" {
		return nil, StateLogChangeOutput{}, fmt.Errorf("summary required")
	}

	path := filepath.Join(s.repoRoot, "CHANGELOG_AGENT.md")
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, StateLogChangeOutput{}, err
	}
	defer f.Close()

	ts := time.Now().Format(time.RFC3339)
	fmt.Fprintf(f, "- %s — %s\n", ts, input.Summary)
	if len(input.Files) > 0 {
		fmt.Fprintf(f, "  - files: %s\n", strings.Join(input.Files, ", "))
	}

	return nil, StateLogChangeOutput{OK: true}, nil
}

func (s *Server) MarkdownLint(ctx context.Context, req *mcp.CallToolRequest, input MarkdownLintInput) (*mcp.CallToolResult, MarkdownLintOutput, error) {
	// Determine the path to lint
	targetPath := s.repoRoot
	if input.Path != nil && *input.Path != "" {
		targetPath = filepath.Join(s.repoRoot, *input.Path)
	}

	// Determine config file path
	configPath := filepath.Join(s.repoRoot, ".markdownlint.json")
	if input.Config != nil && *input.Config != "" {
		configPath = filepath.Join(s.repoRoot, *input.Config)
	}

	// Check if markdownlint is available
	cmd := exec.Command("which", "markdownlint")
	if err := cmd.Run(); err != nil {
		return nil, MarkdownLintOutput{}, fmt.Errorf("markdownlint not found. Please install with: npm install -g markdownlint-cli")
	}

	// Build markdownlint command
	args := []string{}

	// Add config file if it exists
	if _, err := os.Stat(configPath); err == nil {
		args = append(args, "--config", configPath)
	}

	// Add fix flag if requested
	fix := false
	if input.Fix != nil && *input.Fix {
		fix = true
		args = append(args, "--fix")
	}

	// Add target path
	args = append(args, targetPath)

	// Run markdownlint
	cmd = exec.Command("markdownlint", args...)
	output, err := cmd.CombinedOutput()

	var issues []LintIssue
	if err != nil {
		// Parse markdownlint output for issues
		lines := strings.Split(string(output), "\n")
		re := regexp.MustCompile(`^(.+?):(\d+):(\d+)\s+(.+?)\s+(.+)$`)

		for _, line := range lines {
			if matches := re.FindStringSubmatch(line); matches != nil {
				file := strings.TrimPrefix(matches[1], s.repoRoot+"/")
				lineNum, _ := strconv.Atoi(matches[2])
				colNum, _ := strconv.Atoi(matches[3])
				rule := matches[4]
				message := matches[5]

				issues = append(issues, LintIssue{
					File:    file,
					Line:    lineNum,
					Column:  colNum,
					Rule:    rule,
					Message: message,
				})
			}
		}
	}

	// Ensure we always return a non-nil slice
	if issues == nil {
		issues = []LintIssue{}
	}

	return nil, MarkdownLintOutput{
		Issues: issues,
		Fixed:  fix,
		Path:   targetPath,
	}, nil
}

// Template management functions

func (s *Server) TemplateList(ctx context.Context, req *mcp.CallToolRequest, input TemplateListInput) (*mcp.CallToolResult, TemplateListOutput, error) {
	query := `SELECT id, name, description, category, content, created_at, updated_at FROM markdown_templates`
	args := []interface{}{}

	if input.Category != nil && *input.Category != "" {
		query += ` WHERE category = ?`
		args = append(args, *input.Category)
	}

	query += ` ORDER BY category, name`

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, TemplateListOutput{}, err
	}
	defer rows.Close()

	var templates []Template
	for rows.Next() {
		var tmpl Template
		err := rows.Scan(&tmpl.ID, &tmpl.Name, &tmpl.Description, &tmpl.Category, &tmpl.Content, &tmpl.CreatedAt, &tmpl.UpdatedAt)
		if err != nil {
			return nil, TemplateListOutput{}, err
		}

		// Load variables for this template
		varRows, err := s.db.Query(`SELECT name, type, required, default_value, description FROM template_variables WHERE template_id = ?`, tmpl.ID)
		if err != nil {
			return nil, TemplateListOutput{}, err
		}
		defer varRows.Close()

		var variables []TemplateVariable
		for varRows.Next() {
			var variable TemplateVariable
			err := varRows.Scan(&variable.Name, &variable.Type, &variable.Required, &variable.DefaultValue, &variable.Description)
			if err != nil {
				return nil, TemplateListOutput{}, err
			}
			variables = append(variables, variable)
		}
		tmpl.Variables = variables
		templates = append(templates, tmpl)
	}

	if templates == nil {
		templates = []Template{}
	}

	return nil, TemplateListOutput{Templates: templates}, nil
}

func (s *Server) TemplateRegister(ctx context.Context, req *mcp.CallToolRequest, input TemplateRegisterInput) (*mcp.CallToolResult, TemplateRegisterOutput, error) {
	if strings.TrimSpace(input.ID) == "" || strings.TrimSpace(input.Name) == "" || strings.TrimSpace(input.Category) == "" || strings.TrimSpace(input.Content) == "" {
		return nil, TemplateRegisterOutput{}, fmt.Errorf("id, name, category, and content are required")
	}

	tx, err := s.db.Begin()
	if err != nil {
		return nil, TemplateRegisterOutput{}, err
	}
	defer tx.Rollback()

	// Insert template
	description := ""
	if input.Description != nil {
		description = *input.Description
	}

	_, err = tx.Exec(`INSERT INTO markdown_templates (id, name, description, category, content) VALUES (?, ?, ?, ?, ?)`,
		input.ID, input.Name, description, input.Category, input.Content)
	if err != nil {
		return nil, TemplateRegisterOutput{}, err
	}

	// Insert variables
	for _, variable := range input.Variables {
		_, err = tx.Exec(`INSERT INTO template_variables (template_id, name, type, required, default_value, description) VALUES (?, ?, ?, ?, ?, ?)`,
			input.ID, variable.Name, variable.Type, variable.Required, variable.DefaultValue, variable.Description)
		if err != nil {
			return nil, TemplateRegisterOutput{}, err
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, TemplateRegisterOutput{}, err
	}

	return nil, TemplateRegisterOutput{ID: input.ID, Success: true}, nil
}

func (s *Server) TemplateGet(ctx context.Context, req *mcp.CallToolRequest, input TemplateGetInput) (*mcp.CallToolResult, TemplateGetOutput, error) {
	if strings.TrimSpace(input.ID) == "" {
		return nil, TemplateGetOutput{}, fmt.Errorf("template ID is required")
	}

	var tmpl Template
	row := s.db.QueryRow(`SELECT id, name, description, category, content, created_at, updated_at FROM markdown_templates WHERE id = ?`, input.ID)
	err := row.Scan(&tmpl.ID, &tmpl.Name, &tmpl.Description, &tmpl.Category, &tmpl.Content, &tmpl.CreatedAt, &tmpl.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, TemplateGetOutput{}, fmt.Errorf("template not found: %s", input.ID)
		}
		return nil, TemplateGetOutput{}, err
	}

	// Load variables
	varRows, err := s.db.Query(`SELECT name, type, required, default_value, description FROM template_variables WHERE template_id = ?`, input.ID)
	if err != nil {
		return nil, TemplateGetOutput{}, err
	}
	defer varRows.Close()

	var variables []TemplateVariable
	for varRows.Next() {
		var variable TemplateVariable
		err := varRows.Scan(&variable.Name, &variable.Type, &variable.Required, &variable.DefaultValue, &variable.Description)
		if err != nil {
			return nil, TemplateGetOutput{}, err
		}
		variables = append(variables, variable)
	}
	tmpl.Variables = variables

	return nil, TemplateGetOutput{Template: tmpl}, nil
}

func (s *Server) TemplateUpdate(ctx context.Context, req *mcp.CallToolRequest, input TemplateUpdateInput) (*mcp.CallToolResult, TemplateUpdateOutput, error) {
	if strings.TrimSpace(input.ID) == "" {
		return nil, TemplateUpdateOutput{}, fmt.Errorf("template ID is required")
	}

	tx, err := s.db.Begin()
	if err != nil {
		return nil, TemplateUpdateOutput{}, err
	}
	defer tx.Rollback()

	// Build dynamic update query
	setParts := []string{}
	args := []interface{}{}

	if input.Name != nil {
		setParts = append(setParts, "name = ?")
		args = append(args, *input.Name)
	}
	if input.Description != nil {
		setParts = append(setParts, "description = ?")
		args = append(args, *input.Description)
	}
	if input.Category != nil {
		setParts = append(setParts, "category = ?")
		args = append(args, *input.Category)
	}
	if input.Content != nil {
		setParts = append(setParts, "content = ?")
		args = append(args, *input.Content)
	}

	if len(setParts) > 0 {
		setParts = append(setParts, "updated_at = CURRENT_TIMESTAMP")
		args = append(args, input.ID)

		query := fmt.Sprintf("UPDATE markdown_templates SET %s WHERE id = ?", strings.Join(setParts, ", "))
		_, err = tx.Exec(query, args...)
		if err != nil {
			return nil, TemplateUpdateOutput{}, err
		}
	}

	// Update variables if provided
	if len(input.Variables) > 0 {
		// Delete existing variables
		_, err = tx.Exec(`DELETE FROM template_variables WHERE template_id = ?`, input.ID)
		if err != nil {
			return nil, TemplateUpdateOutput{}, err
		}

		// Insert new variables
		for _, variable := range input.Variables {
			_, err = tx.Exec(`INSERT INTO template_variables (template_id, name, type, required, default_value, description) VALUES (?, ?, ?, ?, ?, ?)`,
				input.ID, variable.Name, variable.Type, variable.Required, variable.DefaultValue, variable.Description)
			if err != nil {
				return nil, TemplateUpdateOutput{}, err
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, TemplateUpdateOutput{}, err
	}

	return nil, TemplateUpdateOutput{Updated: true}, nil
}

func (s *Server) TemplateDelete(ctx context.Context, req *mcp.CallToolRequest, input TemplateDeleteInput) (*mcp.CallToolResult, TemplateDeleteOutput, error) {
	if strings.TrimSpace(input.ID) == "" {
		return nil, TemplateDeleteOutput{}, fmt.Errorf("template ID is required")
	}

	result, err := s.db.Exec(`DELETE FROM markdown_templates WHERE id = ?`, input.ID)
	if err != nil {
		return nil, TemplateDeleteOutput{}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, TemplateDeleteOutput{}, err
	}

	return nil, TemplateDeleteOutput{Deleted: rowsAffected > 0}, nil
}

func (s *Server) TemplateApply(ctx context.Context, req *mcp.CallToolRequest, input TemplateApplyInput) (*mcp.CallToolResult, TemplateApplyOutput, error) {
	if strings.TrimSpace(input.TemplateID) == "" {
		return nil, TemplateApplyOutput{}, fmt.Errorf("template ID is required")
	}

	// Get template
	var tmpl Template
	row := s.db.QueryRow(`SELECT id, name, description, category, content, created_at, updated_at FROM markdown_templates WHERE id = ?`, input.TemplateID)
	err := row.Scan(&tmpl.ID, &tmpl.Name, &tmpl.Description, &tmpl.Category, &tmpl.Content, &tmpl.CreatedAt, &tmpl.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, TemplateApplyOutput{}, fmt.Errorf("template not found: %s", input.TemplateID)
		}
		return nil, TemplateApplyOutput{}, err
	}

	// Parse and execute template
	t, err := template.New(tmpl.ID).Parse(tmpl.Content)
	if err != nil {
		return nil, TemplateApplyOutput{}, fmt.Errorf("template parse error: %v", err)
	}

	var result strings.Builder
	err = t.Execute(&result, input.Variables)
	if err != nil {
		return nil, TemplateApplyOutput{}, fmt.Errorf("template execution error: %v", err)
	}

	content := result.String()
	outputPath := ""

	// Write to file if output path specified
	if input.OutputPath != nil && *input.OutputPath != "" {
		fullPath := filepath.Join(s.repoRoot, *input.OutputPath)
		err = os.WriteFile(fullPath, []byte(content), 0644)
		if err != nil {
			return nil, TemplateApplyOutput{}, fmt.Errorf("failed to write file: %v", err)
		}
		outputPath = fullPath
	}

	return nil, TemplateApplyOutput{Content: content, Path: outputPath}, nil
}
