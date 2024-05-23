package email

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReminderInfo_HTML(t *testing.T) {
	ri := &ReminderInfo{Title: "reminder", Content: "this is a reminder email"}
	_, err := ri.HTML()
	assert.NoError(t, err)
}
