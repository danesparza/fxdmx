package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/danesparza/fxdmx/event"
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
