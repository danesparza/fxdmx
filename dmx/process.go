package dmx

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/akualab/dmx"
	"github.com/danesparza/fxdmx/data"
	"github.com/danesparza/fxdmx/event"
)

type PlayTimelineRequest struct {
	ProcessID         string
	RequestedTimeline data.Timeline
}

type timelineProcessMap struct {
	m       map[string]func()
	rwMutex sync.RWMutex
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

	// PlayingTimelines tracks currently playing timelines
	PlayingTimelines timelineProcessMap
}

// HandleAndProcess handles system context calls and channel events to play/stop audio
func (bp *BackgroundProcess) HandleAndProcess(systemctx context.Context) {

	//	Create a map of running timelines and their cancel functions
	bp.PlayingTimelines.m = make(map[string]func())

	//	Loop and respond to channels:
	for {
		select {
		case playReq := <-bp.PlayTimeline:
			//	As we get a request on a channel to play a file...
			//	Spawn a goroutine
			go bp.StartTimelinePlay(systemctx, playReq) // Launch the goroutine

		case stopTL := <-bp.StopTimeline:

			//	Look up the item in the map and call cancel if the item exists (critical section):
			bp.PlayingTimelines.rwMutex.Lock()
			playCancel, exists := bp.PlayingTimelines.m[stopTL]

			if exists {
				//	Call the context cancellation function
				playCancel()

				bp.DB.AddEvent(event.TimelineStopped, fmt.Sprintf("Stopped timeline process %v\n", stopTL), "", bp.HistoryTTL)

				//	Remove ourselves from the map and exit
				delete(bp.PlayingTimelines.m, stopTL)
			}
			bp.PlayingTimelines.rwMutex.Unlock()

		case <-bp.StopAllTimelines:

			//	Loop through all items in the map and call cancel if the item exists (critical section):
			bp.PlayingTimelines.rwMutex.Lock()

			bp.DB.AddEvent(event.AllTimelinesStopped, "Stopping all timeline processes", "", bp.HistoryTTL)

			for stopTL, playCancel := range bp.PlayingTimelines.m {

				//	Call the cancel function
				playCancel()

				//	Remove ourselves from the map
				//	(this is safe to do in a 'range':
				//	https://golang.org/doc/effective_go#for )
				delete(bp.PlayingTimelines.m, stopTL)
			}

			bp.PlayingTimelines.rwMutex.Unlock()

		case <-systemctx.Done():
			bp.DB.AddEvent(event.AllTimelinesStopped, "Stopping timeline processor", "", bp.HistoryTTL)
			return
		}
	}
}

// PlayTimeline plays a timeline
func (bp *BackgroundProcess) StartTimelinePlay(cx context.Context, req PlayTimelineRequest) {
	//	Create a cancelable context from the passed context
	ctx, cancel := context.WithCancel(cx)
	defer cancel()

	//	Add an entry to the map with
	//	- key: instance id
	//	- value: the cancel function (pointer)
	//	(critical section)
	bp.PlayingTimelines.rwMutex.Lock()
	bp.PlayingTimelines.m[req.ProcessID] = cancel
	bp.PlayingTimelines.rwMutex.Unlock()

	//	Process the timeline
	bp.DB.AddEvent(event.TimelineStarted, fmt.Sprintf("Processing timeline %v\n", req.ProcessID), "", bp.HistoryTTL)

	//	First, see if the timeline has a device set on it.
	if strings.TrimSpace(req.RequestedTimeline.USBDevicePath) == "" {
		defaultDevice, err := bp.DB.GetDefaultUSBDev()
		if err != nil {
			bp.DB.AddEvent(event.TimelineError, fmt.Sprintf("An error occurred trying to get the default USB device: %v", err), "", bp.HistoryTTL)
			return
		}

		//	If it doesn't, grab the default and use that.
		req.RequestedTimeline.USBDevicePath = defaultDevice
	}

	// Connect to the DMX controller.
	dmx, e := dmx.NewDMXConnection(req.RequestedTimeline.USBDevicePath)
	if e != nil {
		bp.DB.AddEvent(event.TimelineStarted, fmt.Sprintf("ERROR: Unable to connect to DMX512 interface %v: %v", req.RequestedTimeline.USBDevicePath, e), "", bp.HistoryTTL)
		return
	}
	defer dmx.Close()

	//	Keep a channel state map:
	channelState := map[int]byte{}

	//	Our waitgroup (for sync'ing fade finishes)
	var wg sync.WaitGroup

	//	Iterate through each frame
	for _, frame := range req.RequestedTimeline.Frames {

		select {
		default:

			//	Find out what type of frame this is, and act accordingly:
			switch strings.ToLower(frame.Type) {
			case "scene":
				//	Iterate through each of the channels and set them, then render
				for _, channel := range frame.Channels {
					//	Set dmx value for each channel:
					dmx.SetChannel(channel.Channel, channel.Value)

					//	Track chennel state:
					channelState[channel.Channel] = channel.Value
				}
				dmx.Render()

			case "fade":

				//	Iterate through each of the channels.
				for _, channel := range frame.Channels {

					wg.Add(1)

					//	Find the initial value, and pass that to the fade operation:
					//	(if we can't find it, assume it's 0 and pass that)
					ivalue, prs := channelState[channel.Channel]
					if !prs {
						ivalue = 0
					}

					go func(channelInfo data.ChannelValue, initialValue byte) {
						// Decrement the counter when the goroutine completes.
						defer wg.Done()

						//	Compare with the initial state, then render repeatedly
						//	toward the target value in a for loop (pay attention to the direction).
						if initialValue < channelInfo.Value {
							//	We need to fade up
							for i := initialValue; i < channelInfo.Value; i++ {
								select {
								case <-time.After(1 * time.Millisecond):
									dmx.SetChannel(channelInfo.Channel, byte(i))
									dmx.Render()
								case <-ctx.Done():
									return
								}
							}
						} else {
							//	We need to fade down
							for i := initialValue; i > 0; i-- {
								select {
								case <-time.After(1 * time.Millisecond):
									dmx.SetChannel(channelInfo.Channel, byte(i))
									dmx.Render()
								case <-ctx.Done():
									return
								}
							}
						}

					}(channel, ivalue)

					//	This means we'll need to put a mutex around them map. (for thread safe map interactions)

					//	Also:  Track the new value of the channel
					channelState[channel.Channel] = channel.Value

				}

				// Wait for all fades to complete.
				wg.Wait()

			case "sleep":
				//	Just sleep for the specified number of seconds
				select {
				case <-time.After(time.Duration(frame.SleepTime) * time.Second):
					continue
				case <-ctx.Done():
					return
				}
			}

		case <-ctx.Done():
			// stop
			return
		}
	}

	//	Remove ourselves from the map and exit (critical section)
	bp.PlayingTimelines.rwMutex.Lock()
	delete(bp.PlayingTimelines.m, req.ProcessID)
	bp.PlayingTimelines.rwMutex.Unlock()
}
