package domain

import (
	"time"
)

type TimelineEvent struct {
	Time    time.Time `json:"time"`
	Enabled bool      `json:"enabled"`
}
