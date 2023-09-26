package utils

import (
	"os"
	"time"
)

func CurrentTimeWithLocalTZ() time.Time {
	loc, err := time.LoadLocation(os.Getenv("TZ"))
	if err != nil {
		panic(err)
	}
	return time.Now().In(loc)
}