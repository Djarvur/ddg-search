## Why

Claude Code needs a skill to perform web searches and fetch page content without requiring API keys. The ddg-search project already has working binaries (`ddg-search` and `page-dump`), but they aren't exposed as a Claude Code skill. This change wraps those binaries into a skill that Claude can invoke directly, enabling web search and content fetching capabilities.

## What Changes

- Create new Claude Code skill at `skills/ddg-search/`
- Expose two skill functions:
  - **web-search**: Search the web using DuckDuckGo via the `ddg-search` binary
  - **page-dump**: Fetch and convert web pages to markdown via the `page-dump` binary
- Include install instructions that display when binaries are not found
- Skill configuration in `skill.yaml` with proper tool definitions

## Capabilities

### New Capabilities
- `ddg-search-skill`: Claude Code skill providing web-search and page-dump functionality through tool definitions

### Modified Capabilities
- None - this is a new skill wrapping existing capabilities (ddg-search, url-fetch, html-to-markdown)

## Impact

- **New files**: `skills/ddg-search/skill.yaml`, `skills/ddg-search/SKILL.md`
- **Dependencies**: Requires `ddg-search` and `page-dump` binaries installed via `go install github.com/Djarvur/ddg-search/cmd/...@latest`
- **Systems**: Integrates with Claude Code's skill system
