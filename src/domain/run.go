package domain

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Run struct {
	gorm.Model
	Start time.Time
	End   time.Time
}

func (run *Run) ToString() string {
	if run.Start.Year() == run.End.Year() {
		if run.Start.Month() == run.End.Month() {
			if run.Start.Day() == run.End.Day() {
				return fmt.Sprintf("%s a %s de %s",
					run.Start.Format("15:04:05"),
					run.End.Format("15:04:05"),
					run.Start.Format("2/Jan/2006"))
			}
			return fmt.Sprintf("%s a %s de %s",
				run.Start.Format("2 (15:04)"),
				run.End.Format("2 (15:04)"),
				run.Start.Format("Jan/2006"))
		}
		return fmt.Sprintf("%s a %s de %s",
			run.Start.Format("2/Jan"),
			run.End.Format("2/Jan"),
			run.Start.Format("2006"))
	}
	return fmt.Sprintf("%s a %s",
		run.Start.Format("2/Jan/2006"),
		run.End.Format("2/Jan/2006"))
}
