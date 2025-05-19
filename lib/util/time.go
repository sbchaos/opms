package util

import "time"

func ToISO(t1 time.Time) string {
	return t1.Format(time.RFC3339)
}
