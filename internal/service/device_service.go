package service

import (
    "context"
    "time"
	"fmt"

    "github.com/google/uuid"
    "github.com/lucast-ruiz/devices-api/internal/model"
    "github.com/lucast-ruiz/devices-api/internal/repo"
)


type DeviceService struct {
	repo *repo.DeviceRepository
}

func NewDeviceService(r *repo.DeviceRepository) *DeviceService {
	return &DeviceService{repo: r}
}

func (s *DeviceService) Create(ctx context.Context, name, brand, state string) (*model.Device, error) {
	device := &model.Device{
		ID:        uuid.New().String(),
		Name:      name,
		Brand:     brand,
		State:     model.DeviceState(state),
		CreatedAt: time.Now(),
	}

	err := s.repo.Create(ctx, device)
	if err != nil {
		return nil, err
	}

	return device, nil
}

func (s *DeviceService) GetByID(ctx context.Context, id string) (*model.Device, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *DeviceService) ListAll(ctx context.Context) ([]model.Device, error) {
	return s.repo.ListAll(ctx)
}

func (s *DeviceService) GetByBrand(ctx context.Context, brand string) ([]model.Device, error) {
	return s.repo.GetByBrand(ctx, brand)
}

func (s *DeviceService) GetByState(ctx context.Context, state string) ([]model.Device, error) {
	return s.repo.GetByState(ctx, state)
}

func (s *DeviceService) Update(ctx context.Context, id string, name, brand, state *string) (*model.Device, error) {
	device, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if device == nil {
		return nil, nil
	}

	// Regra: não pode alterar name/brand se o device está "in-use"
	if device.State == model.StateInUse {
		if name != nil && *name != device.Name {
			return nil, fmt.Errorf("cannot change name when device is in-use")
		}
		if brand != nil && *brand != device.Brand {
			return nil, fmt.Errorf("cannot change brand when device is in-use")
		}
	}

	// Apply patch
	if name != nil {
		device.Name = *name
	}
	if brand != nil {
		device.Brand = *brand
	}
	if state != nil {
		device.State = model.DeviceState(*state)
	}

	// Save
	err = s.repo.Update(ctx, device)
	if err != nil {
		return nil, err
	}

	return device, nil
}

func (s *DeviceService) Delete(ctx context.Context, id string) error {
	device, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if device == nil {
		return nil
	}

	// Regra: não pode deletar se in-use
	if device.State == model.StateInUse {
		return fmt.Errorf("cannot delete device that is in-use")
	}

	return s.repo.Delete(ctx, id)
}
