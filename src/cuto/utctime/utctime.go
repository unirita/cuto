package utctime

import "time"

const (
	Default     = "2006-01-02 15:04:05.000"
	NoDelimiter = "20060102150405.000"
	Date8Num    = "20060102"
)

type UTCTime struct {
	tm time.Time
}

// Now creates UTCTime object with current UTC time.
func Now() UTCTime {
	return UTCTime{tm: time.Now().UTC()}
}

// Parse parses value as UTC.
// The layout defines the format by showing how the reference time.
//
// If you need more information of layout, look document for time.Time#Parse.
func Parse(layout, value string) (UTCTime, error) {
	t, err := time.ParseInLocation(layout, value, time.UTC)
	if err != nil {
		return UTCTime{}, err
	}
	return UTCTime{tm: t.UTC()}, nil
}

// ParseLocaltime parses value as localtime.
func ParseLocaltime(layout, value string) (UTCTime, error) {
	t, err := time.ParseInLocation(layout, value, time.Local)
	if err != nil {
		return UTCTime{}, err
	}
	return UTCTime{tm: t.UTC()}, nil
}

// String returns the time formatted using the Default layout.
func (u *UTCTime) String() string {
	return u.Format(Default)
}

// Format returns a textual representation of the time value formatted according to layout.
func (u *UTCTime) Format(layout string) string {
	return u.tm.Format(layout)
}
