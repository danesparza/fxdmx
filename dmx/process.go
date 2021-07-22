package dmx

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/danesparza/fxdmx/data"
)

type PlayTimelineRequest struct {
	ProcessID         string
	RequestedTimeline data.Timeline
}

// BackgroundProcess encapsulates background processing operations
type BackgroundProcess struct {
	DB         *data.Manager
	HistoryTTL time.Duration

	// PlayTimeline signals a timeline should be played
	PlayTimeline chan PlayTimelineRequest

	// StopTimeline signals a running timeline should be stopped
	StopTimeline chan string

	// StopAllTimelines signals all running timlines should be stopped
	StopAllTimelines chan bool
}

type timelineProcessMap struct {
	m       map[string]func()
	rwMutex sync.RWMutex
}

// HandleAndProcess handles system context calls and channel events to play/stop audio
func (bp BackgroundProcess) HandleAndProcess(systemctx context.Context) {

	//	Create a map of running timelines and their cancel functions
	playingTimelines := timelineProcessMap{m: make(map[string]func())}

	//	Loop and respond to channels:
	for {
		select {
		case playReq := <-bp.PlayTimeline:
			//	As we get a request on a channel to play a file...
			//	Spawn a goroutine
			go func(cx context.Context, req PlayTimelineRequest) {

				//	Create a cancelable context from the passed (system) context
				_, cancel := context.WithCancel(cx)
				defer cancel()

				//	Add an entry to the map with
				//	- key: instance id
				//	- value: the cancel function (pointer)
				//	(critical section)
				playingTimelines.rwMutex.Lock()
				playingTimelines.m[req.ProcessID] = cancel
				playingTimelines.rwMutex.Unlock()

				//	Process the timeline
				log.Printf("Processing timeline %v\n", req.ProcessID)

				//	Remove ourselves from the map and exit (critical section)
				playingTimelines.rwMutex.Lock()
				delete(playingTimelines.m, req.ProcessID)
				playingTimelines.rwMutex.Unlock()

			}(systemctx, playReq) // Launch the goroutine

		case stopTL := <-bp.StopTimeline:

			//	Look up the item in the map and call cancel if the item exists (critical section):
			playingTimelines.rwMutex.Lock()
			playCancel, exists := playingTimelines.m[stopTL]

			if exists {
				//	Call the context cancellation function
				playCancel()

				log.Printf("Stopped timeline process %v\n", stopTL)

				//	Remove ourselves from the map and exit
				delete(playingTimelines.m, stopTL)
			}
			playingTimelines.rwMutex.Unlock()

		case <-bp.StopAllTimelines:

			//	Loop through all items in the map and call cancel if the item exists (critical section):
			playingTimelines.rwMutex.Lock()

			log.Printf("Stopping all timeline processes\n")

			for stopTL, playCancel := range playingTimelines.m {

				//	Call the cancel function
				playCancel()

				//	Remove ourselves from the map
				//	(this is safe to do in a 'range':
				//	https://golang.org/doc/effective_go#for )
				delete(playingTimelines.m, stopTL)
			}

			playingTimelines.rwMutex.Unlock()

		case <-systemctx.Done():
			fmt.Println("Stopping timeline processor")
			return
		}
	}
}
