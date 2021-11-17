package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSomething(t *testing.T) {
	date := time.Date(2021, 10, 5, 10, 0, 0, 0, time.UTC)

	dur, _ := time.ParseDuration("30m")

	assert.True(t, CheckEvent(
		date.Add(dur*0),
		date.Add(dur*4),
		date.Add(dur*1),
		date.Add(dur*2),
		"",
	), "")
	assert.True(t, CheckEvent(
		date.Add(dur*0),
		date.Add(dur*2).AddDate(1, 1, 1),
		date.Add(dur*1),
		date.Add(dur*3).AddDate(1, 1, 1),
		"",
	), "")
	assert.True(t, CheckEvent(
		date.Add(dur*1),
		date.Add(dur*4),
		date.Add(dur*0),
		date.Add(dur*2),
		"",
	), "")

	assert.True(t, CheckEvent(
		date.Add(dur*0),
		date.Add(dur*2),
		date.Add(dur*1).AddDate(0, 1, 5),
		date.Add(dur*3).AddDate(0, 1, 5),
		"day",
	), "")

	assert.True(t, CheckEvent(
		date.Add(dur*0),
		date.Add(dur*2),
		date.Add(dur*1).AddDate(0, 0, 14),
		date.Add(dur*3).AddDate(0, 0, 14),
		"week",
	), "")

	assert.True(t, CheckEvent(
		date.Add(dur*0),
		date.Add(dur*2),
		date.Add(dur*1).AddDate(0, 0, 35),
		date.Add(dur*3).AddDate(0, 0, 35),
		"month",
	), "")

	assert.True(t, CheckEvent(
		date.Add(dur*0),
		date.Add(dur*2),
		date.Add(dur*1).AddDate(1, 0, 0),
		date.Add(dur*3).AddDate(1, 0, 0),
		"year",
	), "")
}
