package utils

import "time"

func ParseTime(ts, format string) (time.Time, error) {
	t, err := time.Parse(format, ts)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

var Days = map[string]int{
	"daily":      1,
	"weekly":     7,
	"monthly":    31,
	"quarterly":  92,
	"semiannual": 183,
	"annual":     365,
}

var ValidFrequencies = []string{
	"daily", "weekly", "monthly", "quarterly", "semiannual", "annual",
}
