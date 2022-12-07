package noxast

import (
	"fmt"
	"go/ast"
	"go/token"
	"reflect"
	"strings"

	"github.com/noxworld-dev/opennox-lib/object"
	"github.com/noxworld-dev/opennox-lib/script"
	asm "github.com/noxworld-dev/opennox-lib/script/noxscript/noxasm"
	"github.com/noxworld-dev/opennox-lib/script/noxscript/ns"
	"github.com/noxworld-dev/opennox-lib/script/noxscript/ns/class"
	"github.com/noxworld-dev/opennox-lib/script/noxscript/ns/damage"
	"github.com/noxworld-dev/opennox-lib/script/noxscript/ns/effect"
	"github.com/noxworld-dev/opennox-lib/script/noxscript/ns/enchant"
	"github.com/noxworld-dev/opennox-lib/script/noxscript/ns/spell"
	"github.com/noxworld-dev/opennox-lib/script/noxscript/ns/subclass"
)

type builtinDef struct {
	Name string
	Type reflect.Type
}

var builtins = []*builtinDef{
	asm.BuiltinWall:                {Name: "Wall", Type: reflect.TypeOf(ns.Wall)},
	asm.BuiltinUnused1f:            {Name: "Unused1f", Type: reflect.TypeOf(ns.Unused1f)},
	asm.BuiltinUnused20:            {Name: "Unused20", Type: reflect.TypeOf(ns.Unused20)},
	asm.BuiltinAudioEvent:          {Name: "AudioEvent", Type: reflect.TypeOf(ns.AudioEvent)},
	asm.BuiltinPrint:               {Name: "Print", Type: reflect.TypeOf(ns.Print)},
	asm.BuiltinPrintToAll:          {Name: "PrintToAll", Type: reflect.TypeOf(ns.PrintToAll)},
	asm.BuiltinStopScript:          {Name: "StopScript", Type: reflect.TypeOf(ns.StopScript)},
	asm.BuiltinRandomFloat:         {Name: "RandomFloat", Type: reflect.TypeOf(ns.RandomFloat)},
	asm.BuiltinRandom:              {Name: "Random", Type: reflect.TypeOf(ns.Random)},
	asm.BuiltinCreateObject:        {Name: "CreateObject", Type: reflect.TypeOf(ns.CreateObject)},
	asm.BuiltinGetHost:             {Name: "GetHost", Type: reflect.TypeOf(ns.GetHost)},
	asm.BuiltinObject:              {Name: "Object", Type: reflect.TypeOf(ns.Object)},
	asm.BuiltinUnused50:            {Name: "Unused50", Type: reflect.TypeOf(ns.Unused50)},
	asm.BuiltinUnused58:            {Name: "Unused58", Type: reflect.TypeOf(ns.Unused58)},
	asm.BuiltinUnused59:            {Name: "Unused59", Type: reflect.TypeOf(ns.Unused59)},
	asm.BuiltinUnused5a:            {Name: "Unused5a", Type: reflect.TypeOf(ns.Unused5a)},
	asm.BuiltinUnused5b:            {Name: "Unused5b", Type: reflect.TypeOf(ns.Unused5b)},
	asm.BuiltinUnused5c:            {Name: "Unused5c", Type: reflect.TypeOf(ns.Unused5c)},
	asm.BuiltinUnused5d:            {Name: "Unused5d", Type: reflect.TypeOf(ns.Unused5d)},
	asm.BuiltinUnused5e:            {Name: "Unused5e", Type: reflect.TypeOf(ns.Unused5e)},
	asm.BuiltinGetCharacterData:    {Name: "GetCharacterData", Type: reflect.TypeOf(ns.GetCharacterData)},
	asm.BuiltinEffect:              {Name: "Effect", Type: reflect.TypeOf(ns.Effect)},
	asm.BuiltinWaypoint:            {Name: "Waypoint", Type: reflect.TypeOf(ns.Waypoint)},
	asm.BuiltinWaypointGroup:       {Name: "WaypointGroup", Type: reflect.TypeOf(ns.WaypointGroup)},
	asm.BuiltinObjectGroup:         {Name: "ObjectGroup", Type: reflect.TypeOf(ns.ObjectGroup)},
	asm.BuiltinWallGroup:           {Name: "WallGroup", Type: reflect.TypeOf(ns.WallGroup)},
	asm.BuiltinUnused74:            {Name: "Unused74", Type: reflect.TypeOf(ns.Unused74)},
	asm.BuiltinDestroyEveryChat:    {Name: "DestroyEveryChat", Type: reflect.TypeOf(ns.DestroyEveryChat)},
	asm.BuiltinSetQuestStatus:      {Name: "SetQuestStatus", Type: reflect.TypeOf(ns.SetQuestStatus)},
	asm.BuiltinSetQuestStatusFloat: {Name: "SetQuestStatusFloat", Type: reflect.TypeOf(ns.SetQuestStatusFloat)},
	asm.BuiltinGetQuestStatus:      {Name: "GetQuestStatus", Type: reflect.TypeOf(ns.GetQuestStatus)},
	asm.BuiltinGetQuestStatusFloat: {Name: "GetQuestStatusFloat", Type: reflect.TypeOf(ns.GetQuestStatusFloat)},
	asm.BuiltinResetQuestStatus:    {Name: "ResetQuestStatus", Type: reflect.TypeOf(ns.ResetQuestStatus)},
	asm.BuiltinIsTrigger:           {Name: "IsTrigger", Type: reflect.TypeOf(ns.IsTrigger)},
	asm.BuiltinIsCaller:            {Name: "IsCaller", Type: reflect.TypeOf(ns.IsCaller)},
	asm.BuiltinSetDialog:           {Name: "SetDialog", Type: reflect.TypeOf(ns.SetDialog)},
	asm.BuiltinCancelDialog:        {Name: "CancelDialog", Type: reflect.TypeOf(ns.CancelDialog)},
	asm.BuiltinStoryPic:            {Name: "StoryPic", Type: reflect.TypeOf(ns.StoryPic)},
	asm.BuiltinTellStory:           {Name: "TellStory", Type: reflect.TypeOf(ns.TellStory)},
	asm.BuiltinStartDialog:         {Name: "StartDialog", Type: reflect.TypeOf(ns.StartDialog)},
	asm.BuiltinUnBlind:             {Name: "UnBlind", Type: reflect.TypeOf(ns.UnBlind)},
	asm.BuiltinBlind:               {Name: "Blind", Type: reflect.TypeOf(ns.Blind)},
	asm.BuiltinWideScreen:          {Name: "WideScreen", Type: reflect.TypeOf(ns.WideScreen)},
	asm.BuiltinJournalEntry:        {Name: "JournalEntry", Type: reflect.TypeOf(ns.JournalEntry)},
	asm.BuiltinJournalDelete:       {Name: "JournalDelete", Type: reflect.TypeOf(ns.JournalDelete)},
	asm.BuiltinJournalEdit:         {Name: "JournalEdit", Type: reflect.TypeOf(ns.JournalEdit)},
	asm.BuiltinGetAnswer:           {Name: "GetAnswer", Type: reflect.TypeOf(ns.GetAnswer)},
	asm.BuiltinAutoSave:            {Name: "AutoSave", Type: reflect.TypeOf(ns.AutoSave)},
	asm.BuiltinMusic:               {Name: "Music", Type: reflect.TypeOf(ns.Music)},
	asm.BuiltinStartupScreen:       {Name: "StartupScreen", Type: reflect.TypeOf(ns.StartupScreen)},
	asm.BuiltinIsTalking:           {Name: "IsTalking", Type: reflect.TypeOf(ns.IsTalking)},
	asm.BuiltinGetTrigger:          {Name: "GetTrigger", Type: reflect.TypeOf(ns.GetTrigger)},
	asm.BuiltinGetCaller:           {Name: "GetCaller", Type: reflect.TypeOf(ns.GetCaller)},
	asm.BuiltinMakeFriendly:        {Name: "MakeFriendly", Type: reflect.TypeOf(ns.MakeFriendly)},
	asm.BuiltinMakeEnemy:           {Name: "MakeEnemy", Type: reflect.TypeOf(ns.MakeEnemy)},
	asm.BuiltinBecomePet:           {Name: "BecomePet", Type: reflect.TypeOf(ns.BecomePet)},
	asm.BuiltinBecomeEnemy:         {Name: "BecomeEnemy", Type: reflect.TypeOf(ns.BecomeEnemy)},
	asm.BuiltinUnknownb8:           {Name: "Unknownb8", Type: reflect.TypeOf(ns.Unknownb8)},
	asm.BuiltinUnknownb9:           {Name: "Unknownb9", Type: reflect.TypeOf(ns.Unknownb9)},
	asm.BuiltinSetHalberd:          {Name: "SetHalberd", Type: reflect.TypeOf(ns.SetHalberd)},
	asm.BuiltinDeathScreen:         {Name: "DeathScreen", Type: reflect.TypeOf(ns.DeathScreen)},
	asm.BuiltinNoWallSound:         {Name: "NoWallSound", Type: reflect.TypeOf(ns.NoWallSound)},
	asm.BuiltinIsTrading:           {Name: "IsTrading", Type: reflect.TypeOf(ns.IsTrading)},
	asm.BuiltinClearMessages:       {Name: "ClearMessages", Type: reflect.TypeOf(ns.ClearMessages)},
	asm.BuiltinSetShopkeeperText:   {Name: "SetShopkeeperText", Type: reflect.TypeOf(ns.SetShopkeeperText)},
	asm.BuiltinUnknownc4:           {Name: "Unknownc4", Type: reflect.TypeOf(ns.Unknownc4)},
	asm.BuiltinIsSummoned:          {Name: "IsSummoned", Type: reflect.TypeOf(ns.IsSummoned)},
	asm.BuiltinMusicPushEvent:      {Name: "MusicPushEvent", Type: reflect.TypeOf(ns.MusicPushEvent)},
	asm.BuiltinMusicPopEvent:       {Name: "MusicPopEvent", Type: reflect.TypeOf(ns.MusicPopEvent)},
	asm.BuiltinMusicEvent:          {Name: "MusicEvent", Type: reflect.TypeOf(ns.MusicEvent)},
	asm.BuiltinIsGameBall:          {Name: "IsGameBall", Type: reflect.TypeOf(ns.IsGameBall)},
	asm.BuiltinIsCrown:             {Name: "IsCrown", Type: reflect.TypeOf(ns.IsCrown)},
	asm.BuiltinEndGame:             {Name: "EndGame", Type: reflect.TypeOf(ns.EndGame)},
	asm.BuiltinImmediateBlind:      {Name: "ImmediateBlind", Type: reflect.TypeOf(ns.ImmediateBlind)},
}

