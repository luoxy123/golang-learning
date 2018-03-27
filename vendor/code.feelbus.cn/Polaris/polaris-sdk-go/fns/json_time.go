package fns

import (
	"fmt"
	"time"
)

const timeFormat = "20060102T150405Z"

type JsonTime time.Time

func (jt JsonTime) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("\"%s\"", time.Time(jt).Format(timeFormat))
	return []byte(stamp), nil
}
