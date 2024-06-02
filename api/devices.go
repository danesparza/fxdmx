package api

import (
	"encoding/json"
	"fmt"
	"github.com/danesparza/fxdmx/internal/event"
	"github.com/danesparza/fxdmx/internal/system"
	"net/http"
	"strings"
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

// UpdateDefaultUSBDev godoc
// @Summary Update the default USB device
// @Description Update the default USB device
// @Tags system
// @Accept  json
// @Produce  json
// @Param timeline body api.UpdateDefaultUSBRequest true "The device path to use.  Example: /dev/ttyUSB0"
// @Success 200 {object} api.SystemResponse
// @Failure 400 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /system/defaultusb [put]
func (service Service) UpdateDefaultUSBDev(rw http.ResponseWriter, req *http.Request) {

	//	req.Body is a ReadCloser -- we need to remember to close it:
	defer req.Body.Close()

	//	Decode the request
	request := UpdateDefaultUSBRequest{}
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		sendErrorResponse(rw, err, http.StatusBadRequest)
		return
	}

	//	If we don't have the timeline.id, make sure we indicate that's not valid
	if strings.TrimSpace(request.DevicePath) == "" {
		sendErrorResponse(rw, fmt.Errorf("the device path is required -- it looks something like /dev/ttyUSB0"), http.StatusBadRequest)
		return
	}

	//	Update the timeline:
	updatedDefaultDev, err := service.DB.UpdateDefaultUSBDev(request.DevicePath)
	if err != nil {
		sendErrorResponse(rw, err, http.StatusInternalServerError)
		return
	}

	//	Record the event:
	service.DB.AddEvent(event.ConfigUpdated, fmt.Sprintf("%+v", request), GetIP(req), service.HistoryTTL)

	//	Create our response and send information back:
	response := SystemResponse{
		Message: "Default USB device updated",
		Data:    updatedDefaultDev,
	}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)
}

// GetDefaultUSBDev godoc
// @Summary Get the current default USB device
// @Description Get the current default USB device
// @Tags system
// @Accept  json
// @Produce  json
// @Success 200 {object} api.SystemResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /system/defaultusb [get]
func (service Service) GetDefaultUSBDev(rw http.ResponseWriter, req *http.Request) {

	//	req.Body is a ReadCloser -- we need to remember to close it:
	defer req.Body.Close()

	//	Get all the events:
	defaultDev, err := service.DB.GetDefaultUSBDev()
	if err != nil {
		sendErrorResponse(rw, err, http.StatusInternalServerError)
		return
	}

	//	Create our response and send information back:
	response := SystemResponse{
		Message: "Default USB device",
		Data:    defaultDev,
	}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)
}
