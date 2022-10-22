package ns

import "github.com/noxworld-dev/opennox-lib/script/noxscript/ns/damage"

type builtinObj int

func (obj builtinObj) ScriptID() int {
	return int(obj)
}

type Obj interface {
	Handle
}
type ObjGroup interface {
	Handle
}

const (
	Self  = builtinObj(-2)
	Other = builtinObj(-1)
)

// Object looks up an object by name.
func Object(name string) Obj {
	// header only
	return nil
}

// ObjectGroup looks up object group by name.
func ObjectGroup(name string) ObjGroup {
	// header only
	return nil
}

// GetTrigger returns Self, if valid.
func GetTrigger() Obj {
	// header only
	return nil
}

// GetCaller returns Other, if valid.
func GetCaller() Obj {
	// header only
	return nil
}

// GetObjectX gets object X coordinate.
func GetObjectX(id Obj) float32 {
	// header only
	return 0
}

// GetObjectY gets object Y coordinate.
func GetObjectY(id Obj) float32 {
	// header only
	return 0
}

// GetObjectZ gets object Z coordinate.
func GetObjectZ(id Obj) float32 {
	// header only
	return 0
}

// IsTrigger checks whether object is Self.
func IsTrigger(id Obj) bool {
	// header only
	return false
}

// IsCaller checks whether object is Other.
func IsCaller(id Obj) bool {
	// header only
	return false
}

type Class = string

// HasClass checks whether object has a class.
func HasClass(id Obj, class Class) bool {
	// header only
	return false
}

type Subclass = string

// HasSubclass tests whether an item has a specific subclass.
// The subclass overlaps, so you should probably test for the class first (via HasClass).
func HasSubclass(id Obj, subclass Subclass) bool {
	// header only
	return false
}

type Direction = int

const (
	NW Direction = 0
	N  Direction = 1
	NE Direction = 2
	W  Direction = 3
	E  Direction = 5
	SW Direction = 6
	S  Direction = 7
	SE Direction = 8
)

// GetDirection gets object direction.
//
// See LookWithAngle.
func GetDirection(id Obj) Direction {
	// header only
	return 0
}

// Move an object to a waypoint. The object must be movable or attached to a "Mover".
//
// If the waypoint is linked, the object will continue to move once it reaches the first waypoint.
func Move(id Obj, wp WaypointObj) {
	// header only
}

// GroupMove moves the objects in a group to a waypoint. The objects must be movable or attached to a "Mover".
//
// If the waypoint is linked, the objects will continue to move once they reach the first waypoint.
func GroupMove(group ObjGroup, wp WaypointObj) {
	// header only
}

// LookAtDirection causes object to look in a direction.
func LookAtDirection(id Obj, dir Direction) {
	// header only
}

// GroupLookAtDirection causes objects in a group to look in a direction.
func GroupLookAtDirection(group ObjGroup, dir Direction) {
	// header only
}

// ObjectOn enable an object.
func ObjectOn(id Obj) {
	// header only
}

// ObjectGroupOn enables objects in a group.
func ObjectGroupOn(group ObjGroup) {
	// header only
}

// ObjectOff disables an object.
func ObjectOff(id Obj) {
	// header only
}

// ObjectGroupOff disables objects in a group.
func ObjectGroupOff(group ObjGroup) {
	// header only
}

// ObjectToggle toggles object between enabled and disabled.
func ObjectToggle(id Obj) {
	// header only
}

// ObjectGroupToggle toggles objects in a group between enabled and disabled.
func ObjectGroupToggle(group ObjGroup) {
	// header only
}

// Delete an object.
func Delete(id Obj) {
	// header only
}

// GroupDelete deletes objects in a group.
func GroupDelete(group ObjGroup) {
	// header only
}

// Wander causes an object to wander.
func Wander(id Obj) {
	// header only
}

// GroupWander cause objects in a group to wander.
func GroupWander(group ObjGroup) {
	// header only
}

// GoBackHome cause object to move to its starting location.
func GoBackHome(id Obj) {
	// header only
}

// IsObjectOn gets whether object is enabled.
func IsObjectOn(id Obj) bool {
	// header only
	return false
}

// UnlockDoor unlocks a door. It has no effect if the object is not a door.
func UnlockDoor(id Obj) {
	// header only
}

