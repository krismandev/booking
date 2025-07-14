package utils

import "time"

func FormatTime(time time.Time) string {
	var output string
	layout := "2006-01-02 15:04:05.000 -0700"
	output = time.Format(layout) // "2006-01-02 15:04:05.000 -0700"  // RFC3339 format
	return output
}

func TimeNowString() string {
	now := time.Now().Format(time.RFC3339)
	return now
}
