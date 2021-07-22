package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/danesparza/fxdmx/dmx"
	"github.com/danesparza/fxdmx/event"
	"github.com/gorilla/mux"
	"github.com/rs/xid"
)

// ListAllTimelines godoc
// @Summary List all timelines in the system
// @Description List all timelines in the system
// @Tags timelines
// @Accept  json
// @Produce  json
// @Success 200 {object} api.SystemResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /timelines [get]
func (service Service) ListAllTimelines(rw http.ResponseWriter, req *http.Request) {

	//	Get a list of files
	retval, err := service.DB.GetAllTimelines()
	if err != nil {
		err = fmt.Errorf("error getting a list of timelines: %v", err)
		sendErrorResponse(rw, err, http.StatusInternalServerError)
		return
	}

	//	Construct our response
	response := SystemResponse{
		Message: fmt.Sprintf("%v timeline(s)", len(retval)),
		Data:    retval,
	}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)
}

// CreateTimeline godoc
// @Summary Create a new timeline
// @Description Create a new timeline
// @Tags timelines
// @Accept  json
// @Produce  json
// @Param timeline body api.CreateTimelineRequest true "The timeline to create"
// @Success 200 {object} api.SystemResponse
// @Failure 400 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /timelines [post]
func (service Service) CreateTimeline(rw http.ResponseWriter, req *http.Request) {

	//	req.Body is a ReadCloser -- we need to remember to close it:
	defer req.Body.Close()

	//	Decode the request
	request := CreateTimelineRequest{}
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		sendErrorResponse(rw, err, http.StatusBadRequest)
		return
	}

	//	If we don't have any webhooks associated, make sure we indicate that's not valid
	if len(request.Frames) < 1 {
		sendErrorResponse(rw, fmt.Errorf("at least one frame must be included"), http.StatusBadRequest)
		return
	}

	//	Create the new timeline:
	newTimeline, err := service.DB.AddTimeline(request.Name, request.Frames)
	if err != nil {
		sendErrorResponse(rw, err, http.StatusInternalServerError)
		return
	}

	//	Record the event:
	service.DB.AddEvent(event.TimelineCreated, fmt.Sprintf("%+v", request), GetIP(req), service.HistoryTTL)

	//	Create our response and send information back:
	response := SystemResponse{
		Message: "Timeline created",
		Data:    newTimeline,
	}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)
}

// UpdateTimeline godoc
// @Summary Update a timeline
// @Description Update a timeline
// @Tags timelines
// @Accept  json
// @Produce  json
// @Param timeline body api.UpdateTimelineRequest true "The timeline to update.  Must include timeline.id"
// @Success 200 {object} api.SystemResponse
// @Failure 400 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /timelines [put]
func (service Service) UpdateTimeline(rw http.ResponseWriter, req *http.Request) {

	//	req.Body is a ReadCloser -- we need to remember to close it:
	defer req.Body.Close()

	//	Decode the request
	request := UpdateTimelineRequest{}
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		sendErrorResponse(rw, err, http.StatusBadRequest)
		return
	}

	//	If we don't have the timeline.id, make sure we indicate that's not valid
	if strings.TrimSpace(request.ID) == "" {
		sendErrorResponse(rw, fmt.Errorf("the timeline.id is required"), http.StatusBadRequest)
		return
	}

	//	Make sure the id exists
	timeUpdate, _ := service.DB.GetTimeline(request.ID)
	if timeUpdate.ID != request.ID {
		sendErrorResponse(rw, fmt.Errorf("timeline must already exist"), http.StatusBadRequest)
		return
	}

	//	Only update the name if it's been passed
	if strings.TrimSpace(request.Name) != "" {
		timeUpdate.Name = request.Name
	}

	//	Enabled / disabled is always set
	timeUpdate.Enabled = request.Enabled

	//	Only update frames if we've passed some in
	if len(request.Frames) > 0 {
		timeUpdate.Frames = request.Frames
	}

	//	Update the timeline:
	updatedTimeline, err := service.DB.UpdateTimeline(timeUpdate)
	if err != nil {
		sendErrorResponse(rw, err, http.StatusInternalServerError)
		return
	}

	//	Record the event:
	service.DB.AddEvent(event.TimelineUpdated, fmt.Sprintf("%+v", request), GetIP(req), service.HistoryTTL)

	//	Create our response and send information back:
	response := SystemResponse{
		Message: "Timeline updated",
		Data:    updatedTimeline,
	}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)
}

