package email

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReminderInfo_HTML(t *testing.T) {
	ri := &ReminderInfo{Title: "reminder", Content: "this is a reminder email"}
	_, err := ri.HTML()
	assert.NoError(t, err)
}
