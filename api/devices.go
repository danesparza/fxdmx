package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/danesparza/fxdmx/system"
)

// GetSerialUSBDevices godoc
// @Summary Gets information about currently connected USB serial devices
// @Description Gets information about currently connected USB serial devices
// @Tags system
// @Accept  json
// @Produce  json
// @Success 200 {object} api.SystemResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /system/usbinfo [get]
func (service Service) GetSerialUSBDevices(rw http.ResponseWriter, req *http.Request) {

	//	req.Body is a ReadCloser -- we need to remember to close it:
	defer req.Body.Close()

	//	Get all the devices:
	devices, err := system.GetSerialUSBDeviceInfo()
	if err != nil {
		sendErrorResponse(rw, err, http.StatusInternalServerError)
		return
	}

	//	Create our response and send information back:
	response := SystemResponse{
		Message: fmt.Sprintf("%v devices found", len(devices)),
		Data:    devices,
	}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)
}
