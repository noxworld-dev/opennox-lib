package ns

import "github.com/noxworld-dev/opennox-lib/script"

type WaypointObj interface {
	Handle
	script.Positionable
	script.Enabler
	script.Toggler
}

type WaypointGroupObj interface {
	Handle
	script.EnableSetter
	script.Toggler

	// EachWaypoint calls fnc for all waypoints in the group.
	// If fnc returns false, the iteration stops.
	// If recursive is true, iteration will include items from nested groups.
	EachWaypoint(recursive bool, fnc func(obj WaypointObj) bool)
}

// Waypoint looks up waypoint by name.
func Waypoint(name string) WaypointObj {
	if impl == nil {
		return nil
	}
	return impl.Waypoint(name)
}

// WaypointGroup looks up waypoint group by name.
func WaypointGroup(name string) WaypointGroupObj {
	if impl == nil {
		return nil
	}
	return impl.WaypointGroup(name)
}
