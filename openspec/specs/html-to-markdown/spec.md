## ADDED Requirements

### Requirement: Convert HTML to Markdown

The system SHALL convert HTML documents to markdown format preserving document structure.

#### Scenario: Basic conversion
- **WHEN** valid HTML is provided
- **THEN** system outputs equivalent markdown to stdout

#### Scenario: Empty HTML
- **WHEN** HTML contains no content
- **THEN** system outputs empty string

### Requirement: Preserve headings

The system SHALL convert HTML headings to markdown heading syntax.

#### Scenario: H1 heading
- **WHEN** HTML contains `<h1>Title</h1>`
- **THEN** output contains `# Title`

#### Scenario: H2 heading
- **WHEN** HTML contains `<h2>Subtitle</h2>`
- **THEN** output contains `## Subtitle`

#### Scenario: H3-H6 headings
- **WHEN** HTML contains h3-h6 elements
- **THEN** output contains corresponding markdown heading levels (### through ######)

### Requirement: Preserve links

The system SHALL convert HTML anchors to markdown link syntax.

#### Scenario: Anchor link
- **WHEN** HTML contains `<a href="https://example.com">Example</a>`
- **THEN** output contains `[Example](https://example.com)`

#### Scenario: Relative links
- **WHEN** HTML contains relative links
- **THEN** links are preserved as-is in markdown

### Requirement: Preserve lists

The system SHALL convert HTML lists to markdown list syntax.

#### Scenario: Unordered list
- **WHEN** HTML contains `<ul><li>Item 1</li><li>Item 2</li></ul>`
- **THEN** output contains markdown bullet list

#### Scenario: Ordered list
- **WHEN** HTML contains `<ol><li>First</li><li>Second</li></ol>`
- **THEN** output contains markdown numbered list

### Requirement: Preserve code blocks

The system SHALL convert HTML code elements to markdown code syntax.

#### Scenario: Inline code
- **WHEN** HTML contains `<code>inline code</code>`
- **THEN** output contains `` `inline code` ``

#### Scenario: Code block
- **WHEN** HTML contains `<pre><code>code block</code></pre>`
- **THEN** output contains fenced code block

### Requirement: Preserve paragraphs and line breaks

The system SHALL convert HTML paragraphs and line breaks appropriately.

#### Scenario: Paragraph
- **WHEN** HTML contains `<p>Text</p>`
- **THEN** output contains text with blank line separation

#### Scenario: Line break
- **WHEN** HTML contains `<br>`
- **THEN** output contains markdown line break (two spaces + newline)

### Requirement: Preserve tables

The system SHALL convert HTML tables to markdown table syntax.

#### Scenario: Simple table
- **WHEN** HTML contains a table with headers and rows
- **THEN** output contains markdown table with pipe delimiters

### Requirement: Handle malformed HTML gracefully

The system SHALL handle malformed HTML without crashing.

#### Scenario: Unclosed tags
- **WHEN** HTML has unclosed tags
- **THEN** system still produces readable markdown output

#### Scenario: Invalid HTML entities
- **WHEN** HTML contains invalid entities
- **THEN** system passes them through or converts to text

### Requirement: Strip non-content elements

The system SHALL exclude script, style, and navigation elements from output.

#### Scenario: Script tags
- **WHEN** HTML contains `<script>...</script>`
- **THEN** script content is not included in markdown output

#### Scenario: Style tags
- **WHEN** HTML contains `<style>...</style>`
- **THEN** style content is not included in markdown output
