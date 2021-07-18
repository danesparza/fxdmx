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

	// TimelineError event is when there was an error processing a timeline
	TimelineError = "Timeline error"

	// SystemShutdown event is when the system is shutting down
	SystemShutdown = "System Shutdown"
)
