package main

import (
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"go/token"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/noxworld-dev/opennox-lib/maps"
	"github.com/noxworld-dev/opennox-lib/script/noxscript"
	asm "github.com/noxworld-dev/opennox-lib/script/noxscript/noxasm"
	"github.com/noxworld-dev/opennox-lib/script/noxscript/noxast"
)

func init() {
	cmd := &cobra.Command{
		Use:     "noxscript command",
		Short:   "Tools for working with NoxScript",
		Aliases: []string{"ns"},
	}
	Root.AddCommand(cmd)

	cmdExtr := &cobra.Command{
		Use:   "extract input.map output.obj",
		Short: "Extract binary NoxScript file from a map",
	}
	cmd.AddCommand(cmdExtr)
	cmdExtr.RunE = func(cmd *cobra.Command, args []string) error {
		return cmdNSExtract(cmd, args)
	}

	// TODO: noxscript insert command
	//cmdIns := &cobra.Command{
	//	Use:   "insert input.map input.obj",
	//	Short: "Insert or replace binary NoxScript file in a map",
	//}
	//cmd.AddCommand(cmdIns)
	//cmdIns.RunE = func(cmd *cobra.Command, args []string) error {
	//	return cmdNSInsert(cmd, args)
	//}

	cmdDis := &cobra.Command{
		Use:   "disasm input.obj",
		Short: "Disassemble binary NoxScript file into text assembly",
	}
	cmd.AddCommand(cmdDis)
	cmdDis.RunE = func(cmd *cobra.Command, args []string) error {
		return cmdNSDisasm(cmd, args)
	}

	cmdDecomp := &cobra.Command{
		Use:   "decomp input.obj",
		Short: "Decompile binary NoxScript file into human-readable script",
	}
	cmd.AddCommand(cmdDecomp)
	cmdDecomp.RunE = func(cmd *cobra.Command, args []string) error {
		return cmdNSDecomp(cmd, args)
	}
}

func cmdNSExtract(cmd *cobra.Command, args []string) error {
	switch len(args) {
	case 0:
		return errors.New("at least one map file expected")
	case 1, 2:
		fname := args[0]
		f, err := os.Open(fname)
		if err != nil {
			return err
		}
		defer f.Close()

		s, err := maps.ReadScript(f)
		if err != nil {
			return err
		} else if s == nil || len(s.Data) == 0 {
			err := errors.New("no script in map file")
			if len(args) == 2 {
				return err
			}
			log.Println(err)
			return nil
		}

		out := strings.TrimSuffix(fname, filepath.Ext(fname)) + ".obj"
		if len(args) == 2 {
			out = args[1]
		}
		log.Printf("writing %d bytes to %s\n", len(s.Data), out)
		return os.WriteFile(out, s.Data, 0644)
	default:
		dstDir := ""
		hasDst := false
		if last := args[len(args)-1]; filepath.Ext(last) == "" {
			hasDst, dstDir = true, last
			args = args[:len(args)-1]
		}
		var last error
		for _, fname := range args {
			err := func() error {
				f, err := os.Open(fname)
				if err != nil {
					return err
				}
				defer f.Close()

				s, err := maps.ReadScript(f)
				if err != nil {
					return err
				} else if s == nil || len(s.Data) == 0 {
					return nil
				}
				name := filepath.Base(fname)
				name = strings.TrimSuffix(name, filepath.Ext(name)) + ".obj"
				var out string
				if hasDst {
					out = filepath.Join(dstDir, name)
				} else {
					out = filepath.Join(filepath.Dir(fname), name)
				}
				return os.WriteFile(out, s.Data, 0644)
			}()
			if err != nil {
				log.Printf("%s: %v\n", fname, err)
				last = err
			} else {
				log.Println(fname)
			}
		}
		return last
	}
}

func cmdNSDisasm(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("expected one argument")
	}
	fname := args[0]
	f, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer f.Close()

	var raw []byte
	if filepath.Ext(fname) == ".map" {
		s, err := maps.ReadScript(f)
		if err != nil {
			return err
		}
		raw = s.Data
	} else {
		raw, err = io.ReadAll(f)
		if err != nil {
			return err
		}
	}
	_ = f.Close()

	scr, err := noxscript.ReadScript(bytes.NewReader(raw))
	if err != nil {
		return err
	}
	if len(scr.Strings) != 0 {
		fmt.Println("STRINGS:")
		for i, s := range scr.Strings {
			fmt.Printf("\t%d: %q\n", i, s)
		}
		fmt.Println()
	}
	var last error
	for i, fnc := range scr.Funcs {
		fmt.Printf("func %d: %q\n", i, fnc.Name)
		fmt.Printf("\targs: %d, locals: %d, returns: %d\n",
			fnc.Args, len(fnc.Vars)-fnc.Args, fnc.Return)
		fmt.Println()

		code, err := asm.Decode(fnc.Code)
		if err != nil {
			err = fmt.Errorf("cannot disasm %q: %w", fnc.Name, err)
			log.Println(err)
			last = err
			continue
		}
		_ = asm.Print(os.Stdout, code)
		fmt.Println()
	}
	return last
}

func cmdNSDecomp(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("expected one argument")
	}
	fname := args[0]
	f, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer f.Close()

	var raw []byte
	if filepath.Ext(fname) == ".map" {
		s, err := maps.ReadScript(f)
		if err != nil {
			return err
		}
		raw = s.Data
	} else {
		raw, err = io.ReadAll(f)
		if err != nil {
			return err
		}
	}
	_ = f.Close()

	scr, err := noxscript.ReadScript(bytes.NewReader(raw))
	if err != nil {
		return err
	}
	astf := noxast.Translate(scr)
	return format.Node(os.Stdout, token.NewFileSet(), astf)
}
