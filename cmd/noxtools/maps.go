package main

import (
	"fmt"
	"os"

	"golang.org/x/exp/slices"

	"github.com/noxworld-dev/opennox-lib/maps"
)

func mapReadRawSections(fname string) ([]maps.RawSection, *maps.Header, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()
	rd, err := maps.NewReader(f)
	if err != nil {
		return nil, nil, err
	}
	hdr := rd.Map().Header()
	raw, err := rd.ReadSectionsRaw()
	return raw, &hdr, err
}

func mapWriteRawSections(fname string, hdr *maps.Header, raw []maps.RawSection) error {
	f, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer f.Close()

	wr, err := maps.NewWriter(f, *hdr)
	if err != nil {
		return err
	}
	if err = wr.WriteRawSections(raw); err != nil {
		return err
	}
	if err = wr.Close(); err != nil {
		return err
	}
	return f.Close()
}

func mapRemoveSections(fname string, sections []string, verbose bool) error {
	return mapUpdateSections(fname, func(s *maps.RawSection) (bool, error) {
		if !slices.Contains(sections, s.Name) {
			return false, nil
		}
		*s = maps.RawSection{}
		return true, nil
	}, verbose)
}

func mapUpdateSections(fname string, fnc func(s *maps.RawSection) (bool, error), verbose bool) error {
	raw, hdr, err := mapReadRawSections(fname)
	if err != nil {
		return err
	}
	changed := 0
	out := make([]maps.RawSection, 0, len(raw))
	for _, s := range raw {
		upd, err := fnc(&s)
		if err != nil {
			return err
		}
		if upd {
			changed++
		}
		if !upd || s.Name != "" {
			out = append(out, s)
		}
	}
	if changed == 0 {
		if verbose {
			fmt.Println("no sections updated")
		}
		return nil
	}
	return mapWriteRawSections(fname, hdr, out)
}
