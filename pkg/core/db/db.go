package db

import (
	"database/sql"
	"fmt"
	"net/url"
)

func New(dsn string) (*sql.DB, error) {
	u, err := url.Parse(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DSN: %w", err)
	}
	if u.Scheme != "mysql" {
		return nil, fmt.Errorf("unsupported scheme: %s", u.Scheme)
	}

	credentials := u.User.Username()
	if password, ok := u.User.Password(); ok {
		credentials = fmt.Sprintf("%s:%s", credentials, password)
	}

	dsn = fmt.Sprintf(
		"%s@tcp(%s)/%s?%s",
		credentials,
		u.Host,
		u.Path[1:],
		u.RawQuery,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	return db, nil
}
