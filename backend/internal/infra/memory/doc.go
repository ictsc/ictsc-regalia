// Package memory implements the store, sessions and event bus in process
// memory. It is used only in development fake mode when no database or Redis
// is configured, and loses everything on restart.
package memory
