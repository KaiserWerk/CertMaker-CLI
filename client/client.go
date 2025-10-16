package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/KaiserWerk/CertMaker-CLI/auth"
	"github.com/KaiserWerk/CertMaker-CLI/challenge"
	"github.com/KaiserWerk/CertMaker-CLI/entity"
)

var (
	httpClient = &http.Client{
		Timeout: 5 * time.Minute,
	}
	httpChallengeSolver challenge.Solver
	dnsChallengeSolver  challenge.Solver
)

func requestCertificateWithCSR(csr []byte, days int) ([]byte, error) {
	// days is currently unused
	req, _ := http.NewRequest(http.MethodPost, auth.InstanceURL()+"/certificate/request-with-csr", bytes.NewBuffer(csr))
	auth.SetAuthHeader(req)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		var certResp entity.CertificateResponse
		if err := json.NewDecoder(resp.Body).Decode(&certResp); err != nil {
			return nil, err
		}
		return []byte(certResp.CertificatePEM), nil
		//certificate was issued immediately, return it
	} else if resp.StatusCode == http.StatusAccepted {
		// a challenge needs to be solved
		// TODO
	} else if resp.StatusCode == http.StatusUnauthorized {
		// authentication failed
		return nil, fmt.Errorf("authentication failed")
	} else if resp.StatusCode == http.StatusBadRequest {
		// bad request, possibly invalid CSR
		return nil, fmt.Errorf("bad request, possibly invalid CSR")
	} else {
		// some other error
		return nil, fmt.Errorf("unexpected server response: %s", resp.Status)
	}

}
