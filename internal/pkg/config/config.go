package config

import (
	"net/http"
	"os"
)

// Config defines the configuration object
type Config struct {
	RequestFile *os.File
	Transport   *http.Transport
	ForceSSL    bool
}
