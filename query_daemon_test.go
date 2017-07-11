package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestQueryDaemon_1(t *testing.T) {
	app.qd.Add("198.251.83.150:7777")
	time.Sleep(time.Second)
	assert.Equal(t, app.qd.GetTotal(), int32(1))

	time.Sleep(time.Second)

	app.qd.Remove("198.251.83.150:7777")
	time.Sleep(time.Second)
	assert.Equal(t, app.qd.GetTotal(), int32(0))

	time.Sleep(time.Second)
}

func TestQueryDaemon_2(t *testing.T) {
	app.qd.Add("176.32.37.26:7777")
	app.qd.Add("176.32.39.168:7777")
	app.qd.Add("176.32.39.151:7777")
	app.qd.Add("176.32.39.80:7777")
	app.qd.Add("176.32.36.91:7777")

	time.Sleep(6 * time.Second)

	assert.Equal(t, app.qd.GetTotal(), int32(5))
}

func TestQueryDaemon_3(t *testing.T) {
	app.qd.Add("198.251.83.150:7777")
	time.Sleep(time.Millisecond)
	app.qd.Remove("198.251.83.150:7777")

	app.qd.Add("176.32.37.26:7777")
	app.qd.Add("176.32.39.168:7777")
	app.qd.Add("176.32.39.151:7777")
	app.qd.Add("176.32.39.80:7777")
	app.qd.Add("176.32.36.91:7777")

	time.Sleep(6 * time.Second)

	assert.Equal(t, app.qd.GetTotal(), int32(5))
}
