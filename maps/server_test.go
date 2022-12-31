package maps

import (
	"context"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/noxworld-dev/opennox-lib/ifs"
	"github.com/noxworld-dev/opennox-lib/noxtest"
)

func copyFile(t testing.TB, dst, src string) {
	s, err := ifs.Open(src)
	require.NoError(t, err)
	defer s.Close()
	d, err := os.Create(dst)
	require.NoError(t, err)
	defer d.Close()
	_, err = io.Copy(d, s)
	require.NoError(t, err)
	err = d.Close()
	require.NoError(t, err)
}

func TestMapServer(t *testing.T) {
	dpath := noxtest.DataPath(t, "maps")
	srcdir, err := os.MkdirTemp("", "opennox-map-test-src-*")
	require.NoError(t, err)
	dstdir, err := os.MkdirTemp("", "opennox-map-test-dst-*")
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = os.RemoveAll(srcdir)
		_ = os.RemoveAll(dstdir)
	})
	const mname = "estate"
	for _, name := range []string{
		"estate", "!estate",
		"Estate", "!Estate",
	} {
		name := name
		t.Run(name, func(t *testing.T) {
			// Copy a well-known map to a source folder.
			err = os.MkdirAll(filepath.Join(srcdir, name), 0755)
			require.NoError(t, err)
			smpath := filepath.Join(srcdir, name, name+".map")
			copyFile(t, smpath, filepath.Join(dpath, mname, mname+".map"))

			// Start serving maps from source folder.
			srv := NewServer(srcdir)
			hsrv := &http.Server{Handler: srv}
			l, err := net.Listen("tcp", ":0")
			require.NoError(t, err)
			t.Cleanup(func() {
				_ = l.Close()
			})
			t.Logf("listening on %q", l.Addr())
			go func() {
				err := hsrv.Serve(l)
				if err == http.ErrServerClosed {
					return
				}
				require.NoError(t, err)
			}()
			t.Cleanup(func() {
				_ = hsrv.Close()
			})

			ctx := context.Background()
			cli, err := NewClient(ctx, l.Addr().String())
			require.NoError(t, err)
			defer cli.Close()

			err = cli.DownloadMap(ctx, dstdir, name)
			require.NoError(t, err)
		})
	}
}