// LockDoor locks a door. It has no effect if the object is not a door.
func LockDoor(id Obj) {
	// header only
}

// IsLocked return whether an object is locked. It works with any kind of lock.
func IsLocked(id Obj) bool {
	// header only
	return false
}

// CreateObject creates an object given a type and a starting location.
//
//  Example:
//    spider := CreateObject("SmallAlbinoSpider", Waypoint("SpiderHole"))
func CreateObject(typ string, waypoint Obj) Obj {
	// header only
	return nil
}

type DamageType = damage.Type

// Damage the target with a given source object, amount, and damage type.
func Damage(target Obj, source Obj, amount int, typ DamageType) {
	// header only
}

// GroupDamage damages the target objects with a given source object, amount, and damage type.
func GroupDamage(targetGroup Obj, source Obj, amount int, typ DamageType) {
	// header only
}

// CreateMover creates a Mover for an object.
func CreateMover(id Obj, wp WaypointObj, speed float32) Obj {
	// header only
	return nil
}

// GroupCreateMover creates a Mover for every object in a group.
func GroupCreateMover(group ObjGroup, wp WaypointObj, speed float32) {
	// header only
}

// MoveObject sets an object location.
func MoveObject(id Obj, x float32, y float32) {
	// header only
}

// Raise sets an object's Z coordinate and then let the object fall down.
func Raise(id Obj, z float32) {
	// header only
}

// LookWithAngle sets an object's direction. The direction is an angle represented as an integer between 0 and 255.
// Due east is 0, and the angle increases as the object turns clock-wise.
func LookWithAngle(id Obj, angle int) {
	// header only
}

// PushObjectTo pushes an object to a location.
func PushObjectTo(id Obj, x float32, y float32) {
	// header only
}

// PushObject calculate a unit vector from the object's location to the specified
//  location, and multiply it by the specified magnitude. This vector will be
//  added to the object's location.
func PushObject(id Obj, magnitude float32, x float32, y float32) {
	// header only
}

// GetLastItem returns the object of the last item in the object's inventory. If the inventory is empty, it returns nil.
//
// This is used with GetPreviousItem to iterate through an object's inventory.
func GetLastItem(id Obj) Obj {
	// header only
	return nil
}

// GetPreviousItem returns the object of the previous item in the inventory from the given object.
// If the specified object is not in an inventory, or there are no more items in the inventory, it returns nil.
//
// This is used with GetLastItem to iterate through an object's inventory.
func GetPreviousItem(id Obj) Obj {
	// header only
	return nil
}

// HasItem gets whether the item is in the object's inventory.
func HasItem(holder Obj, item Obj) bool {
	// header only
	return false
}

// GetHolder returns the object that contains the item in its inventory.
func GetHolder(item Obj) Obj {
	// header only
	return nil
}

// Pickup cause object to pickup an item.
func Pickup(id Obj, item Obj) bool {
	// header only
	return false
}

// Drop cause object to drop an item.
func Drop(id Obj, item Obj) bool {
	// header only
	return false
}

// CurrentHealth gets object's health.
func CurrentHealth(id Obj) int {
	// header only
	return 0
}

// MaxHealth gets object's maximum health.
func MaxHealth(id Obj) int {
	// header only
	return 0
}

// RestoreHealth restores object's health.
func RestoreHealth(id Obj, amount int) {
	// header only
}

// IsVisibleTo first checks if the location of the objects are within 512 of each other coordinate-wise. It not, it returns false.
//
// It then checks whether the first object can see the second object.
func IsVisibleTo(object1 Obj, object2 Obj) bool {
	// header only
	return false
}

// LookAtObject sets direction of object so it is looking at another object.
func LookAtObject(id Obj, target Obj) {
	// header only
}

// Walk causes an object to walk to a location.
func Walk(id Obj, x float32, y float32) {
	// header only
}

// GroupWalk causes objects in a group to walk to a location.
func GroupWalk(group ObjGroup, x float32, y float32) {
	// header only
}

// Chat displays a localized string in a speech bubble.
//
// If the string is not in the string database, it will instead print an error message with "MISSING:".
func Chat(id Obj, message StringID) {
	// header only
}

// ChatTimerSeconds displays a localized string in a speech bubble for duration in seconds.
//
// If the string is not in the string database, it will instead print an error message with "MISSING:".
func ChatTimerSeconds(id Obj, message StringID, sec int) {
	// header only
}

