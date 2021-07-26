package data

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/rs/xid"
	"github.com/tidwall/buntdb"
)

type Timeline struct {
	ID            string          `json:"id"`                // Unique Timeline ID
	Enabled       bool            `json:"enabled"`           // Timeline enabled or not
	Created       time.Time       `json:"created"`           // Timeline create time
	Name          string          `json:"name"`              // Timeline name
	USBDevicePath string          `json:"devpath,omitempty"` // The USB device to play the timeline on.  Optional.  If not set, uses the default
	Frames        []TimelineFrame `json:"frames"`            // Frames for the timeline
}

type TimelineFrame struct {
	Type      string         `json:"type"`               // Timeline frame type (scene/sleep/fade) Fade 'fades' between the previous channel state and this frame
	Channels  []ChannelValue `json:"channels,omitempty"` // Channel information to set for the scene (optional) Required if type = scene or fade
	SleepTime int            `json:"sleeptime"`          // Sleep type in seconds (optional) Required if type = sleep
}

type ChannelValue struct {
	Channel int  `json:"channel"` // Unique Fixture ID
	Value   byte `json:"value"`   // Optional fixture name
}

// AddTimeline adds a timeline to the system
func (store Manager) AddTimeline(name string, frames []TimelineFrame) (Timeline, error) {

	//	Our return item
	retval := Timeline{}

	//	If we don't have any frames, return an error
	if len(frames) < 1 {
		return retval, fmt.Errorf("frames must contain at least one item")
	}

	//	Create our new timeline
	newTimeline := Timeline{
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
func (store Manager) UpdateTimeline(updatedTimeline Timeline) (Timeline, error) {

	//	Our return item
	retval := Timeline{}

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
func (store Manager) GetTimeline(id string) (Timeline, error) {
	//	Our return item
	retval := Timeline{}

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

// GetAllTimelines gets all timelines in the system
func (store Manager) GetAllTimelines() ([]Timeline, error) {
	//	Our return item
	retval := []Timeline{}

	//	Set our prefix
	prefix := GetKey("Timeline")

	//	Iterate over our values:
	err := store.systemdb.View(func(tx *buntdb.Tx) error {
		tx.Descend(prefix, func(key, val string) bool {

			if len(val) > 0 {
				//	Create our item:
				item := Timeline{}

				//	Unmarshal data into our item
				bval := []byte(val)
				if err := json.Unmarshal(bval, &item); err != nil {
					return false
				}

				//	Add to the array of returned users:
				retval = append(retval, item)
			}

			return true
		})
		return nil
	})

	//	If there was an error, report it:
	if err != nil {
		return retval, fmt.Errorf("problem getting the list of triggers: %s", err)
	}

	//	Return our data:
	return retval, nil
}

// DeleteTimeline deletes a timeline from the system
func (store Manager) DeleteTimeline(id string) error {

	//	Remove it from the database:
	err := store.systemdb.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete(GetKey("Timeline", id))
		return err
	})

	//	If there was an error removing the data, report it:
	if err != nil {
		return fmt.Errorf("problem removing the timeline: %s", err)
	}

	//	Return our data:
	return nil
}
