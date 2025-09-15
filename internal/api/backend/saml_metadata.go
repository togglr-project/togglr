package apibackend

import (
	"bytes"
	"context"

	"github.com/rom8726/etoggle/internal/domain"
	generatedapi "github.com/rom8726/etoggle/internal/generated/server"
)

func (r *RestAPI) GetSAMLMetadata(ctx context.Context) (generatedapi.GetSAMLMetadataRes, error) {
	metadata, err := r.usersUseCase.GetSSOMetadata(ctx, domain.SSOProviderNameADSaml)
	if err != nil {
		return nil, err
	}

	return &generatedapi.GetSAMLMetadataOK{
		Data: bytes.NewReader(metadata),
	}, nil
}
