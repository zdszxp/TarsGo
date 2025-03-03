// Package local provides a local clock
package local

import (
	gotime "time"

	"github.com/TarsCloud/TarsGo/tars/sync/time"
)

type Time struct{}

func (t *Time) Now() (gotime.Time, error) {
	return gotime.Now(), nil
}

func NewTime(opts ...time.Option) time.Time {
	return new(Time)
}
