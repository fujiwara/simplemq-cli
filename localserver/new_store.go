//go:build !sqlite

package localserver

import (
	"fmt"
	"time"
)

// NewStore creates a new in-memory Store.
// If database is specified, it returns an error because sqlite build tag is not enabled.
func NewStore(visibilityTimeout, messageExpiration time.Duration, database string) (Store, error) {
	if database != "" {
		return nil, fmt.Errorf("--database is specified but sqlite build tag is not enabled; rebuild with -tags sqlite")
	}
	return NewMemoryStore(visibilityTimeout, messageExpiration), nil
}
