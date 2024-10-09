package utils

import (
	"testing"
)

func TestFormatDateTime(t *testing.T) {
	tests := []struct {
		unixTime int64
		expected string
	}{
		{1633024800, "30 Sep 2021 18:00"}, // 30 Sep 2021 18:00 UTC
		{1609459200, "01 Jan 2021 00:00"}, // 01 Jan 2021 00:00 UTC
		{0, "01 Jan 1970 00:00"},          // Unix epoch (01 Jan 1970 00:00 UTC)
		{1672531199, "31 Dec 2022 23:59"}, // 31 Dec 2022 23:59 UTC
	}

	for _, test := range tests {
		result := FormatDateTimeUTC(test.unixTime)
		if result != test.expected {
			t.Errorf("expected %s, got %s for unix time %d", test.expected, result, test.unixTime)
		}
	}
}
