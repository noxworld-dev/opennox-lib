package eval

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/traefik/yaegi/interp"

	eudeval "github.com/noxworld-dev/noxscript/eud/v171/eval"
	nseval3 "github.com/noxworld-dev/noxscript/ns/v3/eval"
	ns4 "github.com/noxworld-dev/noxscript/ns/v4"
	nseval4 "github.com/noxworld-dev/noxscript/ns/v4/eval"

	"github.com/noxworld-dev/opennox-lib/log"
	"github.com/noxworld-dev/opennox-lib/script"
	"github.com/noxworld-dev/opennox-lib/script/eval/imports"
	"github.com/noxworld-dev/opennox-lib/script/eval/stdlib"
)

var (
	Log = log.New("eval")
)

var useSymbols []Exports

type Exports struct {
	Symbols interp.Exports
	Context func(ctx context.Context, g script.Game) context.Context
}

// Register additional libraries for the interpreter.
func Register(m Exports) {
	useSymbols = append(useSymbols, m)
}

func init() {
	script.RegisterVM(script.VMRuntime{
		Name: "eval", Title: "Go interpreter", Log: Log,
		NewMap: func(g script.Game, maps string, name string) (script.VM, error) {
			vm := NewVM(g, maps)
			err := vm.ExecFile(name)
			if os.IsNotExist(err) {
				Log.Println("no go files")
				return vm, nil // still run the empty script for commands
			} else if err != nil {
				return nil, err
			}
			Log.Println("go scripts loaded")
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
	vm := &VM{vm: interp.New(interp.Options{
		GoPath:               gopath,
		BuildTags:            []string{"script"},
		Stdin:                bytes.NewReader(nil),
		Stdout:               stdout,
		Stderr:               stderr,
		Args:                 []string{},
		SourcecodeFilesystem: &modFS{root: dir},
	})}
	vm.addPrinters(stdout, stderr)
	vm.initPackages(g)
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

func (vm *VM) initPackages(g script.Game) {
	ctx := context.Background()

	vm.vm.Use(stdlib.Symbols)
	vm.vm.Use(imports.Symbols)
	vm.vm.Use(nseval3.Symbols)
	vm.vm.Use(nseval4.Symbols)
	vm.vm.Use(eudeval.Symbols)
	ctx = script.WithGame(ctx, g)

	for _, m := range useSymbols {
		if m.Symbols != nil {
			vm.vm.Use(m.Symbols)
		}
		if m.Context != nil {
			ctx = m.Context(ctx, g)
		}
	}

	// Rename map entry below when refactoring!
	// These strange assignments below statically check that types of the function and override are exactly the same.
	getGame := script.Runtime
	getGame = func() script.Game {
		return g
	}
	vm.vm.Use(interp.Exports{
		importPathFor((*script.Game)(nil)) + "/script": {
			"Runtime": reflect.ValueOf(getGame),
		},
	})
	getCtx := context.Background
	getCtx = func() context.Context {
		return ctx
	}
	vm.vm.Use(interp.Exports{
		importPathFor((*context.Context)(nil)) + "/context": {
			"Background": reflect.ValueOf(getCtx),
		},
	})
	// TODO: properly virtualize
	script.SetRuntime(g)
	if v, ok := g.(ns4.Game); ok {
		ns4.SetRuntime(v.NoxScript())
	}
}

func (vm *VM) initMain() {
	if _, err := vm.vm.Eval(`
package main

import (
	"fmt"

	"github.com/noxworld-dev/opennox-lib/types"
	"github.com/noxworld-dev/opennox-lib/object"
	"github.com/noxworld-dev/opennox-lib/script"

	ns3 "github.com/noxworld-dev/noxscript/ns/v3"
	ns4 "github.com/noxworld-dev/noxscript/ns/v4"
	eud "github.com/noxworld-dev/noxscript/eud/v171"
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

func (vm *VM) Exec(s string) (reflect.Value, error) {
	if !vm.defs {
		vm.initMain()
		vm.defs = true
	}
	defer vm.checkExports()
	defer vm.flushPrinters()
	return vm.vm.Eval(s)
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
	if len(vm.curExports) == 0 {
		Log.Println("no exports from go; wrong package name?")
	}
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

func (vm *VM) GetSymbol(name string, typ reflect.Type) (reflect.Value, bool, error) {
	rv, ok := vm.exportByName(name)
	if !ok {
		return reflect.Value{}, false, nil
	} else if rv.Type() != typ {
		return reflect.Value{}, false, fmt.Errorf("unexpected type of %q: expected %v, got %v", name, typ, rv.Type())
	}
	return rv, true, nil
}

func (vm *VM) Close() error {
	vm.flushPrinters()
	*vm = VM{}
	return nil
}
