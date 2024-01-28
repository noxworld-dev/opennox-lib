package noxnet

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/shoenig/test/must"

	"github.com/noxworld-dev/opennox-lib/binenc"
	"github.com/noxworld-dev/opennox-lib/noxnet/xfer"
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
				Unk0:         1,
				Unk1:         32,
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
				MapName: binenc.String{
					Value: "So_Druid.map",
					Junk:  []byte{0x9, 0x0, 0x80, 0x96, 0x98, 0x0, 0x0, 0x0, 0x0, 0x0, 0x57, 0xd2, 0x30, 0x14, 0x1, 0x0, 0x0, 0x0, 0x13},
				},
				CRC: 0x6765031d,
				T:   12561,
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
		{
			name: "text msg global",
			packet: &MsgText{
				NetCode: 935,
				Flags:   TextUTF8,
				PosX:    1472,
				PosY:    2370,
				Size:    13,
				Dur:     0,
				Data:    []byte("hello global\x00"),
			},
		},
		{
			name: "text msg team",
			packet: &MsgText{
				NetCode: 935,
				Flags:   TextUTF8 | TextTeam,
				PosX:    1472,
				PosY:    2370,
				Size:    8,
				Dur:     0,
				Data:    []byte("hi team\x00"),
			},
		},
		{
			name: "text msg payload",
			packet: &MsgText{
				NetCode: 0,
				Flags:   TextUTF8 | TextExt,
				PosX:    0,
				PosY:    0,
				Size:    5,
				Dur:     0,
				Data:    []byte("\x001234"),
			},
		},
		{
			name: "text msg payload 16",
			packet: &MsgText{
				NetCode: 0,
				Flags:   TextExt,
				PosX:    0,
				PosY:    0,
				Size:    5,
				Dur:     0,
				Data:    []byte("\x00\x0012345678"),
			},
		},
		{
			name:   "fade begin",
			packet: &MsgFadeBegin{Out: 1, Menu: 0},
		},
		{
			name:   "fx jiggle",
			packet: &MsgFxJiggle{Val: 17},
		},
		{
			name: "map send start",
			packet: &MsgMapSendStart{
				Unk1:    [3]byte{0, 0, 0},
				MapSize: 208134,
				MapName: binenc.String{Value: "_noxtest.map"},
			},
		},
		{
			name: "map send packet",
			packet: &MsgMapSendPacket{
				Unk:   0,
				Block: 12,
				Data:  []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			},
		},
		{
			name: "stat mult",
			packet: &MsgStatMult{
				Health:   1.1,
				Mana:     1.2,
				Strength: 1.3,
				Speed:    1.4,
			},
		},
		{
			name: "xfer start motd",
			packet: &MsgXfer{&xfer.MsgStart{
				Act:    1,
				Unk1:   0,
				Size:   376,
				Type:   binenc.String{Value: "MOTD"},
				SendID: 0,
				Unk5:   [3]byte{0, 0, 0},
			}},
		},
		{
			name: "xfer accept",
			packet: &MsgXfer{&xfer.MsgAccept{
				RecvID: 0,
				SendID: 0,
			}},
		},
		{
			name: "xfer data motd",
			packet: &MsgXfer{&xfer.MsgData{
				Token:  0,
				RecvID: 0,
				Chunk:  1,
				Data:   []byte("\r\nWelcome to Nox multiplayer!\r\nVisit www.westwood.com for the latest news and updates.\r\n\r\n--------------\r\n\r\nIf you are hosting a game, select a game type and a map \r\nfrom the menu to the right, then click \"GO!\".\r\n\r\n\r\nTo close this message window, click the \"OK\" button.\r\n\r\n\r\n(You can customize this message by editing the file \r\n'motd.txt' found in your Nox game directory)\r\n\x00"),
			}},
		},
		{
			name: "xfer ack",
			packet: &MsgXfer{&xfer.MsgAck{
				Token:  0,
				RecvID: 0,
				Chunk:  1,
			}},
		},
		{
			name: "xfer close",
			packet: &MsgXfer{&xfer.MsgDone{
				RecvID: 0,
			}},
		},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			fname := filepath.Join("testdata", strings.ReplaceAll(c.name, " ", "_")+".dat")
			data, err := os.ReadFile(fname)
			if errors.Is(err, fs.ErrNotExist) {
				data, err = AppendPacket(nil, c.packet)
				must.NoError(t, err)
				err = os.WriteFile(fname, data, 0644)
				must.NoError(t, err)
			}
			must.NoError(t, err)
			p, n, err := DecodeAnyPacket(!c.client, data)
			must.NoError(t, err)
			must.Eq(t, c.packet, p)
			must.EqOp(t, len(data), n)
			buf, err := AppendPacket(nil, p)
			must.NoError(t, err)
			must.Eq(t, data, buf)
			n, err = DecodePacket(data, p)
			must.NoError(t, err)
			must.EqOp(t, len(data), n)
		})
	}
}
