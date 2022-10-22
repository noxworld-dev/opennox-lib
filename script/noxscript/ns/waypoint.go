package ns

type WaypointObj interface {
	Handle
}

type WaypointGroupObj interface {
	Handle
}

// Waypoint looks up waypoint by name.
func Waypoint(name string) WaypointObj {
	// header only
	return nil
}

// WaypointGroup looks up waypoint group by name.
func WaypointGroup(name string) WaypointGroupObj {
	// header only
	return nil
}

// IsWaypointOn gets whether waypoint is enabled.
func IsWaypointOn(wp WaypointObj) bool {
	// header only
	return false
}

// GetWaypointX gets waypoint X coordinate.
func GetWaypointX(wp WaypointObj) float32 {
	// header only
	return 0
}

// GetWaypointY gets waypoint Y coordinate.
func GetWaypointY(wp WaypointObj) float32 {
	// header only
	return 0
}

// MoveWaypoint sets waypoint location.
func MoveWaypoint(wp WaypointObj, x float32, y float32) {
	// header only
}

// WaypointOn enables a waypoint.
func WaypointOn(wp WaypointObj) {
	// header only
}

// WaypointGroupOn enables waypoints in a group.
func WaypointGroupOn(group WaypointGroupObj) {
	// header only
}

// WaypointOff disables a waypoint.
func WaypointOff(wp WaypointObj) {
	// header only
}

// WaypointGroupOff disables waypoints in a group.
func WaypointGroupOff(group WaypointGroupObj) {
	// header only
}

// WaypointToggle toggles waypoint between enabled and disabled.
func WaypointToggle(wp WaypointObj) {
	// header only
}

// WaypointGroupToggle toggles waypoints in group between enabled and disabled.
func WaypointGroupToggle(group WaypointGroupObj) {
	// header only
}
