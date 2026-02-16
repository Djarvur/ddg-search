## Context

The ddg-search project provides two Go binaries:
- `ddg-search`: Performs web searches via DuckDuckGo HTML scraping
- `page-dump`: Fetches URLs and converts HTML to markdown

Claude Code supports custom skills via a `skill.yaml` configuration that defines tools Claude can invoke. This design wraps the existing binaries as Claude Code tools, enabling web search and content fetching without API keys.

## Goals / Non-Goals

**Goals:**
- Create a Claude Code skill that exposes `web-search` and `page-dump` tools
- Display helpful install instructions when binaries are missing
- Keep skill configuration minimal and maintainable

**Non-Goals:**
- Modifying the existing binaries
- Adding new search or fetch capabilities
- Supporting other AI assistants (only Claude Code)

## Decisions

### Skill Location: `skills/ddg-search/`

**Rationale:** Placing the skill in a `skills/` directory at project root follows the pattern of grouping related functionality. The `ddg-search` subdirectory name matches the project and clearly indicates the skill's purpose.

**Alternative considered:** Placing in `.claude/skills/` - rejected because this is a project-specific skill, not a user-level configuration.

### Tool Definitions in skill.yaml

**Rationale:** Define two tools directly in `skill.yaml`:
- `web-search`: Takes a query string, returns search results
- `page-dump`: Takes a URL, returns markdown content

**Alternative considered:** Single combined tool - rejected because separation allows Claude to choose the appropriate action and provides clearer error messages.

### Install Instructions in SKILL.md

**Rationale:** The `SKILL.md` file contains install instructions that Claude displays when invoking the skill if binaries aren't found. This provides just-in-time guidance.

**Alternative considered:** Pre-install validation script - rejected because Claude Code's skill system handles missing tools gracefully by showing the SKILL.md content.

### Binary Discovery via PATH

**Rationale:** Binaries are expected to be in PATH. The skill invokes them directly by name.

**Alternative considered:** Explicit binary paths in config - rejected because PATH-based discovery is more flexible and works across different installation methods.

## Risks / Trade-offs

**Binary not installed** → Skill displays install instructions from SKILL.md, allowing user to build/install

**DuckDuckGo rate limiting** → Existing ddg-search binary already handles rate limiting; skill inherits this behavior

**Binary path issues** → User must ensure binaries are in PATH; SKILL.md documents this requirement
