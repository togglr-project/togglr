package contract

import (
	"context"

	"github.com/rom8726/etoggle/internal/domain"
)

type SegmentsUseCase interface {
	Create(ctx context.Context, segment domain.Segment) (domain.Segment, error)
	GetByID(ctx context.Context, id domain.SegmentID) (domain.Segment, error)
	ListByProjectID(ctx context.Context, projectID domain.ProjectID) ([]domain.Segment, error)
	Update(ctx context.Context, segment domain.Segment) (domain.Segment, error)
	Delete(ctx context.Context, id domain.SegmentID) error
	// ListDesyncFeatureIDs returns feature IDs that have customized rules for the given segment.
	ListDesyncFeatureIDs(ctx context.Context, segmentID domain.SegmentID) ([]domain.FeatureID, error)
}

type SegmentsRepository interface {
	Create(ctx context.Context, segment domain.Segment) (domain.Segment, error)
	GetByID(ctx context.Context, id domain.SegmentID) (domain.Segment, error)
	ListByProjectID(ctx context.Context, projectID domain.ProjectID) ([]domain.Segment, error)
	Update(ctx context.Context, segment domain.Segment) (domain.Segment, error)
	Delete(ctx context.Context, id domain.SegmentID) error
}
