package main

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	_, err := LoadConfig("config-template.yml")
	if err != nil {
		t.Errorf("failed to load config: %v", err)
	}
}
