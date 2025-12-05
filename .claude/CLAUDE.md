# Environment-Specific Tools

This system has modern CLI tools installed. Use these non-interactive tools when executing commands:

## Text Search and File Finding

- **ripgrep (rg)**: Use instead of grep - faster, respects .gitignore
  - Example: `rg "pattern" --type js`
  - Example: `rg "TODO" -g '*.go'` (glob pattern)

- **ast-grep**: Structural search and replace for code
  - Example: `ast-grep --pattern 'console.log($$$)'`
  - Useful for finding code patterns beyond text search

- **fd**: Use instead of find - faster, simpler syntax
  - Example: `fd pattern` or `fd '\.js$'`
  - Example: `fd -e rs` (find Rust files)

## File Viewing and Processing

- **bat**: Use instead of cat - includes syntax highlighting
  - Example: `bat file.txt`
  - For multiple files: `bat dir/*` or `bat dir/*.ext`
  - Prefer over: `for file in ...; do echo "=== $file ==="; cat "$file"; done`

- **eza**: Modern ls replacement
  - Example: `eza -l` (detailed list), `eza --tree` (tree view)
  - Shows git status and file metadata
  - Note: `ls` is NOT aliased, use `eza` explicitly

- **glow**: Markdown renderer for terminal (non-interactive mode)
  - Example: `glow README.md`
  - Renders markdown with formatting

## Data Processing

- **jq**: JSON processor for parsing and filtering
  - Example: `curl api.com/data | jq '.items[]'`
  - Example: `jq '.dependencies' package.json`

- **yq**: YAML/JSON/XML processor (like jq for multiple formats)
  - Example: `yq '.services.web.ports' docker-compose.yml`
  - Can convert between YAML, JSON, and XML

- **sd**: Modern sed replacement for find-and-replace
  - Example: `sd 'old' 'new' file.txt`
  - Simpler syntax than sed

## Code Analysis and Statistics

- **tokei**: Code statistics and line counting
  - Example: `tokei` to see lines of code by language
  - Example: `tokei --sort lines` for detailed breakdown

- **grex**: Generate regular expressions from test cases
  - Example: `grex 'foo123' 'bar456' 'baz789'`
  - Outputs a regex pattern matching the examples

## Development Tools

- **gh**: GitHub CLI for repo management (non-interactive commands)
  - Example: `gh pr list`
  - Example: `gh issue create --title "Bug" --body "Description"`
  - Can be used non-interactively with flags

- **git**: Pre-configured with useful aliases
  - Configured with delta pager and GPG signing enabled

- **mise**: Tool version manager
  - Example: `mise install` to install all tools
  - Example: `mise list` to show installed tools
  - Example: `mise run <task>` to run tasks defined in .mise.toml
  - All tool versions in `.config/mise/config.toml`

- **go**: Go programming language
  - Example: `go build`, `go test`, `go run`
  - GOPATH configured at `~/go`

- **bun**: Fast JavaScript runtime and package manager
  - Example: `bun install`, `bun run dev`, `bun test`
  - Faster than npm for package management

- **air**: Live reload for Go development
  - Can run in background: `air` starts watching
  - Configured via .air.toml if present

## System Utilities

- **dust**: Disk usage visualization
  - Example: `dust` for current directory
  - Example: `dust -d 2` for depth of 2

- **usql**: Universal SQL client
  - Example: `usql postgres://user:pass@localhost/dbname`
  - Supports PostgreSQL, SQLite, and others
  - Can execute queries non-interactively

- **tldr**: Simplified command examples
  - Example: `tldr tar` for practical examples
  - Example: `tldr -u` to update cache
  - Uses tealdeer (tlrc) implementation

## Important Notes

- All listed tools can be used non-interactively in scripts and commands
- Avoid interactive tools not listed here: helix, zellij, btop, fzf, zoxide (zi mode)
- Tool versions and additional tools available in `.config/mise/config.toml`
- Never co-author or add Generated with [Claude Code] to commit messages.
- Before committing, ensure that you are aware of what the diff actually is, so that you can be well informed when writing the message.
- Before creating a PR, ensure that you are aware of all the commits and what their diffs are.