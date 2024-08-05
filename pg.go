package go_pgx_junk

import (
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"io/fs"
)

func LoadFromFS(from fs.FS) (*pgx.ConnConfig, error) {
	connectionBytes, err := fs.ReadFile(from, "pg-connection")
	if err == nil {
		connectionString := string(connectionBytes)
		return pgx.ParseConfig(connectionString)
	}
	if !errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}

	return connectionFromFiles(from)
}

func withFileContent(from fs.FS, fileName string, errs []error, apply func(string) error) []error {
	byteContents, err := fs.ReadFile(from, fileName)
	if err != nil {
		errs = append(errs, err)
		return errs
	}
	if err := apply(string(byteContents)); err != nil {
		errs = append(errs, err)
	}
	return errs
}

func withSafeFileContent(from fs.FS, fileName string, errs []error, apply func(string)) []error {
	return withFileContent(from, fileName, errs, func(s string) error {
		apply(s)
		return nil
	})
}

func withOptionalSafeFileContent(from fs.FS, fileName string, errs []error, apply func(string)) []error {
	if byteContents, err := fs.ReadFile(from, fileName); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return errs
		} else {
			errs = append(errs, err)
			return errs
		}
	} else {
		apply(string(byteContents))
		return errs
	}
}

func connectionFromFiles(from fs.FS) (*pgx.ConnConfig, error) {
	cfg, err := pgx.ParseConfig("")
	if err != nil {
		return nil, err
	}

	var errs []error
	errs = withFileContent(from, "pg-port", errs, func(s string) error {
		_, err := fmt.Sscanf(s, "%d", &cfg.Port)
		return err
	})

	errs = withSafeFileContent(from, "pg-host", errs, func(s string) {
		cfg.Host = s
	})
	errs = withSafeFileContent(from, "pg-database", errs, func(s string) {
		cfg.Database = s
	})
	errs = withOptionalSafeFileContent(from, "pg-user", errs, func(s string) {
		cfg.User = s
	})
	errs = withOptionalSafeFileContent(from, "pg-password", errs, func(s string) {
		cfg.Password = s
	})
	errs = withOptionalSafeFileContent(from, "pg-app", errs, func(s string) {
		cfg.RuntimeParams["application_name"] = s
	})

	return cfg, errors.Join(errs...)
}
