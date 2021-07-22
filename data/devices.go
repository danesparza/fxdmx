package data

// DeviceInfo contains information about a specific USB device.
type DeviceInfo struct {
	DevicePath   string `json:"device"`       // Unique Device path
	ProductName  string `json:"product"`      // Product name from udevadm info
	Manufacturer string `json:"manufacturer"` // Manufacturer name from udevadm info
}
