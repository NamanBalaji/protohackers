package db

import "fmt"

type DB struct {
	store   map[string]string
	version string
}

func New(version string) *DB {
	return &DB{
		store:   make(map[string]string),
		version: version,
	}
}

func (d *DB) Set(key, value string) {
	d.store[key] = value
}

func (d *DB) Get(key string) string {
	if key == "version" {
		return fmt.Sprintf("version=%s", d.version)
	}

	return fmt.Sprintf("%s=%s", key, d.store[key])
}
