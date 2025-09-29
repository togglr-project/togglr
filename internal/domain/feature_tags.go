package domain

import (
	"time"
)

type FeatureTags struct {
	FeatureID FeatureID `db:"feature_id" editable:"true" pk:"true"`
	TagID     TagID     `db:"tag_id"     editable:"true" pk:"true"`
	CreatedAt time.Time `db:"created_at"`
}
