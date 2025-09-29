package domain

import (
	"time"
)

type FeatureTags struct {
	FeatureID FeatureID `db:"feature_id" pk:"true" editable:"true"`
	TagID     TagID     `db:"tag_id" pk:"true" editable:"true"`
	CreatedAt time.Time `db:"created_at"`
}
