package domain

import (
	"time"
)

type SegmentID string

type Segment struct {
	ID          SegmentID
	ProjectID   ProjectID
	Name        string
	Description string
	Conditions  []Condition
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (id SegmentID) String() string {
	return string(id)
}