// ChatTimer displays a localized string in a speech bubble for duration in frames.
//
// If the string is not in the string database, it will instead print an error message with "MISSING:".
func ChatTimer(id Obj, message StringID, frames int) {
	// header only
}

// DestroyChat destroys object's speech bubble.
func DestroyChat(id Obj) {
	// header only
}

// DestroyEveryChat destroys all speech bubbles.
func DestroyEveryChat() {
	// header only
}

// SetOwner makes an object the owner of the target. This will make the target
// friendly to the owner, and it will accredit the target's kills to the owner.
//
// For example, in a multiplayer map, you might have a switch that activates a
// hazard. You can use this so that if the hazard kills anyone, the player who
// activated the hazard gets the credit.
func SetOwner(owner Obj, target Obj) {
	// header only
}

// GroupSetOwner is the same as SetOwner but with an object group as the target.
func GroupSetOwner(owner Obj, targets ObjGroup) {
	// header only
}

// SetOwners is the same as SetOwner but with an object group as the owner.
func SetOwners(owners ObjGroup, target Obj) {
	// header only
}

// GroupSetOwners is the same as SetOwners but with an object group as the target.
func GroupSetOwners(owners ObjGroup, targets ObjGroup) {
	// header only
}

// IsOwnedBy gets whether target is owned by object.
func IsOwnedBy(id Obj, target Obj) bool {
	// header only
	return false
}

// GroupIsOwnedBy gets whether any object in target group is owned by object.
func GroupIsOwnedBy(id Obj, target ObjGroup) bool {
	// header only
	return false
}

// IsOwnedByAny gets whether target is owned by any object in the group.
func IsOwnedByAny(group ObjGroup, target Obj) bool {
	// header only
	return false
}

// GroupIsOwnedByAny gets whether any object in target is owned by any object in the group.
func GroupIsOwnedByAny(group ObjGroup, target ObjGroup) bool {
	// header only
	return false
}

// ClearOwner clears the owner of an object.
func ClearOwner(id Obj) {
	// header only
}

// GetElevatorStatus gets elevator status.
func GetElevatorStatus(id Obj) int {
	// header only
	return 0
}

// CreatureGuard causes a creature to move to a location, guard a nearby location,
// and attack any enemies that move within range of the guarded location.
func CreatureGuard(id Obj, x1 float32, y1 float32, x2 float32, y2 float32, distance float32) {
	// header only
}

// CreatureGroupGuard is the same as CreatureGuard but applies to creatures in a group.
func CreatureGroupGuard(group ObjGroup, x1 float32, y1 float32, x2 float32, y2 float32, distance float32) {
	// header only
}

// CreatureHunt causes creature to hunt.
func CreatureHunt(id Obj) {
	// header only
}

// CreatureGroupHunt causes creatures in a group to hunt.
func CreatureGroupHunt(group ObjGroup) {
	// header only
}

// CreatureIdle causes creature to idle.
func CreatureIdle(id Obj) {
	// header only
}

// CreatureGroupIdle causes creatures in a group to idle.
func CreatureGroupIdle(group ObjGroup) {
	// header only
}

// CreatureFollow causes a creature to follow target, and it won't attack anything unless disrupted or instructed to.
func CreatureFollow(id Obj, target Obj) {
	// header only
}

// CreatureGroupFollow causes the creatures to follow target, and they won't attack anything unless disrupted or instructed to.
func CreatureGroupFollow(group ObjGroup, target Obj) {
	// header only
}

// AggressionLevel sets a creature's aggression level. The most commonly used value is 0.83.
func AggressionLevel(id Obj, level float32) {
	// header only
}

// GroupAggressionLevel sets a group of creature's aggression level. The most commonly used value is 0.83.
func GroupAggressionLevel(group ObjGroup, level float32) {
	// header only
}

// HitLocation causes object to melee attacks a location.
func HitLocation(id Obj, x float32, y float32) {
	// header only
}

// GroupHitLocation causes objects in a group to melee attacks a location.
func GroupHitLocation(group ObjGroup, x float32, y float32) {
	// header only
}

// HitFarLocation causes object to ranged attacks a location.
func HitFarLocation(id Obj, x float32, y float32) {
	// header only
}

