package backup

import (
	"fmt"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
)

func ValidateSchedule(schedule string) error {
	_, err := ParseSchedule(schedule, "")
	return err
}

func ParseSchedule(schedule, timezone string) (cron.Schedule, error) {
	schedule = strings.TrimSpace(schedule)
	timezone = strings.TrimSpace(timezone)
	if timezone != "" {
		if _, err := time.LoadLocation(timezone); err != nil {
			return nil, fmt.Errorf("invalid backup timezone %q: %w", timezone, err)
		}
	}
	if schedule == "manual" {
		return nil, nil
	}
	if timezone != "" {
		if strings.HasPrefix(schedule, "CRON_TZ=") || strings.HasPrefix(schedule, "TZ=") {
			fields := strings.Fields(schedule)
			_, zone, _ := strings.Cut(fields[0], "=")
			if zone != timezone {
				return nil, fmt.Errorf("schedule timezone %q conflicts with backup timezone %q", zone, timezone)
			}
		} else {
			schedule = "CRON_TZ=" + timezone + " " + schedule
		}
	}
	parser := cron.NewParser(cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	parsed, err := parser.Parse(schedule)
	if err != nil {
		return nil, fmt.Errorf("invalid backup schedule: %w", err)
	}
	return parsed, nil
}
