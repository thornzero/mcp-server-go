# Generated Documentation

This directory contains auto-generated documentation from Go source code comments.

## Package Documentation

### github.com/thornzero/project-manager/internal/adrs

```
package adrs // import "github.com/thornzero/project-manager/internal/adrs"

type ADRsHandler struct{ ... }
    func NewADRsHandler(s *server.Server) *ADRsHandler
```

### github.com/thornzero/project-manager/internal/ci

```
package ci // import "github.com/thornzero/project-manager/internal/ci"

type CIHandler struct{ ... }
    func NewCIHandler(s *server.Server) *CIHandler
```

### github.com/thornzero/project-manager/internal/cursorrules

```
package cursorrules // import "github.com/thornzero/project-manager/internal/cursorrules"

/home/thornzero/Repositories/project-manager/internal/cursorrules/cursorrules.go

type CursorRulesHandler struct{ ... }
    func NewCursorRulesHandler(srv *server.Server) *CursorRulesHandler
```

### github.com/thornzero/project-manager/internal/docs

```
package docs // import "github.com/thornzero/project-manager/internal/docs"

Package docs provides MCP tools for accessing Go documentation.

This package implements documentation tools that allow agents to access
godoc-generated documentation for Go packages, functions, and types. It
provides both static documentation generation and dynamic documentation serving
capabilities.

Example usage:

    handler := NewDocsHandler(server)
    result, output, err := handler.DocsGet(ctx, req, types.DocsGetInput{Package: "internal/goals"})

type DocsHandler struct{ ... }
    func NewDocsHandler(s *server.Server) *DocsHandler
```

### github.com/thornzero/project-manager/internal/goals

```
package goals // import "github.com/thornzero/project-manager/internal/goals"

Package goals provides MCP tools for managing project goals and milestones.

This package implements the goals management functionality for the MCP server,
allowing users to create, list, update, and track project goals with priorities
and status tracking.

Example usage:

    handler := NewGoalsHandler(server)
    result, output, err := handler.GoalsList(ctx, req, types.GoalsListInput{Limit: 10})

type GoalsHandler struct{ ... }
    func NewGoalsHandler(s *server.Server) *GoalsHandler
```

### github.com/thornzero/project-manager/internal/logparser

```
package logparser // import "github.com/thornzero/project-manager/internal/logparser"

/home/thornzero/Repositories/project-manager/internal/logparser/logparser.go

type ErrorPatternConfig struct{ ... }
type LogParserHandler struct{ ... }
    func NewLogParserHandler(server interface{ ... }) *LogParserHandler
```

### github.com/thornzero/project-manager/internal/markdown

```
package markdown // import "github.com/thornzero/project-manager/internal/markdown"

/home/thornzero/Repositories/project-manager/internal/markdown/builder.go

func WrapText(s string, width int) []string
type Builder struct{ ... }
    func MCPToolsRuleBuilder(mcpServerPath string) *Builder
    func MCPTroubleshootingGuideBuilder(mcpServerPath string) *Builder
    func MCPUsageGuideBuilder() *Builder
    func NewBuilder() *Builder
type MarkdownHandler struct{ ... }
    func NewMarkdownHandler(s *server.Server) *MarkdownHandler
```

### github.com/thornzero/project-manager/internal/models

```
package models // import "github.com/thornzero/project-manager/internal/models"

type ADR struct{ ... }
type CIRun struct{ ... }
type ChangelogEntry struct{ ... }
type CursorRule struct{ ... }
type Goal struct{ ... }
type MarkdownTemplate struct{ ... }
type PreferredTool struct{ ... }
type TemplateVariable struct{ ... }
```

### github.com/thornzero/project-manager/internal/preferredtools

```
package preferredtools // import "github.com/thornzero/project-manager/internal/preferredtools"

/home/thornzero/Repositories/project-manager/internal/preferredtools/preferredtools.go

type PreferredToolsHandler struct{ ... }
    func NewPreferredToolsHandler(srv *server.Server) *PreferredToolsHandler
```

### github.com/thornzero/project-manager/internal/search

```
package search // import "github.com/thornzero/project-manager/internal/search"

type SearchHandler struct{ ... }
    func NewSearchHandler(s *server.Server) *SearchHandler
```

### github.com/thornzero/project-manager/internal/server

```
package server // import "github.com/thornzero/project-manager/internal/server"

type Server struct{ ... }
    func NewServer(repoRoot string) (*Server, error)
```

