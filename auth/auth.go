package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

func init() {
	authFile := filepath.Join(getHomeDir(), ".auth.json")
	err := Load(authFile)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			panic(err)
		}
	}
}

type AuthInfo struct {
	InstanceURL string `json:"base_url"`
	Token       string `json:"token"`
}

var (
	auth        = AuthInfo{}
	currentFile string
)

func Load(file string) error {
	currentFile = file
	cont, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	return json.Unmarshal(cont, &auth)
}

func Set(instanceUrl, token string) error {
	if instanceUrl == "" {
		return errors.New("cm: missing instance URL")
	}
	if token == "" {
		return errors.New("cm: missing token")
	}

	baseDir := filepath.Dir(currentFile)
	err := os.MkdirAll(baseDir, 0700)
	if err != nil {
		return err
	}

	auth.InstanceURL = instanceUrl
	auth.Token = token

	cont, err := json.Marshal(auth)
	if err != nil {
		return err
	}

	return os.WriteFile(currentFile, cont, 0600)
}

func Remove() error {
	return os.Remove(currentFile)
}

func SetAuthHeader(r *http.Request) {
	r.Header.Set("X-Api-Token", auth.Token)
}

func getHomeDir() string {
	var homedir string
	if runtime.GOOS == "windows" {
		homedir = os.Getenv("USERPROFILE")
	} else {
		homedir = os.Getenv("HOME")
	}

	return filepath.Join(homedir, "KaiserWerk", "CertMaker-CLI")
}
