package ns

// GetQuestStatus gets quest status (int).
func GetQuestStatus(name string) int {
	// header only
	return 0
}

// GetQuestStatusFloat gets quest status (float).
func GetQuestStatusFloat(name string) float32 {
	// header only
	return 0
}

// SetQuestStatus sets quest status (int).
func SetQuestStatus(status int, name string) {
	// header only
}

// SetQuestStatusFloat sets quest status (float).
func SetQuestStatusFloat(status float32, name string) {
	// header only
}

// ResetQuestStatus deletes a quest status.
//
// The name can be a wildcard with an asterisk or with a map name.
func ResetQuestStatus(name string) {
	// header only
}

type EntryType int

const (
	EntryGreen1 = EntryType(0) // green text, no sound
	EntryWhite  = EntryType(1) // white text
	EntryRed    = EntryType(2) // red text with quest label
	EntryGreen2 = EntryType(3) // green text
	EntryGrey   = EntryType(4) // grey text with completed label
	EntryYellow = EntryType(8) // yellow text with hint label
)

// JournalEntry adds entry to player's journal.
//
// If the player object is nil, then it will add the journal entry to all players.
func JournalEntry(obj Obj, message StringID, typ EntryType) {
	// header only
}

// JournalEdit edits entry in player's journal.
//
// If the player object is nil, then it will edit the journal entry for all players.
func JournalEdit(obj Obj, message StringID, typ EntryType) {
	// header only
}

// JournalDelete deletes entry from player's journal.
//
// If the player object is nil, then it will delete the journal entry for all players.
func JournalDelete(obj Obj, message StringID) {
	// header only
}
