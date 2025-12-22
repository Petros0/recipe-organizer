package handler

import (
	"html"
	"regexp"
	"strings"
)

// sanitizeText decodes HTML entities and cleans up text
func sanitizeText(s string) string {
	if s == "" {
		return s
	}
	// Decode HTML entities (e.g., &#39; -> ', &amp; -> &, &nbsp; -> space)
	result := html.UnescapeString(s)
	// Replace non-breaking spaces with regular spaces
	result = strings.ReplaceAll(result, "\u00a0", " ")
	// Collapse multiple whitespace into single space
	spaceRegex := regexp.MustCompile(`\s+`)
	result = spaceRegex.ReplaceAllString(result, " ")
	// Trim leading/trailing whitespace
	result = strings.TrimSpace(result)
	return result
}

// Helper functions for safe type conversion

func getString(obj map[string]interface{}, key string) string {
	val, ok := obj[key]
	if !ok {
		return ""
	}
	if str, ok := val.(string); ok {
		return sanitizeText(str)
	}
	return ""
}

func getStringPtr(obj map[string]interface{}, key string) *string {
	val, ok := obj[key]
	if !ok {
		return nil
	}
	if str, ok := val.(string); ok && str != "" {
		sanitized := sanitizeText(str)
		if sanitized != "" {
			return &sanitized
		}
	}
	return nil
}

func getStringArray(obj map[string]interface{}, key string) []string {
	val, ok := obj[key]
	if !ok {
		return nil
	}

	var result []string
	switch v := val.(type) {
	case []interface{}:
		for _, item := range v {
			if str, ok := item.(string); ok && str != "" {
				sanitized := sanitizeText(str)
				if sanitized != "" {
					result = append(result, sanitized)
				}
			}
		}
	case string:
		if v != "" {
			sanitized := sanitizeText(v)
			if sanitized != "" {
				result = append(result, sanitized)
			}
		}
	}

	return result
}
