package backup

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"strings"
)

func ValidateSchedule(schedule string) error {
	schedule = strings.TrimSpace(schedule)
	if schedule == "manual" {
		return nil
	}
	parser := cron.NewParser(cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	if _, err := parser.Parse(schedule); err != nil {
		return fmt.Errorf("invalid backup schedule: %w", err)
	}
	return nil
}
