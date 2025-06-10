# CLAUDE.md

**回答は日本語で**

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build and Development Commands

```bash
# Build
make build                    # Build binary to bin/reflo
go build -o reflo ./cmd/reflo # Alternative build

# Test
make test                     # Run all tests
go test ./...                # Alternative test
go test -v ./internal/app    # Test specific package

# Run
make start                   # Run with start command
./bin/reflo start           # Run built binary
./bin/reflo end-day         # Generate daily summary
./bin/reflo help            # Show help

# Development
make generate               # Generate mocks
make tidy                  # Clean dependencies
go fmt ./...              # Format code
```

## Architecture Overview

This is a Go CLI application following Clean Architecture with clear separation of concerns:

### Layer Structure
- **CLI Layer** (`internal/cli/`): Command parsing and routing. Supports `start`, `end-day`, `help` commands
- **App Layer** (`internal/app/`): Core business logic implementing the Runner interface
- **Infrastructure**: Logger (JSON to `~/.reflo/YYYY-MM-DD.json`), Timer (Pomodoro-style), Prompt (terminal input)

### Key Patterns
- Dependency injection using functional options
- Interface-based design with mocks for testing
- Context-aware operations with proper cancellation

### Session Flow
1. User runs `reflo start` → declares task goal
2. 25-minute focus timer starts
3. Bell rings on completion → user enters reflection
4. 5-minute break
5. Session logged as JSON with StartTime, EndTime, Goal, and Retro fields

### Testing
- Mock generation via `go.uber.org/mock`
- Test files alongside implementation (`*_test.go`)
- Use `make generate` after interface changes

### Entry Points
- Main: `/cmd/reflo/main.go`
- Runner interface: defines `Start()`, `EndDay()`, `Help()` methods