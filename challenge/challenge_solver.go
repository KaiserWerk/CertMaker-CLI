package challenge

import (
	"context"
	"net/http"

	"github.com/KaiserWerk/CertMaker-CLI/entity"
)

type Solver interface {
	CanSolve(challengeType string) bool
	// Setup prepares the challenge for solving, e.g., by starting a server returning the token
	// or creating DNS records.
	Setup(ctx context.Context, token string, domains []string) error
	// Solve notifies the CertMaker instance that the challenge is ready to be validated. It can also be
	// used for cleanup after the challenge is solved.
	Solve(ctx context.Context, instanceURL, challengeID string, setHeader func(*http.Request)) (*entity.CertificateResponse, error)
}

func normalizeDomains(domains []string) []string {
	// there might be wildcard domains, we need to normalize them by removing the "*." prefix
	// also no double entries allowed
	normalized := make([]string, len(domains))
	seen := make(map[string]struct{})
	idx := 0
	for _, domain := range domains {
		if len(domain) >= 2 && domain[0:2] == "*." {
			domain = domain[2:]
		}
		if _, exists := seen[domain]; !exists {
			seen[domain] = struct{}{}
			normalized[idx] = domain
			idx++
		}
	}
	return normalized[:idx]
}
