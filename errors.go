package main

import "errors"

var (
	ErrNoAuthSelected = errors.New("certctl: no authentication selected from list")
	ErrNoAuthFound    = errors.New("certctl: no entry file found, empty or malformed")
)
