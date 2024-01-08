package certificate

import (
	"crypto/tls"

	"github.com/pkg/errors"
)

// LoadTLSConfig loads a TLS configuration from the given file
// paths and constructs TLS configuration.
func LoadTLSConfig(certFile, keyFile string) (*tls.Config, error) {
	l := &loader{
		tlsLoadX509KeyPair: tls.LoadX509KeyPair,
	}

	return l.loadX509KeyPair(certFile, keyFile)
}

type loader struct {
	tlsLoadX509KeyPair func(certFile, keyFile string) (tls.Certificate, error)
}

func (l *loader) loadX509KeyPair(certFile, keyFile string) (*tls.Config, error) {
	cert, err := l.tlsLoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, errors.Wrap(err, "load x509 key pair")
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
	}, nil
}
