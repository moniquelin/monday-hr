package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	parseDates("2026-02-23")
}

func parseDates(date string) {
	// Parse dates
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		log.Fatal(err)
	}
	year, month, day := t.Date()

	fmt.Println(year, int(month), day)
}