// GroupHitFarLocation causes objects in a group to ranged attacks a location.
func GroupHitFarLocation(group ObjGroup, x float32, y float32) {
	// header only
}

// SetRoamFlag sets roaming flags for object. Default is 0x80.
func SetRoamFlag(id Obj, flags int) {
	// header only
}

// GroupSetRoamFlag sets roaming flags for objects in a group. Default is 0x80.
func GroupSetRoamFlag(group ObjGroup, flags int) {
	// header only
}

// Attack a target.
func Attack(id Obj, target Obj) {
	// header only
}

// GroupAttack causes objects in a group to attack a target.
func GroupAttack(group ObjGroup, target int) {
	// header only
}

// RetreatLevel causes the creature to retreat if its health falls below the specified percentage (0.0 - 1.0).
func RetreatLevel(id Obj, percent float32) {
	// header only
}

// GroupRetreatLevel causes the creatures to retreat if its health falls below the specified percentage (0.0 - 1.0).
func GroupRetreatLevel(group ObjGroup, percent float32) {
	// header only
}

// ResumeLevel causes the creature to stop retreating if its health is above the specified percentage (0.0 - 1.0).
func ResumeLevel(id Obj, percent float32) {
	// header only
}

// GroupResumeLevel causes the creatures to stop retreating if its health is above the specified percentage (0.0 - 1.0).
func GroupResumeLevel(group ObjGroup, percent float32) {
	// header only
}

// RunAway causes creature to run away from target.
func RunAway(id Obj, target Obj, frames int) {
	// header only
}

// GroupRunAway causes creatures to run away from target.
func GroupRunAway(group ObjGroup, target int, frames int) {
	// header only
}

// PauseObject pauses an object temporarily.
func PauseObject(id Obj, frames int) {
	// header only
}

// GroupPauseObject pauses objects of a group temporarily.
func GroupPauseObject(group ObjGroup, frames int) {
	// header only
}

// IsAttackedBy gets whether object1 is being attacked by object2.
func IsAttackedBy(id1 Obj, id2 Obj) bool {
	// header only
	return false
}

// IsSummoned gets whether object is a summoned creature.
func IsSummoned(id Obj) bool {
	// header only
	return false
}

// ZombieStayDown sets zombie to stay down.
func ZombieStayDown(id Obj) {
	// header only
}

// ZombieGroupStayDown sets group of zombies to stay down.
func ZombieGroupStayDown(group ObjGroup) {
	// header only
}

// RaiseZombie raises a zombie. Also clears stay down state.
func RaiseZombie(id Obj) {
	// header only
}

// RaiseZombieGroup raises a zombie group. Also clears stay down state.
func RaiseZombieGroup(group ObjGroup) {
	// header only
}

// IsGameBall gets whether object is a GameBall.
func IsGameBall(id Obj) bool {
	// header only
	return false
}

// IsCrown gets whether object is a Crown.
func IsCrown(id Obj) bool {
	// header only
	return false
}

type ObjectEvent int

const (
	EventEnemySighted    = ObjectEvent(3)
	EventLookingForEnemy = ObjectEvent(4)
	EventDeath           = ObjectEvent(5)
	EventChangeFocus     = ObjectEvent(6)
	EventIsHit           = ObjectEvent(7)
	EventRetreat         = ObjectEvent(8)
	EventCollision       = ObjectEvent(9)
	EventEnemyHeard      = ObjectEvent(10)
	EventEndOfWaypoint   = ObjectEvent(11)
	EventLostEnemy       = ObjectEvent(13)
)

// SetCallback sets a function script to call for an event.
func SetCallback(id Obj, event ObjectEvent, fnc Func) {
	// header only
}

// DeleteObjectTimer delete object after a delay.
func DeleteObjectTimer(id Obj, frames int) {
	// header only
}

// MakeFriendly sets object friendly with host.
func MakeFriendly(id Obj) {
	// header only
}

// MakeEnemy unsets object as friendly.
func MakeEnemy(id Obj) {
	// header only
}

// BecomePet sets object as pet of host.
func BecomePet(id Obj) {
	// header only
}

// BecomeEnemy unsets object as pet of host.
func BecomeEnemy(id Obj) {
	// header only
}

// Frozen sets frozen status of an object.
func Frozen(id Obj, frozen bool) {
	// header only
}
