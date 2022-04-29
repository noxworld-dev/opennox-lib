package cfg

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/noxworld-dev/opennox-lib/noxtest"
)

func TestParseConfig(t *testing.T) {
	path := noxtest.DataPath(t)
	for _, name := range []string{
		"nox.cfg",
		"default.cfg",
	} {
		t.Run(name, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join(path, name))
			if os.IsNotExist(err) {
				t.SkipNow()
			}
			require.NoError(t, err)
			data = bytes.ReplaceAll(data, []byte("\r\n"), []byte("\n"))

			file, err := Parse(bytes.NewReader(data))
			require.NoError(t, err)
			require.Len(t, file.Sections, 2)
			vers, _ := file.Sections[0].Get("Version")
			require.NotZero(t, vers)
			mouse, _ := file.Sections[1].Get("MousePickup")
			require.NotZero(t, mouse)
		})
	}
}

func TestConfig(t *testing.T) {
	const conf = `Version = 65537
# comment
VideoMode = 1024 768 16
---
MousePickup = Left
I+M = ToggleInventory
---
`
	file, err := Parse(strings.NewReader(conf))
	require.NoError(t, err)
	require.Equal(t, &File{
		Sections: []Section{
			{
				{Key: "Version", Value: "65537"},
				{Key: "VideoMode", Value: "1024 768 16", Comment: "comment"},
			},
			{
				{Key: "MousePickup", Value: "Left"},
				{Key: "I+M", Value: "ToggleInventory"},
			},
		},
	}, file)

	var buf bytes.Buffer
	err = file.WriteTo(&buf)
	require.NoError(t, err)
	require.Equal(t, conf, buf.String())

	buf.Reset()
	file.Sections[0].Set("ServerName", "Test")
	err = file.WriteTo(&buf)
	require.NoError(t, err)
	require.Equal(t, `Version = 65537
# comment
VideoMode = 1024 768 16
ServerName = Test
---
MousePickup = Left
I+M = ToggleInventory
---
`, buf.String())

}
