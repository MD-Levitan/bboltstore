package store

import (
	"testing"
)

func TestConfig_setDefault(t *testing.T) {
	config := Config{}
	config.setDefault()
	if config.SessionOptions.Path != DefaultPath {
		t.Errorf("config.SessionOptions.Path should be %s (actual: %s)", DefaultPath, config.SessionOptions.Path)
	}
	if config.SessionOptions.MaxAge != DefaultMaxAge {
		t.Errorf("config.SessionOptions.MaxAge should be %d (actual: %d)", DefaultMaxAge, config.SessionOptions.MaxAge)
	}
	if string(config.DBOptions.BucketName) != DefaultBucketName {
		t.Errorf("config.SessionOptions.BucketName should be %+v (actual: %+v)", DefaultBucketName, config.DBOptions.BucketName)
	}
}