// DeleteTimeline godoc
// @Summary Deletes a timeline in the system
// @Description Deletes a timeline in the system
// @Tags timelines
// @Accept  json
// @Produce  json
// @Param id path string true "The timeline id to delete"
// @Success 200 {object} api.SystemResponse
// @Failure 400 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Failure 503 {object} api.ErrorResponse
// @Router /timelines/{id} [delete]
func (service Service) DeleteTimeline(rw http.ResponseWriter, req *http.Request) {

	//	Get the id from the url (if it's blank, return an error)
	vars := mux.Vars(req)
	if vars["id"] == "" {
		err := fmt.Errorf("requires an id of a timeline to delete")
		sendErrorResponse(rw, err, http.StatusBadRequest)
		return
	}

	//	Delete the timeline
	err := service.DB.DeleteTimeline(vars["id"])
	if err != nil {
		err = fmt.Errorf("error deleting timeline: %v", err)
		sendErrorResponse(rw, err, http.StatusInternalServerError)
		return
	}

	//	Record the event:
	service.DB.AddEvent(event.TimelineDeleted, vars["id"], GetIP(req), service.HistoryTTL)

	//	Construct our response
	response := SystemResponse{
		Message: "Timeline deleted",
		Data:    vars["id"],
	}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)
}

// RequestTimelinePlay godoc
// @Summary Plays a timeline in the system
// @Description Plays a timeline in the system
// @Tags timelines
// @Accept  json
// @Produce  json
// @Param id path string true "The timeline id to play"
// @Success 200 {object} api.SystemResponse
// @Failure 400 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /timelines/play/{id} [post]
func (service Service) RequestTimelinePlay(rw http.ResponseWriter, req *http.Request) {

	//	Get the id from the url (if it's blank, return an error)
	vars := mux.Vars(req)
	if vars["id"] == "" {
		err := fmt.Errorf("requires an id of a timeline to play")
		sendErrorResponse(rw, err, http.StatusBadRequest)
		return
	}

	//	Get the timeline
	timeline, err := service.DB.GetTimeline(vars["id"])
	if err != nil {
		err = fmt.Errorf("error getting timeline: %v", err)
		sendErrorResponse(rw, err, http.StatusInternalServerError)
		return
	}

	//	Send to the channel:
	playRequest := dmx.PlayTimelineRequest{
		ProcessID:         xid.New().String(), // Generate a new id
		RequestedTimeline: timeline,
	}
	service.PlayTimeline <- playRequest

	//	Record the event:
	service.DB.AddEvent(event.TimelineStarted, fmt.Sprintf("Timeline ID: %s / Name: %s", timeline.ID, timeline.Name), GetIP(req), service.HistoryTTL)

	//	Construct our response
	response := SystemResponse{
		Message: "Timeline played",
		Data:    timeline,
	}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)
}

// RequestTimelineStop godoc
// @Summary Stops a specific timeline 'play' process
// @Description Stops a specific timeline 'play' process
// @Tags timelines
// @Accept  json
// @Produce  json
// @Param pid path string true "The process id to stop"
// @Success 200 {object} api.SystemResponse
// @Failure 400 {object} api.ErrorResponse
// @Router /timelines/stop/{pid} [post]
func (service Service) RequestTimelineStop(rw http.ResponseWriter, req *http.Request) {

	//	Get the id from the url (if it's blank, return an error)
	vars := mux.Vars(req)
	if vars["pid"] == "" {
		err := fmt.Errorf("requires a processid of a process to stop")
		sendErrorResponse(rw, err, http.StatusBadRequest)
		return
	}

	//	Send to the channel:
	service.StopTimeline <- vars["pid"]

	//	Record the event:
	service.DB.AddEvent(event.TimelineStopped, vars["pid"], GetIP(req), service.HistoryTTL)

	//	Create our response and send information back:
	response := SystemResponse{
		Message: "Timeline stopping",
		Data:    vars["pid"],
	}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)
}

// RequestAllTimelinesStop godoc
// @Summary Stops all timeline 'play' processes
// @Description Stops all timeline 'play' processes
// @Tags timelines
// @Accept  json
// @Produce  json
// @Success 200 {object} api.SystemResponse
// @Router /timelines/stop [post]
func (service Service) RequestAllTimelinesStop(rw http.ResponseWriter, req *http.Request) {

	//	Send to the channel:
	service.StopAllTimelines <- true

	//	Record the event:
	service.DB.AddEvent(event.AllTimelinesStopped, "Stop all timelines", GetIP(req), service.HistoryTTL)

	//	Create our response and send information back:
	response := SystemResponse{
		Message: "All Timelines stopping",
		Data:    ".",
	}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)
}
