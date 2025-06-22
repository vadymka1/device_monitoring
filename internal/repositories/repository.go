package repositories

import (
	"context"
	"database/sql"

	"device-monitoring/internal/logger"
	"device-monitoring/internal/models"
)

type DeviceRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewDeviceRepository(db *sql.DB, logger *logger.Logger) *DeviceRepository {
	return &DeviceRepository{db: db, logger: logger}
}

func (r *DeviceRepository) GetByID(ctx context.Context, id string) (*models.Device, error) {
	var device models.Device
	err := r.db.QueryRowContext(ctx, "SELECT id, name, ip_address, protocol FROM devices WHERE id=$1", id).
		Scan(&device.ID, &device.Name, &device.IPAddress, &device.Protocol)
	if err != nil {
		r.logger.Error("Error occurred: %v", err)
		return nil, err
	}
	r.logger.Info("Found device: %s", device.ID)
	return &device, nil
}

func (r *DeviceRepository) GetAll(ctx context.Context) ([]models.Device, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name, ip_address, protocol FROM devices")
	if err != nil {
		r.logger.Error("Error occurred: %v", err)
		return nil, err
	}
	defer rows.Close()

	var devices []models.Device
	for rows.Next() {
		var d models.Device
		if err := rows.Scan(&d.ID, &d.Name, &d.IPAddress, &d.Protocol); err != nil {
			r.logger.Error("Error occurred: %v", err)
			return nil, err
		}
		devices = append(devices, d)
	}
	return devices, nil
}

func (r *DeviceRepository) SaveStatus(ctx context.Context, status *models.DeviceStatus) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO device_statuses (id, device_id, status, hw_version, sw_version, fw_version, checksum, updated_at)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, NOW())`,
		status.DeviceID, status.Status, status.HWVersion, status.SWVersion, status.FWVersion, status.Checksum)
	return err
}

func (r *DeviceRepository) GetStatusByDeviceID(ctx context.Context, deviceID string) (*models.DeviceStatus, error) {
	var status models.DeviceStatus
	err := r.db.QueryRowContext(ctx, `
		SELECT id, device_id, status, hw_version, sw_version, fw_version, checksum 
		FROM device_statuses 
		WHERE device_id=$1 
		ORDER BY updated_at DESC 
		LIMIT 1`, deviceID).
		Scan(&status.ID, &status.DeviceID, &status.Status, &status.HWVersion, &status.SWVersion, &status.FWVersion, &status.Checksum)
	if err != nil {
		r.logger.Error("Error occurred: %v", err)
		return nil, err
	}
	return &status, nil
}

func (r *DeviceRepository) GetAllDevices(ctx context.Context) ([]models.Device, error) {
	r.logger.Info("Fetching all devices from DB")
	rows, err := r.db.QueryContext(ctx, "SELECT id, name, ip_address, protocol FROM devices")
	if err != nil {
		r.logger.Error("Query failed: %v", err)
		return nil, err
	}
	defer rows.Close()

	var devices []models.Device
	for rows.Next() {
		var d models.Device
		if err := rows.Scan(&d.ID, &d.Name, &d.IPAddress, &d.Protocol); err != nil {
			r.logger.Error("Row scan failed: %v", err)
			return nil, err
		}
		devices = append(devices, d)
	}
	return devices, nil
}
