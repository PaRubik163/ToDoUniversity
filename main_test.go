package main

import (
	"testing"
	"time"
)

func TestParseHomework(t *testing.T) {
	tz := time.FixedZone("MSK", 3*60*60)
	now := time.Date(2026, time.September, 11, 12, 0, 0, 0, tz)
	subject, title, deadline, err := parseHomework("Предмет: Разработка безопасного ПО\nЗадание: Практическая 1\nДедлайн: 6 октября", tz, now)
	if err != nil {
		t.Fatal(err)
	}
	if subject != "Разработка безопасного ПО" || title != "Практическая 1" {
		t.Fatalf("unexpected fields: %q, %q", subject, title)
	}
	want := time.Date(2026, time.October, 6, 0, 0, 0, 0, tz)
	if !deadline.Equal(want) {
		t.Fatalf("deadline = %v, want %v", deadline, want)
	}
}

func TestParseHomeworkMovesPastDateToNextYear(t *testing.T) {
	tz := time.UTC
	now := time.Date(2026, time.October, 7, 0, 0, 0, 0, tz)
	_, _, deadline, err := parseHomework("Предмет: Математика\nЗадание: Лабораторная\nДедлайн: 6 октября", tz, now)
	if err != nil {
		t.Fatal(err)
	}
	if deadline.Year() != 2027 {
		t.Fatalf("year = %d, want 2027", deadline.Year())
	}
}

func TestParseHomeworkRequiresAllFields(t *testing.T) {
	_, _, _, err := parseHomework("Предмет: Математика\nДедлайн: 6 октября", time.UTC, time.Now())
	if err == nil {
		t.Fatal("expected error")
	}
}
