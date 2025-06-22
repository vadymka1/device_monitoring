-- +goose Up
INSERT INTO devices (id, name, ip_address, protocol) VALUES
                                                         ('11111111-1111-1111-1111-111111111111', 'Router A', '192.168.0.101', 'http'),
                                                         ('22222222-2222-2222-2222-222222222222', 'Camera B', '192.168.0.102', 'http'),
                                                         ('33333333-3333-3333-3333-333333333333', 'Switch C', '192.168.0.103', 'grpc');

INSERT INTO device_statuses (id, device_id, status, hw_version, sw_version, fw_version, checksum, updated_at)
VALUES
  (gen_random_uuid(), '11111111-1111-1111-1111-111111111111', 'healthy', '1.0', '2.1', '3.3', 'abc123', NOW()),
  (gen_random_uuid(), '22222222-2222-2222-2222-222222222222', 'unreachable', '', '', '', '', NOW()),
  (gen_random_uuid(), '33333333-3333-3333-3333-333333333333', 'healthy', '1.0', '2.1', '3.3', 'abc123', NOW());



DELETE FROM device_statuses
WHERE device_id IN (
  '11111111-1111-1111-1111-111111111111',
  '22222222-2222-2222-2222-222222222222'
);

-- +goose Down
DELETE FROM devices WHERE id IN (
  '11111111-1111-1111-1111-111111111111',
  '22222222-2222-2222-2222-222222222222',
  '33333333-3333-3333-3333-333333333333'
);
