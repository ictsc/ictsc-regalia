// Package github reads competition content from the master repository on
// GitHub and verifies the GitHub Actions OIDC tokens that ask for a refresh.
// Content is always read at one full commit SHA through the REST API and
// validated before it is returned; see docs/content-repository-contract.md.
package github
