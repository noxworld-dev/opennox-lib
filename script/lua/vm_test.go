package lua

import (
	"testing"

	"github.com/noxworld-dev/opennox-lib/script/scripttest"
)

func TestLUA(t *testing.T) {
	vm := NewVM(scripttest.Game{T: t}, "")
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
