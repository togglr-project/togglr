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
}

type SegmentsRepository interface {
	Create(ctx context.Context, segment domain.Segment) (domain.Segment, error)
	GetByID(ctx context.Context, id domain.SegmentID) (domain.Segment, error)
	ListByProjectID(ctx context.Context, projectID domain.ProjectID) ([]domain.Segment, error)
	Update(ctx context.Context, segment domain.Segment) (domain.Segment, error)
	Delete(ctx context.Context, id domain.SegmentID) error
}
