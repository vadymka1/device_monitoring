package services

import (
	"context"
	"log"
	"sync"
	"time"

	"device-monitoring/internal/models"
)

type DeviceRepo interface {
	GetAll(ctx context.Context) ([]models.Device, error)
	SaveStatus(ctx context.Context, status *models.DeviceStatus) error
}

type DeviceMonitor interface {
	CheckDevice(ctx context.Context, d *models.Device) (*models.DeviceStatus, error)
}

type MonitorRunner struct {
	repo    DeviceRepo
	monitor DeviceMonitor
}

func NewMonitorRunner(repo DeviceRepo, monitor DeviceMonitor) *MonitorRunner {
	return &MonitorRunner{repo: repo, monitor: monitor}
}

func (r *MonitorRunner) Start(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			devices, err := r.repo.GetAll(ctx)
			if err != nil {
				log.Printf("failed to load devices: %v", err)
				continue
			}
			var wg sync.WaitGroup
			for _, d := range devices {
				wg.Add(1)
				go func(dev models.Device) {
					defer wg.Done()
					status, err := r.monitor.CheckDevice(ctx, &dev)
					if err != nil {
						log.Printf("monitor error for device %s: %v", dev.ID, err)
						return
					}
					log.Printf("✓ Monitored %s -> %s", dev.ID, status.Status)
				}(d)
			}
			wg.Wait()
		case <-ctx.Done():
			return
		}
	}
}
