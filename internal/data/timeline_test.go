package data_test

import (
	data2 "github.com/danesparza/fxdmx/internal/data"
	"os"
	"testing"
)

func TestTimeline_AddTimeline_ValidTimeline_Successful(t *testing.T) {

	//	Arrange
	systemdb := getTestFiles()

	db, err := data2.NewManager(systemdb)
	if err != nil {
		t.Fatalf("NewManager failed: %s", err)
	}
	defer func() {
		db.Close()
		os.RemoveAll(systemdb)
	}()

	testTimelineFrames := []data2.TimelineFrame{
		{
			Type: "scene",
			Channels: []data2.ChannelValue{
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
			Channels: []data2.ChannelValue{
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
	newTimeline, err := db.AddTimeline("unittest_timeline1", "", testTimelineFrames)

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

func TestTimeline_AddTimeline_NoFrames_ReturnsError(t *testing.T) {

	//	Arrange
	systemdb := getTestFiles()

	db, err := data2.NewManager(systemdb)
	if err != nil {
		t.Fatalf("NewManager failed: %s", err)
	}
	defer func() {
		db.Close()
		os.RemoveAll(systemdb)
	}()

	testTimelineFrames := []data2.TimelineFrame{} // No items

	//	Act
	_, err = db.AddTimeline("unittest_timeline1", "", testTimelineFrames)

	//	Assert
	if err == nil {
		t.Errorf("AddTimeline - Should return error, but got none")
	}
}

func TestTimeline_GetTimeline_ValidTimeline_Successful(t *testing.T) {

	//	Arrange
	systemdb := getTestFiles()

	db, err := data2.NewManager(systemdb)
	if err != nil {
		t.Fatalf("NewManager failed: %s", err)
	}
	defer func() {
		db.Close()
		os.RemoveAll(systemdb)
	}()

	testTimeline1 := data2.Timeline{Name: "Timeline1", Frames: []data2.TimelineFrame{
		{
			Type: "scene",
			Channels: []data2.ChannelValue{
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
			Channels: []data2.ChannelValue{
				{Channel: 1, Value: 255},
				{Channel: 8, Value: 255},
			},
		},
		{
			Type:      "sleep",
			SleepTime: 10,
		},
	}}

	testTimeline2 := data2.Timeline{Name: "Timeline2", Frames: []data2.TimelineFrame{
		{
			Type: "fade",
			Channels: []data2.ChannelValue{
				{Channel: 1, Value: 255},
				{Channel: 8, Value: 255},
			},
		},
		{
			Type:      "sleep",
			SleepTime: 10,
		},
	}}

	testTimeline3 := data2.Timeline{Name: "Timeline3", Frames: []data2.TimelineFrame{
		{
			Type: "scene",
			Channels: []data2.ChannelValue{
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
			Channels: []data2.ChannelValue{
				{Channel: 1, Value: 255},
				{Channel: 8, Value: 255},
			},
		},
		{
			Type:      "sleep",
			SleepTime: 10,
		},
	}}

	//	Act
	db.AddTimeline(testTimeline1.Name, "", testTimeline1.Frames)
	newTimeline2, _ := db.AddTimeline(testTimeline2.Name, "", testTimeline2.Frames)
	db.AddTimeline(testTimeline3.Name, "", testTimeline3.Frames)

	gotTimeline, err := db.GetTimeline(newTimeline2.ID)

	//	Log the file details:
	t.Logf("Timeline: %+v", gotTimeline)

	//	Assert
	if err != nil {
		t.Errorf("GetTimeline - Should get timeline without error, but got: %s", err)
	}

	if gotTimeline.Name != newTimeline2.Name {
		t.Errorf("GetTimeline failed: Should get valid name but got: %v", gotTimeline.Name)
	}

	if gotTimeline.Frames[0].Type != testTimeline2.Frames[0].Type {
		t.Errorf("GetTimeline failed: Frames don't match what I expected: %+v", gotTimeline)
	}
}

func TestTimeline_GetAllTimelines_Successful(t *testing.T) {

	//	Arrange
	systemdb := getTestFiles()

	db, err := data2.NewManager(systemdb)
	if err != nil {
		t.Fatalf("NewManager failed: %s", err)
	}
	defer func() {
		db.Close()
		os.RemoveAll(systemdb)
	}()

	testTimeline1 := data2.Timeline{Name: "Timeline1", Frames: []data2.TimelineFrame{
		{
			Type: "scene",
			Channels: []data2.ChannelValue{
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
			Channels: []data2.ChannelValue{
				{Channel: 1, Value: 255},
				{Channel: 8, Value: 255},
			},
		},
		{
			Type:      "sleep",
			SleepTime: 10,
		},
	}}

	testTimeline2 := data2.Timeline{Name: "Timeline2", Frames: []data2.TimelineFrame{
		{
			Type: "fade",
			Channels: []data2.ChannelValue{
				{Channel: 1, Value: 255},
				{Channel: 8, Value: 255},
			},
		},
		{
			Type:      "sleep",
			SleepTime: 10,
		},
	}}

	testTimeline3 := data2.Timeline{Name: "Timeline3", Frames: []data2.TimelineFrame{
		{
			Type: "scene",
			Channels: []data2.ChannelValue{
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
			Channels: []data2.ChannelValue{
				{Channel: 1, Value: 255},
				{Channel: 8, Value: 255},
			},
		},
		{
			Type:      "sleep",
			SleepTime: 10,
		},
	}}

	//	Act
	db.AddTimeline(testTimeline1.Name, "", testTimeline1.Frames)
	newTimeline2, _ := db.AddTimeline(testTimeline2.Name, "", testTimeline2.Frames)
	db.AddTimeline(testTimeline3.Name, "", testTimeline3.Frames)

	gotTimelines, err := db.GetAllTimelines()

	//	Assert
	if err != nil {
		t.Errorf("GetAllTimelines - Should get all timelines without error, but got: %s", err)
	}

	if len(gotTimelines) < 2 {
		t.Errorf("GetAllTimelines failed: Should get all items but got: %v", len(gotTimelines))
	}

	if gotTimelines[1].Name != newTimeline2.Name {
		t.Errorf("GetAllTimelines failed: Should get an item with the correct details: %+v", gotTimelines[1])
	}
}

func TestTimeline_UpdateTimeline_ValidTimelines_Successful(t *testing.T) {

	//	Arrange
	systemdb := getTestFiles()

	db, err := data2.NewManager(systemdb)
	if err != nil {
		t.Fatalf("NewManager failed: %s", err)
	}
	defer func() {
		db.Close()
		os.RemoveAll(systemdb)
	}()

	testTimeline1 := data2.Timeline{Name: "Timeline1", Frames: []data2.TimelineFrame{
		{
			Type: "scene",
			Channels: []data2.ChannelValue{
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
			Channels: []data2.ChannelValue{
				{Channel: 1, Value: 255},
				{Channel: 8, Value: 255},
			},
		},
		{
			Type:      "sleep",
			SleepTime: 10,
		},
	}}

	testTimeline2 := data2.Timeline{Name: "Timeline2", Frames: []data2.TimelineFrame{
		{
			Type: "fade",
			Channels: []data2.ChannelValue{
				{Channel: 1, Value: 255},
				{Channel: 8, Value: 255},
			},
		},
		{
			Type:      "sleep",
			SleepTime: 10,
		},
	}}

	testTimeline3 := data2.Timeline{Name: "Timeline3", Frames: []data2.TimelineFrame{
		{
			Type: "scene",
			Channels: []data2.ChannelValue{
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
			Channels: []data2.ChannelValue{
				{Channel: 1, Value: 255},
				{Channel: 8, Value: 255},
			},
		},
		{
			Type:      "sleep",
			SleepTime: 10,
		},
	}}

	//	Act
	db.AddTimeline(testTimeline1.Name, "", testTimeline1.Frames)
	newTimeline2, _ := db.AddTimeline(testTimeline2.Name, "", testTimeline2.Frames)
	db.AddTimeline(testTimeline3.Name, "", testTimeline3.Frames)

	//	Update the 2nd trigger:
	newTimeline2.Enabled = false
	_, err = db.UpdateTimeline(newTimeline2) //	Update the 2nd timeline

	gotTimeline, _ := db.GetTimeline(newTimeline2.ID) // Refetch to verify

	//	Assert
	if err != nil {
		t.Errorf("UpdateTimeline - Should update timeline without error, but got: %s", err)
	}

	if gotTimeline.Enabled != false {
		t.Errorf("UpdateTimeline failed: Should get an item that has been disabled but got: %+v", gotTimeline)
	}

}

func TestTimeline_DeleteTimeline_ValidTimeline_Successful(t *testing.T) {

	//	Arrange
	systemdb := getTestFiles()

	db, err := data2.NewManager(systemdb)
	if err != nil {
		t.Fatalf("NewManager failed: %s", err)
	}
	defer func() {
		db.Close()
		os.RemoveAll(systemdb)
	}()

	testTimeline1 := data2.Timeline{Name: "Timeline1", Frames: []data2.TimelineFrame{
		{
			Type: "scene",
			Channels: []data2.ChannelValue{
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
			Channels: []data2.ChannelValue{
				{Channel: 1, Value: 255},
				{Channel: 8, Value: 255},
			},
		},
		{
			Type:      "sleep",
			SleepTime: 10,
		},
	}}

	testTimeline2 := data2.Timeline{Name: "Timeline2", Frames: []data2.TimelineFrame{
		{
			Type: "fade",
			Channels: []data2.ChannelValue{
				{Channel: 1, Value: 255},
				{Channel: 8, Value: 255},
			},
		},
		{
			Type:      "sleep",
			SleepTime: 10,
		},
	}}

	testTimeline3 := data2.Timeline{Name: "Timeline3", Frames: []data2.TimelineFrame{
		{
			Type: "scene",
			Channels: []data2.ChannelValue{
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
			Channels: []data2.ChannelValue{
				{Channel: 1, Value: 255},
				{Channel: 8, Value: 255},
			},
		},
		{
			Type:      "sleep",
			SleepTime: 10,
		},
	}}

	//	Act
	db.AddTimeline(testTimeline1.Name, "", testTimeline1.Frames)
	newTimeline2, _ := db.AddTimeline(testTimeline2.Name, "", testTimeline2.Frames)
	db.AddTimeline(testTimeline3.Name, "", testTimeline3.Frames)

	err = db.DeleteTimeline(newTimeline2.ID) //	Delete the 2nd timeline

	gotTimelines, _ := db.GetAllTimelines()

	//	Assert
	if err != nil {
		t.Errorf("DeleteTimeline - Should delete timeline without error, but got: %s", err)
	}

	if len(gotTimelines) != 2 {
		t.Errorf("DeleteTimeline failed: Should remove an item but got: %v", len(gotTimelines))
	}

	if gotTimelines[1].Name == newTimeline2.Name {
		t.Errorf("DeleteTimeline failed: Should get an item with different details than the removed item but got: %+v", gotTimelines[1])
	}

}
