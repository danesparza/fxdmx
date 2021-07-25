package system

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"time"
)

// DeviceInfo contains information about a specific USB device.
type DeviceInfo struct {
	DevicePath   string `json:"device"`       // Unique Device path
	ProductName  string `json:"product"`      // Product name from udevadm info
	Manufacturer string `json:"manufacturer"` // Manufacturer name from udevadm info
}

// GetSerialUSBDeviceInfo gets a list of serial USB devices in the system
func GetSerialUSBDeviceInfo() ([]DeviceInfo, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	retval := []DeviceInfo{}

	//	First, find a list of active serial USB devices in the system
	matches, _ := filepath.Glob("/dev/ttyUSB*")

	//	For each device found...
	for _, match := range matches {
		//	Get information from udevadm..
		devinfoCmd := exec.CommandContext(ctx, "udevadm", "info", "--a", "--name", match)

		out, err := devinfoCmd.Output()
		if err != nil {
			//	If we get an error, return it:
			return retval, fmt.Errorf("error running udevadm to get information about USB serial devices: %v", err)
		}

		//	Find product
		reProd := regexp.MustCompile(`ATTRS{product}=="(?P<prod>.*)"`)
		prodStrings := reProd.FindStringSubmatch(string(out))
		prodString := "Not found"
		if len(prodStrings) > 1 {
			prodString = prodStrings[1]
		}

		//	Find manufacturer
		reMan := regexp.MustCompile(`ATTRS{manufacturer}=="(?P<man>.*)"`)
		manStrings := reMan.FindStringSubmatch(string(out))
		manString := "Not found"
		if len(manStrings) > 1 {
			manString = manStrings[1]
		}

		dev := DeviceInfo{
			DevicePath:   match,
			ProductName:  prodString,
			Manufacturer: manString,
		}

		retval = append(retval, dev)
	}

	//	Return what we found
	return retval, nil
}
