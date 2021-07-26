package event

const (
	// SystemStartup event is when the system has started up
	SystemStartup = "System startup"

	// TimelineCreated event is when a timeline has been created
	TimelineCreated = "Timeline created"

	// TimelineUpdated event is when a timeline has been updated
	TimelineUpdated = "Timeline updated"

	// TimelineDeleted event is when a timeline has been removed
	TimelineDeleted = "Timeline deleted"

	// TimelineStarted event is when a timeline sequence has been started
	TimelineStarted = "Timeline started"

	// TimelineStopped event is when a timeline sequence has been stopped
	TimelineStopped = "Timeline stopped"

	// AllTimelinesStopped event is when all timelines have been stopped
	AllTimelinesStopped = "All Timelines stopped"

	// TimelineError event is when there was an error processing a timeline
	TimelineError = "Timeline error"

	// ConfigUpdated event is when the system configuration has been updated (specifically the default usb device has been udpated)
	ConfigUpdated = "Config updated"

	// SystemShutdown event is when the system is shutting down
	SystemShutdown = "System Shutdown"
)
