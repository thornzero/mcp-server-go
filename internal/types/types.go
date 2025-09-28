package types

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

// Preferred Tool types
type PreferredTool struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Language    string `json:"language"`
	UseCase     string `json:"use_case"`
	Priority    int    `json:"priority"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type PreferredToolsListInput struct {
	Category string `json:"category,omitempty"`
	Language string `json:"language,omitempty"`
	Limit    int    `json:"limit,omitempty"`
}

type PreferredToolsListOutput struct {
	Tools []PreferredTool `json:"tools"`
}

type PreferredToolsAddInput struct {
	Name        string `json:"name"`
	Category    string `json:"category"`
	Description string `json:"description,omitempty"`
	Language    string `json:"language,omitempty"`
	UseCase     string `json:"use_case,omitempty"`
	Priority    int    `json:"priority,omitempty"`
}

type PreferredToolsAddOutput struct {
	ID      uint `json:"id"`
	Success bool `json:"success"`
}

type PreferredToolsUpdateInput struct {
	ID          uint   `json:"id"`
	Name        string `json:"name,omitempty"`
	Category    string `json:"category,omitempty"`
	Description string `json:"description,omitempty"`
	Language    string `json:"language,omitempty"`
	UseCase     string `json:"use_case,omitempty"`
	Priority    int    `json:"priority,omitempty"`
}

type PreferredToolsUpdateOutput struct {
	Success bool `json:"success"`
}

type PreferredToolsDeleteInput struct {
	ID uint `json:"id"`
}

type PreferredToolsDeleteOutput struct {
	Success bool `json:"success"`
}

// Cursor Rules types
type CursorRule struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Content     string `json:"content"`
	Tags        string `json:"tags"`
	Source      string `json:"source"`
	IsActive    bool   `json:"is_active"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type CursorRulesListInput struct {
	Category string `json:"category,omitempty"`
	Tags     string `json:"tags,omitempty"`
	Source   string `json:"source,omitempty"`
	Active   *bool  `json:"active,omitempty"`
	Limit    int    `json:"limit,omitempty"`
}

type CursorRulesListOutput struct {
	Rules []CursorRule `json:"rules"`
}

type CursorRulesAddInput struct {
	Name        string `json:"name"`
	Category    string `json:"category"`
	Description string `json:"description,omitempty"`
	Content     string `json:"content"`
	Tags        string `json:"tags,omitempty"`
	Source      string `json:"source,omitempty"`
	IsActive    *bool  `json:"is_active,omitempty"`
}

type CursorRulesAddOutput struct {
	ID      uint `json:"id"`
	Success bool `json:"success"`
}

type CursorRulesUpdateInput struct {
	ID          uint   `json:"id"`
	Name        string `json:"name,omitempty"`
	Category    string `json:"category,omitempty"`
	Description string `json:"description,omitempty"`
	Content     string `json:"content,omitempty"`
	Tags        string `json:"tags,omitempty"`
	Source      string `json:"source,omitempty"`
	IsActive    *bool  `json:"is_active,omitempty"`
}

type CursorRulesUpdateOutput struct {
	Success bool `json:"success"`
}

type CursorRulesDeleteInput struct {
	ID uint `json:"id"`
}

type CursorRulesDeleteOutput struct {
	Success bool `json:"success"`
}

type CursorRulesSuggestInput struct {
	Language string `json:"language,omitempty"`
	Category string `json:"category,omitempty"`
	Tags     string `json:"tags,omitempty"`
	Limit    int    `json:"limit,omitempty"`
}

type CursorRulesSuggestOutput struct {
	Suggestions []CommunityRule `json:"suggestions"`
}

type CommunityRule struct {
	Name        string `json:"name"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Tags        string `json:"tags"`
	Source      string `json:"source"`
	URL         string `json:"url"`
}

type CursorRulesInstallInput struct {
	RuleName string `json:"rule_name"`
	URL      string `json:"url,omitempty"`
}

type CursorRulesInstallOutput struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
