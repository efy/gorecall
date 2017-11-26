package server

import (
	"net/http"
	"net/http/cgi"
	"net/http/fcgi"
	"time"
)

var (
	HTTP = &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	FCGI = fcgi.Serve

	CGI = cgi.Serve
)
