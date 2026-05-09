package shared

import "time"

const DateLayout = "2006-01-02"

func ParseDate(value string) (time.Time, error) {
	return time.Parse(DateLayout, value)
}

func AddDays(t time.Time, days int) time.Time {
	return t.AddDate(0, 0, days)
}

func DaysUntilEndOfYear(t time.Time) int {
	endOfYear := time.Date(
		t.Year(),
		12,
		31,
		0,
		0,
		0,
		0,
		t.Location(),
	)

	return int(endOfYear.Sub(t).Hours()/24) + 1
}
