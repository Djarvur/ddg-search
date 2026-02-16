package search

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/Djarvur/ddg-search/internal/config"
	"github.com/PuerkitoBio/goquery"
)

// Parser extracts search results from DuckDuckGo HTML responses.
type Parser struct{}

// NewParser creates a new HTML parser.
func NewParser() *Parser {
	return &Parser{}
}

// Parse extracts search results from the HTML response body.
func (p *Parser) Parse(html string, maxResults int) ([]config.Result, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var results []config.Result

	// Try finding results by .result__a directly (more reliable)
	doc.Find(".result__a").Each(func(_ int, s *goquery.Selection) {
		if maxResults > 0 && len(results) >= maxResults {
			return
		}

		result := config.Result{}
		result.Title = strings.TrimSpace(s.Text())

		href, exists := s.Attr("href")
		if exists {
			result.URL = p.extractURL(href)
		}

		// Try to find snippet in nearby element
		parent := s.Parent()
		if parent.HasClass("result__title") {
			grandparent := parent.Parent()
			snippet := grandparent.Find(".result__snippet")
			result.Snippet = strings.TrimSpace(snippet.Text())
		}

		// Skip empty results
		if result.Title == "" || result.URL == "" {
			return
		}

		results = append(results, result)
	})

	return results, nil
}

// IsEmptyResults checks if the response indicates no results (potential rate limiting).
func (p *Parser) IsEmptyResults(html string) bool {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return true
	}

	// Check for "No results" message
	noResults := doc.Find(".no-results")

	return noResults.Length() > 0
}

// IsRateLimitPage checks if the response is a rate limit/captcha page.
func (p *Parser) IsRateLimitPage(html string) bool {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return false
	}

	// Check for captcha or rate limit indicators
	body := doc.Find("body").Text()
	lowerBody := strings.ToLower(body)

	indicators := []string{
		"captcha",
		"rate limit",
		"too many requests",
		"blocked",
		"automated",
		"bots use duckduckgo",
		"challenge",
		"anomaly",
	}

	for _, indicator := range indicators {
		if strings.Contains(lowerBody, indicator) {
			return true
		}
	}

	return false
}

// FindRateLimitIndicator returns the first rate limit indicator found in the HTML.
// Returns empty string if no indicator is found or if search results exist (not a rate limit page).
// Currently disabled due to false positives when indicator words appear in search result snippets.
func (p *Parser) FindRateLimitIndicator(_ string) string {
	// Disabled: searching for indicator words in body text causes false positives
	// when those words appear in legitimate search result snippets.
	return ""
}

// extractURL extracts the actual URL from DuckDuckGo's redirect URL.
func (p *Parser) extractURL(redirectURL string) string {
	// DuckDuckGo redirect URLs look like: //duckduckgo.com/l/?uddg=encoded_url
	if !strings.Contains(redirectURL, "uddg=") {
		return redirectURL
	}

	u, err := url.Parse(redirectURL)
	if err != nil {
		return redirectURL
	}

	// Handle relative URLs
	if !u.IsAbs() {
		u.Scheme = "https"
		u.Host = "duckduckgo.com"
	}

	uddg := u.Query().Get("uddg")
	if uddg == "" {
		return redirectURL
	}

	decoded, err := url.QueryUnescape(uddg)
	if err != nil {
		return uddg
	}

	return decoded
}
