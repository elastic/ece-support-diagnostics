package config

import (
	"path/filepath"

	"github.com/elastic/beats/libbeat/logp"
)

func (c *Config) setupLogging() {
	config := logp.Config{
		Beat:       "ece-support-diag",
		JSON:       false,
		Level:      logp.DebugLevel,
		ToStderr:   false,
		ToSyslog:   false,
		ToFiles:    true,
		ToEventLog: false,
		Files: logp.FileConfig{
			Path:        filepath.Dir(c.DiagnosticLogFilePath()),
			Name:        filepath.Base(c.DiagnosticLogFilePath()),
			MaxSize:     20000000,
			MaxBackups:  0,
			Permissions: 0644,
			// Interval:       4 * time.Hour,
			RedirectStderr: true,
		},
	}
	logp.Configure(config)
}
