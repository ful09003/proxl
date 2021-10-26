// Package criteria provides the underlying primative functions which Cards uses to evaluate various metrics.
// These functions consist of things such as regex scorers, label-length scorers, or label-exlusion scorers.
// Each type of primative function SHOULD live in their own separate files (e.g. regex_scorer.go) along with their associated test files.
// Debug logging may be utilized here, though other logging activity or logging to non-debug severities SHOULD NOT happen in this package.
// Methods which make sense to share should be placed into criteria_helpers.go.

package criteria