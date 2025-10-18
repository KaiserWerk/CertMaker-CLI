package challenge

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/KaiserWerk/CertMaker-CLI/entity"
)

type DNSChallengeSolver struct {
}

func (d *DNSChallengeSolver) CanSolve(challengeType string) bool {
	return challengeType == "dns-01"
}

func (d *DNSChallengeSolver) Setup(ctx context.Context, token string, domains []string) error {
	if len(domains) == 0 {
		return ErrNoDomainsProvided
	}
	domains = normalizeDomains(domains)

	// display the token and domains (with the certmaker subdomains) to the user with the
	// instructions to create the necessary DNS TXT records.
	// the token is the same for all domains.
	fmt.Println("To complete the DNS-01 challenge, please create the following DNS TXT records:")
	for _, domain := range domains {
		fmt.Printf(" - _certmaker_challenge.%s\n", domain)
	}
	fmt.Printf("with the value: %s\n\n", token)
	fmt.Println("After creating the DNS records, please wait a moment for the changes to propagate, then press Enter to continue...")
	fmt.Scanln() // wait for user to press Enter
	// might want to add a timeout or a retry mechanism
	// to check if the DNS records are actually propagated before proceeding. Might...
	return nil
}

func (d *DNSChallengeSolver) Solve(ctx context.Context, instanceURL, challengeID string, setHeader func(*http.Request)) (*entity.CertificateResponse, error) {
	solveURL := fmt.Sprintf("%s/api/v1/dns-01/%s/solve", instanceURL, challengeID)

	httpClient := &http.Client{Timeout: 1 * time.Minute}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, solveURL, nil)
	if err != nil {
		return nil, err
	}
	setHeader(req)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("certmaker-sdk: expected status code 201, got %d", resp.StatusCode)
	}

	var response entity.CertificateResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	return &response, err
}
