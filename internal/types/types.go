package types

// Goal management inputs and outputs
type GoalsListInput struct {
	Limit int `json:"limit,omitempty" jsonschema:"Maximum number of goals to return (0 = no limit)"`
}

type GoalsListOutput struct {
	Goals []Goal `json:"goals" jsonschema:"List of active goals"`
}

type Goal struct {
	ID        int    `json:"id" jsonschema:"Unique goal identifier"`
	Title     string `json:"title" jsonschema:"Goal title or description"`
	Priority  int    `json:"priority" jsonschema:"Goal priority (lower number = higher priority)"`
	Status    string `json:"status" jsonschema:"Current goal status (active, paused, done)"`
	Notes     string `json:"notes" jsonschema:"Additional notes or details about the goal"`
	UpdatedAt string `json:"updated_at" jsonschema:"Last update timestamp"`
}

type GoalsAddInput struct {
	Title    string  `json:"title" jsonschema:"Goal title or description (required)"`
	Priority *int    `json:"priority,omitempty" jsonschema:"Goal priority (lower number = higher priority, defaults to 0)"`
	Notes    *string `json:"notes,omitempty" jsonschema:"Additional notes or context for the goal"`
}

type GoalsAddOutput struct {
	ID int `json:"id" jsonschema:"ID of the created goal"`
}

type GoalsUpdateInput struct {
	ID       int     `json:"id" jsonschema:"Goal ID to update (required)"`
	Status   *string `json:"status,omitempty" jsonschema:"New status (active, paused, done)"`
	Notes    *string `json:"notes,omitempty" jsonschema:"Updated notes or context"`
	Priority *int    `json:"priority,omitempty" jsonschema:"Updated priority (lower number = higher priority)"`
}

type GoalsUpdateOutput struct {
	Updated int `json:"updated" jsonschema:"Number of rows updated"`
}

// ADR management inputs and outputs
type ADRsListInput struct {
	Query *string `json:"query,omitempty" jsonschema:"Search query to filter ADRs by title or content"`
}

type ADRsListOutput struct {
	ADRs []ADR `json:"adrs" jsonschema:"List of ADRs"`
}

type ADR struct {
	ID        string `json:"id" jsonschema:"ADR identifier (e.g., ADR-001)"`
	Title     string `json:"title" jsonschema:"ADR title or subject"`
	Content   string `json:"content" jsonschema:"Full content of the ADR document"`
	UpdatedAt string `json:"updated_at" jsonschema:"Last modification timestamp"`
}

type ADRsGetInput struct {
	ID string `json:"id" jsonschema:"ADR ID to retrieve (e.g., ADR-001)"`
}

type ADRsGetOutput struct {
	ID      string `json:"id" jsonschema:"ADR identifier"`
	Title   string `json:"title" jsonschema:"ADR title"`
	Content string `json:"content" jsonschema:"Full content of the ADR document"`
}

// CI management inputs and outputs
type CIRunTestsInput struct {
	Scope *string `json:"scope,omitempty" jsonschema:"Test scope to run specific package or directory (e.g., ./cmd/jukebox)"`
}

type CIRunTestsOutput struct {
	Status string `json:"status" jsonschema:"Test execution status (success, failure, running)"`
	Output string `json:"output" jsonschema:"Test output and results"`
}

type CILastFailureInput struct {
	RandomString *string `json:"random_string,omitempty" jsonschema:"Dummy parameter for no-parameter tools (optional)"`
}

type CILastFailureOutput struct {
	Status    string  `json:"status" jsonschema:"Last test failure status"`
	Scope     *string `json:"scope,omitempty" jsonschema:"Test scope that failed (if available)"`
	StartedAt *string `json:"started_at,omitempty" jsonschema:"Test start timestamp (if available)"`
}

// Repository search inputs and outputs
type RepoSearchInput struct {
	Q    string  `json:"q" jsonschema:"Search query pattern (supports regex)"`
	Path *string `json:"path,omitempty" jsonschema:"Path to search within (file or directory)"`
	Max  *int    `json:"max,omitempty" jsonschema:"Maximum number of results to return"`
}

type RepoSearchOutput struct {
	Results []SearchResult `json:"results" jsonschema:"List of search results"`
}

type SearchResult struct {
	File  string `json:"file" jsonschema:"File path where match was found"`
	Line  string `json:"line" jsonschema:"Line number where match was found"`
	Match string `json:"match" jsonschema:"Matching text content"`
}

