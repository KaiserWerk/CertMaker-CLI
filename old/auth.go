package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
)

type authInfo struct {
	BaseURL string `json:"base_url"`
	Token   string `json:"token"`
}

type authenticator struct {
	file string
	info *authInfo
}

func newAuthenticator(file string) *authenticator {
	return &authenticator{
		file: file,
		info: &authInfo{},
	}
}

func (auth *authenticator) isAuthenticated() bool {
	return auth.info != nil && auth.info.Token != "" && auth.info.BaseURL != ""
}

func (auth *authenticator) load() error {
	cont, err := os.ReadFile(auth.file)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ErrNoAuthFound
		}
		return nil
	}

	return json.Unmarshal(cont, &auth.info)
}

func (auth *authenticator) set(baseUrl, token string) error {

	if baseUrl == "" || token == "" {
		return errors.New("missing required fields")
	}

	auth.info = &authInfo{}
	auth.info.BaseURL = baseUrl
	auth.info.Token = token

	cont, err := json.Marshal(auth.info)
	if err != nil {
		return err
	}

	return os.WriteFile(auth.file, cont, 0644)
}

func (auth *authenticator) newRequest(r *http.Request) error {
	if auth.info.Token == "" {
		return errors.New("no authentication token found")
	}

	r.Header.Set("Authorization", "Bearer "+auth.info.Token)
	return nil
}

func (auth *authenticator) clear() error {
	auth.info = nil
	return os.Remove(auth.file)
}

func (auth *authenticator) authInfo() (string, string) {
	return auth.info.BaseURL, auth.info.Token
}
