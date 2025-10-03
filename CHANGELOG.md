# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Process lock mechanism to prevent multiple server instances
- Debug logging system with `PROJECT_MANAGER_DEBUG` environment variable
- MCP server tools for project management

### Changed

- Fixed PID file location to use build directory instead of `.agent` directory
- Improved server startup and shutdown handling
- Enhanced error handling in integration tests

### Fixed

- Nil pointer dereference in markdown tests
- Integration test validation errors
- CI test hanging issues
- Broken pipe errors in markdown tools

## [1.0.0] - 2025-10-03

### Initial Release

- Initial release of Project Manager MCP server
- Goals management system
- ADR (Architecture Decision Records) management
- CI/CD integration tools
- Repository search functionality
- Markdown template system
- Cursor rules management
- Documentation generation
- State logging and changelog generation

### Development History

- Testing MCP tools to verify they work
- Created comprehensive MCP testing report documenting 4 working tools and 5 tools with issues
- Fixed goals_update tool validation issues
- Implemented markdown linting tool for consistent formatting validation
- MCP Testing Report Round 2 - 7/9 tools working (78% success rate)
- Fixed repo_search tool validation issues
- Fixed ci_last_failure tool parameter validation
- Implemented markdown_lint MCP tool with auto-fix functionality
- Research phase: Go markdown libraries and advanced correction methods
- Implemented comprehensive markdown template system with 6 new MCP tools
