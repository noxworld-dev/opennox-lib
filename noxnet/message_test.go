package noxnet

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodePacket(t *testing.T) {
	var cases = []struct {
		name   string
		packet Message
		client bool
	}{
		{
			name: "server info",
			packet: &MsgServerInfo{
				PlayersCur: 1,
				PlayersMax: 32,
				Unk2:       [5]byte{0x0f, 0x0f, 0xff, 0xff, 0xff},
				MapName:    "BluDeath",
				Status1:    0x02,
				Status2:    0x00,
				Unk19:      [7]byte{0x00, 0x55, 0x00, 0x9a, 0x03, 0x01, 0x00},
				Flags:      0x2107,
				Unk27:      [2]byte{0x03, 0x10},
				Unk29:      [8]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
				Unk37:      [4]byte{0xc0, 0x00, 0xd4, 0x00},
				Token:      0x12345678,
				Unk45:      [20]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xef, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
				Unk65:      [4]byte{0x50, 0xec, 0x98, 0x06},
				ServerName: "User Game",
			},
		},
		{
			name: "server join",
			packet: &MsgServerJoin{
				PlayerName: "Игрок",
				Serial:     "1234567890123456789012",
				Version:    0x1039a,
			},
		},
		{
			name: "server accept",
			packet: &MsgServerAccept{
				Unk0:   [2]byte{0x00, 0x01},
				ID:     1,
				XorKey: 0x9e,
			},
		},
		{
			name:   "client accept",
			client: true,
			packet: &MsgClientAccept{
				Unk0:         [2]byte{0x1, 0x20},
				PlayerName:   "Denn",
				PlayerClass:  1,
				Unk70:        [29]byte{0x73, 0x4d, 0x22, 0xda, 0x9a, 0x6e, 0xda, 0x9a, 0x6e, 0xda, 0x9a, 0x6e, 0xda, 0x9a, 0x6e, 0x1f, 0x1f, 0x8, 0x17, 0x6, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
				ScreenWidth:  1024,
				ScreenHeight: 768,
				Serial:       "1234567890123456789012",
				Unk129:       [26]byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			},
		},
		{
			name: "timestamp",
			packet: &MsgTimestamp{
				T: 12561,
			},
		},
		{
			name: "timestamp full",
			packet: &MsgFullTimestamp{
				T: 12561,
			},
		},
		{
			name: "join data",
			packet: &MsgJoinData{
				NetCode: 96,
				Unk2:    0,
			},
		},
		{
			name: "use map",
			packet: &MsgUseMap{
				MapName:     "So_Druid.map",
				CRC:         0x6765031d,
				T:           12561,
				MapNameJunk: []byte{0x9, 0x0, 0x80, 0x96, 0x98, 0x0, 0x0, 0x0, 0x0, 0x0, 0x57, 0xd2, 0x30, 0x14, 0x1, 0x0, 0x0, 0x0, 0x13},
			},
		},
		{
			name: "player input",
			packet: &MsgPlayerInput{
				Inputs: []PlayerInput{
					&PlayerInput1{Code: CCOrientation, Val: 130},
				},
			},
		},
		{
			name: "player mouse",
			packet: &MsgMouse{
				X: 3103,
				Y: 2963,
			},
		},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			fname := filepath.Join("testdata", strings.ReplaceAll(c.name, " ", "_")+".dat")
			data, err := os.ReadFile(fname)
			if errors.Is(err, fs.ErrNotExist) {
				data, err = AppendPacket(nil, c.packet)
				require.NoError(t, err)
				err = os.WriteFile(fname, data, 0644)
				require.NoError(t, err)
			}
			require.NoError(t, err)
			p, n, err := DecodeAnyPacket(!c.client, data)
			require.NoError(t, err)
			require.Equal(t, c.packet, p)
			require.Equal(t, int(len(data)), int(n))
			buf, err := AppendPacket(nil, p)
			require.NoError(t, err)
			require.Equal(t, data, buf)
		})
	}
}
