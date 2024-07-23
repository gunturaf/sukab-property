package infra

import (
	"net/url"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DBAttributes struct {
	Host         string
	Port         string
	Username     string
	Password     string
	DatabaseName string
	EnableSSL    bool
}

// NewDB returns "DB".
// "DB" in Go-lang is an abstraction for
// database connection pooler.
func NewDB(attr DBAttributes) *sqlx.DB {
	u := url.URL{
		Host:     attr.Host + ":" + attr.Port,
		User:     url.UserPassword(attr.Username, attr.Password),
		Scheme:   "postgresql",
		Path:     attr.DatabaseName,
		RawQuery: "sslmode=disable",
	}

	if attr.EnableSSL {
		u.RawQuery = ""
	}

	db, err := sqlx.Connect("postgres", u.String())
	if err != nil {
		// we just let the app to crash because we need actual db conn to proceed.
		panic(err)
	}
	return db
}
