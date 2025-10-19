package client

import (
	"bytes"
	"context"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"time"

	"github.com/KaiserWerk/CertMaker-CLI/auth"
	"github.com/KaiserWerk/CertMaker-CLI/challenge"
	"github.com/KaiserWerk/CertMaker-CLI/entity"
)

var (
	httpClient = &http.Client{
		Timeout: 1 * time.Minute,
	}
)

func RequestCertificateWithCSR(csr []byte, days int, challengeType string) ([]byte, error) {
	// TODO: days is currently unused
	req, _ := http.NewRequest(http.MethodPost, auth.InstanceURL()+"/certificate/request-with-csr", bytes.NewBuffer(csr))
	auth.SetAuthHeader(req)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusCreated:
		var certResp entity.CertificateResponse
		if err := json.NewDecoder(resp.Body).Decode(&certResp); err != nil {
			return nil, err
		}
		return []byte(certResp.CertificatePEM), nil
		// certificate was issued without a challenge, return it
	case http.StatusAccepted:
		// a challenge needs to be solved

		// parse the CSR to extract domains
		block, _ := pem.Decode(csr)
		if block == nil || block.Type != "CERTIFICATE REQUEST" {
			return nil, fmt.Errorf("failed to parse CSR PEM")
		}

		csrData, err := x509.ParseCertificateRequest(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse CSR: %w", err)
		}

		var challengeResp entity.CertificateResponse
		if err := json.NewDecoder(resp.Body).Decode(&challengeResp); err != nil {
			return nil, err
		}

		// no challenges provided? No way to get a certificate then
		if len(challengeResp.Challenges) == 0 {
			return nil, fmt.Errorf("no challenges provided by server")
		}

		// determine the challenge to solve
		var challengeToSolve entity.ChallengeResponse
		for _, ch := range challengeResp.Challenges {
			if ch.ChallengeType == challengeType {
				challengeToSolve = ch
				break
			}
		}

		var solver challenge.Solver
		switch challengeType {
		case "http-01":
			solver = &challenge.HTTP01ChallengeSolver{}
		case "dns-01":
			solver = &challenge.DNS01ChallengeSolver{}
		default:
			return nil, fmt.Errorf("unsupported challenge type: %s", challengeType)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		err = solver.Setup(ctx, challengeToSolve.ChallengeToken, csrData.DNSNames)
		if err != nil {
			return nil, fmt.Errorf("failed to set up challenge solver: %w", err)
		}
		certResp, err := solver.Solve(ctx, auth.InstanceURL(), challengeToSolve.ChallengeID, auth.SetAuthHeader)
		if err != nil {
			return nil, fmt.Errorf("failed to solve challenge: %w", err)
		}

		return []byte(certResp.CertificatePEM), nil
	case http.StatusUnauthorized:
		// authentication failed
		return nil, fmt.Errorf("authentication failed")
	case http.StatusBadRequest:
		// bad request, possibly invalid CSR
		return nil, fmt.Errorf("bad request, possibly invalid CSR")
	default:
		// some other error
		return nil, fmt.Errorf("unexpected server response: %s", resp.Status)
	}
}

func RequestCertificateWithSimpleRequest(domains []string, ips []string, emails []string, days int, challengeType string) ([]byte, []byte, error) {

	sr := entity.SimpleRequest{
		Domains:        domains,
		IPs:            ips,
		EmailAddresses: emails,
		Days:           days,
	}
	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(sr); err != nil {
		return nil, nil, err
	}
	req, _ := http.NewRequest(http.MethodPost, auth.InstanceURL()+"/certificate/request-with-simple-request", &b)
	auth.SetAuthHeader(req)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusCreated:
		var certResp entity.CertificateResponse
		if err := json.NewDecoder(resp.Body).Decode(&certResp); err != nil {
			return nil, nil, err
		}
		return []byte(certResp.CertificatePEM), nil, nil
		//certificate was issued immediately, return it
	case http.StatusAccepted:
		// a challenge needs to be solved

		var challengeResp entity.CertificateResponse
		if err := json.NewDecoder(resp.Body).Decode(&challengeResp); err != nil {
			return nil, nil, err
		}

		if len(challengeResp.Challenges) == 0 {
			return nil, nil, fmt.Errorf("no challenges provided by server")
		}

		var challengeToSolve entity.ChallengeResponse
		for _, ch := range challengeResp.Challenges {
			if ch.ChallengeType == challengeType {
				challengeToSolve = ch
				break
			}
		}

		var solver challenge.Solver
		switch challengeType {
		case "http-01":
			solver = &challenge.HTTP01ChallengeSolver{}
		case "dns-01":
			solver = &challenge.DNS01ChallengeSolver{}
		default:
			return nil, nil, fmt.Errorf("unsupported challenge type: %s", challengeType)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		err = solver.Setup(ctx, challengeToSolve.ChallengeToken, domains)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to set up challenge solver: %w", err)
		}
		certResp, err := solver.Solve(ctx, auth.InstanceURL(), challengeToSolve.ChallengeID, auth.SetAuthHeader)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to solve challenge: %w", err)
		}

		return []byte(certResp.CertificatePEM), []byte(certResp.PrivateKeyPEM), nil
	case http.StatusUnauthorized:
		// authentication failed, user should do re-auth
		return nil, nil, fmt.Errorf("authentication failed")
	case http.StatusBadRequest:
		// bad request, possibly encoding or format error
		return nil, nil, fmt.Errorf("bad request, possibly invalid CSR")
	default:
		// some other error
		return nil, nil, fmt.Errorf("unexpected server response: %s", resp.Status)
	}
}
