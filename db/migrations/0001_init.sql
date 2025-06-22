-- +goose Up
CREATE TABLE IF NOT EXISTS devices (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    ip_address TEXT NOT NULL,
    protocol TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS device_statuses (
    id UUID PRIMARY KEY,
    device_id UUID REFERENCES devices(id) ON DELETE CASCADE,
    status TEXT NOT NULL,
    hw_version TEXT,
    sw_version TEXT,
    fw_version TEXT,
    checksum TEXT,
    updated_at TIMESTAMP DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS device_statuses;
DROP TABLE IF EXISTS devices;
