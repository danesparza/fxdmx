package data

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/rs/xid"
	"github.com/tidwall/buntdb"
)

type TimeLine struct {
	ID      string          `json:"id"`      // Unique Timeline ID
	Enabled bool            `json:"enabled"` // Timeline enabled or not
	Created time.Time       `json:"created"` // Timeline create time
	Name    string          `json:"name"`    // Scene name
	Frames  []TimeLineFrame `json:"frames"`  // Frames for the timeline
}

type TimeLineFrame struct {
	Type      string         `json:"type"`      // Timeline frame type (scene/sleep/fade) Fade 'fades' between the previous channel state and this frame
	Channels  []ChannelValue `json:"channels"`  // Channel information to set for the scene (optional) Required if type = scene or fade
	SleepTime int            `json:"sleeptime"` // Sleep type in seconds (optional) Required if type = sleep
}

type ChannelValue struct {
	Channel int  `json:"channel"` // Unique Fixture ID
	Value   byte `json:"value"`   // Optional fixture name
}

// AddTrigger adds a trigger to the system
func (store Manager) AddTimeline(name string, frames []TimeLineFrame) (TimeLine, error) {

	//	Our return item
	retval := TimeLine{}

	newTimeline := TimeLine{
		ID:      xid.New().String(), // Generate a new id
		Created: time.Now(),
		Enabled: true,
		Name:    name,
		Frames:  frames,
	}

	//	Serialize to JSON format
	encoded, err := json.Marshal(newTimeline)
	if err != nil {
		return retval, fmt.Errorf("problem serializing the data: %s", err)
	}

	//	Save it to the database:
	err = store.systemdb.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(GetKey("Timeline", newTimeline.ID), string(encoded), &buntdb.SetOptions{})
		return err
	})

	//	If there was an error saving the data, report it:
	if err != nil {
		return retval, fmt.Errorf("problem saving the timeline: %s", err)
	}

	//	Set our retval:
	retval = newTimeline

	//	Return our data:
	return retval, nil
}
