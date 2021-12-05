package internal

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const (
	ENV_PREFIX    = "SIMGR_"
	OPTIONS_PORT  = "port"
	DEFAULT_PORT  = 8080
	OPTIONS_QUIET = "quiet"
	DEFAULT_QUIET = false
	OPTIONS_URL   = "url"
)

type Config struct {
	Port  int
	Quiet bool
	URL   *url.URL
}

func envInt(name string, def int) (int, error) {
	if val := os.Getenv(name); val != "" {
		intVal, err := strconv.Atoi(val)
		if err != nil {
			return 0, fmt.Errorf("Invalid value for %s: %w", name, err)
		}
		return intVal, nil
	}
	return def, nil
}

func envBool(name string, def bool) (bool, error) {
	if val := os.Getenv(name); val != "" {
		lowerVal := strings.ToLower(val)
		// Accept "y", "yes", "t", "true", "s", "si", "sí"...
		for _, pref := range []string{"y", "s", "t"} {
			if strings.HasPrefix(lowerVal, pref) {
				return true, nil
			}
		}
	}
	return def, nil
}

func envName(name string) string {
	return ENV_PREFIX + strings.ToUpper(name)
}

func NewConfig() (c Config, err error) {
	defPort, err := envInt(envName(OPTIONS_PORT), DEFAULT_PORT)
	if err != nil {
		return c, err
	}
	defVerbose, err := envBool(envName(OPTIONS_QUIET), DEFAULT_QUIET)
	if err != nil {
		return c, err
	}
	defURL := os.Getenv(envName(OPTIONS_URL))
	var strURL string
	flag.IntVar(&c.Port, OPTIONS_PORT, defPort, "HTTP Listen Port Number")
	flag.BoolVar(&c.Quiet, OPTIONS_QUIET, defVerbose, "Disable regular logging")
	flag.StringVar(&strURL, OPTIONS_URL, "", "URL to proxy API requests to")
	flag.Parse()
	if c.Port <= 1024 || c.Port > 65535 {
		return c, fmt.Errorf("Invalid port number %d, must be between 1025 and 65535", c.Port)
	}
	if strURL == "" {
		strURL = defURL
	}
	if strURL == "" {
		return c, fmt.Errorf("URL must be specified either via CLI -url or env %s", envName(OPTIONS_URL))
	}
	c.URL, err = url.Parse(strURL)
	if err != nil {
		return c, err
	}
	return c, nil
}
