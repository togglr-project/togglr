package productinfo

import (
	"context"
	"fmt"
	"time"

	"github.com/togglr-project/togglr/internal/contract"
	"github.com/togglr-project/togglr/internal/domain"
)

type ProductInfoService struct {
	productInfoRepo contract.ProductInfoRepository
}

func New(productInfoRepo contract.ProductInfoRepository) *ProductInfoService {
	return &ProductInfoService{
		productInfoRepo: productInfoRepo,
	}
}

func (s *ProductInfoService) GetProductInfo(ctx context.Context) (domain.ProductInfo, error) {
	clientID, err := s.productInfoRepo.GetClientID(ctx)
	if err != nil {
		return domain.ProductInfo{}, fmt.Errorf("failed to get client ID: %w", err)
	}

	// For now, we'll use current time as created_at since the repository doesn't return it
	// In a real implementation, you might want to add a method to get the creation time
	return domain.ProductInfo{
		ClientID:  clientID,
		CreatedAt: time.Now(), // This should ideally come from the database
	}, nil
}
