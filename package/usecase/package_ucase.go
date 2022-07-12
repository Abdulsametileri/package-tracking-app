package usecase

import (
	"context"
	"encoding/json"

	"github.com/Abdulsametileri/package-tracking-app/domain"
)

type packageUsecase struct {
	pc domain.PackageClient
}

func NewPackageUsecase(pClient domain.PackageClient) domain.PackageUsecase {
	return &packageUsecase{pc: pClient}
}

func (p *packageUsecase) TrackByVehicleID(ctx context.Context, id string) (*domain.Package, error) {
	bytes, err := p.pc.ConsumeByVehicleID(ctx, id)
	if err != nil {
		return nil, err
	}

	var res domain.Package
	err = json.Unmarshal(bytes, &res)
	return &res, err
}
