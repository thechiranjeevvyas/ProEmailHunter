package extractor

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"
)

// ANSI color codes
const (
	REDCOLOR    = "\033[91m"
	GREENCOLOR  = "\033[92m"
	YELLOWCOLOR = "\033[93m"
	CYANCOLOR   = "\033[96m"
	BLUECOLOR   = "\033[94m"
	RESETCOLOR  = "\033[0m"
)

var excludePatterns = []string{
	".jpg", ".png", ".gif", ".webp", ".ico", ".mp4", ".pdf", ".eot",
	".doc", ".docx", ".xls", ".xlsx", ".woff", ".woff2", ".css", ".json",
	".xml", ".rss", ".svg", ".yaml", ".yml", ".csv", ".dockerfile", ".cfg",
	".lock", ".js", ".md", ".toml",
}

// Improved Email regex pattern - more strict to avoid matching image filenames
var emailRegex = regexp.MustCompile(`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`)

type ExtractionResult struct {
	Emails map[string]bool
	Links  map[string]bool
	Error  error
	URL    string
}

func shouldExclude(link string) bool {
	lowerLink := strings.ToLower(link)
	for _, pattern := range excludePatterns {
		if strings.Contains(lowerLink, pattern) {
			return true
		}
	}
	return false
}

func isValidEmail(email string) bool {
	if strings.Contains(strings.ToLower(email), "example") ||
		strings.Contains(strings.ToLower(email), "email") ||
		len(email) < 5 {
		return false
	}
	if strings.Contains(email, ".png") || strings.Contains(email, ".jpg") ||
		strings.Contains(email, ".webp") || strings.Contains(email, ".gif") {
		return false
	}
	return true
}

func extractEmailsAndLinks(targetURL, baseURL string, timeout time.Duration, verbose bool) ExtractionResult {
	if verbose {
		fmt.Printf("%s[PROCESSING]%s %s\n", YELLOWCOLOR, RESETCOLOR, targetURL)
	}

	result := ExtractionResult{
		Emails: make(map[string]bool),
		Links:  make(map[string]bool),
		URL:    targetURL,
	}

	client := &http.Client{Timeout: timeout}
	resp, err := client.Get(targetURL)
	if err != nil {
		if verbose {
			fmt.Printf("%s[ERROR]%s %v\n", REDCOLOR, RESETCOLOR, err)
		}
		result.Error = err
		return result
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if verbose {
			fmt.Printf("%s[ERROR]%s HTTP %d for %s\n", REDCOLOR, RESETCOLOR, resp.StatusCode, targetURL)
		}
		return result
	}

	body, _ := io.ReadAll(resp.Body)
	html := string(body)

	hrefPatterns := []string{`href="([^"]*)"`, `href='([^']*)'`}

	for _, pattern := range hrefPatterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllStringSubmatch(html, -1)

		for _, match := range matches {
			if len(match) < 2 {
				continue
			}
			link := match[1]
			if shouldExclude(link) {
				continue
			}

			if strings.HasPrefix(strings.ToLower(link), "mailto:") {
				emailPart := link[7:]
				emails := emailRegex.FindAllString(emailPart, -1)
				for _, email := range emails {
					if email != "" && isValidEmail(email) {
						result.Emails[email] = true
					}
				}
				continue
			}

			emailsInText := emailRegex.FindAllString(link, -1)
			for _, email := range emailsInText {
				if email != "" && isValidEmail(email) {
					result.Emails[email] = true
				}
			}

			absoluteLink, err := url.Parse(link)
			if err != nil {
				continue
			}
			base, _ := url.Parse(targetURL)
			resolvedLink := base.ResolveReference(absoluteLink).String()

			baseParsed, _ := url.Parse(baseURL)
			resolvedParsed, _ := url.Parse(resolvedLink)

			if resolvedParsed.Host == baseParsed.Host {
				result.Links[resolvedLink] = true
			}
		}
	}

	allEmailsInPage := emailRegex.FindAllString(html, -1)
	for _, email := range allEmailsInPage {
		if email != "" && isValidEmail(email) {
			result.Emails[email] = true
		}
	}

	return result
}

// ProcessDomain runs the scan logic
func ProcessDomain(domainURL string, maxWorkers int, timeout time.Duration, verbose bool) {
	allEmails := make(map[string]bool)
	foundEmails := make(map[string][]string)

	mainResult := extractEmailsAndLinks(domainURL, domainURL, timeout, verbose)
	for email := range mainResult.Emails {
		allEmails[email] = true
		foundEmails[email] = append(foundEmails[email], domainURL)
	}

	links := make([]string, 0, len(mainResult.Links))
	for link := range mainResult.Links {
		links = append(links, link)
	}

	if len(links) > 0 {
		var wg sync.WaitGroup
		results := make(chan ExtractionResult, len(links))
		semaphore := make(chan struct{}, maxWorkers)

		for _, link := range links {
			wg.Add(1)
			go func(l string) {
				defer wg.Done()
				semaphore <- struct{}{}
				defer func() { <-semaphore }()
				result := extractEmailsAndLinks(l, domainURL, timeout, verbose)
				results <- result
			}(link)
		}

		go func() {
			wg.Wait()
			close(results)
		}()

		for result := range results {
			for email := range result.Emails {
				allEmails[email] = true
				foundEmails[email] = append(foundEmails[email], result.URL)
			}
		}
	}

	if len(allEmails) > 0 {
		fmt.Printf("\n%s%s%s\n", BLUECOLOR, strings.Repeat("═", 80), RESETCOLOR)
		fmt.Printf("%s🎯 DOMAIN:%s %s\n", GREENCOLOR, RESETCOLOR, domainURL)
		fmt.Printf("%s📧 FOUND:%s %d unique email(s)\n", GREENCOLOR, RESETCOLOR, len(allEmails))
		fmt.Printf("%s%s%s\n", BLUECOLOR, strings.Repeat("═", 80), RESETCOLOR)

		i := 1
		for email := range allEmails {
			sources := foundEmails[email]
			sourceInfo := sources[0]
			if len(sources) > 1 {
				sourceInfo = fmt.Sprintf("%s (+%d more pages)", sources[0], len(sources)-1)
			}
			fmt.Printf("%s%d.%s %s:: %s%s%s\n",
				CYANCOLOR, i, RESETCOLOR, sourceInfo, GREENCOLOR, email, RESETCOLOR)
			i++
		}
		fmt.Printf("%s%s%s\n\n", BLUECOLOR, strings.Repeat("═", 80), RESETCOLOR)
	} else {
		fmt.Printf("\n%s❌ NO EMAILS FOUND:%s %s\n\n", REDCOLOR, RESETCOLOR, domainURL)
	}
}
