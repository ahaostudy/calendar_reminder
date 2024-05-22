.PHNOY: format
format:
	gofumpt -w -extra . && goimports  -local github.com/ahaostudy/calendar_reminder -w .