// Change logging inputs and outputs
type StateLogChangeInput struct {
	Summary string   `json:"summary" jsonschema:"Brief summary of the change made (required)"`
	Files   []string `json:"files,omitempty" jsonschema:"List of files that were modified"`
}

type StateLogChangeOutput struct {
	OK bool `json:"ok" jsonschema:"Whether the change was logged successfully"`
}

// Markdown linting inputs and outputs
type MarkdownLintInput struct {
	Path   *string `json:"path,omitempty" jsonschema:"Path to lint (file or directory, defaults to current directory)"`
	Fix    *bool   `json:"fix,omitempty" jsonschema:"Whether to automatically fix formatting issues"`
	Config *string `json:"config,omitempty" jsonschema:"Path to markdownlint configuration file"`
}

type MarkdownLintOutput struct {
	Issues []LintIssue `json:"issues" jsonschema:"List of linting issues found"`
	Fixed  bool        `json:"fixed" jsonschema:"Whether any issues were automatically fixed"`
	Path   string      `json:"path" jsonschema:"Path that was linted"`
}

type LintIssue struct {
	File    string `json:"file" jsonschema:"File path where the issue was found"`
	Line    int    `json:"line" jsonschema:"Line number of the issue"`
	Column  int    `json:"column" jsonschema:"Column number of the issue"`
	Rule    string `json:"rule" jsonschema:"Markdownlint rule that was violated"`
	Message string `json:"message" jsonschema:"Description of the linting issue"`
}

// Template system inputs and outputs
type TemplateListInput struct {
	Category *string `json:"category,omitempty" jsonschema:"Filter templates by category (optional)"`
}

type TemplateListOutput struct {
	Templates []Template `json:"templates" jsonschema:"List of available templates"`
}

type Template struct {
	ID          string             `json:"id" jsonschema:"Unique template identifier"`
	Name        string             `json:"name" jsonschema:"Template name"`
	Description string             `json:"description" jsonschema:"Template description"`
	Category    string             `json:"category" jsonschema:"Template category"`
	Content     string             `json:"content" jsonschema:"Template markdown content"`
	Variables   []TemplateVariable `json:"variables" jsonschema:"Template variables definition"`
	CreatedAt   string             `json:"created_at" jsonschema:"Template creation timestamp"`
	UpdatedAt   string             `json:"updated_at" jsonschema:"Template last update timestamp"`
}

type TemplateVariable struct {
	Name         string `json:"name" jsonschema:"Variable name"`
	Type         string `json:"type" jsonschema:"Variable type (string, date, list, number, boolean)"`
	Required     bool   `json:"required" jsonschema:"Whether this variable is required"`
	DefaultValue string `json:"default_value" jsonschema:"Default value for the variable"`
	Description  string `json:"description" jsonschema:"Description of what this variable represents"`
}

type TemplateRegisterInput struct {
	ID          string             `json:"id" jsonschema:"Unique template identifier (required)"`
	Name        string             `json:"name" jsonschema:"Template name (required)"`
	Description *string            `json:"description,omitempty" jsonschema:"Template description (optional)"`
	Category    string             `json:"category" jsonschema:"Template category (required)"`
	Content     string             `json:"content" jsonschema:"Template markdown content (required)"`
	Variables   []TemplateVariable `json:"variables,omitempty" jsonschema:"Template variables definition (optional)"`
}

type TemplateRegisterOutput struct {
	ID      string `json:"id" jsonschema:"Registered template identifier"`
	Success bool   `json:"success" jsonschema:"Whether registration was successful"`
}

type TemplateGetInput struct {
	ID string `json:"id" jsonschema:"Template identifier to retrieve"`
}

type TemplateGetOutput struct {
	Template Template `json:"template" jsonschema:"Retrieved template"`
}

type TemplateUpdateInput struct {
	ID          string             `json:"id" jsonschema:"Template identifier to update (required)"`
	Name        *string            `json:"name,omitempty" jsonschema:"Updated template name (optional)"`
	Description *string            `json:"description,omitempty" jsonschema:"Updated template description (optional)"`
	Category    *string            `json:"category,omitempty" jsonschema:"Updated template category (optional)"`
	Content     *string            `json:"content,omitempty" jsonschema:"Updated template content (optional)"`
	Variables   []TemplateVariable `json:"variables,omitempty" jsonschema:"Updated template variables (optional)"`
}

type TemplateUpdateOutput struct {
	Updated bool `json:"updated" jsonschema:"Whether template was successfully updated"`
}

type TemplateDeleteInput struct {
	ID string `json:"id" jsonschema:"Template identifier to delete"`
}

