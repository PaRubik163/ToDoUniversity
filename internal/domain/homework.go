// Package domain contains business concepts and rules independent of storage
// and delivery mechanisms.
package domain

import "time"

// Homework is a task owned by a Telegram user.
type Homework struct {
	ID           int64
	Subject      string
	Title        string
	Priority     string
	Note         string
	Deadline     time.Time
	ReminderDays []int
	Completed    bool
	CompletedAt  *time.Time
}

// CalendarDays returns the number of calendar-day boundaries between dates in
// the configured time zone. It intentionally ignores the time of day.
func CalendarDays(from, to time.Time, tz *time.Location) int {
	a := time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, tz)
	b := time.Date(to.Year(), to.Month(), to.Day(), 0, 0, 0, 0, tz)
	return int(b.Sub(a).Hours() / 24)
}
