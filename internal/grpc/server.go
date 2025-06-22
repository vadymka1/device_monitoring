package grpc

import (
	"context"

	"device-monitoring/internal/logger"
	"device-monitoring/internal/repositories"
	"device-monitoring/internal/services"

	proto "device-monitoring/internal/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MonitoringServer struct {
	proto.UnimplementedMonitoringServer

	repo    *repositories.DeviceRepository
	monitor *services.MonitorService
	logger  *logger.Logger
}

func NewMonitoringServer(repo *repositories.DeviceRepository, monitor *services.MonitorService, log *logger.Logger) *MonitoringServer {
	return &MonitoringServer{
		repo:    repo,
		monitor: monitor,
		logger:  log,
	}
}

func (s *MonitoringServer) GetDevice(ctx context.Context, req *proto.DeviceRequest) (*proto.Device, error) {
	device, err := s.repo.GetByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "device not found: %v", err)
	}
	return &proto.Device{
		Id:        device.ID,
		Name:      device.Name,
		IpAddress: device.IPAddress,
		Protocol:  device.Protocol,
	}, nil
}

func (s *MonitoringServer) GetDeviceStatus(ctx context.Context, req *proto.DeviceRequest) (*proto.DeviceStatus, error) {
	status, err := s.repo.GetStatusByDeviceID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &proto.DeviceStatus{
		Id:        status.ID,
		DeviceId:  status.DeviceID,
		Status:    status.Status,
		HwVersion: status.HWVersion,
		SwVersion: status.SWVersion,
		FwVersion: status.FWVersion,
		Checksum:  status.Checksum,
	}, nil
}

func (s *MonitoringServer) CheckDeviceStatus(ctx context.Context, req *proto.DeviceRequest) (*proto.DeviceStatusResponse, error) {
	device, err := s.repo.GetByID(ctx, req.Id)
	if err != nil {
		s.logger.Error("Device not found: %v", err)
		return nil, status.Errorf(codes.NotFound, "device not found")
	}

	_, err = s.monitor.CheckDevice(ctx, device)
	if err != nil {
		s.logger.Error("CheckDevice failed: %v", err)
		return nil, status.Errorf(codes.Internal, "check failed")
	}

	statusData, err := s.repo.GetStatusByDeviceID(ctx, req.Id)
	if err != nil {
		s.logger.Error("Failed to get status after check: %v", err)
		return nil, status.Errorf(codes.Internal, "status fetch failed")
	}

	return &proto.DeviceStatusResponse{
		Status:    statusData.Status,
		HwVersion: statusData.HWVersion,
		SwVersion: statusData.SWVersion,
		FwVersion: statusData.FWVersion,
		Checksum:  statusData.Checksum,
	}, nil
}
