package database

import (
	"net/url"
	"os"
)

func GetConnectionString() string {
	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD")),
		Host:     os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT"),
		Path:     os.Getenv("DB_NAME"),
		RawQuery: "sslmode=disable",
	}
	return u.String()
}
