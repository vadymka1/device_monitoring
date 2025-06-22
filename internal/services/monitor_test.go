package services

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"device-monitoring/internal/logger"
	"device-monitoring/internal/models"

	"github.com/stretchr/testify/assert"
)

type mockRepo struct {
	Saved   *models.DeviceStatus
	SaveErr error
}

func (m *mockRepo) SaveStatus(ctx context.Context, status *models.DeviceStatus) error {
	m.Saved = status
	return m.SaveErr
}

func (m *mockRepo) GetStatusByDeviceID(ctx context.Context, deviceID string) (*models.DeviceStatus, error) {
	return nil, nil
}

func TestCheckDevice_Healthy(t *testing.T) {
	health := models.HealthPayload{
		HWVersion: "v1.1",
		SWVersion: "v2.2",
		FWVersion: "v3.3",
		Checksum:  "abc123",
	}
	body, _ := json.Marshal(health)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}))
	defer ts.Close()

	device := &models.Device{
		ID:        "dev123",
		Name:      "Device A",
		IPAddress: ts.Listener.Addr().String(),
		Protocol:  "http",
	}

	repo := &mockRepo{}
	log := logger.New()
	service := NewMonitorService(repo, log)

	status, err := service.CheckDevice(context.Background(), device)
	assert.NoError(t, err)
	assert.Equal(t, "unhealthy", status.Status)
	assert.NotNil(t, repo.Saved)
}

func TestSaveStatusFailure(t *testing.T) {
	ctx := context.Background()
	repo := &mockRepo{SaveErr: errors.New("db failure")}
	log := logger.New()
	monitor := NewMonitorService(repo, log)

	status := &models.DeviceStatus{
		DeviceID: "test-device",
		Status:   "healthy",
	}

	err := monitor.repo.SaveStatus(ctx, status)
	assert.Error(t, err)
	assert.EqualError(t, err, "db failure")
}

func TestChecksumGeneration(t *testing.T) {
	status := &models.DeviceStatus{
		DeviceID:  "dev001",
		HWVersion: "HW-1",
		SWVersion: "SW-2",
		FWVersion: "FW-3",
	}
	checksum, err := generateChecksum(status)
	assert.NoError(t, err)
	assert.NotEmpty(t, checksum)
	assert.Len(t, checksum, 64)
}
