package framework

import "time"

func MostRecent(dates ...*time.Time) *time.Time {
	var mostRecent *time.Time

	for _, date := range dates {
		if date != nil {
			if mostRecent == nil {
				mostRecent = date
			} else if date.After(*mostRecent) {
				mostRecent = date
			}
		}
	}

	return mostRecent
}

func ToPTR[Source any](src Source) *Source {
	return &src
}
