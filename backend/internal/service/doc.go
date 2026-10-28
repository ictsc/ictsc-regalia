// Package service carries out the use cases behind the API: answering,
// grading, ranking and its visibility, redeploys, content refresh, sign-in
// and sessions. It works through core.Store and the interfaces in
// dependencies.go (Discord, GitHub content, SState, event bus), so it never
// talks to a database or an external service directly.
package service
