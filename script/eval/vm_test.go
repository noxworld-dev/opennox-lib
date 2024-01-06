package eval

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/shoenig/test/must"

	"github.com/noxworld-dev/opennox-lib/script"
	"github.com/noxworld-dev/opennox-lib/script/scripttest"
)

func TestScriptDir(t *testing.T) {
	dir, err := os.MkdirTemp("", "opennox_eval_*")
	must.NoError(t, err)
	defer func() {
		_ = os.RemoveAll(dir)
	}()

	writeFile := func(path string, src string) {
		path = filepath.Join(dir, path)
		err = os.MkdirAll(filepath.Dir(path), 0755)
		must.NoError(t, err)
		err = os.WriteFile(path, []byte(src), 0644)
		must.NoError(t, err)
	}

	writeFile("pkg1/script.go", `
package pkg1

import (
	"pkg2"
	"pkg3"
)

var Cnt int

func init() { println("init one"); Cnt++ }
func One() { println("one"); pkg2.Two(); pkg3.Three(); Cnt++ }
func OnFrame() { println("frame"); Cnt++ }

func Panic1() {
	panic2()
}
func panic2() {
	func(){
		panic3()
	}()
}
func panic3() {
	panic("test")
}
`)

	writeFile("pkg2/script.go", `
package pkg2

func init() { println("init two") }
func Two() { println("two") }
`)

	writeFile("pkg1/vendor/pkg3/script.go", `
package pkg3

func init() { println("init three") }
func Three() { println("three") }
`)
	g := &scripttest.Game{T: t}
	vm := NewVM(g, dir)
	err = vm.ExecFile("pkg1")
	must.NoError(t, err)

	rv, ok := vm.exportByName(`One`)
	must.True(t, ok)
	fnc, ok := rv.Interface().(func())
	must.True(t, ok)
	fnc()

	_, err = vm.Exec(`println("builtin")`)
	must.NoError(t, err)
	_, err = vm.Exec(`fmt.Println("fmt")`)
	must.NoError(t, err)
	_, err = vm.Exec(`Game.Global().Print("global")`)
	must.NoError(t, err)
	vm.OnFrame()

	rv, ok = vm.exportByName(`Cnt`)
	must.True(t, ok)
	p, ok := rv.Interface().(int)
	must.True(t, ok)
	must.EqOp(t, 3, p)

	one, err := script.GetVMSymbol[func()](vm, "One")
	must.NoError(t, err)
	one()

	cnt, err := script.GetVMSymbolPtr[int](vm, "Cnt")
	must.NoError(t, err)
	must.EqOp(t, 4, *cnt)

	panic1, err := script.GetVMSymbol[func()](vm, "Panic1")
	must.NoError(t, err)
	func() {
		defer func() {
			recover()
		}()
		panic1()
	}()
	must.EqOp(t, `info: init two
info: init three
info: init one
info: one
info: two
info: three
info: builtin
info: fmt
info: frame
info: one
info: two
info: three
error: /src/pkg1/script.go:24:2: panic: pkg1.panic3(...)
error: /src/pkg1/script.go:20:3: panic: pkg1.panic2.func(...)
error: /src/pkg1/script.go:19:2: panic: pkg1.panic2.func(...)
error: /src/pkg1/script.go:16:2: panic: pkg1.Panic1(...)
`, g.Log.String())
}
