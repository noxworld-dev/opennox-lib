package eval

import (
	"bytes"
	"os"
	"reflect"
	"strings"

	"github.com/traefik/yaegi/interp"

	"github.com/noxworld-dev/opennox-lib/log"
	"github.com/noxworld-dev/opennox-lib/script"
	"github.com/noxworld-dev/opennox-lib/script/eval/imports"
)

var (
	Log = log.New("eval")
)

func init() {
	script.RegisterVM(script.VMRuntime{
		Name: "eval", Title: "Go interpreter", Log: Log,
		NewMap: func(g script.Game, maps string, name string) (script.VM, error) {
			vm := NewVM(g, maps)
			err := vm.ExecFile(name)
			if os.IsNotExist(err) {
				return vm, nil // still run the empty script for commands
			} else if err != nil {
				return nil, err
			}
			return vm, nil
		},
	})
}

type printer struct {
	vm  *VM
	buf bytes.Buffer
	p   script.Printer
}

func (p *printer) Write(data []byte) (int, error) {
	p.buf.Write(data)
	for p.buf.Len() > 0 {
		data := p.buf.Bytes()
		i := bytes.IndexByte(data, '\n')
		if i < 0 {
			break
		}
		p.p.Print(string(p.buf.Next(i + 1)))
	}
	return len(data), nil
}

var _ script.VM = (*VM)(nil)

type VM struct {
	g        script.Game
	vm       *interp.Interpreter
	printers []*printer
	defs     bool

	curExports map[string]reflect.Value

	onFrame    func()
	onEvent    func(typ script.EventType)
	onEventStr func(typ string)
}

func NewVM(g script.Game, dir string) *VM {
	stdout := &printer{p: g.Console(false)}
	stderr := &printer{p: g.Console(true)}
	vm := &VM{g: g, vm: interp.New(interp.Options{
		GoPath:               "/",
		BuildTags:            []string{"script"},
		Stdin:                bytes.NewReader(nil),
		Stdout:               stdout,
		Stderr:               stderr,
		Args:                 []string{},
		SourcecodeFilesystem: &modFS{root: dir},
	})}
	vm.addPrinters(stdout, stderr)
	vm.initPackages()
	return vm
}

func (vm *VM) addPrinters(list ...*printer) {
	for _, p := range list {
		p.vm = vm
		vm.printers = append(vm.printers, p)
	}
}

func (vm *VM) flushPrinters() {
	for _, p := range vm.printers {
		if p.buf.Len() == 0 {
			continue
		}
		p.p.Print(p.buf.String())
		p.buf.Reset()
	}
}

func importPathFor(v any) string {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.PkgPath()
}

func (vm *VM) initPackages() {
	vm.vm.Use(imports.Symbols)

	// Rename map entry below when refactoring!
	// These strange assignments below statically check that types of the function and override are exactly the same.
	getGame := script.Runtime
	getGame = vm.runtime
	vm.vm.Use(interp.Exports{
		importPathFor((*script.Game)(nil)) + "/script": {
			"Runtime": reflect.ValueOf(getGame),
		},
	})
}

func (vm *VM) runtime() script.Game {
	if vm.g == nil {
		return script.Runtime()
	}
	return vm.g
}

func (vm *VM) initMain() {
	if _, err := vm.vm.Eval(`
package main

import (
	"fmt"
	"github.com/noxworld-dev/opennox-lib/types"
	"github.com/noxworld-dev/opennox-lib/object"
	"github.com/noxworld-dev/opennox-lib/script"
)

var Game = script.Runtime()
`); err != nil {
		panic(err)
	}
}

func (vm *VM) exportByName(name string) (reflect.Value, bool) {
	rv, ok := vm.curExports[name]
	return rv, ok
}

func maybeSetExport[T any](vm *VM, p *T, name string) {
	v, ok := vm.exportByName(name)
	if !ok || v.Kind() != reflect.TypeOf(p).Elem().Kind() {
		return
	}
	*p, _ = v.Interface().(T)
}

func (vm *VM) checkExports() {
	if vm.onFrame == nil {
		maybeSetExport(vm, &vm.onFrame, "OnFrame")
	}
	if vm.onEvent == nil && vm.onEventStr == nil {
		maybeSetExport(vm, &vm.onEvent, "OnEvent")
		maybeSetExport(vm, &vm.onEventStr, "OnEvent")
	}
}

func (vm *VM) Exec(s string) error {
	if !vm.defs {
		vm.initMain()
		vm.defs = true
	}
	defer vm.checkExports()
	defer vm.flushPrinters()
	_, err := vm.vm.Eval(s)
	return err
}

func (vm *VM) ExecFile(pkg string) error {
	defer vm.checkExports()
	defer vm.flushPrinters()
	_, err := vm.vm.EvalPath(pkg)
	if err != nil {
		if strings.Contains(err.Error(), "no Go files") {
			return os.ErrNotExist
		}
		return err
	}
	vm.curExports = vm.vm.Symbols(pkg)[pkg]
	return nil
}

func (vm *VM) OnFrame() {
	if vm.onFrame == nil {
		return
	}
	defer vm.flushPrinters()
	vm.onFrame()
}

func (vm *VM) OnEvent(typ script.EventType) {
	if vm.onEvent == nil && vm.onEventStr == nil {
		return
	}
	defer vm.flushPrinters()
	if vm.onEvent != nil {
		vm.onEvent(typ)
	} else if vm.onEventStr != nil {
		vm.onEventStr(string(typ))
	}
}

func (vm *VM) Close() error {
	vm.flushPrinters()
	*vm = VM{}
	return nil
}
