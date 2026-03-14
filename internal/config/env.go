package config

import "os"

type EnvKey string

func (key EnvKey) GetValue() string {
	return os.Getenv((string(key)))
}

func (key EnvKey) IsSet() bool {
	_, set := os.LookupEnv(string(key))
	return set
}
