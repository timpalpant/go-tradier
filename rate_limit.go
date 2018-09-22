package tradier

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Extract quota violation expiration from body message.
func parseQuotaViolationExpiration(body string) time.Time {
	if !strings.HasPrefix(body, "Quota Violation") {
		return time.Time{}
	}

	parts := strings.Fields(body)
	ms, err := strconv.ParseInt(parts[len(parts)-1], 10, 64)
	if err != nil {
		return time.Time{}
	}

	return time.Unix(ms/1000, 0)
}

// Attempt to extract the rate-limit expiry time if we have exceeded the rate limit available.
func getRateLimitExpiration(h http.Header) time.Time {
	rateLimit, err := strconv.ParseInt(h.Get(rateLimitAvailable), 10, 64)
	if err != nil {
		return time.Time{}
	}

	rateLimitExpiry, err := strconv.ParseInt(h.Get(rateLimitExpiry), 10, 64)
	if err != nil {
		return time.Time{}
	}

	if rateLimit == 0 && rateLimitExpiry > 0 {
		return time.Unix(rateLimitExpiry, 0)
	}

	return time.Time{}
}
