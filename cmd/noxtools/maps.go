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
	raw, hdr, err := mapReadRawSections(fname)
	if err != nil {
		return err
	}
	removed := 0
	out := make([]maps.RawSection, 0, len(raw))
	for _, s := range raw {
		if slices.Contains(sections, s.Name) {
			removed++
		} else {
			out = append(out, s)
		}
	}
	if removed == 0 {
		if verbose {
			fmt.Println("no sections to remove")
		}
		return nil
	}
	return mapWriteRawSections(fname, hdr, out)
}