type TemplateDeleteOutput struct {
	Deleted bool `json:"deleted" jsonschema:"Whether template was successfully deleted"`
}

type TemplateApplyInput struct {
	TemplateID string                 `json:"template_id" jsonschema:"Template identifier to apply (required)"`
	Variables  map[string]interface{} `json:"variables" jsonschema:"Variable values to substitute in template"`
	OutputPath *string                `json:"output_path,omitempty" jsonschema:"Output file path (optional, auto-generated if not provided)"`
}

type TemplateApplyOutput struct {
	Content string `json:"content" jsonschema:"Generated markdown content"`
	Path    string `json:"path" jsonschema:"Output file path where content was written"`
}

// Preferred Tool types
type PreferredTool struct {
	ID          uint   `json:"id" jsonschema:"Unique tool identifier"`
	Name        string `json:"name" jsonschema:"Tool name"`
	Category    string `json:"category" jsonschema:"Tool category"`
	Description string `json:"description" jsonschema:"Tool description"`
	Language    string `json:"language" jsonschema:"Programming language"`
	UseCase     string `json:"use_case" jsonschema:"Use case description"`
	Priority    int    `json:"priority" jsonschema:"Tool priority (higher = more preferred)"`
	CreatedAt   string `json:"created_at" jsonschema:"Creation timestamp"`
	UpdatedAt   string `json:"updated_at" jsonschema:"Last update timestamp"`
}

type PreferredToolsListInput struct {
	Category string `json:"category,omitempty" jsonschema:"Filter by tool category (optional)"`
	Language string `json:"language,omitempty" jsonschema:"Filter by programming language (optional)"`
	Limit    int    `json:"limit,omitempty" jsonschema:"Maximum number of results to return (optional)"`
}

type PreferredToolsListOutput struct {
	Tools []PreferredTool `json:"tools" jsonschema:"List of preferred tools"`
}

type PreferredToolsAddInput struct {
	Name        string `json:"name" jsonschema:"Tool name (required)"`
	Category    string `json:"category" jsonschema:"Tool category (required)"`
	Description string `json:"description,omitempty" jsonschema:"Tool description (optional)"`
	Language    string `json:"language,omitempty" jsonschema:"Programming language (optional)"`
	UseCase     string `json:"use_case,omitempty" jsonschema:"Use case description (optional)"`
	Priority    int    `json:"priority,omitempty" jsonschema:"Tool priority (optional, higher = more preferred)"`
}

type PreferredToolsAddOutput struct {
	ID      uint `json:"id" jsonschema:"Created tool identifier"`
	Success bool `json:"success" jsonschema:"Whether tool was successfully added"`
}

type PreferredToolsUpdateInput struct {
	ID          uint   `json:"id" jsonschema:"Tool identifier to update (required)"`
	Name        string `json:"name,omitempty" jsonschema:"Updated tool name (optional)"`
	Category    string `json:"category,omitempty" jsonschema:"Updated tool category (optional)"`
	Description string `json:"description,omitempty" jsonschema:"Updated tool description (optional)"`
	Language    string `json:"language,omitempty" jsonschema:"Updated programming language (optional)"`
	UseCase     string `json:"use_case,omitempty" jsonschema:"Updated use case description (optional)"`
	Priority    int    `json:"priority,omitempty" jsonschema:"Updated tool priority (optional)"`
}

type PreferredToolsUpdateOutput struct {
	Success bool `json:"success" jsonschema:"Whether tool was successfully updated"`
}

type PreferredToolsDeleteInput struct {
	ID uint `json:"id" jsonschema:"Tool identifier to delete"`
}

type PreferredToolsDeleteOutput struct {
	Success bool `json:"success" jsonschema:"Whether tool was successfully deleted"`
}

// Cursor Rules types
type CursorRule struct {
	ID          uint   `json:"id" jsonschema:"Unique rule identifier"`
	Name        string `json:"name" jsonschema:"Rule name"`
	Category    string `json:"category" jsonschema:"Rule category"`
	Description string `json:"description" jsonschema:"Rule description"`
	Content     string `json:"content" jsonschema:"Rule content (MDC format)"`
	Tags        string `json:"tags" jsonschema:"Comma-separated tags"`
	Source      string `json:"source" jsonschema:"Rule source (local, community, custom)"`
	IsActive    bool   `json:"is_active" jsonschema:"Whether rule is active"`
	CreatedAt   string `json:"created_at" jsonschema:"Creation timestamp"`
	UpdatedAt   string `json:"updated_at" jsonschema:"Last update timestamp"`
}

