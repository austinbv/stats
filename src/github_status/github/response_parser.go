package github

import (
	"net/http"
	"net/url"
	"time"
	"strconv"
	"regexp"
)

type GitHubHeader struct {
	Next                      *url.URL
	RateLimitRemaining        int
	RateLimitReset            time.Time
}

func ParseHeader(header http.Header) GitHubHeader {
	return GitHubHeader{
		RateLimitRemaining: getRateLimitRemaining(header),
		RateLimitReset: getRateLimitResetTime(header),
		Next: getNextPageLink(header),
	}
}

func getNextPageLink(header http.Header) *url.URL {
	re := regexp.MustCompile(`<(.*?)>; rel="next"`)
	next_link_match := re.FindStringSubmatch(header.Get("Link"))

	next_page := ""
	if len(next_link_match) != 0 {
		next_page = next_link_match[1]
	}

	next_page_url, _ := url.Parse(next_page)
	return next_page_url
}

func getRateLimitRemaining(header http.Header) int {
	rate_limit, _ := strconv.Atoi(header.Get("X-RateLimit-Remaining"))
	return rate_limit
}

func getRateLimitResetTime(header http.Header) time.Time{
	reset, _ := strconv.Atoi(header.Get("X-RateLimit-Reset"))
	return time.Unix(int64(reset), 0)
}
