package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"device-monitoring/internal/logger"

	"device-monitoring/internal/models"
)

type DeviceStatusRepository interface {
	SaveStatus(ctx context.Context, status *models.DeviceStatus) error
	GetStatusByDeviceID(ctx context.Context, deviceID string) (*models.DeviceStatus, error)
}

type MonitorService struct {
	repo   DeviceStatusRepository
	logger *logger.Logger
}

func NewMonitorService(repo DeviceStatusRepository, logger *logger.Logger) *MonitorService {
	return &MonitorService{repo: repo, logger: logger}
}

func (s *MonitorService) CheckDevice(ctx context.Context, device *models.Device) (*models.DeviceStatus, error) {
	//address := fmt.Sprintf("%s:80", device.IPAddress)
	timeout := 3 * time.Second
	//conn, err := net.DialTimeout("tcp", address, timeout)
	//if err != nil {
	//	return s.saveStatus(ctx, device.ID, "unreachable", nil)
	//}
	//conn.Close()

	const maxRetries = 3
	var lastErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		s.logger.Info("Attempt %d to check device: %s", attempt, device.ID)

		healthURL := fmt.Sprintf("http://localhost:8080/health/%s", device.ID)
		req, err := http.NewRequestWithContext(ctx, "GET", healthURL, nil)
		if err != nil {
			s.logger.Warn("Request creation failed on attempt %d: %v", attempt, err)
			lastErr = err
			time.Sleep(time.Duration(attempt) * time.Second)
			continue
		}

		client := &http.Client{Timeout: timeout}
		resp, err := client.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			s.logger.Warn("Device check failed on attempt %d: %v", attempt, err)
			lastErr = err
			time.Sleep(time.Duration(attempt) * time.Second)
			continue
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			s.logger.Error("Reading response failed: %v", err)
			return s.saveStatus(ctx, device.ID, "error", nil)
		}

		var payload models.HealthPayload
		if err := json.Unmarshal(body, &payload); err != nil {
			s.logger.Error("Unmarshal health response failed: %v", err)
			return s.saveStatus(ctx, device.ID, "error", nil)
		}

		status := &models.DeviceStatus{
			DeviceID:  device.ID,
			Status:    "healthy",
			HWVersion: payload.HWVersion,
			SWVersion: payload.SWVersion,
			FWVersion: payload.FWVersion,
			Checksum:  payload.Checksum,
		}

		checksum, err := generateChecksum(status)
		if err == nil {
			payload.Checksum = checksum
			s.logger.Info("Checksum generated: %s", payload.Checksum)
		} else {
			s.logger.Error("Checksum generation failed: %v", err)
		}

		s.logger.Info("Device %s is healthy", device.ID)
		return s.saveStatus(ctx, device.ID, "healthy", status)
	}

	s.logger.Error("Device %s is unhealthy after %d attempts. Last error: %v", device.ID, maxRetries, lastErr)
	return s.saveStatus(ctx, device.ID, "unhealthy", nil)

}

func (s *MonitorService) saveStatus(ctx context.Context, deviceID, state string, data *models.DeviceStatus) (*models.DeviceStatus, error) {
	status := &models.DeviceStatus{
		DeviceID: deviceID,
		Status:   state,
	}
	if data != nil {
		status.HWVersion = data.HWVersion
		status.SWVersion = data.SWVersion
		status.FWVersion = data.FWVersion
		status.Checksum = data.Checksum
	}
	if err := s.repo.SaveStatus(ctx, status); err != nil {
		return nil, err
	}
	return status, nil
}

func generateChecksum(status *models.DeviceStatus) (string, error) {
	if status == nil {
		return "", fmt.Errorf("status is nil")
	}

	data := fmt.Sprintf("%s:%s:%s:%s", status.DeviceID, status.HWVersion, status.SWVersion, status.FWVersion)
	h := sha256.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil)), nil
}
