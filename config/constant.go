package config

import "time"

const RedisTimeDuration time.Duration = 1000 * 1000 * 1000 * 60 * 5

// Format
const DateFormat string = "2006-01-02"
const DateTimeFormat string = "2006-01-02 15:04:05"

// Error
const Error string = "something wrong"
