package config

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

func NewHttpClient() *http.Client {
	var tr = &http.Transport{
		// Disable Certificate Checking
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		ResponseHeaderTimeout: 15 * time.Second,
		// Connection timeout = 5s
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		// TLS Handshake Timeout = 5s
		TLSHandshakeTimeout: 5 * time.Second,
	}
	// HTTP Timeout = 10s
	return &http.Client{Timeout: 10 * time.Second, Transport: tr}
}
