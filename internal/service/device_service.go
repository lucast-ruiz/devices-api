package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lucast-ruiz/devices-api/internal/model"
)

type DeviceService struct {
    repo DeviceRepo
}


func NewDeviceService(r DeviceRepo) *DeviceService {
    return &DeviceService{repo: r}
}

func (s *DeviceService) Create(ctx context.Context, name, brand, state string) (*model.Device, error) {
	if name == "" || brand == "" {
		return nil, fmt.Errorf("name and brand are required")
	}
	if !model.IsValidState(state) {
		return nil, fmt.Errorf("invalid state value")
	}

	device := &model.Device{
		ID:        uuid.New().String(),
		Name:      name,
		Brand:     brand,
		State:     model.DeviceState(state),
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, device); err != nil {
		return nil, err
	}

	return device, nil
}

func (s *DeviceService) GetByID(ctx context.Context, id string) (*model.Device, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *DeviceService) ListAll(ctx context.Context, limit, offset int) ([]model.Device, error) {
	return s.repo.ListAll(ctx, limit, offset)
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

	// Valida estado, se enviado
	if state != nil {
		if !model.IsValidState(*state) {
			return nil, fmt.Errorf("invalid state value")
		}
		device.State = model.DeviceState(*state)
	}

	// Apply patch
	if name != nil {
		device.Name = *name
	}
	if brand != nil {
		device.Brand = *brand
	}

	// Save
	if err := s.repo.Update(ctx, device); err != nil {
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

type DeviceRepo interface {
    Create(ctx context.Context, d *model.Device) error
    GetByID(ctx context.Context, id string) (*model.Device, error)
    ListAll(ctx context.Context, limit, offset int) ([]model.Device, error)
    GetByBrand(ctx context.Context, brand string) ([]model.Device, error)
    GetByState(ctx context.Context, state string) ([]model.Device, error)
    Update(ctx context.Context, d *model.Device) error
    Delete(ctx context.Context, id string) error
}

