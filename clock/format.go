package clock

import "time"

func Format(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

func Parse(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, nil
	}
	return time.ParseInLocation("2006-01-02 15:04:05", s, time.FixedZone("CST", int(Offset.Seconds())))
}

func Since(now, t time.Time) time.Duration {
	if t.IsZero() || now.IsZero() {
		return 0
	}
	return now.Sub(t)
}

func AddSeconds(t time.Time, sec int) time.Time {
	if sec == 0 {
		return t
	}
	return t.Add(time.Duration(sec) * time.Second)
}

func Zone() *time.Location {
	return time.FixedZone("CST", int(Offset.Seconds()))
}
