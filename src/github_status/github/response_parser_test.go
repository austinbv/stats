package github

import (
	"testing"
	"net/http"
	"time"
	"strconv"
	"net/url"
)

func TestParseHeader_returns_a_GitHubHeader_with_the_remaining_rate_limit(t *testing.T) {
	header := http.Header{}
	header.Add("X-RateLimit-Remaining", "1")

	if h := ParseHeader(header).RateLimitRemaining; h != 1 {
		t.Errorf("Expected %v to equal %v", h, 1)
	}
}

func TestParseHeader_returns_a_GitHubHeader_with_the_reset_time(t *testing.T) {
	header := http.Header{}
	reset_time := 1385779257546
	header.Add("X-RateLimit-Reset", strconv.Itoa(reset_time))

	expected := ParseHeader(header).RateLimitReset
	actual := time.Unix(int64(reset_time), 0)

	if expected != actual {
		t.Errorf("Expected %v to equal %v", expected, actual)
	}
}

func TestParseHeader_returns_a_GitHubHeader_with_the_next_page_link(t *testing.T) {
	header := http.Header{}
	header.Add("Link", `<https://api.github.com/resource?page=2>; rel="next", <https://api.github.com/resource?page=5>; rel="last"`)

	expected := ParseHeader(header).Next
	actual, _ := url.Parse("https://api.github.com/resource?page=2")

	if *expected != *actual {
		t.Errorf("Expected %v to equal %v", expected, actual)
	}
}
