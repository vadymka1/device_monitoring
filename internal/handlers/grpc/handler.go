package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	monitoringpb "device-monitoring/internal/proto"
	"device-monitoring/internal/services"
)

type MonitoringServer struct {
	service *services.DeviceService
	monitor *services.MonitorService
}

func NewMonitoringServer(service *services.DeviceService, monitor *services.MonitorService) *MonitoringServer {
	return &MonitoringServer{service: service, monitor: monitor}
}

func (h *MonitoringServer) CheckDeviceStatus(ctx context.Context, req *monitoringpb.DeviceRequest) (*monitoringpb.DeviceStatusResponse, error) {
	device, err := h.service.GetDeviceByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "device not found")
	}
	statusData, err := h.monitor.CheckDevice(ctx, device)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "check failed")
	}
	return &monitoringpb.DeviceStatusResponse{
		Status:    statusData.Status,
		HwVersion: statusData.HWVersion,
		SwVersion: statusData.SWVersion,
		FwVersion: statusData.FWVersion,
		Checksum:  statusData.Checksum,
	}, nil
}