type CursorRulesListInput struct {
	Category string `json:"category,omitempty" jsonschema:"Filter by rule category (optional)"`
	Tags     string `json:"tags,omitempty" jsonschema:"Filter by tags (optional)"`
	Source   string `json:"source,omitempty" jsonschema:"Filter by rule source (optional)"`
	Active   *bool  `json:"active,omitempty" jsonschema:"Filter by active status (optional)"`
	Limit    int    `json:"limit,omitempty" jsonschema:"Maximum number of results to return (optional)"`
}

type CursorRulesListOutput struct {
	Rules []CursorRule `json:"rules" jsonschema:"List of cursor rules"`
}

type CursorRulesAddInput struct {
	Name        string `json:"name" jsonschema:"Rule name (required)"`
	Category    string `json:"category" jsonschema:"Rule category (required)"`
	Description string `json:"description,omitempty" jsonschema:"Rule description (optional)"`
	Content     string `json:"content" jsonschema:"Rule content in MDC format (required)"`
	Tags        string `json:"tags,omitempty" jsonschema:"Comma-separated tags (optional)"`
	Source      string `json:"source,omitempty" jsonschema:"Rule source (optional, defaults to local)"`
	IsActive    *bool  `json:"is_active,omitempty" jsonschema:"Whether rule is active (optional, defaults to true)"`
}

type CursorRulesAddOutput struct {
	ID      uint `json:"id" jsonschema:"Created rule identifier"`
	Success bool `json:"success" jsonschema:"Whether rule was successfully added"`
}

type CursorRulesUpdateInput struct {
	ID          uint   `json:"id" jsonschema:"Rule identifier to update (required)"`
	Name        string `json:"name,omitempty" jsonschema:"Updated rule name (optional)"`
	Category    string `json:"category,omitempty" jsonschema:"Updated rule category (optional)"`
	Description string `json:"description,omitempty" jsonschema:"Updated rule description (optional)"`
	Content     string `json:"content,omitempty" jsonschema:"Updated rule content (optional)"`
	Tags        string `json:"tags,omitempty" jsonschema:"Updated comma-separated tags (optional)"`
	Source      string `json:"source,omitempty" jsonschema:"Updated rule source (optional)"`
	IsActive    *bool  `json:"is_active,omitempty" jsonschema:"Updated active status (optional)"`
}

type CursorRulesUpdateOutput struct {
	Success bool `json:"success" jsonschema:"Whether rule was successfully updated"`
}

type CursorRulesDeleteInput struct {
	ID uint `json:"id" jsonschema:"Rule identifier to delete"`
}

type CursorRulesDeleteOutput struct {
	Success bool `json:"success" jsonschema:"Whether rule was successfully deleted"`
}

type CursorRulesSuggestInput struct {
	Language string `json:"language,omitempty" jsonschema:"Programming language for rule suggestions (optional)"`
	Category string `json:"category,omitempty" jsonschema:"Rule category for suggestions (optional)"`
	Tags     string `json:"tags,omitempty" jsonschema:"Tags for rule suggestions (optional)"`
	Limit    int    `json:"limit,omitempty" jsonschema:"Maximum number of suggestions to return (optional)"`
}

type CursorRulesSuggestOutput struct {
	Suggestions []CommunityRule `json:"suggestions" jsonschema:"List of suggested community rules"`
}

type CommunityRule struct {
	Name        string `json:"name" jsonschema:"Community rule name"`
	Category    string `json:"category" jsonschema:"Rule category"`
	Description string `json:"description" jsonschema:"Rule description"`
	Tags        string `json:"tags" jsonschema:"Comma-separated tags"`
	Source      string `json:"source" jsonschema:"Rule source (community)"`
	URL         string `json:"url" jsonschema:"URL to the community rule"`
}

type CursorRulesInstallInput struct {
	RuleName string `json:"rule_name" jsonschema:"Name of the rule to install (required)"`
	URL      string `json:"url,omitempty" jsonschema:"URL to the rule (optional)"`
}

type CursorRulesInstallOutput struct {
	Success bool   `json:"success" jsonschema:"Whether rule was successfully installed"`
	Message string `json:"message" jsonschema:"Installation result message"`
}

// Setup MCP Tools types
type SetupMCPToolsInput struct {
	ProjectPath string `json:"project_path" jsonschema:"Path to the project directory (required)"`
}

