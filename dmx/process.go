package dmx

import (
	"context"
	"fmt"
	"time"

	"github.com/danesparza/fxdmx/data"
)

// BackgroundProcess encapsulates background processing operations
type BackgroundProcess struct {
	DB         *data.Manager
	HistoryTTL time.Duration

	// PlayTimeline signals a timeline should be played
	PlayTimeline chan data.Timeline

	// StopTimeline signals a running timeline should be stopped
	StopTimeline chan string

	// StopAllTimelines signals all running timlines should be stopped
	StopAllTimelines chan bool
}

// HandleAndProcess handles system context calls and channel events to play/stop audio
func (bp BackgroundProcess) HandleAndProcess(systemctx context.Context) {

	//	Create a map of running timelines and their cancel functions
	// playingAudio := audioProcessMap{m: make(map[string]func())}

	//	Loop and respond to channels:
	for {
		select {
		case playReq := <-bp.PlayTimeline:
			//	As we get a request on a channel to play a file...
			//	Spawn a goroutine
			go func(cx context.Context, req data.Timeline) {
				/*
					//	Create a cancelable context from the passed (system) context
					ctx, cancel := context.WithCancel(cx)
					defer cancel()

					//	Add an entry to the map with
					//	- key: instance id
					//	- value: the cancel function (pointer)
					//	(critical section)
					playingAudio.rwMutex.Lock()
					playingAudio.m[req.ProcessID] = cancel
					playingAudio.rwMutex.Unlock()

					//	Create the command with context and play the audio
					playCommand := exec.CommandContext(ctx, "mpg123", playReq.FilePath)

					if err := playCommand.Run(); err != nil {
						//	Log an error playing a file
						fmt.Printf("error playing %v: %v", playReq.FilePath, err)
					}

					//	Remove ourselves from the map and exit (critical section)
					playingAudio.rwMutex.Lock()
					delete(playingAudio.m, req.ProcessID)
					playingAudio.rwMutex.Unlock()
				*/

			}(systemctx, playReq) // Launch the goroutine

		/*
			case stopTL := <-bp.StopTimeline:

					//	Look up the item in the map and call cancel if the item exists (critical section):
					playingAudio.rwMutex.Lock()
					playCancel, exists := playingAudio.m[stopFile]

					if exists {
						//	Call the context cancellation function
						playCancel()

						//	Remove ourselves from the map and exit
						delete(playingAudio.m, stopFile)
					}
					playingAudio.rwMutex.Unlock()
		*/
		case <-bp.StopAllTimelines:
			/*
				//	Loop through all items in the map and call cancel if the item exists (critical section):
				playingAudio.rwMutex.Lock()

				for stopFile, playCancel := range playingAudio.m {

					//	Call the cancel function
					playCancel()

					//	Remove ourselves from the map
					//	(this is safe to do in a 'range':
					//	https://golang.org/doc/effective_go#for )
					delete(playingAudio.m, stopFile)
				}

				playingAudio.rwMutex.Unlock()
			*/

		case <-systemctx.Done():
			fmt.Println("Stopping timeline processor")
			return
		}
	}
}
