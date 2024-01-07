package config

import (
	"os"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

const (
	// DefaultGracefulShutdownTimeout is the default graceful shutdown timeout.
	DefaultGracefulShutdownTimeout = time.Second * 30
	// MaxGracefulShutdownTimeout is the maximum graceful shutdown timeout.
	MaxGracefulShutdownTimeout = time.Minute * 5
	// DefaultOpenStreamTimeout is the default timeout for opening a stream.
	DefaultOpenStreamTimeout = time.Second * 30
	// MaxOpenStreamTimeout is the maximum timeout for opening a stream.
	MaxOpenStreamTimeout = time.Minute * 5
)

// Config contains broker server application configuration.
type Config struct {
	SubscriberPort          int           `yaml:"subscriberPort"`
	PublisherPort           int           `yaml:"publisherPort"`
	GracefulShutdownTimeout time.Duration `yaml:"gracefulShutdownTimeout"`
	OpenStreamTimeout       time.Duration `yaml:"openStreamTimeout"`
}

// LoadConfig loads the configuration from the given path.
func LoadConfig(path string) (Config, error) {
	reader := &reader{
		osReadFile:    os.ReadFile,
		yamlUnmarshal: yaml.Unmarshal,
	}
	return reader.readConfig(path)
}

type reader struct {
	osReadFile    func(string) ([]byte, error)
	yamlUnmarshal func([]byte, interface{}) error
}

func (r *reader) readConfig(path string) (Config, error) {
	data, err := r.osReadFile(path)
	if err != nil {
		return Config{}, errors.Wrap(err, "read file")
	}

	var config Config
	if err := r.yamlUnmarshal(data, &config); err != nil {
		return Config{}, errors.Wrap(err, "unmarshal yaml")
	}

	config.GracefulShutdownTimeout = clampDuration(
		config.GracefulShutdownTimeout,
		DefaultGracefulShutdownTimeout,
		MaxGracefulShutdownTimeout,
	)
	config.OpenStreamTimeout = clampDuration(
		config.OpenStreamTimeout,
		DefaultOpenStreamTimeout,
		MaxOpenStreamTimeout,
	)

	return config, nil
}

func clampDuration(
	duration, defaultDuration, maxDuration time.Duration,
) time.Duration {
	if duration == 0 {
		return defaultDuration
	}
	if duration > maxDuration {
		return maxDuration
	}
	return duration
}
