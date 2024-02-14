package config

import (
	"errors"
	"fmt"
	"os"
)

// Provider is an interface for getting a value from a secret store.
type Provider interface {
	Get(key string) (string, error)
}

// Configuration ...
type Configuration struct {
	provider Provider
}

// New ...
func New(provider Provider) *Configuration {
	c := &Configuration{
		provider: provider,
	}

	return c
}

// Get returns the value from environment variable `<key>`. When an environment variable `<key>_SECURE` exists
// the provider is used for getting the value.
func (c *Configuration) Get(key string) (string, error) {
	res := os.Getenv(key)

	if c.provider != nil {
		if secret, ok := os.LookupEnv(key + "_SECURE"); ok {
			if c.provider == nil {
				return "", errors.New("provider is nil")
			}

			providedSecret, err := c.provider.Get(secret)

			if err != nil {
				return "", fmt.Errorf("failed to get secret: %w", err)
			}

			res = providedSecret
		}
	}

	if res == "" {
		return "", errors.New("empty value")
	}

	return res, nil
}