type SetupMCPToolsOutput struct {
	Success      bool     `json:"success" jsonschema:"Whether setup was successful"`
	ProjectPath  string   `json:"project_path" jsonschema:"Resolved project path"`
	RulesDir     string   `json:"rules_dir" jsonschema:"Path to .cursor/rules directory"`
	FilesCreated []string `json:"files_created" jsonschema:"List of files created"`
	Message      string   `json:"message" jsonschema:"Setup result message"`
}

// Log parsing inputs and outputs
type LogParseInput struct {
	FilePath string `json:"file_path,omitempty" jsonschema:"Path to log file (optional, auto-detects if not provided)"`
	Format   string `json:"format,omitempty" jsonschema:"Output format: json, summary, detailed, ai-friendly (default: ai-friendly)"`
}

type LogParseOutput struct {
	File            string          `json:"file" jsonschema:"Log file name"`
	AnalysisTime    string          `json:"analysis_time" jsonschema:"When the analysis was performed"`
	Statistics      FileStatistics  `json:"statistics" jsonschema:"Basic file statistics"`
	ErrorCounts     ErrorCounts     `json:"error_counts" jsonschema:"Counts of different error types"`
	CriticalIssues  CriticalIssues  `json:"critical_issues" jsonschema:"Critical problems requiring immediate attention"`
	ErrorPatterns   []ErrorPattern  `json:"error_patterns" jsonschema:"Detailed error pattern analysis"`
	RecentErrors    []string        `json:"recent_errors" jsonschema:"Most recent error messages"`
	MissingFiles    []string        `json:"missing_files" jsonschema:"Files that were referenced but not found"`
	Recommendations []string        `json:"recommendations" jsonschema:"Priority-based recommendations for fixing issues"`
	Context         AnalysisContext `json:"context" jsonschema:"Context information for AI analysis"`
}

type FileStatistics struct {
	Lines    int    `json:"lines" jsonschema:"Total number of lines in the log file"`
	Size     string `json:"size" jsonschema:"Human-readable file size"`
	Modified string `json:"modified" jsonschema:"File modification date"`
}

type ErrorCounts struct {
	NetworkIssues    int `json:"network_issues" jsonschema:"Internet connectivity failures"`
	FileErrors       int `json:"file_errors" jsonschema:"Missing files and file system errors"`
	MemoryLeaks      int `json:"memory_leaks" jsonschema:"Potential listener leaks"`
	ComposerErrors   int `json:"composer_errors" jsonschema:"Composer context errors"`
	MCPErrors        int `json:"mcp_errors" jsonschema:"MCP server-related issues"`
	JavaScriptErrors int `json:"javascript_errors" jsonschema:"JavaScript type, reference, and syntax errors"`
	PermissionErrors int `json:"permission_errors" jsonschema:"Access denied and permission issues"`
	ConnectionErrors int `json:"connection_errors" jsonschema:"Connection refused and timeout errors"`
}

type CriticalIssues struct {
	DiskSpace        int `json:"disk_space" jsonschema:"Disk space errors (ENOSPC)"`
	SyntaxErrors     int `json:"syntax_errors" jsonschema:"JavaScript syntax errors"`
	PermissionDenied int `json:"permission_denied" jsonschema:"Permission denied errors"`
}

type ErrorPattern struct {
	Pattern     string   `json:"pattern" jsonschema:"Regex pattern that matched"`
	Description string   `json:"description" jsonschema:"Human-readable description of the error"`
	Severity    string   `json:"severity" jsonschema:"Severity level: low, medium, high, critical"`
	Count       int      `json:"count" jsonschema:"Number of occurrences"`
	Recent      []string `json:"recent" jsonschema:"Recent examples of this error"`
}

type AnalysisContext struct {
	TotalErrors     int    `json:"total_errors" jsonschema:"Total number of errors found"`
	MostCommonIssue string `json:"most_common_issue" jsonschema:"Description of the most frequent issue"`
	SeverityLevel   string `json:"severity_level" jsonschema:"Overall severity assessment"`
	Environment     string `json:"environment" jsonschema:"Development environment assessment"`
}

// Changelog generation inputs and outputs
type ChangelogGenerateInput struct {
	Format string `json:"format,omitempty" jsonschema:"Output format: markdown, json (default: markdown)"`
	Limit  int    `json:"limit,omitempty" jsonschema:"Maximum number of entries to include (0 = no limit)"`
}

type ChangelogGenerateOutput struct {
	Content string `json:"content" jsonschema:"Generated changelog content"`
	Path    string `json:"path" jsonschema:"Path where changelog was written"`
	Entries int    `json:"entries" jsonschema:"Number of entries included"`
}
