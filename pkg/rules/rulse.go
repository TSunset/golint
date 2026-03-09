package rules

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

var sensitiveKeywords = []string{
	"password",
	"passwd",
	"pwd",
	"secret",
	"token",
	"api_key",
	"apikey",
	"access_key",
	"private_key",
}

var specialCharsRegex = regexp.MustCompile(`[@#$%^&*()=+\[\]{}<>\\|/~` + "`" + `]`)
var repeatedPunctRegex = regexp.MustCompile(`[!?.,:;]{2,}`)
var suspiciousPatternRegex = regexp.MustCompile(`(?i)(password|passwd|pwd|token|api[_-]?key|secret)\s*[:=]`)

func CheckMessage(msg string) []string {
	var issues []string

	if msg == "" {
		return issues
	}

	if !startsWithLowercase(msg) {
		issues = append(issues, "log message must start with a lowercase letter")
	}

	if !isEnglishOnly(msg) {
		issues = append(issues, "log message must be in English only")
	}

	if containsSpecialCharsOrEmoji(msg) {
		issues = append(issues, "log message must not contain special characters or emoji")
	}

	if containsSensitiveData(msg) {
		issues = append(issues, "log message must not contain sensitive data")
	}

	return issues
}

func startsWithLowercase(msg string) bool {
	r, _ := utf8.DecodeRuneInString(strings.TrimSpace(msg))
	if r == utf8.RuneError {
		return false
	}
	if unicode.IsLetter(r) {
		return unicode.IsLower(r)
	}
	return false
}

func isEnglishOnly(msg string) bool {
	for _, r := range msg {
		if r > unicode.MaxASCII {
			return false
		}
	}
	return true
}

func containsSpecialCharsOrEmoji(msg string) bool {
	if specialCharsRegex.MatchString(msg) {
		return true
	}
	if repeatedPunctRegex.MatchString(msg) {
		return true
	}

	for _, r := range msg {
		if unicode.IsSymbol(r) {
			return true
		}
	}

	return false
}

func containsSensitiveData(msg string) bool {
	lower := strings.ToLower(msg)

	if suspiciousPatternRegex.MatchString(lower) {
		return true
	}

	for _, kw := range sensitiveKeywords {
		if strings.Contains(lower, kw) {
			return true
		}
	}

	return false
}
