package lua

import (
	"image"
	"testing"
	"time"

	"github.com/noxworld-dev/opennox-lib/object"
	"github.com/noxworld-dev/opennox-lib/script"
	"github.com/noxworld-dev/opennox-lib/types"
)

type fakeObject struct {
	pos types.Pointf
	id  string
}

func (v *fakeObject) String() string {
	panic("implement me")
}

func (v *fakeObject) GetObject() script.Object {
	panic("implement me")
}

func (v *fakeObject) Class() object.Class {
	panic("implement me")
}

func (v *fakeObject) Owner() script.Object {
	panic("implement me")
}

func (v *fakeObject) SetOwner(owner script.ObjectWrapper) {
	panic("implement me")
}

func (v *fakeObject) Destroy() {
	panic("implement me")
}

func (v *fakeObject) OnTriggerActivate(fnc func(obj script.Object)) {
	panic("implement me")
}

func (v *fakeObject) OnTriggerDeactivate(fnc func()) {
	panic("implement me")
}

func (v *fakeObject) ObjectType() script.ObjectType {
	panic("implement me")
}

func (v *fakeObject) Z() float32 {
	panic("implement me")
}

func (v *fakeObject) SetZ(z float32) {
	panic("implement me")
}

func (v *fakeObject) Push(vec types.Pointf, force float32) {
	panic("implement me")
}

func (v *fakeObject) PushTo(p types.Pointf) {
	panic("implement me")
}

func (v *fakeObject) IsEnabled() bool {
	panic("implement me")
}

func (v *fakeObject) ID() string {
	return v.id
}

func (v *fakeObject) Type() script.ObjectType {
	panic("implement me")
}

func (v *fakeObject) Pos() types.Pointf {
	return v.pos
}

func (v *fakeObject) SetPos(p types.Pointf) {
	v.pos = p
}

func (v *fakeObject) Enable(val bool) {
	panic("implement me")
}

func (v *fakeObject) Delete() {
	panic("implement me")
}

type testPrinter struct {
	t   testing.TB
	lvl string
}

func (p testPrinter) Print(text string) {
	p.t.Logf("%s: %s", p.lvl, text)
}

type fakeGame struct {
	t testing.TB
}

func (g fakeGame) BlindPlayers(blind bool) {
	panic("implement me")
}

func (g fakeGame) CinemaPlayers(v bool) {
	panic("implement me")
}

func (g fakeGame) OnPlayerJoin(fnc func(p script.Player)) {
	panic("implement me")
}

func (g fakeGame) OnPlayerLeave(fnc func(p script.Player)) {
	panic("implement me")
}

func (g fakeGame) ObjectGroupByID(id string) *script.ObjectGroup {
	panic("implement me")
}

func (g fakeGame) WaypointGroupByID(id string) *script.WaypointGroup {
	panic("implement me")
}

func (g fakeGame) WallAt(pos types.Pointf) script.Wall {
	panic("implement me")
}

func (g fakeGame) WallNear(pos types.Pointf) script.Wall {
	panic("implement me")
}

func (g fakeGame) WallAtGrid(pos image.Point) script.Wall {
	panic("implement me")
}

func (g fakeGame) WallGroupByID(id string) *script.WallGroup {
	panic("implement me")
}

func (g fakeGame) Console(error bool) script.Printer {
	lvl := "info"
	if error {
		lvl = "error"
	}
	return testPrinter{t: g.t, lvl: lvl}
}

func (g fakeGame) Frame() int {
	return 1
}

func (g fakeGame) Time() time.Duration {
	return time.Second
}

func (g fakeGame) Players() []script.Player {
	panic("implement me")
}

func (g fakeGame) HostPlayer() script.Player {
	panic("implement me")
}

func (g fakeGame) WaypointByID(id string) script.Waypoint {
	panic("implement me")
}

func (g fakeGame) WallByPos(pos image.Point) script.Wall {
	panic("implement me")
}

func (g fakeGame) AudioEffect(name string, pos script.Positioner) {
	panic("implement me")
}

func (g fakeGame) Global() script.Printer {
	return testPrinter{t: g.t, lvl: "global"}
}

func (g fakeGame) ObjectTypeByID(id string) script.ObjectType {
	return nil
}

func (g fakeGame) ObjectByID(id string) script.Object {
	switch id {
	case "Frog":
		return &fakeObject{pos: types.Pointf{1, 2}, id: id}
	case "Cookie":
		return &fakeObject{pos: types.Pointf{3, 4}, id: id}
	}
	return nil
}

func TestLUA(t *testing.T) {
	vm := NewVM(fakeGame{t: t}, "")
	err := vm.Exec(`
Nox = require("Nox.Map.Script.v0")
print(Nox)
Nox.API = 0
print(Nox)

frog = Nox.Object("Frog")
print(frog.id, frog.x, frog.y)
print(frog:Pos())

frog.x = 5
print(frog:Pos())

function OnStart()
	cookie = Nox.Object("Cookie")
	print("cookie", frog:SetPos(cookie):Pos())
end

frame = 0
function OnFrame()
	frame = frame + 1
	print("frame", frame)
end
`)
	if err != nil {
		t.Fatal(err)
	}
	vm.OnEvent("Start")
	defer vm.OnEvent("End")
	for i := 0; i < 3; i++ {
		vm.OnFrame()
	}
}
