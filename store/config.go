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
	// ReaperOptions represents options for a reaper routine
	ReaperOptions ROptions
}

// Options represents options for a database.
type Options struct {
	// BucketName represents the name of the bucket which contains sessions.
	BucketName []byte
	// If FreeDB is true, then after closing DB and store, DB will be deleted.
	FreeDB bool
}

type ROptions struct {
	//StartRoutine represents the status of reaper
	StartRoutine bool
	// BucketName represents the name of the bucket which contains sessions.
	BucketName []byte
	// BatchSize represents the maximum number of sessions which the reaper
	// process at one time.
	BatchSize int
	// CheckInterval represents the interval between the reaper's invocation.
	CheckInterval time.Duration
}

func NewDefaultConfig() *Config {
	return &Config{
		SessionOptions: sessions.Options{Path: DefaultPath, MaxAge: DefaultMaxAge},
		DBOptions:      Options{BucketName: []byte(DefaultBucketName), FreeDB: false},
		ReaperOptions: ROptions{StartRoutine: false, BucketName: []byte(DefaultBucketName),
			BatchSize: DefaultBatchSize, CheckInterval: DefaultCheckInterval},
	}
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
	if c.ReaperOptions.BatchSize == 0 {
		c.ReaperOptions.BatchSize = DefaultBatchSize
	}
	if c.ReaperOptions.BucketName == nil {
		c.ReaperOptions.BucketName = []byte(DefaultBucketName)
	}
	if c.ReaperOptions.CheckInterval == 0 {
		c.ReaperOptions.CheckInterval = DefaultCheckInterval
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
