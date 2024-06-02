package api

import (
	"encoding/json"
	"fmt"
	data2 "github.com/danesparza/fxdmx/internal/data"
	"github.com/danesparza/fxdmx/internal/dmx"
	"net/http"
	"time"
)

// Service encapsulates API service operations
type Service struct {
	DB         *data2.Manager
	StartTime  time.Time
	HistoryTTL time.Duration

	// PlayTimeline signals a timeline should be played
	PlayTimeline chan dmx.PlayTimelineRequest

	// StopTimeline signals a timeline should stop playing
	StopTimeline chan string

	//	StopAllTimelines signals all timelines should stop playing
	StopAllTimelines chan bool
}

// CreateTimelineRequest is a request to create a new timeline
type CreateTimelineRequest struct {
	Name          string                `json:"name"`    // The timeline name
	USBDevicePath string                `json:"devpath"` // The usb device path to use for this timeline
	Frames        []data2.TimelineFrame `json:"frames"`  // The frame sequence to progress through
}

// UpdateTimelineRequest is a request to update a timeline
type UpdateTimelineRequest struct {
	ID            string                `json:"id"`      // Unique Timeline ID
	Enabled       bool                  `json:"enabled"` // Timeline enabled or not
	Name          string                `json:"name"`    // The timeline name
	USBDevicePath string                `json:"devpath"` // The usb device path to use for this timeline
	Frames        []data2.TimelineFrame `json:"frames"`  // The frame sequence to progress through
}

// UpdateDefaultUSBRequest is a request to update the default USB device to use
type UpdateDefaultUSBRequest struct {
	DevicePath string `json:"devicepath"` // Unique USB device path
}

// SystemResponse is a response for a system request
type SystemResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// ErrorResponse represents an API response
type ErrorResponse struct {
	Message string `json:"message"`
}

// Used to send back an error:
func sendErrorResponse(rw http.ResponseWriter, err error, code int) {
	//	Our return value
	response := ErrorResponse{
		Message: "Error: " + err.Error()}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(code)
	json.NewEncoder(rw).Encode(response)
}

// ShowUI redirects to the /ui/ url path
func ShowUI(rw http.ResponseWriter, req *http.Request) {
	// http.Redirect(rw, req, "/ui/", 301)
	fmt.Fprintf(rw, "Hello, world - UI")
}

// GetIP gets a requests IP address by reading off the forwarded-for
// header (for proxies) and falls back to use the remote address.
func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}
