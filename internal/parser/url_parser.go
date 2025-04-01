package parser

import (
	"errors"
	"net/url"
	"regexp"
	"strings"
)

type URLParser struct {
	urlRegex *regexp.Regexp
}

type URLParseResult struct {
	OriginalURL string
	Normalized  string
	Domain      string
	Path        string
	Params      map[string]string
	IsValid     bool
}

func NewURLParser() *URLParser {
	regex := regexp.MustCompile(`^(http|https)://[a-zA-Z0-9]+([\-\.]{1}[a-zA-Z0-9]+)*\.[a-zA-Z]{2,}(:[0-9]{1,5})?(/.*)?$`)

	return &URLParser{
		urlRegex: regex,
	}
}

func (parser *URLParser) Parse(rawURL string) (*URLParseResult, error) {
	rawURL = strings.TrimSpace(rawURL)

	if rawURL == "" {
		return nil, errors.New("empty URL provided")
	}

	prefixIsMissing := !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://")
	if prefixIsMissing {
		rawURL = "https://" + rawURL
	}

	isValid := parser.urlRegex.MatchString(rawURL)
	if !isValid {
		return nil, errors.New("invalid URL format")
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	params := make(map[string]string)
	for k, v := range parsedURL.Query() {
		if len(v) > 0 {
			params[k] = v[0]
		}
	}

	normalized := strings.TrimSuffix(strings.ToLower(rawURL), "/")

	return &URLParseResult{
		OriginalURL: rawURL,
		Normalized:  normalized,
		Domain:      parsedURL.Host,
		Path:        parsedURL.Path,
		Params:      params,
		IsValid:     true,
	}, nil
}
func (parser *URLParser) ParseLogEntry(logEntry string) (map[string]string, error) {
	// Example log format: [timestamp] "GET /abc123 HTTP/1.1" 301 "Mozilla/5.0 ..." "192.168.1.1" "referrer"
	data := make(map[string]string)

	shortCodeRegex := regexp.MustCompile(`GET /([a-zA-Z0-9]+)`)
	matches := shortCodeRegex.FindStringSubmatch(logEntry)
	if len(matches) > 1 {
		data["shortcode"] = matches[1]
	}

	ipRegex := regexp.MustCompile(`"(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})"`)
	ipMatches := ipRegex.FindStringSubmatch(logEntry)
	if len(ipMatches) > 1 {
		data["ip"] = ipMatches[1]
	}

	uaRegex := regexp.MustCompile(`"(Mozilla[^"]*)"`)
	uaMatches := uaRegex.FindStringSubmatch(logEntry)
	if len(uaMatches) > 1 {
		data["user_agent"] = uaMatches[1]
	}

	tsRegex := regexp.MustCompile(`\[(.*?)\]`)
	tsMatches := tsRegex.FindStringSubmatch(logEntry)
	if len(tsMatches) > 1 {
		data["timestamp"] = tsMatches[1]
	}

	return data, nil
}
