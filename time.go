package tradier

import (
	"strconv"
	"time"
)

// DateTime wraps time.Time and adds flexible implementations for unmarshaling
// JSON in the different forms it appears in the Tradier API.
type DateTime struct {
	time.Time
}

func (d *DateTime) Set(s string) error {
	if s == "null" {
		return nil
	}

	t, err := time.Parse(time.RFC3339, s)
	if err == nil {
		*d = DateTime{t}
		return nil
	}

	// Date and time.
	t, err = time.Parse("2006-01-02T15:04:05", s)
	if err == nil {
		*d = DateTime{t}
		return nil
	}

	// Just the date
	t, err = time.Parse("2006-01-02", s)
	if err == nil {
		*d = DateTime{t}
		return nil
	}

	// Just the hour
	t, err = time.Parse("15:04", s)
	if err == nil {
		*d = DateTime{t}
		return nil
	}

	// Milliseconds since the Unix epoch.
	t, err = ParseTimeMs(s)
	if err == nil {
		*d = DateTime{t}
		return nil
	}

	return err
}

func (d *DateTime) UnmarshalJSON(b []byte) error {
	if b[0] == '"' && b[len(b)-1] == '"' {
		b = b[1 : len(b)-1]
	}
	s := string(b)

	return d.Set(s)
}

func ParseTimeMs(tsMs string) (time.Time, error) {
	msecs, err := strconv.ParseInt(tsMs, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	secs := msecs / 1000
	nsecs := 1000000 * (msecs % 1000)
	t := time.Unix(secs, nsecs)
	return t, nil
}
