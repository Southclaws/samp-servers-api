package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestQueryDaemon_Setup(t *testing.T) {
	app.qd.Add("198.251.83.150:7777")
	time.Sleep(time.Second)
	assert.Equal(t, app.qd.Total, 1)

	time.Sleep(time.Second)

	app.qd.Remove("198.251.83.150:7777")
	time.Sleep(time.Second)
	assert.Equal(t, app.qd.Total, 0)

	time.Sleep(time.Second)
}
