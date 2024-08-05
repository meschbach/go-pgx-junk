package go_pgx_junk

import (
	"github.com/jackc/pgx/v5"
	"io/fs"
)

type Config struct {
	// DirConfig, if set and exists, points to a set of files to load the Postgres configuration from.  Accepts either a
	// single DSN/URL or a set of files describe various attributes.
	// DSN format:
	//   * $base/pg-connection: In either DSN key=value format or postgres URL format.
	// Attribute format:
	//   * pg-host
	//   * pg-port
	//   * pg-database
	//   * pg-user
	//   * pg-password
	DirConfig string
	// PKIConfig, if set and exists, points to the k8s tls secrets.  Client certificate and secret is optional however
	// at least the certificate authority is expected.
	// Format:
	//   * ca.crt (require) Certificate authority common
	//   * tls.crt - Certificate this client should authorize as
	//   * tls.key - Key for the client
	PKIConfig string
}

func Load(c Config, from fs.SubFS) (*pgx.ConnConfig, error) {
	/*
	 * Load base configuration
	 */
	var connectionConfig *pgx.ConnConfig
	if len(c.DirConfig) > 0 {
		if conn, err := pgx.ParseConfig(""); err == nil {
			connectionConfig = conn
		} else {
			return nil, err
		}
	} else {
		configDir, err := from.Sub(c.DirConfig)
		if err != nil {
			return nil, err
		}
		if conn, err := LoadFromFS(configDir); err == nil {
			connectionConfig = conn
		} else {
			return nil, err
		}
	}
	/*
	 * Load PKI
	 */
	if len(c.PKIConfig) > 0 {
		pkiDir, err := from.Sub(c.PKIConfig)
		if err != nil {
			return nil, err
		}
		if config, err := loadPKI(pkiDir); err != nil {
			return nil, err
		} else {
			connectionConfig.TLSConfig = config
		}
	}

	return connectionConfig, nil
}
