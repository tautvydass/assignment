package certificate

import (
	"crypto/tls"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoader_loadX509KeyPair(t *testing.T) {
	t.Run("error_loading_key_pair", func(t *testing.T) {
		l := &loader{
			tlsLoadX509KeyPair: func(certFile, keyFile string) (tls.Certificate, error) {
				return tls.Certificate{}, assert.AnError
			},
		}

		got, err := l.loadX509KeyPair("certFile", "keyFile")
		require.True(t, errors.Is(err, assert.AnError))
		assert.Empty(t, got)
	})
	t.Run("happy_path", func(t *testing.T) {
		cert := tls.Certificate{
			Certificate: [][]byte{
				[]byte("test certificate"),
			},
		}
		l := &loader{
			tlsLoadX509KeyPair: func(certFile, keyFile string) (tls.Certificate, error) {
				return cert, nil
			},
		}

		got, err := l.loadX509KeyPair("certFile", "keyFile")
		require.NoError(t, err)
		assert.Equal(t, &tls.Config{
			Certificates: []tls.Certificate{cert},
		}, got)
	})
}
