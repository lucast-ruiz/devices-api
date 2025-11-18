package repo

import (
	"context"
	"database/sql"

	"github.com/lucast-ruiz/devices-api/internal/model"
)

type DeviceRepository struct {
	db *sql.DB
}

func NewDeviceRepository(db *sql.DB) *DeviceRepository {
	return &DeviceRepository{db}
}

func (r *DeviceRepository) Create(ctx context.Context, d *model.Device) error {
	query := `
		INSERT INTO devices (id, name, brand, state, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.ExecContext(ctx, query, d.ID, d.Name, d.Brand, d.State, d.CreatedAt)
	return err
}

func (r *DeviceRepository) GetByID(ctx context.Context, id string) (*model.Device, error) {
	query := `
		SELECT id, name, brand, state, created_at
		FROM devices
		WHERE id = $1
	`

	row := r.db.QueryRowContext(ctx, query, id)

	var d model.Device
	err := row.Scan(&d.ID, &d.Name, &d.Brand, &d.State, &d.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &d, nil
}

func (r *DeviceRepository) GetByBrand(ctx context.Context, brand string) ([]model.Device, error) {
	query := `SELECT id, name, brand, state, created_at FROM devices WHERE brand = $1`

	rows, err := r.db.QueryContext(ctx, query, brand)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []model.Device
	for rows.Next() {
		var d model.Device
		if err := rows.Scan(&d.ID, &d.Name, &d.Brand, &d.State, &d.CreatedAt); err != nil {
			return nil, err
		}
		devices = append(devices, d)
	}

	return devices, nil
}

func (r *DeviceRepository) GetByState(ctx context.Context, state string) ([]model.Device, error) {
	query := `SELECT id, name, brand, state, created_at FROM devices WHERE state = $1`

	rows, err := r.db.QueryContext(ctx, query, state)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []model.Device
	for rows.Next() {
		var d model.Device
		if err := rows.Scan(&d.ID, &d.Name, &d.Brand, &d.State, &d.CreatedAt); err != nil {
			return nil, err
		}
		devices = append(devices, d)
	}

	return devices, nil
}

func (r *DeviceRepository) Update(ctx context.Context, d *model.Device) error {
	query := `
		UPDATE devices
		SET name = $1, brand = $2, state = $3
		WHERE id = $4
	`
	_, err := r.db.ExecContext(ctx, query, d.Name, d.Brand, d.State, d.ID)
	return err
}

func (r *DeviceRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM devices WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *DeviceRepository) ListAll(ctx context.Context) ([]model.Device, error) {
	query := `SELECT id, name, brand, state, created_at FROM devices`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []model.Device

	for rows.Next() {
		var d model.Device
		if err := rows.Scan(&d.ID, &d.Name, &d.Brand, &d.State, &d.CreatedAt); err != nil {
			return nil, err
		}
		devices = append(devices, d)
	}

	return devices, nil
}