type callBuiltinFunc func(args []ast.Expr) ast.Expr

func (t *translator) callBuiltin(ind asm.Builtin) (rt reflect.Type, callExp callBuiltinFunc) {
	frames := func(x ast.Expr) ast.Expr {
		return callSel(t.imports.script, "Frames", x)
	}
	seconds := func(x ast.Expr) ast.Expr {
		var dt ast.Expr
		if v, ok := asInt(x); ok {
			switch v {
			default:
				dt = &ast.BinaryExpr{
					X: x, Op: token.MUL, Y: sel(t.imports.time, "Second"),
				}
			case 1:
				dt = sel(t.imports.time, "Second")
			case 0:
				dt = x
			}
		} else {
			dt = &ast.BinaryExpr{
				X:  x,
				Op: token.MUL,
				Y:  sel(t.imports.time, "Second"),
			}
		}
		return callSel(t.imports.script, "Time", dt)
	}
	switch ind {
	case asm.BuiltinGetObjectX, asm.BuiltinGetObjectY,
		asm.BuiltinGetWaypointX, asm.BuiltinGetWaypointY:
		if ind == asm.BuiltinGetObjectX || ind == asm.BuiltinGetObjectY {
			rt = reflect.TypeOf((*func(ns.Obj) float32)(nil)).Elem()
		} else {
			rt = reflect.TypeOf((*func(ns.WaypointObj) float32)(nil)).Elem()
		}
		name := "X"
		if ind == asm.BuiltinGetObjectY || ind == asm.BuiltinGetWaypointY {
			name = "Y"
		}
		callExp = func(args []ast.Expr) ast.Expr {
			return sel(callSel(args[0], "Pos"), name)
		}
	case asm.BuiltinGetObjectZ:
		rt = reflect.TypeOf((*func(ns.Obj) float32)(nil)).Elem()
		callExp = func(args []ast.Expr) ast.Expr {
			return callSel(args[0], "Z")
		}
	case asm.BuiltinGetDirection:
		rt = reflect.TypeOf((*func(ns.Obj) ns.Direction)(nil)).Elem()
		callExp = func(args []ast.Expr) ast.Expr {
			return callSel(args[0], "Direction")
		}
	case asm.BuiltinHasClass:
		rt = reflect.TypeOf((*func(ns.Obj, class.Class) bool)(nil)).Elem()
		callExp = func(args []ast.Expr) ast.Expr {
			if str, ok := asStr(args[1]); ok {
				if c, err := object.ParseClass(str); err == nil {
					ce := callSel(args[0], "Class")
					return callSel(ce, "Has", sel(t.imports.object, c.GoString()))
				}
			}
			if s, ok := asStr(args[1]); ok {
				args[1] = sel(t.imports.class, s)
			}
			return callSel(args[0], "HasClass", args[1:]...)
		}
	case asm.BuiltinHasSubclass:
		rt = reflect.TypeOf((*func(ns.Obj, subclass.SubClass) bool)(nil)).Elem()
		callExp = func(args []ast.Expr) ast.Expr {
			if s, ok := asStr(args[1]); ok {
				args[1] = sel(t.imports.subclass, s)
			}
			return callSel(args[0], "HasSubclass", args[1:]...)
		}
	case asm.BuiltinHasEnchant:
		rt = reflect.TypeOf((*func(ns.Obj, enchant.Enchant) bool)(nil)).Elem()
		callExp = func(args []ast.Expr) ast.Expr {
			if s, ok := asStr(args[1]); ok && strings.HasPrefix(s, "ENCHANT_") {
				args[1] = sel(t.imports.enchant, s[8:])
			}
			return callSel(args[0], "HasEnchant", args[1:]...)
		}
	case asm.BuiltinMove, asm.BuiltinGroupMove:
		if ind == asm.BuiltinMove {
			rt = reflect.TypeOf((*func(ns.Obj, ns.WaypointObj))(nil)).Elem()
		} else {
			rt = reflect.TypeOf((*func(ns.ObjGroup, ns.WaypointObj))(nil)).Elem()
		}
		callExp = func(args []ast.Expr) ast.Expr {
			return callSel(args[0], "Move", args[1:]...)
		}
	case asm.BuiltinObjectOn, asm.BuiltinObjectOff,
		asm.BuiltinObjectGroupOn, asm.BuiltinObjectGroupOff,
		asm.BuiltinWaypointOn, asm.BuiltinWaypointOff,
		asm.BuiltinWaypointGroupOn, asm.BuiltinWaypointGroupOff,
		asm.BuiltinWallOpen, asm.BuiltinWallClose,
		asm.BuiltinWallGroupOpen, asm.BuiltinWallGroupClose:
		switch ind {
		case asm.BuiltinObjectOn, asm.BuiltinObjectOff:
			rt = reflect.TypeOf((*func(ns.Obj))(nil)).Elem()
		case asm.BuiltinObjectGroupOn, asm.BuiltinObjectGroupOff:
			rt = reflect.TypeOf((*func(ns.ObjGroup))(nil)).Elem()
		case asm.BuiltinWaypointOn, asm.BuiltinWaypointOff:
			rt = reflect.TypeOf((*func(ns.WaypointObj))(nil)).Elem()
		case asm.BuiltinWaypointGroupOn, asm.BuiltinWaypointGroupOff:
			rt = reflect.TypeOf((*func(ns.WaypointGroupObj))(nil)).Elem()
		case asm.BuiltinWallOpen, asm.BuiltinWallClose:
			rt = reflect.TypeOf((*func(ns.WallObj))(nil)).Elem()
		case asm.BuiltinWallGroupOpen, asm.BuiltinWallGroupClose:
			rt = reflect.TypeOf((*func(ns.WallGroupObj))(nil)).Elem()
		default:
			panic(ind)
		}
		val := t.types.true
		switch ind {
		case asm.BuiltinObjectOff, asm.BuiltinObjectGroupOff,
			asm.BuiltinWaypointOff, asm.BuiltinWaypointGroupOff,
			asm.BuiltinWallClose, asm.BuiltinWallGroupClose:
			val = t.types.false
		}
		callExp = func(args []ast.Expr) ast.Expr {
			return callSel(args[0], "Enable", val)
		}
	case asm.BuiltinFrozen:
		rt = reflect.TypeOf((*func(ns.Obj, bool))(nil)).Elem()
		callExp = func(args []ast.Expr) ast.Expr {
			return callSel(args[0], "Freeze", args[1])
		}
	case asm.BuiltinPauseObject, asm.BuiltinGroupPauseObject:
		if ind == asm.BuiltinPauseObject {
			rt = reflect.TypeOf((*func(ns.Obj, int))(nil)).Elem()
		} else {
			rt = reflect.TypeOf((*func(ns.ObjGroup, int))(nil)).Elem()
		}
		callExp = func(args []ast.Expr) ast.Expr {
			return callSel(args[0], "Pause", frames(args[1]))
		}
	case asm.BuiltinObjectToggle, asm.BuiltinObjectGroupToggle,
		asm.BuiltinWaypointToggle, asm.BuiltinWaypointGroupToggle,
		asm.BuiltinWallToggle, asm.BuiltinWallGroupToggle:
		switch ind {
		case asm.BuiltinObjectToggle:
			rt = reflect.TypeOf((*func(ns.Obj))(nil)).Elem()
		case asm.BuiltinObjectGroupToggle:
			rt = reflect.TypeOf((*func(ns.ObjGroup))(nil)).Elem()
		case asm.BuiltinWaypointToggle:
			rt = reflect.TypeOf((*func(ns.WaypointObj))(nil)).Elem()
		case asm.BuiltinWaypointGroupToggle:
			rt = reflect.TypeOf((*func(ns.WaypointGroupObj))(nil)).Elem()
		case asm.BuiltinWallToggle:
			rt = reflect.TypeOf((*func(ns.WallObj))(nil)).Elem()
		case asm.BuiltinWallGroupToggle:
			rt = reflect.TypeOf((*func(ns.WallGroupObj))(nil)).Elem()
		default:
			panic(ind)
		}
		callExp = func(args []ast.Expr) ast.Expr {
			return callSel(args[0], "Toggle")
		}
	case asm.BuiltinWallBreak, asm.BuiltinWallGroupBreak:
		switch ind {
		case asm.BuiltinWallBreak:
			rt = reflect.TypeOf((*func(ns.WallObj))(nil)).Elem()
		case asm.BuiltinWallGroupBreak:
			rt = reflect.TypeOf((*func(ns.WallGroupObj))(nil)).Elem()
		default:
			panic(ind)
		}
		callExp = func(args []ast.Expr) ast.Expr {
			return callSel(args[0], "Destroy")
		}
	case asm.BuiltinIsObjectOn, asm.BuiltinIsWaypointOn:
		switch ind {
		case asm.BuiltinIsObjectOn:
			rt = reflect.TypeOf((*func(ns.Obj) bool)(nil)).Elem()
		case asm.BuiltinIsWaypointOn:
			rt = reflect.TypeOf((*func(ns.WaypointObj) bool)(nil)).Elem()
		default:
			panic(ind)
		}
		callExp = func(args []ast.Expr) ast.Expr {
			return callSel(args[0], "IsEnabled")
		}
	case asm.BuiltinCurrentHealth, asm.BuiltinMaxHealth:
		rt = reflect.TypeOf((*func(ns.Obj) int)(nil)).Elem()
		name := "CurrentHealth"
		if ind == asm.BuiltinMaxHealth {
			name = "MaxHealth"
		}
		callExp = func(args []ast.Expr) ast.Expr {
			return callSel(args[0], name)
		}
	case asm.BuiltinRestoreHealth:
		rt = reflect.TypeOf((*func(ns.Obj, int))(nil)).Elem()
		callExp = func(args []ast.Expr) ast.Expr {
			return callSel(args[0], "RestoreHealth", args[1])
		}
	case asm.BuiltinLookAtDirection, asm.BuiltinGroupLookAtDirection:
		if ind == asm.BuiltinLookAtDirection {
			rt = reflect.TypeOf((*func(ns.Obj, ns.Direction))(nil)).Elem()
		} else {
			rt = reflect.TypeOf((*func(ns.ObjGroup, ns.Direction))(nil)).Elem()
		}
		callExp = func(args []ast.Expr) ast.Expr {
			if d, ok := asInt(args[1]); ok {
				var name string
				switch d {
				case 0:
					name = "NW"
				case 1:
					name = "N"
				case 2:
					name = "NE"
				case 3:
					name = "W"
				case 5:
					name = "E"
				case 6:
					name = "SW"
				case 7:
					name = "S"
				case 8:
					name = "SE"
				}
				if name != "" {
					args[1] = sel(t.imports.ns, name)
				}
			}
			return callSel(args[0], "LookAtDirection", args[1])
		}
	case asm.BuiltinLookWithAngle:
		rt = reflect.TypeOf((*func(ns.Obj, int))(nil)).Elem()
		callExp = func(args []ast.Expr) ast.Expr {
			return callSel(args[0], "LookWithAngle", args[1])
		}
	case asm.BuiltinLookAtObject:
		rt = reflect.TypeOf((*func(ns.Obj, script.Positioner))(nil)).Elem()
		callExp = func(args []ast.Expr) ast.Expr {
			return callSel(args[0], "LookAtObject", args[1])
		}
	case asm.BuiltinDelete, asm.BuiltinGroupDelete:
		if ind == asm.BuiltinDelete {
			rt = reflect.TypeOf((*func(ns.Obj))(nil)).Elem()
		} else {
			rt = reflect.TypeOf((*func(ns.ObjGroup))(nil)).Elem()
		}
		callExp = func(args []ast.Expr) ast.Expr {
			return callSel(args[0], "Delete")
		}
	case asm.BuiltinDeleteObjectTimer:
		rt = reflect.TypeOf((*func(ns.Obj, int))(nil)).Elem()
		callExp = func(args []ast.Expr) ast.Expr {
			return callSel(args[0], "DeleteAfter", frames(args[1]))
		}
	case asm.BuiltinCreatureIdle, asm.BuiltinCreatureGroupIdle,
		asm.BuiltinWander, asm.BuiltinGroupWander,
		asm.BuiltinCreatureHunt, asm.BuiltinCreatureGroupHunt,
		asm.BuiltinGoBackHome:
		switch ind {
		case asm.BuiltinWander, asm.BuiltinGoBackHome:
			rt = reflect.TypeOf((*func(ns.Obj))(nil)).Elem()
		default:
			rt = reflect.TypeOf((*func(ns.ObjGroup))(nil)).Elem()
		}
		var name string
		switch ind {
		case asm.BuiltinCreatureIdle, asm.BuiltinCreatureGroupIdle:
			name = "Idle"
		case asm.BuiltinWander, asm.BuiltinGroupWander:
			name = "Wander"
		case asm.BuiltinCreatureHunt, asm.BuiltinCreatureGroupHunt:
			name = "Hunt"
		case asm.BuiltinGoBackHome:
			name = "Return"
		default:
			panic(ind)
		}
		callExp = func(args []ast.Expr) ast.Expr {
			return callSel(args[0], name)
		}
	case asm.BuiltinLockDoor, asm.BuiltinUnlockDoor:
		rt = reflect.TypeOf((*func(ns.Obj))(nil)).Elem()
		val := t.types.true
		if ind == asm.BuiltinUnlockDoor {
			val = t.types.false
		}
		callExp = func(args []ast.Expr) ast.Expr {
			return callSel(args[0], "Lock", val)
		}
	case asm.BuiltinIsLocked:
		rt = reflect.TypeOf((*func(ns.Obj) bool)(nil)).Elem()
		callExp = func(args []ast.Expr) ast.Expr {
			return callSel(args[0], "IsLocked")
		}
	case asm.BuiltinMoveObject, asm.BuiltinMoveWaypoint:
		if ind == asm.BuiltinMoveObject {
			rt = reflect.TypeOf((*func(ns.Obj, float32, float32))(nil)).Elem()
		} else {
			rt = reflect.TypeOf((*func(ns.WaypointObj, float32, float32))(nil)).Elem()
		}
		callExp = func(args []ast.Expr) ast.Expr {
			pos, _ := t.asPos(args[1], args[2])
			return callSel(args[0], "SetPos", pos)
		}
	case asm.BuiltinPushObjectTo:
		rt = reflect.TypeOf((*func(ns.Obj, float32, float32))(nil)).Elem()
		callExp = func(args []ast.Expr) ast.Expr {
			pos, _ := t.asPos(args[1], args[2])
			return callSel(args[0], "ApplyForce", pos)
		}
	case asm.BuiltinPushObject:
		rt = reflect.TypeOf((*func(ns.Obj, float32, float32, float32))(nil)).Elem()
		callExp = func(args []ast.Expr) ast.Expr {
			pos, _ := t.asPos(args[2], args[3])
			return callSel(args[0], "PushTo", pos, args[1])
		}
	case asm.BuiltinRaise:
		rt = reflect.TypeOf((*func(ns.Obj, float32))(nil)).Elem()
		callExp = func(args []ast.Expr) ast.Expr {
			return callSel(args[0], "SetZ", args[1])
		}
	case asm.BuiltinDamage, asm.BuiltinGroupDamage:
		if ind == asm.BuiltinDamage {
			rt = reflect.TypeOf((*func(ns.Obj, ns.Obj, int, damage.Type))(nil)).Elem()
		} else {
			rt = reflect.TypeOf((*func(ns.ObjGroup, ns.Obj, int, damage.Type))(nil)).Elem()
		}
		callExp = func(args []ast.Expr) ast.Expr {
			if i, ok := asInt(args[3]); ok {
				args[3] = sel(t.imports.damage, damage.Type(i).String())
			}
			return callSel(args[0], "Damage", args[1:]...)
		}
	case asm.BuiltinIsAttackedBy:
		rt = reflect.TypeOf((*func(ns.Obj, ns.Obj) bool)(nil)).Elem()
		callExp = func(args []ast.Expr) ast.Expr {
			return callSel(args[0], "IsAttackedBy", args[1])
		}
	case asm.BuiltinHasItem, asm.BuiltinPickup, asm.BuiltinDrop:
		rt = reflect.TypeOf((*func(ns.Obj, ns.Obj) bool)(nil)).Elem()
		var name string
		switch ind {
		case asm.BuiltinHasItem:
			name = "HasItem"
		case asm.BuiltinPickup:
			name = "Pickup"
		case asm.BuiltinDrop:
			name = "Drop"
		default:
			panic(ind)
		}
		callExp = func(args []ast.Expr) ast.Expr {
			return callSel(args[0], name, args[1])
		}
	case asm.BuiltinGetHolder, asm.BuiltinGetLastItem, asm.BuiltinGetPreviousItem:
		rt = reflect.TypeOf((*func(ns.Obj) ns.Obj)(nil)).Elem()
		var name string
		switch ind {
		case asm.BuiltinGetHolder:
			name = "GetHolder"
		case asm.BuiltinGetLastItem:
			name = "GetLastItem"
		case asm.BuiltinGetPreviousItem:
			name = "GetPreviousItem"
		default:
			panic(ind)
		}
		callExp = func(args []ast.Expr) ast.Expr {
			return callSel(args[0], name)
		}
	case asm.BuiltinAttack, asm.BuiltinGroupAttack,
		asm.BuiltinCreatureFollow, asm.BuiltinCreatureGroupFollow:
		switch ind {
		case asm.BuiltinAttack, asm.BuiltinCreatureFollow:
			rt = reflect.TypeOf((*func(ns.Obj, script.Positioner))(nil)).Elem()
		default:
			rt = reflect.TypeOf((*func(ns.ObjGroup, script.Positioner))(nil)).Elem()
		}
		var name string
		switch ind {
		case asm.BuiltinAttack, asm.BuiltinGroupAttack:
			name = "Attack"
		case asm.BuiltinCreatureFollow, asm.BuiltinCreatureGroupFollow:
			name = "Follow"
		default:
			panic(ind)
		}
		callExp = func(args []ast.Expr) ast.Expr {
			return callSel(args[0], name, args[1])
		}
	case asm.BuiltinCreatureGuard, asm.BuiltinCreatureGroupGuard:
		if ind == asm.BuiltinCreatureGuard {
			rt = reflect.TypeOf((*func(ns.Obj, float32, float32, float32, float32, float32))(nil)).Elem()
		} else {
			rt = reflect.TypeOf((*func(ns.ObjGroup, float32, float32, float32, float32, float32))(nil)).Elem()
		}
		callExp = func(args []ast.Expr) ast.Expr {
			p1, _ := t.asPos(args[1], args[2])
			p2, _ := t.asPos(args[3], args[4])
			return callSel(args[0], "Guard", p1, p2, args[5])
		}
	case asm.BuiltinRunAway, asm.BuiltinGroupRunAway:
		if ind == asm.BuiltinRunAway {
			rt = reflect.TypeOf((*func(ns.Obj, script.Positioner, int))(nil)).Elem()
		} else {
			rt = reflect.TypeOf((*func(ns.ObjGroup, script.Positioner, int))(nil)).Elem()
		}
		callExp = func(args []ast.Expr) ast.Expr {
			return callSel(args[0], "Flee", args[1], frames(args[2]))
		}
	case asm.BuiltinWalk, asm.BuiltinGroupWalk:
		if ind == asm.BuiltinWalk {
			rt = reflect.TypeOf((*func(ns.Obj, float32, float32))(nil)).Elem()
		} else {
			rt = reflect.TypeOf((*func(ns.ObjGroup, float32, float32))(nil)).Elem()
		}
		callExp = func(args []ast.Expr) ast.Expr {
			pos, _ := t.asPos(args[1], args[2])
			return callSel(args[0], "WalkTo", pos)
		}
	case asm.BuiltinHitLocation, asm.BuiltinGroupHitLocation:
		if ind == asm.BuiltinHitLocation {
			rt = reflect.TypeOf((*func(ns.Obj, float32, float32))(nil)).Elem()
		} else {
			rt = reflect.TypeOf((*func(ns.ObjGroup, float32, float32))(nil)).Elem()
		}
		callExp = func(args []ast.Expr) ast.Expr {
			pos, _ := t.asPos(args[1], args[2])
			return callSel(args[0], "HitMelee", pos)
		}
	case asm.BuiltinHitFarLocation, asm.BuiltinGroupHitFarLocation:
		if ind == asm.BuiltinHitFarLocation {
			rt = reflect.TypeOf((*func(ns.Obj, float32, float32))(nil)).Elem()
		} else {
			rt = reflect.TypeOf((*func(ns.ObjGroup, float32, float32))(nil)).Elem()
		}
		callExp = func(args []ast.Expr) ast.Expr {
			pos, _ := t.asPos(args[1], args[2])
			return callSel(args[0], "HitRanged", pos)
		}
	case asm.BuiltinSetCallback:
		rt = reflect.TypeOf((*func(ns.Obj, ns.ObjectEvent, ns.Func))(nil)).Elem()
		callExp = func(args []ast.Expr) ast.Expr {
			if fp, ok := asInt(args[2]); ok && fp >= 0 && fp < len(t.funcs) {
				args[2] = t.funcs[fp]
			}
			return callSel(args[0], "OnEvent", args[1:]...)
		}
	case asm.BuiltinChat:
		rt = reflect.TypeOf((*func(ns.Obj, string))(nil)).Elem()
		callExp = func(args []ast.Expr) ast.Expr {
			return callSel(args[0], "Chat", args[1])
		}
	case asm.BuiltinChatTimer, asm.BuiltinChatTimerSeconds:
		rt = reflect.TypeOf((*func(ns.Obj, string, int))(nil)).Elem()
		wrap := frames
		if ind == asm.BuiltinChatTimerSeconds {
			wrap = seconds
		}
		callExp = func(args []ast.Expr) ast.Expr {
			return callSel(args[0], "ChatTimer", args[1], wrap(args[2]))
		}
	case asm.BuiltinDestroyChat:
		rt = reflect.TypeOf((*func(ns.Obj))(nil)).Elem()
		callExp = func(args []ast.Expr) ast.Expr {
			return callSel(args[0], "DestroyChat")
		}
	case asm.BuiltinGetElevatorStatus:
		rt = reflect.TypeOf((*func(ns.Obj) int)(nil)).Elem()
		callExp = func(args []ast.Expr) ast.Expr {
			return callSel(args[0], "GetElevatorStatus")
		}
	case asm.BuiltinRetreatLevel, asm.BuiltinGroupRetreatLevel,
		asm.BuiltinResumeLevel, asm.BuiltinGroupResumeLevel,
		asm.BuiltinAggressionLevel, asm.BuiltinGroupAggressionLevel:
		if ind == asm.BuiltinRetreatLevel || ind == asm.BuiltinResumeLevel || ind == asm.BuiltinAggressionLevel {
			rt = reflect.TypeOf((*func(ns.Obj, float32))(nil)).Elem()
		} else {
			rt = reflect.TypeOf((*func(ns.ObjGroup, float32))(nil)).Elem()
		}
		name := "RetreatLevel"
		if ind == asm.BuiltinResumeLevel || ind == asm.BuiltinGroupResumeLevel {
			name = "ResumeLevel"
		} else if ind == asm.BuiltinAggressionLevel || ind == asm.BuiltinGroupAggressionLevel {
			name = "AggressionLevel"
		}
		callExp = func(args []ast.Expr) ast.Expr {
			return callSel(args[0], name, args[1])
		}
	case asm.BuiltinSetRoamFlag, asm.BuiltinGroupSetRoamFlag:
		if ind == asm.BuiltinSetRoamFlag {
			rt = reflect.TypeOf((*func(ns.Obj, int))(nil)).Elem()
		} else {
			rt = reflect.TypeOf((*func(ns.ObjGroup, int))(nil)).Elem()
		}
		callExp = func(args []ast.Expr) ast.Expr {
			return callSel(args[0], "SetRoamFlag", args[1])
		}
	case asm.BuiltinZombieStayDown, asm.BuiltinZombieGroupStayDown,
		asm.BuiltinRaiseZombie, asm.BuiltinRaiseZombieGroup:
		if ind == asm.BuiltinZombieStayDown || ind == asm.BuiltinRaiseZombie {
			rt = reflect.TypeOf((*func(ns.Obj))(nil)).Elem()
		} else {
			rt = reflect.TypeOf((*func(ns.ObjGroup))(nil)).Elem()
		}
		name := "ZombieStayDown"
		if ind == asm.BuiltinRaiseZombie || ind == asm.BuiltinRaiseZombieGroup {
			name = "RaiseZombie"
		}
		callExp = func(args []ast.Expr) ast.Expr {
			return callSel(args[0], name)
		}
	case asm.BuiltinCreateMover:
		rt = reflect.TypeOf((*func(ns.Obj, ns.WaypointObj, float32) ns.Obj)(nil)).Elem()
		callExp = func(args []ast.Expr) ast.Expr {
			return callSel(args[0], "CreateMover", args[1:]...)
		}
	case asm.BuiltinGroupCreateMover:
		rt = reflect.TypeOf((*func(ns.ObjGroup, ns.WaypointObj, float32))(nil)).Elem()
		callExp = func(args []ast.Expr) ast.Expr {
			return callSel(args[0], "CreateMover", args[1:]...)
		}
	case asm.BuiltinIsVisibleTo:
		rt = reflect.TypeOf((*func(ns.Obj, ns.Obj) bool)(nil)).Elem()
		callExp = func(args []ast.Expr) ast.Expr {
			return callSel(args[0], "CanSee", args[1])
		}
	case asm.BuiltinSetOwner, asm.BuiltinGroupSetOwner,
		asm.BuiltinSetOwners, asm.BuiltinGroupSetOwners:
		var name string
		switch ind {
		case asm.BuiltinSetOwner:
			name = "SetOwner"
			rt = reflect.TypeOf((*func(ns.Obj, ns.Obj))(nil)).Elem()
		case asm.BuiltinSetOwners:
			name = "SetOwners"
			rt = reflect.TypeOf((*func(ns.ObjGroup, ns.Obj))(nil)).Elem()
		case asm.BuiltinGroupSetOwner:
			name = "SetOwner"
			rt = reflect.TypeOf((*func(ns.Obj, ns.ObjGroup))(nil)).Elem()
		case asm.BuiltinGroupSetOwners:
			name = "SetOwners"
			rt = reflect.TypeOf((*func(ns.ObjGroup, ns.ObjGroup))(nil)).Elem()
		default:
			panic(ind)
		}
		callExp = func(args []ast.Expr) ast.Expr {
			// Important: caller is the second argument!
			return callSel(args[1], name, args[0])
		}
	case asm.BuiltinIsOwnedBy, asm.BuiltinGroupIsOwnedBy,
		asm.BuiltinIsOwnedByAny, asm.BuiltinGroupIsOwnedByAny:
		var name string
		switch ind {
		case asm.BuiltinIsOwnedBy:
			name = "HasOwner"
			rt = reflect.TypeOf((*func(ns.Obj, ns.Obj) bool)(nil)).Elem()
		case asm.BuiltinIsOwnedByAny:
			name = "HasOwnerIn"
			rt = reflect.TypeOf((*func(ns.ObjGroup, ns.Obj) bool)(nil)).Elem()
		case asm.BuiltinGroupIsOwnedBy:
			name = "HasOwner"
			rt = reflect.TypeOf((*func(ns.Obj, ns.ObjGroup) bool)(nil)).Elem()
		case asm.BuiltinGroupIsOwnedByAny:
			name = "HasOwnerIn"
			rt = reflect.TypeOf((*func(ns.ObjGroup, ns.ObjGroup) bool)(nil)).Elem()
		default:
			panic(ind)
		}
		callExp = func(args []ast.Expr) ast.Expr {
			// Important: caller is the second argument!
			return callSel(args[1], name, args[0])
		}
	case asm.BuiltinClearOwner:
		rt = reflect.TypeOf((*func(ns.Obj))(nil)).Elem()
		callExp = func(args []ast.Expr) ast.Expr {
			return callSel(args[0], "ClearOwner")
		}
	case asm.BuiltinEffect:
		rt = reflect.TypeOf((*func(effect.Effect, float32, float32, float32, float32))(nil)).Elem()
		callExp = func(args []ast.Expr) ast.Expr {
			if s, ok := asStr(args[0]); ok {
				args[0] = sel(t.imports.effects, s)
			}
			p1, _ := t.asPos(args[1], args[2])
			p2, _ := t.asPos(args[3], args[4])
			return callSel(t.imports.ns, "Effect", args[0], p1, p2)
		}
	case asm.BuiltinGetGold, asm.BuiltinGetScore:
		var name string
		switch ind {
		case asm.BuiltinGetGold:
			name = "GetGold"
		case asm.BuiltinGetScore:
			name = "GetScore"
		default:
			panic(ind)
		}
		rt = reflect.TypeOf((*func(ns.Obj) int)(nil)).Elem()
		callExp = func(args []ast.Expr) ast.Expr {
			return callSel(args[0], name)
		}
	case asm.BuiltinChangeGold, asm.BuiltinChangeScore:
		var name string
		switch ind {
		case asm.BuiltinChangeGold:
			name = "ChangeGold"
		case asm.BuiltinChangeScore:
			name = "ChangeScore"
		default:
			panic(ind)
		}
		rt = reflect.TypeOf((*func(ns.Obj, int))(nil)).Elem()
		callExp = func(args []ast.Expr) ast.Expr {
			return callSel(args[0], name, args[1])
		}
	case asm.BuiltinGiveXp:
		rt = reflect.TypeOf((*func(ns.Obj, float32))(nil)).Elem()
		callExp = func(args []ast.Expr) ast.Expr {
			return callSel(args[0], "GiveXp", args[1])
		}
	case asm.BuiltinSecondTimer, asm.BuiltinFrameTimer,
		asm.BuiltinSecondTimerWithArg, asm.BuiltinFrameTimerWithArg:
		var wrap func(ast.Expr) ast.Expr
		switch ind {
		case asm.BuiltinSecondTimer:
			wrap = seconds
			rt = reflect.TypeOf((*func(int, ns.Func) ns.Timer)(nil)).Elem()
		case asm.BuiltinFrameTimer:
			wrap = frames
			rt = reflect.TypeOf((*func(int, ns.Func) ns.Timer)(nil)).Elem()
		case asm.BuiltinSecondTimerWithArg:
			wrap = seconds
			rt = reflect.TypeOf((*func(int, ns.Func, any) ns.Timer)(nil)).Elem()
		case asm.BuiltinFrameTimerWithArg:
			wrap = frames
			rt = reflect.TypeOf((*func(int, ns.Func, any) ns.Timer)(nil)).Elem()
		default:
			panic(ind)
		}
		callExp = func(args []ast.Expr) ast.Expr {
			args[0] = wrap(args[0])
			if len(args) > 2 {
				args[1], args[2] = args[2], args[1]
			}
			if fp, ok := asInt(args[1]); ok && fp >= 0 && fp < len(t.funcs) {
				args[1] = t.funcs[fp]
			}
			return callSel(t.imports.ns, "NewTimer", args...)
		}
	case asm.BuiltinIntToString:
		rt = reflect.TypeOf((*func(int) string)(nil)).Elem()
		callExp = func(args []ast.Expr) ast.Expr {
			return callSel(t.imports.strconv, "Itoa", args[0])
		}
	case asm.BuiltinFloatToString:
		rt = reflect.TypeOf((*func(float32) string)(nil)).Elem()
		callExp = func(args []ast.Expr) ast.Expr {
			return callSel(t.imports.strconv, "FormatFloat",
				call(ast.NewIdent("float64"), args[0]),
				&ast.BasicLit{Kind: token.CHAR, Value: "'g'"},
				intLit(-1), intLit(32),
			)
		}
	case asm.BuiltinCancelTimer:
		rt = reflect.TypeOf((*func(ns.Timer) bool)(nil)).Elem()
		callExp = func(args []ast.Expr) ast.Expr {
			return callSel(args[0], "Cancel")
		}
	case asm.BuiltinDistance:
		rt = reflect.TypeOf((*func(float32, float32, float32, float32) float32)(nil)).Elem()
		callExp = func(args []ast.Expr) ast.Expr {
			p1, _ := t.asPos(args[0], args[1])
			p2, _ := t.asPos(args[2], args[3])
			return call(ast.NewIdent("float32"), callSel(callSel(p1, "Sub", p2), "Len"))
		}
	case asm.BuiltinCastSpellLocationLocation:
		rt = reflect.TypeOf((*func(spell.Spell, float32, float32, float32, float32))(nil)).Elem()
		callExp = func(args []ast.Expr) ast.Expr {
			if s, ok := asStr(args[0]); ok && strings.HasPrefix(s, "SPELL_") {
				args[0] = sel(t.imports.spell, s[6:])
			}
			p1, _ := t.asPos(args[1], args[2])
			p2, _ := t.asPos(args[3], args[4])
			return callSel(t.imports.ns, "CastSpell", args[0], p1, p2)
		}
	case asm.BuiltinCastSpellLocationObject:
		rt = reflect.TypeOf((*func(spell.Spell, float32, float32, script.Positioner))(nil)).Elem()
		callExp = func(args []ast.Expr) ast.Expr {
			if s, ok := asStr(args[0]); ok && strings.HasPrefix(s, "SPELL_") {
				args[0] = sel(t.imports.spell, s[6:])
			}
			p, _ := t.asPos(args[1], args[2])
			return callSel(t.imports.ns, "CastSpell", args[0], p, args[3])
		}
	case asm.BuiltinCastSpellObjectLocation:
		rt = reflect.TypeOf((*func(spell.Spell, script.Positioner, float32, float32))(nil)).Elem()
		callExp = func(args []ast.Expr) ast.Expr {
			if s, ok := asStr(args[0]); ok && strings.HasPrefix(s, "SPELL_") {
				args[0] = sel(t.imports.spell, s[6:])
			}
			p, _ := t.asPos(args[2], args[3])
			return callSel(t.imports.ns, "CastSpell", args[0], args[1], p)
		}
	case asm.BuiltinCastSpellObjectObject:
		rt = reflect.TypeOf((*func(spell.Spell, script.Positioner, script.Positioner))(nil)).Elem()
		callExp = func(args []ast.Expr) ast.Expr {
			if s, ok := asStr(args[0]); ok && strings.HasPrefix(s, "SPELL_") {
				args[0] = sel(t.imports.spell, s[6:])
			}
			return callSel(t.imports.ns, "CastSpell", args[0], args[1], args[2])
		}
	case asm.BuiltinAwardSpell, asm.BuiltinGroupAwardSpell:
		if ind == asm.BuiltinAwardSpell {
			rt = reflect.TypeOf((*func(ns.Obj, spell.Spell) bool)(nil)).Elem()
		} else {
			rt = reflect.TypeOf((*func(ns.ObjGroup, spell.Spell) bool)(nil)).Elem()
		}
		callExp = func(args []ast.Expr) ast.Expr {
			if s, ok := asStr(args[1]); ok && strings.HasPrefix(s, "SPELL_") {
				args[1] = sel(t.imports.spell, s[6:])
			}
			return callSel(args[0], "AwardSpell", args[1])
		}
	case asm.BuiltinEnchant, asm.BuiltinGroupEnchant:
		if ind == asm.BuiltinEnchant {
			rt = reflect.TypeOf((*func(ns.Obj, enchant.Enchant, float32))(nil)).Elem()
		} else {
			rt = reflect.TypeOf((*func(ns.ObjGroup, enchant.Enchant, float32))(nil)).Elem()
		}
		callExp = func(args []ast.Expr) ast.Expr {
			if s, ok := asStr(args[1]); ok && strings.HasPrefix(s, "ENCHANT_") {
				args[1] = sel(t.imports.enchant, s[8:])
			}
			return callSel(args[0], "Enchant", args[1], seconds(args[2]))
		}
	case asm.BuiltinEnchantOff:
		rt = reflect.TypeOf((*func(ns.Obj, enchant.Enchant))(nil)).Elem()
		callExp = func(args []ast.Expr) ast.Expr {
			if s, ok := asStr(args[1]); ok && strings.HasPrefix(s, "ENCHANT_") {
				args[1] = sel(t.imports.enchant, s[8:])
			}
			return callSel(args[0], "EnchantOff", args[1])
		}
	case asm.BuiltinTrapSpells:
		rt = reflect.TypeOf((*func(ns.Obj, spell.Spell, spell.Spell, spell.Spell))(nil)).Elem()
		callExp = func(args []ast.Expr) ast.Expr {
			if s, ok := asStr(args[1]); ok && strings.HasPrefix(s, "SPELL_") {
				args[1] = sel(t.imports.spell, s[6:])
			} else if s == "NULL" {
				args[1] = stringLit("")
			}
			if s, ok := asStr(args[2]); ok && strings.HasPrefix(s, "SPELL_") {
				args[2] = sel(t.imports.spell, s[6:])
			} else if s == "NULL" {
				args[2] = stringLit("")
			}
			if s, ok := asStr(args[3]); ok && strings.HasPrefix(s, "SPELL_") {
				args[3] = sel(t.imports.spell, s[6:])
			} else if s == "NULL" {
				args[3] = stringLit("")
			}
			return callSel(args[0], "TrapSpells", args[1:]...)
		}
	default:
		var fnc ast.Expr
		if ind >= 0 && int(ind) < len(builtins) {
			fnc = t.builtins[ind]
			if fnc == nil {
				panic(fmt.Errorf("no builtin: %v", ind))
			}
		} else {
			fnc = ast.NewIdent(fmt.Sprintf("builtin_overflow_%d", uint32(ind)))
		}
		rt, _ = getType(fnc).(reflect.Type)
		callExp = func(args []ast.Expr) ast.Expr {
			switch ind {
			case asm.BuiltinSetDialog:
				if s, ok := asStr(args[1]); ok {
					name := s
					switch name {
					case "NORMAL":
						name = "Normal"
					case "NEXT":
						name = "Next"
					case "YESNO":
						name = "YesNo"
					case "YESNONEXT":
						name = "YesNoNext"
					case "FALSE":
						name = "False"
					}
					args[1] = sel(t.imports.ns, "Dialog"+name)
				}
				if fp, ok := asInt(args[2]); ok && fp >= 0 && fp < len(t.funcs) {
					args[2] = t.funcs[fp]
				}
				if fp, ok := asInt(args[3]); ok && fp >= 0 && fp < len(t.funcs) {
					args[3] = t.funcs[fp]
				}
			case asm.BuiltinAudioEvent, asm.BuiltinTellStory:
				if s, ok := asStr(args[0]); ok {
					args[0] = sel(t.imports.audio, s)
				}
			case asm.BuiltinEnchant, asm.BuiltinEnchantOff, asm.BuiltinHasEnchant:
				if s, ok := asStr(args[1]); ok && strings.HasPrefix(s, "ENCHANT_") {
					args[1] = sel(t.imports.enchant, strings.TrimPrefix(s, "ENCHANT_"))
				}
			}
			return &ast.CallExpr{Fun: fnc, Args: args}
		}
	}
	return
}
