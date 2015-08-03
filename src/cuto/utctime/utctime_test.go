package utctime

import (
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	var err error
	time.Local, err = time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic("Setup failed: " + err.Error())
	}
	os.Exit(m.Run())
}

func TestParse(t *testing.T) {
	value := "2015-07-30 15:44:07.123"
	u, err := Parse(Default, value)
	if err != nil {
		t.Fatalf("Error occured: %s", err)
	}
	if u.tm.Hour() != 15 {
		t.Errorf("Parsed value wrong: %s", u.String())
	}
	if u.tm.Location().String() != "UTC" {
		t.Errorf("Location wrong: %s", u.tm.Location())
	}
}

func TestParseLocaltime(t *testing.T) {
	value := "2015-07-30 15:44:07.123"
	u, err := ParseLocaltime(Default, value)
	if err != nil {
		t.Fatalf("Error occured: %s", err)
	}
	if u.tm.Hour() != 6 {
		t.Errorf("Parsed value wrong: %s", u.String())
	}
	if u.tm.Location().String() != "UTC" {
		t.Errorf("Location wrong: %s", u.tm.Location())
	}
}

func TestTimeString(t *testing.T) {
	u := UTCTime{}
	var err error
	u.tm, err = time.ParseInLocation(Default, "2015-07-30 15:44:07.000", time.UTC)
	if err != nil {
		t.Fatalf("Error occured: %s", err)
	}
	if u.String() != "2015-07-30 15:44:07.000" {
		t.Errorf("u.String() => %s, want %s", u.String(), "2015-07-30 15:44:07.000")
	}
}

func TestTimeFormat(t *testing.T) {
	u := UTCTime{}
	var err error
	u.tm, err = time.ParseInLocation(Default, "2015-07-30 15:44:07.000", time.UTC)
	if err != nil {
		t.Fatalf("Error occured: %s", err)
	}
	if u.Format("20060102150405") != "20150730154407" {
		t.Errorf("u.Format() => %s, want %s", u.Format("20060102150405"), "20150730154407")
	}
}

func TestTimeFormatInLocal(t *testing.T) {
	u := UTCTime{}
	var err error
	u.tm, err = time.ParseInLocation(Default, "2015-07-30 15:44:07.000", time.UTC)
	if err != nil {
		t.Fatalf("Error occured: %s", err)
	}
	if u.FormatLocaltime("20060102150405") != "20150731004407" {
		t.Errorf("u.FormatLocaltime() => %s, want %s", u.FormatLocaltime("20060102150405"), "20150730154407")
	}
}
