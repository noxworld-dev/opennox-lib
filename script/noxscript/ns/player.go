package ns

// GetHost gets host's player object.
func GetHost() Obj {
	// header only
	return nil
}

// GetCharacterData gets information about the loaded character.
func GetCharacterData(field int) int {
	// header only
	return 0
}

// StringID is an ID of a localized string in the string database.
type StringID = string

// Print displays a localized string on the screen of Other.
//
// If the string is not in the string database, it will instead print an error message with "MISSING:".
func Print(message StringID) {
	// header only
}

// PrintToAll displays a localized string to everyone.
//
// If the string is not in the string database, it will instead print an error message with "MISSING:".
func PrintToAll(message StringID) {
	// header only
}

// ClearMessages clears messages on player's screen.
func ClearMessages(player Obj) {
	// header only
}

// UnBlind the host.
func UnBlind() {
	// header only
}

// Blind the host.
func Blind() {
	// header only
}

// WideScreen enables or disables cinematic wide-screen effect.
func WideScreen(enable bool) {
	// header only
}

// GetGold gets amount of gold for player object.
func GetGold(player Obj) int {
	// header only
	return 0
}

// ChangeGold changes amount of gold for player object.
func ChangeGold(player Obj, delta int) {
	// header only
}

// GiveXp grants experience to a player.
func GiveXp(player Obj, xp float32) {
	// header only
}

// IsTalking gets whether host is talking.
func IsTalking() bool {
	// header only
	return false
}

type HalberdLevel int

const (
	OblivionHalberd   = HalberdLevel(0)
	OblivionHeart     = HalberdLevel(1)
	OblivionWierdling = HalberdLevel(2)
	OblivionOrb       = HalberdLevel(3)
)

// SetHalberd upgrades host's oblivion staff.
func SetHalberd(upgrade HalberdLevel) {
	// header only
}

// ImmediateBlind immediately blinds the host.
func ImmediateBlind() {
	// header only
}

type PlayerClass int

const (
	Warrior  = PlayerClass(0)
	Wizard   = PlayerClass(1)
	Conjurer = PlayerClass(2)
)

// EndGame end of game for a specific class.
func EndGame(class PlayerClass) {
	// header only
}

// GetScore gets player's score.
func GetScore(player Obj) int {
	// header only
	return 0
}

// ChangeScore changes player's score.
func ChangeScore(player Obj, score int) {
	// header only
}

// IsTrading returns whether the host is currently talking to shopkeeper.
func IsTrading() bool {
	// header only
	return false
}
