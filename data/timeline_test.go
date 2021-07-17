package data_test

import (
	"os"
	"testing"

	"github.com/danesparza/fxdmx/data"
)

func TestTimeline_AddTimeline_ValidTimeline_Successful(t *testing.T) {

	//	Arrange
	systemdb := getTestFiles()

	db, err := data.NewManager(systemdb)
	if err != nil {
		t.Fatalf("NewManager failed: %s", err)
	}
	defer func() {
		db.Close()
		os.RemoveAll(systemdb)
	}()

	testTimelineFrames := []data.TimeLineFrame{
		{
			Type: "scene",
			Channels: []data.ChannelValue{
				{Channel: 2, Value: 255},
				{Channel: 3, Value: 140},
				{Channel: 4, Value: 25},
				{Channel: 9, Value: 255},
				{Channel: 10, Value: 140},
				{Channel: 11, Value: 25},
			},
		},
		{
			Type: "fade",
			Channels: []data.ChannelValue{
				{Channel: 1, Value: 255},
				{Channel: 8, Value: 255},
			},
		},
		{
			Type:      "sleep",
			SleepTime: 10,
		},
	}

	//	Act
	newTimeline, err := db.AddTimeline("unittest_timeline1", testTimelineFrames)

	//	Assert
	if err != nil {
		t.Errorf("AddTimeline - Should add timeline without error, but got: %s", err)
	}

	if newTimeline.Created.IsZero() {
		t.Errorf("AddTimeline failed: Should have set an item with the correct datetime: %+v", newTimeline)
	}

	if newTimeline.Enabled != true {
		t.Errorf("AddTimeline failed: Should have enabled the timeline by default: %+v", newTimeline)
	}

	if newTimeline.Frames[0].Channels[0].Value != 255 {
		t.Errorf("AddTimeline failed: Should have added channels correctly: %+v", newTimeline)
	}

}
