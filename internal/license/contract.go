package license

import (
	"context"

	"github.com/rom8726/etoggle/internal/domain"
)

type LicensesRepository interface {
	GetLastByExpiresAt(ctx context.Context) (domain.License, error)
	GetByID(ctx context.Context, id string) (domain.License, error)
	Create(ctx context.Context, license domain.License) (domain.License, error)
}

type ProductInfoRepository interface {
	GetClientID(ctx context.Context) (string, error)
}

type ProjectsRepository interface {
	Count(ctx context.Context) (uint, error)
}
