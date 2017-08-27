package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApp_GetServers(t *testing.T) {
	servers, err := app.GetServers("1", "", "", "")
	if err != nil {
		t.Errorf("App.GetServers() error = %v", err)
	}

	expected := []ServerCore{
		{"127.0.0.1", "Scavenge and Survive Official", 4, 32, "Scavenge & Survive by Southclaws", "English", false},
	}

	assert.ObjectsAreEqual(expected, servers)
}
