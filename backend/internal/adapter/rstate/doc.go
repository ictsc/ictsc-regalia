// Package rstate asks rstate (github.com/ictsc/rstate), the Terraform-based
// redeploy tool used in the secondary qualifier, to redeploy a team's problem
// and reads back its job state. rstate has no callbacks, so the service polls
// Observe; see docs/rstate-integration.md.
package rstate
