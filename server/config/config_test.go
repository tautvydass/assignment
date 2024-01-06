package config

import (
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestLoadConfig(t *testing.T) {
	var (
		path  = "test_path"
		tests = map[string]struct {
			osReadFile    func(string) ([]byte, error)
			yamlUnmarshal func([]byte, interface{}) error
			want          Config
			wantErr       error
		}{
			"error_reading_file": {
				osReadFile: func(string) ([]byte, error) {
					return nil, assert.AnError
				},
				wantErr: errors.Wrap(assert.AnError, "read file"),
			},
			"error_unmarshalling_file_content_to_yaml": {
				osReadFile: func(string) ([]byte, error) {
					return make([]byte, 0), nil
				},
				yamlUnmarshal: func([]byte, interface{}) error {
					return assert.AnError
				},
				wantErr: errors.Wrap(assert.AnError, "unmarshal yaml"),
			},
			"happy_path_with_defaults": {
				osReadFile: func(string) ([]byte, error) {
					return make([]byte, 0), nil
				},
				want: Config{
					GracefulShutdownTimeout: DefaultGracefulShutdownTimeout,
				},
			},
			"happy_path_with_out_of_bounds_graceful_shutdown_timeout": {
				osReadFile: func(string) ([]byte, error) {
					return yaml.Marshal(Config{
						GracefulShutdownTimeout: MaxGracefulShutdownTimeout + 1,
					})
				},
				want: Config{
					GracefulShutdownTimeout: MaxGracefulShutdownTimeout,
				},
			},
			"happy_path": {
				osReadFile: func(string) ([]byte, error) {
					return yaml.Marshal(Config{
						SubscriberPort:          1111,
						PublisherPort:           2222,
						GracefulShutdownTimeout: time.Second,
					})
				},
				want: Config{
					SubscriberPort:          1111,
					PublisherPort:           2222,
					GracefulShutdownTimeout: time.Second,
				},
			},
		}
	)

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			reader := &reader{
				osReadFile:    tc.osReadFile,
				yamlUnmarshal: tc.yamlUnmarshal,
			}
			if reader.yamlUnmarshal == nil {
				reader.yamlUnmarshal = yaml.Unmarshal
			}

			got, err := reader.readConfig(path)
			if tc.wantErr != nil {
				require.EqualError(t, err, tc.wantErr.Error())
				assert.Empty(t, got)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}
