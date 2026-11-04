// Package postgres implements core.Store on PostgreSQL. Changes that span
// several rows run in one transaction, and competition-wide changes (rule,
// ranking freeze, final reveal, score recalculation) are recorded in
// audit_log. The schema lives in db/migrations.
package postgres
