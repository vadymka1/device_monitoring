package services

import (
	"context"

	"device-monitoring/internal/logger"
	"device-monitoring/internal/models"
	"device-monitoring/internal/repositories"
)

type DeviceService struct {
	repo    repositories.DeviceRepository
	monitor *MonitorService
	logger  *logger.Logger
}

func NewDeviceService(repo repositories.DeviceRepository, monitor *MonitorService, logger *logger.Logger) *DeviceService {
	return &DeviceService{repo: repo, monitor: monitor, logger: logger}
}

func (s *DeviceService) GetDeviceByID(ctx context.Context, id string) (*models.Device, error) {
	s.logger.Info("Service: GetDeviceByID %s", id)
	return s.repo.GetByID(ctx, id)
}

func (s *DeviceService) GetStatusByDeviceID(ctx context.Context, deviceID string) (*models.DeviceStatus, error) {
	s.logger.Info("Service: GetStatusByDeviceID %s", deviceID)
	return s.repo.GetStatusByDeviceID(ctx, deviceID)
}

func (s *DeviceService) GetAllDevices(ctx context.Context) ([]models.Device, error) {
	s.logger.Info("Service: GetAllDevices")
	return s.repo.GetAllDevices(ctx)
}
