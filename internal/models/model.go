package models

type Device struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	IPAddress string `json:"ip_address"`
	Protocol  string `json:"protocol"`
}

type DeviceStatus struct {
	ID        string `json:"id"`
	DeviceID  string `json:"device_id"`
	Status    string `json:"status"`
	HWVersion string `json:"hw_version"`
	SWVersion string `json:"sw_version"`
	FWVersion string `json:"fw_version"`
	Checksum  string `json:"checksum"`
}

type HealthPayload struct {
	HWVersion string `json:"hw_version"`
	SWVersion string `json:"sw_version"`
	FWVersion string `json:"fw_version"`
	Checksum  string `json:"checksum"`
}
