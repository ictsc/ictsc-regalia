// Package config reads the runtime settings from ICTSC_* environment
// variables and validates them all at start-up, so a misconfigured server
// refuses to start instead of failing on the first request.
package config
