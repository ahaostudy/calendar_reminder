package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseDate(t *testing.T) {
	d := "2003-10-19"
	expect := time.Date(2003, 10, 19, 0, 0, 0, 0, time.UTC)

	date, err := ParseDate(d)
	assert.NoError(t, err)
	assert.Equal(t, date, expect)
}

func TestGetDateStartAndEnd(t *testing.T) {
	date := time.Date(2003, 10, 19, 0, 0, 0, 0, time.UTC)

	start, end := GetDateStartAndEnd(date)
	assert.Equal(t, time.Date(2003, 10, 19, 0, 0, 0, 0, time.UTC), start)
	assert.Equal(t, time.Date(2003, 10, 19, 23, 59, 59, 0, time.UTC), end)
}
