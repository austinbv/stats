package github

import (
	"testing"
	"net/http"
	"time"
	"strconv"
	"net/url"
	"github.com/stretchr/testify/assert"
)

func TestParseHeader_returns_a_GitHubHeader_with_the_remaining_rate_limit(t *testing.T) {
	header := http.Header{}
	header.Add("X-RateLimit-Remaining", "1")

	assert.Equal(t, ParseHeader(header).RateLimitRemaining, 1)
}

func TestParseHeader_returns_a_GitHubHeader_with_the_reset_time(t *testing.T) {
	header := http.Header{}
	reset_time := 1385779257546
	header.Add("X-RateLimit-Reset", strconv.Itoa(reset_time))

	expected := ParseHeader(header).RateLimitReset
	actual := time.Unix(int64(reset_time), 0)

	assert.Equal(t, expected, actual)
}

func TestParseHeader_returns_a_GitHubHeader_with_the_next_page_link(t *testing.T) {
	header := http.Header{}
	header.Add("Link", `<https://api.github.com/resource?page=2>; rel="next", <https://api.github.com/resource?page=5>; rel="last"`)

	expected := ParseHeader(header).Next
	actual, error := url.Parse("https://api.github.com/resource?page=2")

	assert.Nil(t, error)
	assert.Equal(t, *expected, *actual)
}
