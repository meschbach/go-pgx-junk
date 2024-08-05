package go_pgx_junk

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/fs"
)

func loadPKI(from fs.FS) (*tls.Config, error) {
	caBytes, err := fs.ReadFile(from, "ca.crt")
	if err != nil {
		return nil, err
	}

	trustedCertificateAuthorities, err := x509.ParseCertificates(caBytes)
	if err != nil {
		return nil, err
	}
	trustedCertificateAuthorityPool := x509.NewCertPool()
	for _, cert := range trustedCertificateAuthorities {
		trustedCertificateAuthorityPool.AddCert(cert)
	}
	cfg := &tls.Config{
		InsecureSkipVerify: false,
		RootCAs:            trustedCertificateAuthorityPool,
	}

	serviceCertificate, err := fs.ReadFile(from, "tls.crt")
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return nil, err
		}
	} else {
		serviceKey, err := fs.ReadFile(from, "tls.key")
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				return nil, errors.New("tls certificate found but not corresponding key")
			} else {
				return nil, err
			}
		} else {
			cert, err := tls.X509KeyPair(serviceCertificate, serviceKey)
			if err != nil {
				return nil, err
			}
			cfg.Certificates = append(cfg.Certificates, cert)
		}
	}

	return cfg, nil
}
