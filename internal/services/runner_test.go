package services

import (
	"context"
	"testing"
	"time"

	"device-monitoring/internal/models"

	"github.com/stretchr/testify/assert"
)

type MockRunnerRepo struct{}

func (m *MockRunnerRepo) GetAll(ctx context.Context) ([]models.Device, error) {
	return []models.Device{
		{ID: "1", Name: "router", IPAddress: "localhost", Protocol: "http"},
	}, nil
}

func (m *MockRunnerRepo) GetByID(ctx context.Context, id string) (*models.Device, error) {
	return nil, nil
}

func (m *MockRunnerRepo) SaveStatus(ctx context.Context, status *models.DeviceStatus) error {
	return nil
}

func (m *MockRunnerRepo) GetStatusByDeviceID(ctx context.Context, id string) (*models.DeviceStatus, error) {
	return nil, nil
}

type MockMonitor struct{}

func (m *MockMonitor) CheckDevice(ctx context.Context, d *models.Device) (*models.DeviceStatus, error) {
	return &models.DeviceStatus{DeviceID: d.ID, Status: "healthy"}, nil
}

func TestMonitorRunner(t *testing.T) {

	repo := &MockRunnerRepo{}
	monitor := &MockMonitor{}
	runner := NewMonitorRunner(repo, monitor)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	done := make(chan struct{})
	go func() {
		runner.Start(ctx, 500*time.Millisecond)
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("runner did not finish in time")
	}

	assert.True(t, true, "runner completed without panic")
}
