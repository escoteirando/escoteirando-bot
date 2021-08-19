package domain

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Schedule struct {
	gorm.Model
	Id         string `gorm:"primaryKey"`
	LastRun    time.Time
	NextRun    time.Time
	Interval   time.Duration
	RoundStart bool
}

func (schedule Schedule) ToString() string {
	return fmt.Sprintf("Schedule '%s' LastRun=%v NextRun=%v", schedule.Id, schedule.LastRun, schedule.NextRun)
}