### github.com/thornzero/project-manager/internal/setup

```
package setup // import "github.com/thornzero/project-manager/internal/setup"

/home/thornzero/Repositories/project-manager/internal/setup/setup.go

type SetupHandler struct{ ... }
    func NewSetupHandler(s *server.Server) *SetupHandler
```

### github.com/thornzero/project-manager/internal/state

```
package state // import "github.com/thornzero/project-manager/internal/state"

type StateHandler struct{ ... }
    func NewStateHandler(s *server.Server) *StateHandler
```

### github.com/thornzero/project-manager/internal/templates

```
package templates // import "github.com/thornzero/project-manager/internal/templates"

type TemplatesHandler struct{ ... }
    func NewTemplatesHandler(s *server.Server) *TemplatesHandler
```

### github.com/thornzero/project-manager/internal/types

```
package types // import "github.com/thornzero/project-manager/internal/types"

type ADR struct{ ... }
type ADRsGetInput struct{ ... }
type ADRsGetOutput struct{ ... }
type ADRsListInput struct{ ... }
type ADRsListOutput struct{ ... }
type AnalysisContext struct{ ... }
type CILastFailureInput struct{ ... }
type CILastFailureOutput struct{ ... }
type CIRunTestsInput struct{ ... }
type CIRunTestsOutput struct{ ... }
type ChangelogGenerateInput struct{ ... }
type ChangelogGenerateOutput struct{ ... }
type CommunityRule struct{ ... }
type CriticalIssues struct{ ... }
type CursorRule struct{ ... }
type CursorRulesAddInput struct{ ... }
type CursorRulesAddOutput struct{ ... }
type CursorRulesDeleteInput struct{ ... }
type CursorRulesDeleteOutput struct{ ... }
type CursorRulesInstallInput struct{ ... }
type CursorRulesInstallOutput struct{ ... }
type CursorRulesListInput struct{ ... }
type CursorRulesListOutput struct{ ... }
type CursorRulesSuggestInput struct{ ... }
type CursorRulesSuggestOutput struct{ ... }
type CursorRulesUpdateInput struct{ ... }
type CursorRulesUpdateOutput struct{ ... }
type DocsGenerateInput struct{ ... }
type DocsGenerateOutput struct{ ... }
type DocsGetInput struct{ ... }
type DocsGetOutput struct{ ... }
type DocsListInput struct{ ... }
type DocsListOutput struct{ ... }
type ErrorCounts struct{ ... }
type ErrorPattern struct{ ... }
type FileStatistics struct{ ... }
type Goal struct{ ... }
type GoalsAddInput struct{ ... }
type GoalsAddOutput struct{ ... }
type GoalsListInput struct{ ... }
type GoalsListOutput struct{ ... }
type GoalsUpdateInput struct{ ... }
type GoalsUpdateOutput struct{ ... }
type LintIssue struct{ ... }
type LogParseInput struct{ ... }
type LogParseOutput struct{ ... }
type MarkdownLintInput struct{ ... }
type MarkdownLintOutput struct{ ... }
type PreferredTool struct{ ... }
type PreferredToolsAddInput struct{ ... }
type PreferredToolsAddOutput struct{ ... }
type PreferredToolsDeleteInput struct{ ... }
type PreferredToolsDeleteOutput struct{ ... }
type PreferredToolsListInput struct{ ... }
type PreferredToolsListOutput struct{ ... }
type PreferredToolsUpdateInput struct{ ... }
type PreferredToolsUpdateOutput struct{ ... }
type RepoSearchInput struct{ ... }
type RepoSearchOutput struct{ ... }
type SearchResult struct{ ... }
type SetupMCPToolsInput struct{ ... }
type SetupMCPToolsOutput struct{ ... }
type StateLogChangeInput struct{ ... }
type StateLogChangeOutput struct{ ... }
type Template struct{ ... }
type TemplateApplyInput struct{ ... }
type TemplateApplyOutput struct{ ... }
type TemplateDeleteInput struct{ ... }
type TemplateDeleteOutput struct{ ... }
type TemplateGetInput struct{ ... }
type TemplateGetOutput struct{ ... }
type TemplateListInput struct{ ... }
type TemplateListOutput struct{ ... }
type TemplateRegisterInput struct{ ... }
type TemplateRegisterOutput struct{ ... }
type TemplateUpdateInput struct{ ... }
type TemplateUpdateOutput struct{ ... }
type TemplateVariable struct{ ... }
```

