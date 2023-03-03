package eval

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/noxworld-dev/opennox-lib/script/scripttest"
)

func TestScriptDir(t *testing.T) {
	dir, err := os.MkdirTemp("", "opennox_eval_*")
	require.NoError(t, err)
	defer func() {
		_ = os.RemoveAll(dir)
	}()

	writeFile := func(path string, src string) {
		path = filepath.Join(dir, path)
		err = os.MkdirAll(filepath.Dir(path), 0755)
		require.NoError(t, err)
		err = os.WriteFile(path, []byte(src), 0644)
		require.NoError(t, err)
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

	vm := NewVM(scripttest.Game{T: t}, dir)
	err = vm.ExecFile("pkg1")
	require.NoError(t, err)

	rv, ok := vm.exportByName(`One`)
	require.True(t, ok)
	fnc, ok := rv.Interface().(func())
	require.True(t, ok)
	fnc()

	err = vm.Exec(`println("builtin")`)
	require.NoError(t, err)
	err = vm.Exec(`fmt.Println("fmt")`)
	require.NoError(t, err)
	err = vm.Exec(`Game.Global().Print("global")`)
	require.NoError(t, err)
	vm.OnFrame()

	rv, ok = vm.exportByName(`Cnt`)
	require.True(t, ok)
	p, ok := rv.Interface().(int)
	require.True(t, ok)
	require.Equal(t, int(3), p)
}
