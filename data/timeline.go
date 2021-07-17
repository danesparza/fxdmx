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

// AddTimeline adds a timeline to the system
func (store Manager) AddTimeline(name string, frames []TimeLineFrame) (TimeLine, error) {

	//	Our return item
	retval := TimeLine{}

	//	If we don't have any frames, return an error
	if len(frames) < 1 {
		return retval, fmt.Errorf("frames must contain at least one item")
	}

	//	Create our new timeline
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

// UpdateTimeline updates a timeline in the system
func (store Manager) UpdateTimeline(updatedTimeline TimeLine) (TimeLine, error) {

	//	Our return item
	retval := TimeLine{}

	//	Serialize to JSON format
	encoded, err := json.Marshal(updatedTimeline)
	if err != nil {
		return retval, fmt.Errorf("problem serializing the data: %s", err)
	}

	//	Save it to the database:
	err = store.systemdb.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(GetKey("Timeline", updatedTimeline.ID), string(encoded), &buntdb.SetOptions{})
		return err
	})

	//	If there was an error saving the data, report it:
	if err != nil {
		return retval, fmt.Errorf("problem saving the timeline: %s", err)
	}

	//	Set our retval:
	retval = updatedTimeline

	//	Return our data:
	return retval, nil
}

// GetTimeline gets information about a single timeline in the system based on its id
func (store Manager) GetTimeline(id string) (TimeLine, error) {
	//	Our return item
	retval := TimeLine{}

	//	Find the item:
	err := store.systemdb.View(func(tx *buntdb.Tx) error {

		val, err := tx.Get(GetKey("Timeline", id))
		if err != nil {
			return err
		}

		if len(val) > 0 {
			//	Unmarshal data into our item
			if err := json.Unmarshal([]byte(val), &retval); err != nil {
				return err
			}
		}

		//	If we get to this point and there is no error...
		return nil
	})

	//	If there was an error, report it:
	if err != nil {
		return retval, fmt.Errorf("problem getting the timeline: %s", err)
	}

	//	Return our data:
	return retval, nil
}
