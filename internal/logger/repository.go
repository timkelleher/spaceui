package logger

import "strings"

type LogRepository struct {
	items []string
	limit int
}

func NewLogRepository(limit int) LogRepository {
	return LogRepository{limit: limit}
}

func (lr *LogRepository) Push(item string) {
	lr.items = append([]string{item}, lr.items...)
	if len(lr.items) > lr.limit {
		lr.items = lr.items[0:lr.limit]
	}
}

func (lr *LogRepository) Logs() string {
	return strings.Join(lr.items, "\n")
}
