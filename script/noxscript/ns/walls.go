package ns

import "github.com/noxworld-dev/opennox-lib/script"

// NoWallSound enables or disables wall sounds.
func NoWallSound(noWallSound bool) {
	if impl == nil {
		return
	}
	impl.NoWallSound(noWallSound)
}

// WallObj is a pointer to a wall.
type WallObj interface {
	Handle
	script.Enabler
	script.Toggler

	// Destroy breaks a wall.
	Destroy()
}

// Wall gets a wall by its grid coordinates.
func Wall(x int, y int) WallObj {
	if impl == nil {
		return nil
	}
	return impl.Wall(x, y)
}

// WallGroupObj is a group of walls.
type WallGroupObj interface {
	Handle
	script.EnableSetter
	script.Toggler

	// Destroy breaks walls in a group.
	Destroy()

	// EachWall calls fnc for all walls in the group.
	// If fnc returns false, the iteration stops.
	// If recursive is true, iteration will include items from nested groups.
	EachWall(recursive bool, fnc func(obj WallObj) bool)
}

// WallGroup lookups wall group by name.
func WallGroup(name string) WallGroupObj {
	if impl == nil {
		return nil
	}
	return impl.WallGroup(name)
}
