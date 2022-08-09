package config

import (
	"testing"
)

func TestLoad(t *testing.T) {
	c, err := LoadConfig()
	if err != nil {
		t.Logf("Encountered an Error: %v\n", err)
		t.Fail()
	}
	t.Log(c)
}

func TestGetUsername(t *testing.T) {
	username, err := GetUserName()
	if err != nil {
		t.Logf("Ecountered an Error: %v", err)
		t.Fail()
	}
	t.Log(username)
}
