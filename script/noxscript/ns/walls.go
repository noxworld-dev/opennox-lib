package ns

// NoWallSound enables or disables wall sounds.
func NoWallSound(noWallSound bool) {
	// header only
}

// WallObj is a pointer to a wall.
type WallObj interface {
	Handle
}

// Wall gets a wall by its grid coordinates.
func Wall(x int, y int) WallObj {
	// header only
	return nil
}

// WallOpen opens a wall.
func WallOpen(wall WallObj) {
	// header only
}

// WallClose closes a wall.
func WallClose(wall WallObj) {
	// header only
}

// WallToggle toggles a wall between opened and closed.
func WallToggle(wall WallObj) {
	// header only
}

// WallBreak breaks a wall.
func WallBreak(wall WallObj) {
	// header only
}

// WallGroupObj is a group of walls.
type WallGroupObj interface {
	Handle
}

// WallGroup lookups wall group by name.
func WallGroup(name string) WallGroupObj {
	// header only
	return nil
}

// WallGroupOpen opens walls in a group.
func WallGroupOpen(wallGroup WallGroupObj) {
	// header only
}

// WallGroupClose closes walls in a group.
func WallGroupClose(wallGroup WallGroupObj) {
	// header only
}

// WallGroupToggle toggles walls in a group between opened and closed.
func WallGroupToggle(wallGroup WallGroupObj) {
	// header only
}

// WallGroupBreak breaks walls in a group.
func WallGroupBreak(wallGroup WallGroupObj) {
	// header only
}
