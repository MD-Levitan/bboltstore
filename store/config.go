package store

import (
	"time"

	"github.com/gorilla/sessions"
)

// Config represents a config for a session store.
type Config struct {
	// SessionOptions represents options for a session.
	SessionOptions sessions.Options
	// DBOptions represents options for a database.
	DBOptions Options
}

// Options represents options for a database.
type Options struct {
	// BucketName represents the name of the bucket which contains sessions.
	BucketName []byte
}

// setDefault sets default to the config.
func (c *Config) setDefault() {
	if c.SessionOptions.Path == "" {
		c.SessionOptions.Path = DefaultPath
	}
	if c.SessionOptions.MaxAge == 0 {
		c.SessionOptions.MaxAge = DefaultMaxAge
	}
	if c.DBOptions.BucketName == nil {
		c.DBOptions.BucketName = []byte(DefaultBucketName)
	}
}

// Defaults for sessions.Options
const (
	DefaultPath   = "/"
	DefaultMaxAge = 60 * 60 * 24 * 30 // 30days
)

// Defaults for store.Options
const (
	DefaultBucketName = "sessions"
)

// Defaults for reaper.Options
const (
	DefaultBatchSize     = 100
	DefaultCheckInterval = time.Minute
)
