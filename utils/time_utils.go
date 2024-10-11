package utils

import "time"

// formatDateTime formats a Unix timestamp to the format "02 Jan 2006 15:04".
func FormatDateTimeUTC(unixTime int64) string {
	t := time.Unix(unixTime, 0).UTC()
	return t.Format("02 Jan 2006 15:04")
}

func CalculateRemainingDays(unixTimestamp int64) int {
	targetTime := time.Unix(unixTimestamp, 0)
	now := time.Now()
	duration := targetTime.Sub(now)
	daysRemaining := int(duration.Hours() / 24)

	return daysRemaining
}
