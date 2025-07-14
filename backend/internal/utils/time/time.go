package time

import "time"

func DefaultDate() time.Time {
	return time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
}
