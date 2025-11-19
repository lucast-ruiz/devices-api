package service

import (
    "context"
    "testing"
    "time"

    "github.com/lucast-ruiz/devices-api/internal/model"
)


type mockRepo struct {
    CreateFn   func(ctx context.Context, d *model.Device) error
    GetByIDFn  func(ctx context.Context, id string) (*model.Device, error)
    UpdateFn   func(ctx context.Context, d *model.Device) error
    DeleteFn   func(ctx context.Context, id string) error
}

func (m *mockRepo) Create(ctx context.Context, d *model.Device) error {
    if m.CreateFn != nil {
        return m.CreateFn(ctx, d)
    }
    return nil
}

func (m *mockRepo) GetByID(ctx context.Context, id string) (*model.Device, error) {
    if m.GetByIDFn != nil {
        return m.GetByIDFn(ctx, id)
    }
    return nil, nil
}

func (m *mockRepo) ListAll(ctx context.Context, limit, offset int) ([]model.Device, error) {
    return nil, nil
}

func (m *mockRepo) GetByBrand(ctx context.Context, brand string) ([]model.Device, error) {
    return nil, nil
}

func (m *mockRepo) GetByState(ctx context.Context, state string) ([]model.Device, error) {
    return nil, nil
}

func (m *mockRepo) Update(ctx context.Context, d *model.Device) error {
    if m.UpdateFn != nil {
        return m.UpdateFn(ctx, d)
    }
    return nil
}

func (m *mockRepo) Delete(ctx context.Context, id string) error {
    if m.DeleteFn != nil {
        return m.DeleteFn(ctx, id)
    }
    return nil
}

//
// TESTES DAS REGRAS DE NEGÓCIO
//

func TestUpdate_CannotChangeNameWhenInUse(t *testing.T) {
    newName := "NewName"

    repo := &mockRepo{
        GetByIDFn: func(ctx context.Context, id string) (*model.Device, error) {
            return &model.Device{
                ID:        id,
                Name:      "Original",
                Brand:     "BrandX",
                State:     model.StateInUse,
                CreatedAt: time.Now(),
            }, nil
        },
    }

    svc := NewDeviceService(repo)

    _, err := svc.Update(context.Background(), "1", &newName, nil, nil)
    if err == nil {
        t.Fatalf("expected error, got nil")
    }
}

func TestUpdate_CannotChangeBrandWhenInUse(t *testing.T) {
    newBrand := "BrandY"

    repo := &mockRepo{
        GetByIDFn: func(ctx context.Context, id string) (*model.Device, error) {
            return &model.Device{
                ID:        id,
                Name:      "Device",
                Brand:     "BrandX",
                State:     model.StateInUse,
                CreatedAt: time.Now(),
            }, nil
        },
    }

    svc := NewDeviceService(repo)

    _, err := svc.Update(context.Background(), "1", nil, &newBrand, nil)
    if err == nil {
        t.Fatalf("expected error, got nil")
    }
}

func TestDelete_CannotDeleteInUse(t *testing.T) {
    repo := &mockRepo{
        GetByIDFn: func(ctx context.Context, id string) (*model.Device, error) {
            return &model.Device{
                ID:        id,
                Name:      "X",
                Brand:     "Y",
                State:     model.StateInUse,
                CreatedAt: time.Now(),
            }, nil
        },
    }

    svc := NewDeviceService(repo)

    err := svc.Delete(context.Background(), "1")
    if err == nil {
        t.Fatalf("expected error on deleting in-use device")
    }
}

func TestUpdate_InvalidState(t *testing.T) {
    invalid := "broken"

    repo := &mockRepo{
        GetByIDFn: func(ctx context.Context, id string) (*model.Device, error) {
            return &model.Device{
                ID:        id,
                Name:      "A",
                Brand:     "B",
                State:     model.StateAvailable,
                CreatedAt: time.Now(),
            }, nil
        },
    }

    svc := NewDeviceService(repo)

    _, err := svc.Update(context.Background(), "1", nil, nil, &invalid)
    if err == nil {
        t.Fatalf("expected error for invalid state")
    }
}

func TestUpdate_CreatedAtIsNotModified(t *testing.T) {
    created := time.Now().Add(-1 * time.Hour)

    saved := &model.Device{}

    repo := &mockRepo{
        GetByIDFn: func(ctx context.Context, id string) (*model.Device, error) {
            return &model.Device{
                ID:        id,
                Name:      "A",
                Brand:     "B",
                State:     model.StateAvailable,
                CreatedAt: created,
            }, nil
        },
        UpdateFn: func(ctx context.Context, d *model.Device) error {
            saved = d
            return nil
        },
    }

    svc := NewDeviceService(repo)

    newName := "Updated"
    _, err := svc.Update(context.Background(), "1", &newName, nil, nil)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if !saved.CreatedAt.Equal(created) {
        t.Fatalf("expected created_at to remain unchanged")
    }
}

func TestUpdate_DeviceNotFound(t *testing.T) {
    repo := &mockRepo{
        GetByIDFn: func(ctx context.Context, id string) (*model.Device, error) {
            return nil, nil
        },
    }

    svc := NewDeviceService(repo)

    dev, err := svc.Update(context.Background(), "1", nil, nil, nil)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if dev != nil {
        t.Fatalf("expected nil device when not found")
    }
}

func TestCreate_RequiresNameAndBrand(t *testing.T) {
    repo := &mockRepo{}
    svc := NewDeviceService(repo)

    _, err := svc.Create(context.Background(), "", "Brand", "available")
    if err == nil {
        t.Fatalf("expected error for missing name")
    }

    _, err = svc.Create(context.Background(), "Device", "", "available")
    if err == nil {
        t.Fatalf("expected error for missing brand")
    }
}
