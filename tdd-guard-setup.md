TDD Guard
npm version CI Security License: MIT

Automated Test-Driven Development enforcement for Claude Code.

Overview
TDD Guard ensures Claude Code follows Test-Driven Development principles. When your agent tries to skip tests or over-implement, TDD Guard blocks the action and explains what needs to happen instead.

TDD Guard Demo
Click to watch TDD Guard in action

Features
Test-First Enforcement - Blocks implementation without failing tests
Minimal Implementation - Prevents code beyond current test requirements
Lint Integration - Enforces refactoring using your linting rules
Multi-Language Support - TypeScript, JavaScript, Python, PHP, Go, Rust, and Storybook
Customizable Rules - Adjust validation rules to match your TDD style
Flexible Validation - Choose faster or more capable models for your needs
Session Control - Toggle on and off mid-session
Requirements
Node.js 22+
Claude Code or Anthropic API key
Test framework (Jest, Vitest, Storybook, pytest, PHPUnit, Go 1.24+, or Rust with cargo/cargo-nextest)
Quick Start
1. Install TDD Guard
Using npm:

npm install -g tdd-guard
Or using Homebrew:

brew install tdd-guard
2. Add Test Reporter
TDD Guard needs to capture test results from your test runner. Choose your language below:

JavaScript/TypeScript
Python (pytest)
PHP (PHPUnit)
Go
Install the tdd-guard-go reporter:

go install github.com/nizos/tdd-guard/reporters/go/cmd/tdd-guard-go@latest
Pipe go test -json output to the reporter:

go test -json ./... 2>&1 | tdd-guard-go -project-root /Users/username/projects/my-app
For Makefile integration:

test:
	go test -json ./... 2>&1 | tdd-guard-go -project-root /Users/username/projects/my-app
Note: The reporter acts as a filter that passes test output through unchanged while capturing results for TDD Guard. See the Go reporter configuration for more details.

Rust
3. Configure Claude Code Hooks
TDD Guard uses hooks to validate operations and provide convenience features like quick toggle commands and automatic session management.

Choose either interactive or manual setup below:

Interactive Setup
Manual Configuration
If you prefer to edit settings files directly, add all three hooks to your chosen settings file. See Settings File Locations to choose the appropriate file:

{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "Write|Edit|MultiEdit|TodoWrite",
        "hooks": [
          {
            "type": "command",
            "command": "tdd-guard"
          }
        ]
      }
    ],
    "UserPromptSubmit": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "tdd-guard"
          }
        ]
      }
    ],
    "SessionStart": [
      {
        "matcher": "startup|resume|clear",
        "hooks": [
          {
            "type": "command",
            "command": "tdd-guard"
          }
        ]
      }
    ]
  }
}
