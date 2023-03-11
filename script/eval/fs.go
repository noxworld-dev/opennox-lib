package eval

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func allowedFile(mode fs.FileMode) bool {
	if mode&fs.ModeSymlink != 0 {
		return false
	} else if !mode.IsDir() && !mode.IsRegular() {
		return false
	}
	return true
}

var _ fs.StatFS = (*modFS)(nil)

// modFS is a modified FS tailored for serving OpenNox mod/map packages.
//
// It has the following assumptions/constrains:
// - File system is read-only.
// - Files outside root cannot be accessed.
// - Symlinks and other non-regular files are ignored.
// - All paths must start with /src/ - this is to make Yaegi happy (it assumes packages in GOPATH must be in src).
type modFS struct {
	root string // folder mapped to /src/ in virtual FS
}

type modFileInfo struct {
	m *modFS // to correctly print path
	fs.FileInfo
}

func (fi modFileInfo) Name() string {
	p := fi.m.virtPath(fi.FileInfo.Name())
	return p
}

type modDirEntry struct {
	m *modFS // to correctly print path
	fs.DirEntry
}

func (e modDirEntry) Name() string {
	p := e.m.virtPath(e.DirEntry.Name())
	return p
}

func (e modDirEntry) Info() (fs.FileInfo, error) {
	fi, err := e.DirEntry.Info()
	if err != nil {
		return nil, err
	}
	return modFileInfo{m: e.m, FileInfo: fi}, nil
}

type modFile struct {
	m *modFS // to correctly print path in Stat
	f fs.ReadDirFile
}

func (f modFile) Read(data []byte) (int, error) {
	return f.f.Read(data)
}

func (f modFile) Close() error {
	return f.f.Close()
}

func (f modFile) Stat() (fs.FileInfo, error) {
	fi, err := f.f.Stat()
	if err != nil {
		return nil, err
	}
	return modFileInfo{m: f.m, FileInfo: fi}, nil
}

func (f modFile) ReadDir(n int) ([]fs.DirEntry, error) {
	list, err := f.f.ReadDir(n)
	if err != nil {
		return nil, err
	}
	out := make([]fs.DirEntry, 0, len(list))
	for i := 0; i < len(list); i++ {
		fi := list[i]
		if allowedFile(fi.Type()) {
			out = append(out, modDirEntry{m: f.m, DirEntry: fi})
		}
	}
	return out, nil
}

func (m *modFS) virtPath(p string) string {
	if !filepath.IsAbs(p) {
		p = filepath.Clean(p)
		return p
	}
	if !strings.HasPrefix(p, m.root) {
		return filepath.Base(p)
	}
	p = p[len(m.root):]
	p = filepath.Clean(p)
	p = filepath.Join(rootPref, p)
	return p
}

func (m *modFS) realPath(p string) (string, bool) {
	if filepath.IsAbs(p) {
		if !strings.HasPrefix(p, rootPref) {
			return "", false
		}
		p = p[len(rootPref):]
	}
	p = filepath.Clean(p)
	p = filepath.Join(m.root, p)
	return p, true
}

func (m *modFS) Stat(orig string) (fs.FileInfo, error) {
	p, ok := m.realPath(orig)
	if !ok {
		return nil, fs.ErrNotExist
	}
	fi, err := os.Stat(p)
	if err != nil {
		return nil, err
	}
	// do not expose other methods except the base ones, fix file name
	return modFileInfo{m: m, FileInfo: fi}, nil
}

func (m *modFS) Open(orig string) (fs.File, error) {
	p, ok := m.realPath(orig)
	if !ok {
		return nil, fs.ErrNotExist
	}
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	fi, err := f.Stat()
	if err != nil {
		_ = f.Close()
		return nil, err
	}
	if !allowedFile(fi.Mode()) {
		_ = f.Close()
		return nil, fs.ErrNotExist
	}
	// do not expose other methods except the base ones
	return modFile{m: m, f: f}, nil
}
