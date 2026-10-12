package dsn

import (
	"net"
	"net/url"
	"strconv"

	"interval_attribution/internal/app/config"
)

func FromConfig(cfg config.DatabaseConfig) string {
	return (&url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.User, cfg.Password),
		Host:   net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port)),
		Path:   cfg.Name,
		RawQuery: url.Values{
			"sslmode": []string{"disable"},
		}.Encode(),
	}).String()
}
