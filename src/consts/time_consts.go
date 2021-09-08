package consts

import "time"

const (
	Day   = time.Duration(24) * time.Hour
	Week  = time.Duration(7) * Day
	Month = time.Duration(30) * Day
)
