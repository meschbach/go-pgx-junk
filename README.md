# Postgres PGX Library Junk

## Usage
Grab it via:
```bash
go get -u github.com/meschbach/go-pgx-junk
```

### Loading configuration

```go
func setup() (*pgx.ConnConfig, error) {
	return go_pgx_junk.Load(go_pgx_junk.Config{
		DirConfig: "/secrets/k8s/pg-config",
		PKIConfig: "/secrets/k8s/pg-client"
    })
}
```

See [Config for options and details](config.go#L8).
